package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/williamfotso/acc/assignment/notion/types"
)

const BASE_URL = "https://api.notion.com/v1"

func sendRequest(req interface{}, method, url, notion_id string) (respBody []byte, err error) {

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var final_url string
	if notion_id != "" {
		final_url = fmt.Sprintf("%s/%s/%s", BASE_URL, url, notion_id)
	} else {
		final_url = fmt.Sprintf("%s/%s", BASE_URL, url)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, final_url, bytes.NewBuffer(jsonData))

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_KEY"))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	// Send request
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("notion API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// AddAssignmentToNotion adds an assignment to Notion efficiently
func AddAssignmentToNotion(assign, type_info, course_info map[string]string) (string, error) {

	// Create a single rich text object for reuse
	richTextObj := types.RichText{
		Type: "text",
		Text: &types.TextContent{
			Content: "",
			Link:    nil,
		},
		Annotations: &types.TextAnnotation{
			Bold:          false,
			Italic:        false,
			Strikethrough: false,
			Underline:     false,
			Code:          false,
			Color:         "default",
		},
		PlainText: "",
		Href:      nil,
	}

	// Create the request with strongly typed fields
	req := types.PageRequest{}
	req.Parent.Type = "database_id"
	req.Parent.DatabaseID = "17e40a21a7e381a18a85ccc380a0beec"
	// Set deadline
	req.Properties = &types.Properties{
		Deadline: types.Deadline{
			ID:   "_UjC",
			Type: "date",
			Date: &types.DateObject{
				Start: assign["deadline"], // 2025-06-05T00:00:00.000Z
			},
		},
		Courses: types.Courses{
			ID:   "w%3FC%3B",
			Type: "relation",
			Relation: []types.Relation{
				{
					ID: course_info["notion_id"],
				},
			},
		},
		Type: types.Type{
			ID:     "S~Ce",
			Type:   "select",
			Select: type_info,
		},
		Status: types.Status{
			ID:   "%5Bm%5Cs",
			Type: "status",
			Status: &types.StatusObject{
				ID:    "3aa77cf8-c39e-4c7b-b7d2-ab15ae43ff23",
				Name:  "Not started",
				Color: "default",
			},
		},
		TODO: types.TODO{
			ID:   "%5DJfC",
			Type: "rich_text",
		},
		AssignmentName: types.AssignmentName{
			ID:   "title",
			Type: "title",
		},
	}
	// Set TODO
	todo_obj := types.TODO{
		ID:   "%5DJfC",
		Type: "rich_text",
	}
	todoText := richTextObj
	todoText.Text.Content = assign["todo"]
	todoText.PlainText = assign["todo"]
	todo_obj.RichText = []types.RichText{todoText}
	req.Properties.TODO = todo_obj

	// Set title
	assignment_name_obj := types.AssignmentName{
		ID:   "title",
		Type: "title",
	}
	titleText := richTextObj
	titleText.Text.Content = assign["title"]
	titleText.PlainText = assign["title"]
	assignment_name_obj.Title = []types.RichText{titleText}
	req.Properties.AssignmentName = assignment_name_obj

	resp, err := sendRequest(req, "POST", "pages", "")
	if err != nil {
		return "", err
	}

	// Parse response
	type NotionResponse struct {
		ID string `json:"id"`
	}

	var notionResp NotionResponse
	if err := json.Unmarshal(resp, &notionResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return notionResp.ID, nil
}

func UpdateAssignementToNotion(assign map[string]string, col string, value string, type_info map[string]string) (err error) {
	var req interface{}

	switch col {

	case "course_code":
		courseReq := types.UpdateCourseCodeRequest{}
		courseReq.Properties = types.PropertiesWithRequiredCourseCode{}
		courseReq.Properties.Courses = types.Courses{
			ID:   "w%3FC%3B",
			Type: "relation",
			Relation: []types.Relation{
				{
					ID: value,
				},
			},
		}
		req = courseReq

	case "deadline":
		deadlineReq := types.UpdateDeadlineRequest{}
		deadlineReq.Properties = types.PropertiesWithRequiredDeadline{}

		dateObj := types.DateObject{
			Start: value,
		}

		deadlineReq.Properties.Deadline = types.Deadline{
			ID:   "_UjC",
			Type: "date",
			Date: &dateObj,
		}

		req = deadlineReq

	case "link":
		linkReq := types.UpdateLinkRequest{}
		linkReq.Properties = types.PropertiesWithRequiredLink{}

		linkReq.Properties.Link = types.Link{
			ID:   "jgPD",
			Type: "url",
			URL:  value,
		}

		req = linkReq

	case "title":
		titleReq := types.UpdateTitleRequest{}
		titleReq.Properties = types.PropertiesWithRequiredName{}

		richTextObj := types.RichText{
			Type: "text",
			Text: &types.TextContent{
				Content: value,
				Link:    nil,
			},
			Annotations: &types.TextAnnotation{
				Bold:          false,
				Italic:        false,
				Strikethrough: false,
				Underline:     false,
				Code:          false,
				Color:         "default",
			},
			PlainText: value,
			Href:      nil,
		}

		titleReq.Properties.AssignmentName = types.AssignmentName{
			Title: []types.RichText{richTextObj},
		}

		req = titleReq

	case "todo":
		todoReq := types.UpdateTODORequest{}
		todoReq.Properties = types.PropertiesWithRequiredTODO{}

		richTextObj := types.RichText{
			Type: "text",
			Text: &types.TextContent{
				Content: value,
				Link:    nil,
			},
			Annotations: &types.TextAnnotation{
				Bold:          false,
				Italic:        false,
				Strikethrough: false,
				Underline:     false,
				Code:          false,
				Color:         "default",
			},
			PlainText: value,
			Href:      nil,
		}

		todoReq.Properties.TODO = types.TODO{
			ID:       "%5DJfC",
			Type:     "rich_text",
			RichText: []types.RichText{richTextObj},
		}

		req = todoReq

	case "type":

		typeReq := types.UpdateTypeRequest{}
		typeReq.Properties = types.PropertiesWithRequiredType{}

		typeReq.Properties.Type = types.Type{
			ID:     "S~Ce",
			Type:   "select",
			Select: type_info,
		}

		req = typeReq

	}

	if req == nil {
		return fmt.Errorf("invalid column type: %s", col)
	}

	_, err = sendRequest(req, "PATCH", "pages", assign["notion_id"])

	return err
}

func DeleteAssignementFromNotion(assign map[string]string) (err error) {

	req := types.DeletionRequest{}
	req.Archived = true

	_, err = sendRequest(req, "PATCH", "pages", assign["notion_id"])

	return err
}
