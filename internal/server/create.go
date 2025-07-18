package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"gorm.io/gorm"
	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/core/models/course"
	"github.com/williamfotso/acc/internal/core/models/user"
	"github.com/williamfotso/acc/internal/types"
)

func WebhookCreateHandler(w http.ResponseWriter, r *http.Request, payload types.NotionWebhookPayload, u *user.User) {
	

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

	// 1. Get the page id
	page_id := payload.Entity.Id

	// 2. Get the page properties
	new_page, err := assignment.GetPage(page_id, u.NotionAPIKey)
	if err != nil {
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error getting page properties: %s", err))
		return
	}

	var pageResp types.PageRequest
	if err := json.Unmarshal(new_page, &pageResp); err != nil {
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error unmarshalling page: %s", err))
		return
	}

	var properties = pageResp.Properties

	course_notion_id := properties.Courses.Relation[0].ID
	course_code := course.Get_Course_byNotionID(course_notion_id, db).Code

	deadline, err := time.Parse(time.DateOnly ,properties.Deadline.Date.Start)
	if err != nil {

		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error parsing deadline: %s", err))
		return
	}

	PrintLog(fmt.Sprintf("User id : %d,\n page_id : %s\n ",u.ID, page_id))

	aVal := assignment.Assignment{
                UserID:		u.ID,
                CourseCode:	course_code,
                Title:		properties.AssignmentName.Title[0].PlainText,
                TypeName:	properties.Type.Select["name"],
                Deadline:	deadline,
                Todo:		properties.TODO.RichText[0].PlainText,
                StatusName:	properties.Status.Status.Name,
                Link:		properties.Link.URL,
		NotionID:	page_id}

        result := tx.Create(&aVal)
        if result.Error != nil {
		tx.Rollback()
                PrintERROR(w, http.StatusConflict, fmt.Sprintf("Error creating assignment in database",err))
                return
        }
	
	tx.Commit()
	/*notification_id := fmt.Sprintf("%s-created", assignment["notion_id"])
	title := fmt.Sprintf("%s: %s", assignment["course_code"], assignment["title"])
	subtitle := fmt.Sprintf("Created at %s", time.Now().Format(time.Stamp))
	message := "New assignment created"

	args := []string{
		"-group", notification_id,
		"-title", title,
		"-subtitle", subtitle,
		"-message", message,
		"-sound", "Frog",
		"-timeout", "5", // Notification stays for 30 seconds
	}

	err = notifier.UseNotifier(args)
	if err != nil {
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error sending notification: %s", err))
	}

	time.Sleep(10 * time.Second) // Wait for the notification to be sent

	err = notifier.UseNotifier([]string{"-remove", notification_id})
	if err != nil {
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error removing notification: %s", err))
	}*/

	PrintLog(fmt.Sprintf("Assignment created: %+v %s", aVal.NotionID, aVal.Title))
}
