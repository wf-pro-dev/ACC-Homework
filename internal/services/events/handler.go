// cmd/root.go
package events

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/services/notifier"
	"github.com/williamfotso/acc/internal/storage/local"
)

type AssignmentResponse struct {
	ID         string `json:"id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Title      string `json:"title"`
	Todo       string `json:"todo"`
	Deadline   string `json:"deadline"`
	Link       string `json:"link"`
	CourseCode string `json:"course_code"`
	TypeName   string `json:"type"`
	StatusName string `json:"status"`
	NotionID   string `json:"notion_id"`
}

func HandleAssignmentCreate(data json.RawMessage, message string) {

	// Handle assignment creation
	userID, err := local.GetCurrentUserID()
	if err != nil {
		return
	}

	db, err := local.GetLocalDB(userID)
	if err != nil {
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var ar AssignmentResponse
	if err := json.Unmarshal(data, &ar); err != nil {
		log.Printf("Error unmarshalling assignment: %v", err)
		return
	}

	id, err := strconv.Atoi(ar.ID)
	if err != nil {
		log.Printf("Error converting ID to int: %v", err)
		return
	}

	deadline, err := time.Parse(time.RFC3339, ar.Deadline)
	if err != nil {
		log.Printf("Error parsing deadline: %v", err)
		return
	}

	a := assignment.LocalAssignment{
		RemoteID:   uint(id),
		Title:      ar.Title,
		Todo:       ar.Todo,
		Deadline:   deadline,
		Link:       ar.Link,
		CourseCode: ar.CourseCode,
		TypeName:   ar.TypeName,
		StatusName: ar.StatusName,
		NotionID:   ar.NotionID,
		SyncStatus: assignment.SyncStatusSynced,
	}
	if err := tx.Create(&a).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating assignment: %v", err)
		return
	}

	Notify("created", message, a.ToMap())

	tx.Commit()

}

func HandleAssignmentUpdate(data json.RawMessage, message string) {
	// Similar to handleAssignmentCreate but with update logic

	userID, err := local.GetCurrentUserID()
	if err != nil {
		return
	}

	db, err := local.GetLocalDB(userID)
	if err != nil {
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var update struct {
		ID     string `json:"id"`
		Column string `json:"column"`
		Value  string `json:"value"`
	}

	if err := json.Unmarshal(data, &update); err != nil {
		log.Printf("Error parsing update: %v", err)
		return
	}

	if err := tx.Model(&assignment.LocalAssignment{}).Where("remote_id = ?", update.ID).Update(update.Column, update.Value).Error; err != nil {
		fmt.Printf("Error updating assignment %s with %s = %s\n", update.ID, update.Column, update.Value)
		tx.Rollback()
		panic(err)
	}

	var a assignment.LocalAssignment
	err = tx.Model(&assignment.LocalAssignment{}).Where("remote_id = ?", update.ID).First(&a).Error
	if err != nil {
		log.Printf("Error getting assignment: %v", err)
		return
	}

	tx.Commit()

	Notify("updated", message, a.ToMap())

}

func HandleAssignmentDelete(data json.RawMessage, message string) {
	// Similar to handleAssignmentCreate but with delete logic
	userID, err := local.GetCurrentUserID()
	if err != nil {
		return
	}

	db, err := local.GetLocalDB(userID)
	if err != nil {
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var ar AssignmentResponse
	if err := json.Unmarshal(data, &ar); err != nil {
		log.Printf("Error unmarshalling assignment: %v", err)
		return
	}

	var a assignment.LocalAssignment
	if err := tx.Model(&assignment.LocalAssignment{}).Where("notion_id = ?", ar.NotionID).First(&a).Error; err != nil {
		tx.Rollback()
		log.Printf("Error getting assignment: %v", err)
		return
	}

	if err := tx.Where("notion_id = ?", ar.NotionID).Delete(&assignment.LocalAssignment{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting assignment: %v", err)
		return
	}

	tx.Commit()

	Notify("deleted", message, a.ToMap())

}

func Notify(action, message string, assignment map[string]string) {

	notification_id := fmt.Sprintf("%s-%s", assignment["notion_id"], action)
	title := fmt.Sprintf("%s: %s", assignment["course_code"], assignment["title"])
	subtitle := fmt.Sprintf("%s at %s", action, time.Now().Format(time.Stamp))

	args := []string{
		"-group", notification_id,
		"-title", title,
		"-subtitle", subtitle,
		"-message", message,
		"-sound", "Frog",
		"-timeout", "60", // Notification stays for 30 seconds
	}

	err := notifier.UseNotifier(args)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	time.Sleep(15 * time.Second) // Wait for the notification to be sent

	err = notifier.UseNotifier([]string{"-remove", notification_id})
	if err != nil {
		log.Printf("Error removing notification: %v", err)
	}

}
