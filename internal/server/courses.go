package server

import (
	"encoding/json"
	"net/http"
	"fmt"

	"github.com/williamfotso/acc/internal/core/models/course"
	"gorm.io/gorm"
)
func GetCourseHandler(w http.ResponseWriter, r *http.Request) {
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

	var courses []course.Course
	if err := db.Where("user_id = ?", userID).Find(&courses).Error ; err != nil {
		PrintERROR(w, http.StatusInternalServerError, fmt.Sprintf("Error getting assignment for user id = %d : %s",userID, err))
                return
	}
	
	var coursesMap []map[string]string
	for _, a := range courses {
		coursesMap = append(coursesMap, a.ToMap())
	}

	w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
                "message": "User's Assignments retrieved successfully",
                "courses":    coursesMap,
        })
}

func CreateCourseHandler(w http.ResponseWriter, r *http.Request) {

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
        	Name		string    `json:"name"`
        	Code		string    `json:"code"`
       		Duration	string    `json:"duration"`
        	RoomNumber	string    `json:"room_number"`
    	}

    	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        	PrintERROR(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
        	return
    	}

    	// Validate all required fields
    	if input.Name == "" || input.Code == "" || input.Duration == "" || input.RoomNumber == "" {
        	PrintERROR(w, http.StatusBadRequest, "Missing required fields")
        	return
    	}


    	cVal := course.Course{
        	UserID:     	userID,
		Name:		input.Name,	
		Code:		input.Code,	
		Duration:	input.Duration,
		RoomNumber:	input.RoomNumber,	
    	}	

	err := tx.Create(&cVal).Error
        if err != nil {
                PrintERROR(w, http.StatusConflict, fmt.Sprintf("Error creating assignment in database",err))
                return
        }

	aObj := &cVal

	c, err := course.Get_Course_byId(aObj.ID, tx)
	if err != nil {
                PrintERROR(w,http.StatusInternalServerError,fmt.Sprintf("failed to getting assignment: %s", err))
                return
        }

	notion_id, err := c.Add_Notion()
	if err != nil {
		tx.Rollback()
                PrintERROR(w, http.StatusInternalServerError, fmt.Sprintf("Error creating assignment in notion",err))
                return
        }

	c.NotionID = notion_id
	err = tx.Save(&c).Error
	if err != nil {
		tx.Rollback()
                PrintERROR(w, http.StatusInternalServerError, fmt.Sprintf("Error updating new assignment",err))
                return
        }


	// Convert to map safely
        courseMap := c.ToMap()
        if courseMap == nil {
		tx.Rollback()
                PrintERROR(w, http.StatusInternalServerError, "Failed to process assignment data")
                return
        }

	tx.Commit()

        // Return response
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
                "message": "Assignment created successfully",
                "course":  courseMap,
        })

}
/*
func UpdateCourseHandler(w http.ResponseWriter, r *http.Request) {
	
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

	c, err := course.Get_Course_byId(uint(int_id), db)
	if err != nil {
                PrintERROR(w,http.StatusInternalServerError,fmt.Sprintf("failed to getting assignment: %s", err))
                return
        }

	if err := db.Exec(fmt.Sprintf("UPDATE courses SET %s = ?, updated_at = ? WHERE id = ?",updateData.Column), 
        	 updateData.Value, time.Now().Format(time.RFC3339), c.ID).Error; err != nil {
    		 PrintERROR(w, http.StatusInternalServerError,
                                        fmt.Sprintf("Error updating assignment in database: %s", err))
		return
	}


	err = c.Update_Notion(updateData.Column,updateData.Value)
        if err != nil {
                PrintERROR(w, http.StatusInternalServerError, fmt.Sprintf("Error updating assignment in notion",err))
                return
        }

}*/

