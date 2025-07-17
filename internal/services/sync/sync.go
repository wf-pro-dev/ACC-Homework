package sync

import (
	"strconv"
	"time"

	"github.com/williamfotso/acc/internal/core/models"
	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/core/models/course"
	"github.com/williamfotso/acc/internal/services/client"
	"gorm.io/gorm"
)

func Sync(db *gorm.DB) error {
	db = db.Debug()

	if err := SyncAssignment(db); err != nil {
		return err
	}
	if err := SyncCourse(db); err != nil {
		return err
	}
	return nil
}

func SyncAssignment(db *gorm.DB) error {
	var pending []assignment.LocalAssignment
	if err := db.Model(&assignment.LocalAssignment{}).Where("sync_status = ?", assignment.SyncStatusPending).Find(&pending).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, local := range pending {

			var remote map[string]string
			var err error

			if local.NotionID == "" {
				remote, err = client.CreateAssignment(local.ToMap())
				if err != nil {
					return err
				}

			} else {

				var update models.LocalUpdate
				if err = tx.Model(&models.LocalUpdate{}).Where("entity = ? AND entity_id = ?", models.Assignment, local.RemoteID).First(&update).Error; err != nil {
					return err
				}

				id := strconv.Itoa(int(local.RemoteID))

				err = client.SendUpdate(id, update.Column, update.Value)
				if err != nil {
					return err
				}

				if tx.Delete(&update).Error != nil {
					return err
				}

				remote = local.ToMap()
			}

			// Update local record
			if err := tx.Model(&local).Updates(map[string]interface{}{
				"notion_id":   remote["notion_id"],
				"sync_status": assignment.SyncStatusSynced,
				"updated_at":  time.Now(),
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func SyncCourse(db *gorm.DB) error {
	var pending []course.LocalCourse
	if err := db.Model(&course.LocalCourse{}).Where("sync_status = ?", course.SyncStatusPending).Find(&pending).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, local := range pending {

			var remote map[string]string
			var err error

			if local.NotionID == "" {
				remote, err = client.CreateCourse(local.ToMap())
				if err != nil {
					return err
				}
			}

			// Update local record
			if err := tx.Model(&local).Updates(map[string]interface{}{
				"notion_id":   remote["notion_id"],
				"sync_status": course.SyncStatusSynced,
				"updated_at":  time.Now(),
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
