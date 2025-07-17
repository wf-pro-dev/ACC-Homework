package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/williamfotso/acc/internal/core/models"
	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/services/network"
	"github.com/williamfotso/acc/internal/storage/local"
	"gorm.io/gorm"
)

func GetAssignments() ([]map[string]string, error) {

	var response struct {
		Message     string              `json:"message"`
		Assignments []map[string]string `json:"assignments"`
		Error       string              `json:"error,omitempty"`
	}

	isOnline := network.IsOnline()

	if isOnline {

		new_client, err := NewClient()
		if err != nil {
			return nil, err
		}

		resp, err := new_client.Get("https://newsroom.dedyn.io/acc-homework/assignment/get")

		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		if response.Error != "" {
			return nil, fmt.Errorf(response.Error)
		}

		if response.Assignments == nil {
			return make([]map[string]string, 0), nil
		}

	}

	return response.Assignments, nil

}

func CreateAssignment(assignmentData map[string]string) (map[string]string, error) {

	userID, err := local.GetCurrentUserID()
	if err != nil {
		return nil, err
	}

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
			return nil, errors.New(response.Error)
		}

		if response.Assignment == nil {
			return nil, fmt.Errorf("no assignment data in response")
		}

		a.NotionID = response.Assignment["notion_id"].(string)
		remote_id, err := strconv.Atoi(response.Assignment["id"].(string))
		if err != nil {
			return nil, fmt.Errorf("error formating remote_id: %s", err)
		}
		fmt.Printf("Remote ID: %d\n", remote_id)
		a.RemoteID = uint(remote_id)
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

	userID, err := local.GetCurrentUserID()
	if err != nil {
		return err
	}

	db, err := local.GetLocalDB(userID)
	if err != nil {
		return err
	}

	db = db.Debug()

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var local_assignment assignment.LocalAssignment
	err = tx.First(&local_assignment, "remote_id = ?", id).Error
	if err != nil {
		return err
	}

	err = tx.Model(&local_assignment).Update(column, value).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	isOnline := network.IsOnline()

	if isOnline {

		err = SendUpdate(id, column, value)
		if err != nil {
			tx.Rollback()
			return err
		}

	} else {

		int_id, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("error formating id: %s", err)
		}

		err = tx.Model(&local_assignment).Update("sync_status", assignment.SyncStatusPending).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		update := models.LocalUpdate{
			Entity:   models.Assignment,
			EntityID: uint(int_id),
			Column:   column,
			Value:    value,
		}

		if err := tx.Create(&update).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("local update failed: %w", err)
		}

	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func SendUpdate(id, column, value string) error {

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
		"https://newsroom.dedyn.io/acc-homework/assignment/update",
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

func MigrateAssignments(db *gorm.DB) error {

	count := 0

	remoteAssignments, err := GetAssignments()
	if err != nil {
		fmt.Printf("ERROR : %s", err)
		return err
	}

	for _, ra := range remoteAssignments {

		deadline, err := time.Parse(time.DateOnly, ra["deadline"])
		if err != nil {

			return fmt.Errorf("Error formating deadline : %s", err)
		}

		remote_id, err := strconv.Atoi(ra["id"])
		if err != nil {
			return fmt.Errorf("Error formating remote_id : %s", err)
		}

		localAssignment := assignment.LocalAssignment{
			RemoteID:   uint(remote_id),
			Title:      ra["title"],
			Todo:       ra["todo"],
			Deadline:   deadline,
			Link:       ra["link"],
			CourseCode: ra["course_code"],
			TypeName:   ra["type"],
			StatusName: ra["status"],
			NotionID:   ra["notion_id"],
			SyncStatus: assignment.SyncStatusSynced,
		}

		if err := db.First(&localAssignment, "remote_id = ?", remote_id).Error; err == nil {
			continue
		}

		if err := db.Create(&localAssignment).Error; err != nil {
			count++
			return err
		}
		count++
	}

	return nil
}
