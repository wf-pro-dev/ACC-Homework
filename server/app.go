package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/williamfotso/acc/assignment"
	"github.com/williamfotso/acc/assignment/notion"
	"github.com/williamfotso/acc/assignment/notion/types"
	"github.com/williamfotso/acc/course"
	"github.com/williamfotso/acc/crud"
)

func StartServer() {

	http.HandleFunc("/notion-webhooks/test", testHandler)
	http.HandleFunc("/notion-webhooks", notionWebhookHandler)
	log.Println("Server listening on :8080...") // Changed from fmt
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s test: %s", time.Now().Format(time.Stamp), r.URL.Path)

	var payload struct {
		Test string `json:"test"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("%s test: %s", time.Now().Format(time.Stamp), payload.Test)
}

func notionWebhookHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Verify it's a POST request
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decode the payload
	var payload NotionWebhookPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Printf("%s Error decoding payload: %s", time.Now().Format(time.Stamp), err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// 3. Get the page id
	page_id := payload.Entity.Id

	// 4. Get the author of the change
	var author_type string
	if len(payload.Authors) > 0 {
		author_type = payload.Authors[0].Type
		if author_type == "bot" {
			log.Printf("%s Author is a bot, skipping", time.Now().Format(time.Stamp))
			return
		}
	}

	// 5. Get the database
	db, err := crud.GetDB()
	if err != nil {
		log.Printf("%s Error getting database: %s", time.Now().Format(time.Stamp), err)
		http.Error(w, "Error getting database", http.StatusInternalServerError)
		return
	}

	// 6. Loop through the properties
	for _, property_id := range payload.Data.Properties {

		// Get the timestamp
		timestamp := time.Now().Format(time.Stamp)

		// 6. Get the updated property
		property, err := notion.GetPageProperties(page_id, property_id)
		if err != nil {
			log.Printf("%s Fetching properties: %s", timestamp, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// 8. Get the column name from the property id
		column := types.COLUMNS[property_id]

		// 9. Get the new value from the property
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
					log.Printf("%s Error getting course: %s", timestamp, err)
					http.Error(w, "Error getting course", http.StatusInternalServerError)
					return
				}

				value = course.Code

			} else {
				value = ""
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
		}

		// Log the update
		log.Printf("%s page_id %s property_id %s column %s value %s",
			timestamp, page_id, property_id, column, value)

		if value != "" {
			// Update the assignment in the database
			assignment := assignment.GetAssignmentsbyNotionID(page_id, db)
			if assignment == nil {
				log.Printf("%s Error getting assignment: %s", timestamp, err)
				http.Error(w, "Error getting assignment", http.StatusInternalServerError)
				return
			}

			err = crud.PutHanlder(assignment.GetID(db), column, "assignements", value, db)
			if err != nil {
				log.Printf("%s Error updating assignment in database: %s", timestamp, err)
				http.Error(w, "Error updating assignment in database", http.StatusInternalServerError)
				return
			}

		}
	}

	fmt.Print("\n\n")
}
