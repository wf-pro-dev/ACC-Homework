package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/williamfotso/acc/assignment/notion"
	"github.com/williamfotso/acc/assignment/notion/types"
	"github.com/williamfotso/acc/course"
	"github.com/williamfotso/acc/crud"
	"github.com/williamfotso/acc/notifier"
)

func WebhookCreateHandler(w http.ResponseWriter, r *http.Request, payload NotionWebhookPayload) {

	// 1. Get the page id
	page_id := payload.Entity.Id

	// 2. Get the page properties
	new_page, err := notion.GetPage(page_id)
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

	db, err := crud.GetDB()
	if err != nil {
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error getting database: %s", err))
		return
	}

	var assignment = make(map[string]string)
	var properties = pageResp.Properties

	assignment["notion_id"] = page_id
	assignment["title"] = properties.AssignmentName.Title[0].PlainText

	course_notion_id := properties.Courses.Relation[0].ID
	course_code := course.GetCoursebyNotionID(course_notion_id, db).Code
	assignment["course_code"] = course_code

	assignment["type"] = properties.Type.Select["name"]
	assignment["deadline"] = properties.Deadline.Date.Start
	assignment["todo"] = properties.TODO.RichText[0].PlainText
	assignment["link"] = properties.Link.URL

	switch properties.Status.Status.Name {
	case "Done":
		assignment["status"] = "done"
	case "In progress":
		assignment["status"] = "start"
	default:
		assignment["status"] = "default"
	}

	err = crud.PostHandler(assignment, "assignements", db)
	if err != nil {
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error creating assignment: %s", err))
		return
	}

	notification_id := assignment["notion_id"]
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
	time.Sleep(5 * time.Second) // Wait for the notification to be sent

	err = notifier.UseNotifier([]string{"-remove", notification_id})
	if err != nil {
		PrintERROR(w, http.StatusInternalServerError,
			fmt.Sprintf("Error removing notification: %s", err))
	}

	PrintLog(fmt.Sprintf("Assignment created: %+v %s", assignment["notion_id"], assignment["title"]))
}
