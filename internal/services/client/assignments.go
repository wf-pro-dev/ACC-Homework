package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/services/network"
	"github.com/williamfotso/acc/internal/storage/local"
)

func CreateAssignment(assignmentData map[string]string) (map[string]string, error) {

	userID := uint(1)
	db, err := local.GetLocalDB(userID)
	if err != nil {
		return nil, err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	deadline, err := time.Parse(time.RFC3339, assignmentData["deadline"])
	if err != nil {
		return nil, err
	}

	// Create local assignment
	a := assignment.LocalAssignment{
		Title:      assignmentData["title"],
		Todo:       assignmentData["todo"],
		Deadline:   deadline,
		Link:       assignmentData["link"],
		CourseCode: assignmentData["course_code"],
		TypeName:   assignmentData["type_name"],
		StatusName: assignmentData["status_name"],
	}

	isOnline := network.IsOnline()

	if isOnline {

		new_client, err := NewClient()
		if err != nil {
			return nil, err
		}

		jsonData, _ := json.Marshal(assignmentData)

		resp, err := new_client.Post(
			"https://newsroom.dedyn.io/acc-homework/assignment",
			"application/json",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
		}

		var response struct {
			Message    string                 `json:"message"`
			Assignment map[string]interface{} `json:"assignment"`
			Error      string                 `json:"error,omitempty"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		if response.Error != "" {
			return nil, fmt.Errorf(response.Error)
		}

		if response.Assignment == nil {
			return nil, fmt.Errorf("no assignment data in response")
		}

		a.NotionID = response.Assignment["notion_id"].(string)
		a.SyncStatus = assignment.SyncStatusSynced

	} else {
		a.SyncStatus = assignment.SyncStatusPending
	}

	if err := tx.Create(&a).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("local create failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return a.ToMap(), nil
}

func UpdateAssignment(id, column, value string) error {
	new_client, err := NewClient()
	if err != nil {
		return err
	}

	updateData := map[string]interface{}{
		"id":     id,
		"value":  value,
		"column": column,
	}

	jsonData, _ := json.Marshal(updateData)

	resp, err := new_client.Post(
		"http://localhost:3000/acc-homework/assignment/update",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
