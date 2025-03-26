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
)

// TypeInfo represents the type of assignment
type TypeInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// CourseInfo represents the course relation
type CourseID struct {
	id string
}

// AddAssignmentToNotion adds an assignment to Notion efficiently
func AddAssignmentToNotion(assign, type_info, course_info map[string]string) (string, error) {

	// Use struct to define the request structure for better type safety and readability
	type TextContent struct {
		Content string      `json:"content"`
		Link    interface{} `json:"link"`
	}

	type TextAnnotation struct {
		Bold          bool   `json:"bold"`
		Italic        bool   `json:"italic"`
		Strikethrough bool   `json:"strikethrough"`
		Underline     bool   `json:"underline"`
		Code          bool   `json:"code"`
		Color         string `json:"color"`
	}

	type RichText struct {
		Type        string         `json:"type"`
		Text        TextContent    `json:"text"`
		Annotations TextAnnotation `json:"annotations"`
		PlainText   string         `json:"plain_text"`
		Href        interface{}    `json:"href"`
	}

	type DateObject struct {
		Start    string      `json:"start"`
		End      interface{} `json:"end"`
		TimeZone interface{} `json:"time_zone"`
	}

	type StatusObject struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	type PageRequest struct {
		Cover  interface{} `json:"cover"`
		Icon   interface{} `json:"icon"`
		Parent struct {
			Type       string `json:"type"`
			DatabaseID string `json:"database_id"`
		} `json:"parent"`
		Archived   bool `json:"archived"`
		InTrash    bool `json:"in_trash"`
		Properties struct {
			Deadline struct {
				ID   string     `json:"id"`
				Type string     `json:"type"`
				Date DateObject `json:"date"`
			} `json:"Deadline"`
			Courses struct {
				ID       string              `json:"id"`
				Type     string              `json:"type"`
				Relation []map[string]string `json:"relation"`
				HasMore  bool                `json:"has_more"`
			} `json:"Courses"`
			Type struct {
				ID     string            `json:"id"`
				Type   string            `json:"type"`
				Select map[string]string `json:"select"`
			} `json:"Type"`
			Status struct {
				ID     string       `json:"id"`
				Type   string       `json:"type"`
				Status StatusObject `json:"status"`
			} `json:"Status"`
			TODO struct {
				ID       string     `json:"id"`
				Type     string     `json:"type"`
				RichText []RichText `json:"rich_text"`
			} `json:"TODO"`
			AssignmentName struct {
				ID    string     `json:"id"`
				Type  string     `json:"type"`
				Title []RichText `json:"title"`
			} `json:"Assignment Name"`
		} `json:"properties"`
	}

	// Create a single rich text object for reuse
	richTextObj := RichText{
		Type: "text",
		Text: TextContent{
			Content: "",
			Link:    nil,
		},
		Annotations: TextAnnotation{
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
	req := PageRequest{}
	req.Parent.Type = "database_id"
	req.Parent.DatabaseID = "17e40a21a7e381a18a85ccc380a0beec"
	// Set deadline
	req.Properties.Deadline.ID = "_UjC"
	req.Properties.Deadline.Type = "date"
	req.Properties.Deadline.Date.Start = assign["deadline"]

	// Set course
	req.Properties.Courses.ID = "w%3FC%3B"
	req.Properties.Courses.Type = "relation"
	req.Properties.Courses.Relation = []map[string]string{course_info}

	// Set type
	req.Properties.Type.ID = "S~Ce"
	req.Properties.Type.Type = "select"
	req.Properties.Type.Select = type_info

	// Set status
	req.Properties.Status.ID = "%5Bm%5Cs"
	req.Properties.Status.Type = "status"
	req.Properties.Status.Status = StatusObject{
		ID:    "3aa77cf8-c39e-4c7b-b7d2-ab15ae43ff23",
		Name:  "Not started",
		Color: "default",
	}

	// Set TODO
	req.Properties.TODO.ID = "%5DJfC"
	req.Properties.TODO.Type = "rich_text"
	todoText := richTextObj
	todoText.Text.Content = assign["todo"]
	todoText.PlainText = assign["todo"]
	req.Properties.TODO.RichText = []RichText{todoText}

	// Set title
	req.Properties.AssignmentName.ID = "title"
	req.Properties.AssignmentName.Type = "title"
	titleText := richTextObj
	titleText.Text.Content = assign["title"]
	titleText.PlainText = assign["title"]
	req.Properties.AssignmentName.Title = []RichText{titleText}

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
	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_KEY"))
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

	// Update the .info file more efficiently
	if err := updateInfoFile(notionResp.ID); err != nil {
		return notionResp.ID, fmt.Errorf("created page but failed to update .info file: %w", err)
	}

	return notionResp.ID, nil
}

// updateInfoFile efficiently updates the .info file with the new page ID
func updateInfoFile(pageID string) error {
	// Try to open the file in append mode first
	file, err := os.OpenFile(".info", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .info file: %w", err)
	}
	defer file.Close()

	// Get file info to check if it's empty
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if fileInfo.Size() == 0 {
		// File is empty, write the page ID
		_, err = file.WriteString(pageID + "\n")
		return err
	}

	// File exists with content - seek to the end of the last line
	_, err = file.Seek(-1, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek in file: %w", err)
	}

	// Read the last byte to check if it's a newline
	lastByte := make([]byte, 1)
	_, err = file.Read(lastByte)
	if err != nil {
		return fmt.Errorf("failed to read last byte: %w", err)
	}

	// Seek to the end for writing
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek to end: %w", err)
	}

	// Append the page ID (with a space if not at the start of a line)
	if lastByte[0] == '\n' {
		_, err = file.WriteString(pageID + "\n")
	} else {
		_, err = file.WriteString(" " + pageID + "\n")
	}

	return err
}
