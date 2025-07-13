package server

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/types"
)

func WebhookDeleteHandler(w http.ResponseWriter, r *http.Request, payload types.NotionWebhookPayload) {

	
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

	a, err := assignment.Get_Assignment_byNotionID(payload.Entity.Id, tx)

	if err != nil {
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error getting assignment: %s", err))
		return
	}

	err = tx.Delete(&a).Error 
	if err != nil {
		tx.Rollback()
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error deleting assignment: %s", err))
		return
	}

	/*notification_id := fmt.Sprintf("%s-deleted", assignment.NotionID)
	title := fmt.Sprintf("%s: %s", assignment.CourseCode, assignment.Title)
	subtitle := fmt.Sprintf("Deleted at %s", time.Now().Format(time.Stamp))
	message := "Assignment deleted"

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

	tx.Commit()

	PrintLog(fmt.Sprintf("Assignment deleted: %s %s", payload.Entity.Id, a.Title))
}
