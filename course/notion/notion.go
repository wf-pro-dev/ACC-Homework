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

	"github.com/williamfotso/acc/course/notion/types"
)

var NOTION_API_KEY = os.Getenv("NOTION_API_KEY")

const BASE_URL = "https://api.notion.com/v1"
const DATABASE_ID = "17e40a21a7e381129891fdcdaaa5dbec"

// AddAssignmentToNotion adds an assignment to Notion efficiently
func AddCourseToNotion(course map[string]string) (string, error) {

	// Create a single rich text object for reuse
	createRichTextObj := func() types.RichText {
		text := types.RichText{
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
		return text
	}

	// Create the request with strongly typed fields
	req := types.PageRequest{}
	req.Parent.Type = "database_id"
	req.Parent.DatabaseID = DATABASE_ID

	properties := types.Properties{
		Name: types.Title{
			Type: "title",
		},
		Code: types.Text{
			Type: "rich_text",
		},
		RoomNumber: types.Text{
			Type: "rich_text",
		},
		Duration: types.Text{
			Type: "rich_text",
		},
	}

	req.Properties = &properties

	// Set name
	name_obj := types.Title{
		Type: "title",
	}
	nameText := createRichTextObj()
	nameText.Text.Content = course["name"]
	nameText.PlainText = course["name"]
	name_obj.Title = []types.RichText{nameText}
	req.Properties.Name = name_obj

	// Set code
	code_obj := types.Text{
		Type: "rich_text",
	}
	codeText := createRichTextObj()
	codeText.Text.Content = course["code"]
	codeText.PlainText = course["code"]
	code_obj.RichText = []types.RichText{codeText}
	req.Properties.Code = code_obj

	// Set name
	room_number_obj := types.Text{
		Type: "rich_text",
	}
	roomNumberText := createRichTextObj()
	roomNumberText.Text.Content = course["room_number"]
	roomNumberText.PlainText = course["room_number"]
	room_number_obj.RichText = []types.RichText{roomNumberText}
	req.Properties.RoomNumber = room_number_obj

	// Set duration
	duration_obj := types.Text{
		Type: "rich_text",
	}
	durationText := createRichTextObj()
	durationText.Text.Content = course["duration"]
	durationText.PlainText = course["duration"]
	duration_obj.RichText = []types.RichText{durationText}
	req.Properties.Duration = duration_obj

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Set up the HTTP client with timeouts for better reliability
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request with context for potential cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.notion.com/v1/pages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+NOTION_API_KEY)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Notion-Version", "2022-06-28")

	// Send request
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("notion API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	type NotionResponse struct {
		ID string `json:"id"`
	}

	var notionResp NotionResponse
	if err := json.Unmarshal(respBody, &notionResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return notionResp.ID, nil
}
