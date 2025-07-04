package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/williamfotso/acc/assignment"
	"github.com/williamfotso/acc/crud"
	"github.com/williamfotso/acc/notifier"
)

func WebhookDeleteHandler(w http.ResponseWriter, r *http.Request, payload NotionWebhookPayload) {

	db, err := crud.GetDB()
	if err != nil {
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error getting database: %s", err))
		return
	}

	assignment := assignment.GetAssignmentsbyNotionID(payload.Entity.Id, db)
	if assignment == nil {
		err = fmt.Errorf("assignment not found")
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error getting assignment: %s", err))
		return
	}

	err = crud.DeleteHandler("assignements", "notion_id", payload.Entity.Id, db)
	if err != nil {
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error deleting assignment: %s", err))
		return
	}

	notification_id := assignment.NotionID
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

	PrintLog(fmt.Sprintf("Assignment deleted: %s %s", payload.Entity.Id, assignment.Title))
}
