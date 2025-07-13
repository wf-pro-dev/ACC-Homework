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
