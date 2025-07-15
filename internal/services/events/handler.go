// cmd/root.go
package events

import (
	"encoding/json"
	"fmt"

	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/storage/local"
)

func HandleAssignmentCreate(data json.RawMessage) {
	// Handle assignment creation
	userID, err := local.GetCurrentUserID()
	if err != nil {
		return
	}

	db, err := local.GetLocalDB(userID)
	if err != nil {
		return
	}

	var assignment assignment.LocalAssignment
	if err := json.Unmarshal(data, &assignment); err == nil {
		db.Create(&assignment)
	}
}

func HandleAssignmentUpdate(data json.RawMessage) {
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
		panic(err)
	}

	if err := tx.Model(&assignment.LocalAssignment{}).Where("remote_id = ?", update.ID).Update(update.Column, update.Value).Error; err != nil {
		fmt.Printf("Error updating assignment %s with %s = %s\n", update.ID, update.Column, update.Value)
		tx.Rollback()
		panic(err)
	}

	tx.Commit()

}

func HandleAssignmentDelete(data json.RawMessage) {
	// Similar to handleAssignmentCreate but with delete logic

	userID, err := local.GetCurrentUserID()
	if err != nil {
		return
	}

	db, err := local.GetLocalDB(userID)
	if err != nil {
		return
	}

	var assignment assignment.LocalAssignment
	if err := json.Unmarshal(data, &assignment); err == nil {
		db.Delete(&assignment)
	}
}
