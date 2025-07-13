package server

import (
	"encoding/json"
	"net/http"
	"fmt"
	"time"
	"strconv"

	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/core/models"
	"gorm.io/gorm"
)

func CreateAssignmentHandler(w http.ResponseWriter, r *http.Request) {

	userIDVal := r.Context().Value("user_id")
        if userIDVal == nil {
                PrintERROR(w, http.StatusUnauthorized, "User ID not found in context")
                return
        }
	
	userID, ok := userIDVal.(uint)
        if !ok {
                PrintERROR(w, http.StatusUnauthorized, "Invalid user ID format")
                return
        }

	dbVal := r.Context().Value("db")
	if dbVal == nil {
		PrintERROR(w, http.StatusInternalServerError, "Database connection not found")
		return
	}

	db, ok := dbVal.(*gorm.DB)
	if !ok {
		PrintERROR(w, http.StatusInternalServerError, "Invalid database connection")
		return
	}

	tx := db.Begin()
	defer func() {
    		if r := recover(); r != nil {
        		tx.Rollback()
    		}
	}()

	
	var input struct {
        	CourseCode string    `json:"course_code"`
        	Title      string    `json:"title"`
       		TypeName   string    `json:"type_name"`
        	Deadline   string    `json:"deadline"`
        	Todo       string    `json:"todo"`
        	StatusName string    `json:"status_name"`
    	}

    	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        	PrintERROR(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
        	return
    	}

    	// Validate all required fields
    	if input.CourseCode == "" || input.Title == "" || input.TypeName == "" || input.Deadline == "" {
        	PrintERROR(w, http.StatusBadRequest, "Missing required fields")
        	return
    	}

    	deadline, err := time.Parse(time.RFC3339, input.Deadline)
    	if err != nil {
        	PrintERROR(w, http.StatusBadRequest, "Invalid deadline format")
        	return
    	}

    	aVal := assignment.Assignment{
        	UserID:     userID,
        	CourseCode: input.CourseCode,
        	Title:      input.Title,
        	TypeName:   input.TypeName,
        	Deadline:   deadline,
        	Todo:       input.Todo,
        	StatusName: input.StatusName,
        	Link:       "https://acconline.austincc.edu/ultra/stream",
    	}	

	result := tx.Create(&aVal)
        if result.Error != nil {
                PrintERROR(w, http.StatusConflict, fmt.Sprintf("Error creating assignment in database",err))
                return
        }

	aObj := &aVal

	a, err := assignment.Get_Assignment_byId(aObj.ID, tx)
	if err != nil {
                PrintERROR(w,http.StatusInternalServerError,fmt.Sprintf("failed to getting assignment: %s", err))
                return
        }

	notion_id, err := a.Add_Notion()
	if err != nil {
		tx.Rollback()
                PrintERROR(w, http.StatusInternalServerError, fmt.Sprintf("Error creating assignment in notion",err))
                return
        }

	a.NotionID = notion_id
	err = tx.Save(&a).Error
	if err != nil {
		tx.Rollback()
                PrintERROR(w, http.StatusInternalServerError, fmt.Sprintf("Error updating new assignment",err))
                return
        }


	// Convert to map safely
        assignmentMap := a.ToMap()
        if assignmentMap == nil {
		tx.Rollback()
                PrintERROR(w, http.StatusInternalServerError, "Failed to process assignment data")
                return
        }

	tx.Commit()

        // Return response
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
                "message": "User retrieved successfully",
                "assignment":    assignmentMap,
        })

}
func UpdateAssignmentHandler(w http.ResponseWriter, r *http.Request) {
	
	dbVal := r.Context().Value("db")
        if dbVal == nil {
                PrintERROR(w, http.StatusInternalServerError, "Database connection not found")
                return
        }


        db, ok := dbVal.(*gorm.DB)
        if !ok {
                PrintERROR(w, http.StatusInternalServerError, "Invalid database connection")
                return
        }


	var updateData struct {
		ID 	string		`json:"id"`
        	Value	string		`json:"value"`
        	Column	string		`json:"column`
    	}

	err := json.NewDecoder(r.Body).Decode(&updateData)
        if err != nil {
                PrintERROR(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body %s",err))
                return
        }

	int_id, err := strconv.Atoi(updateData.ID)
	if err != nil {
		PrintERROR(w,http.StatusInternalServerError,fmt.Sprintf("failed to convert assignment ID to int: %s", err))
                return
        }

	a, err := assignment.Get_Assignment_byId(uint(int_id), db)
	if err != nil {
                PrintERROR(w,http.StatusInternalServerError,fmt.Sprintf("failed to getting assignment: %s", err))
                return
        }

	if err := db.Exec(fmt.Sprintf("UPDATE assignments SET %s = ?, updated_at = ? WHERE id = ?",updateData.Column), 
        	 updateData.Value, time.Now().Format(time.RFC3339), a.ID).Error; err != nil {
    		 PrintERROR(w, http.StatusInternalServerError,
                                        fmt.Sprintf("Error updating assignment in database: %s", err))
		return
	}

	value := updateData.Value

	if updateData.Column == "course_code" {
		value = a.Course.NotionID
	}

	var obj map[string]string

	if updateData.Column == "status_name" {
		var status = models.Get_AssignmentStatus_byName(value, db)
		obj = status.ToMap()
	} else {
		var t = models.Get_AssignmentType_byName(value, db)
		obj = t.ToMap()
	}

	err = a.Update_Notion(updateData.Column,value , obj)
        if err != nil {
                PrintERROR(w, http.StatusInternalServerError, fmt.Sprintf("Error updating assignment in notion",err))
                return
        }

}
