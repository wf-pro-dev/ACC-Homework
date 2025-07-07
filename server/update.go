package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/williamfotso/acc/assignment"
	"github.com/williamfotso/acc/assignment/notion"
	"github.com/williamfotso/acc/assignment/notion/types"
	"github.com/williamfotso/acc/course"
	"github.com/williamfotso/acc/database"
	"github.com/williamfotso/acc/notifier"
)

func WebhookUpdateHandler(w http.ResponseWriter, r *http.Request, payload NotionWebhookPayload) {

	// 3. Get the page id
	page_id := payload.Entity.Id

	// 5. Get the database
	db, err := database.GetDB()
	if err != nil {
		PrintERROR(w, http.StatusInternalServerError, fmt.Sprintf("Error getting database: %s", err))
		return
	}

	// 6. Loop through the properties
	for _, property_id := range payload.Data.Properties {

		// 7. Get the updated property
		property, err := notion.GetPageProperties(page_id, property_id)
		if err != nil {
			PrintERROR(w, http.StatusInternalServerError, fmt.Sprintf("Fetching properties: %s", err))
			return
		}

		// For debugging
		// var property_map map[string]interface{}
		// json.Unmarshal(property, &property_map)

		// 8. Get the column name from the property id
		column := types.COLUMNS[property_id]
		if column == "" {
			PrintERROR(w, http.StatusInternalServerError,
				fmt.Sprintf("Column not found for id: %s", property_id))
			return
		}

		// 9. Get the new value from the property
		value := GetValue(w, property, column, db)

		// Log the update
		PrintLog(fmt.Sprintf("page_id %s property_id %s column %s value %s",
			page_id, property_id, column, value))

		if value != "" {
			// Update the assignment in the database
			assignment := assignment.GetAssignmentsbyNotionID(page_id, db)
			if assignment == nil {
				PrintERROR(w, http.StatusInternalServerError,
					fmt.Sprintf("Error getting assignment: %s", err))
			}

			err = database.PutHanlder(assignment.GetID(db), column, "assignements", value, db)
			if err != nil {
				PrintERROR(w, http.StatusInternalServerError,
					fmt.Sprintf("Error updating assignment in database: %s", err))
			}

			notification_id := fmt.Sprintf("%s-%s-%s", assignment.GetNotionID(), column, value)
			title := fmt.Sprintf("%s: %s", assignment.GetCourseCode(), assignment.GetTitle())
			subtitle := fmt.Sprintf("Updated at %s", time.Now().Format(time.Stamp))
			message := fmt.Sprintf("%s is now %s", column, value)

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
			}
		}
	}

	fmt.Print("\n\n")
}

func GetValue(w http.ResponseWriter, property []byte, column string, db *sql.DB) string {
	var value string
	switch column {

	case "course_code":
		var coursesType struct {
			Courses []struct {
				Relation struct {
					ID string `json:"id"`
				} `json:"relation"`
			} `json:"results"`
		}
		json.Unmarshal(property, &coursesType)

		if len(coursesType.Courses) > 0 {
			course := course.GetCoursebyNotionID(coursesType.Courses[0].Relation.ID, db)

			if course == nil {
				value = ""
				err := fmt.Errorf("course not found")
				PrintERROR(w, http.StatusInternalServerError,
					fmt.Sprintf("Error getting course: %s", err))
			} else {
				value = course.Code
			}
		}

	case "deadline":
		var dateType struct {
			Date struct {
				Start string `json:"start"`
			} `json:"date"`
		}

		json.Unmarshal(property, &dateType)
		value = dateType.Date.Start

	case "link":
		var linkType struct {
			URL string `json:"url"`
		}
		json.Unmarshal(property, &linkType)
		value = linkType.URL

	case "todo":
		var todoType struct {
			Results []struct {
				RichText struct {
					PlainText string `json:"plain_text"`
				} `json:"rich_text"`
			} `json:"results"`
		}

		json.Unmarshal(property, &todoType)
		value = todoType.Results[0].RichText.PlainText

	case "title":
		var titleType struct {
			Results []struct {
				Title struct {
					PlainText string `json:"plain_text"`
				} `json:"title"`
			} `json:"results"`
		}
		json.Unmarshal(property, &titleType)
		value = titleType.Results[0].Title.PlainText

	case "type":

		var selectType struct {
			Select map[string]string `json:"select"`
		}

		json.Unmarshal(property, &selectType)
		value = selectType.Select["name"]

	case "status":
		var statusType struct {
			Status struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"status"`
		}
		json.Unmarshal(property, &statusType)
		switch statusType.Status.Name {
		case "Done":
			value = "done"
		case "In progress":
			value = "start"
		default:
			value = "default"
		}
	default:
		value = ""
	}

	return value

}
