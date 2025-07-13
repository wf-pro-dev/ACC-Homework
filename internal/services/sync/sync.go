package sync

import  (

	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/core/models/courses"
)

type SyncService struct {
    localDB  *gorm.DB // User's local SQLite DB
    remoteDB *gorm.DB // Central PostgreSQL DB
}

func (s *SyncService) SyncAssignment() error {
    var pending []assignment.LocalAssignment
    if err := s.localDB.Model(&assignment.LocalAssignment{}).Where("sync_status = ?", assignment.SyncStatusPending).Find(&pending).Error; err != nil {
        return err
    }

    return s.remoteDB.Transaction(func(tx *gorm.DB) error {
        for _, local := range pending {
            remote := assignment.Assignment{
                Title:       local.Title,
                Todo:        local.Todo,
                Deadline:    local.Deadline,
                Link:        local.Link,
                CourseCode:  local.CourseCode,
                TypeName:    local.TypeName,
                StatusName:  local.StatusName,
                UserID:      local.UserID,
            }

            if err := tx.Create(&remote).Error; err != nil {
                return err
            }

            // Update local record
            if err := s.localDB.Model(&local).Updates(map[string]interface{}{
                "notion_id":   remote.NotionID,
                "sync_status": assignment.yncStatusSynced,
                "updated_at":  time.Now(),
            }).Error; err != nil {
                return err
            }
        }
        return nil
    })
}


func (s *SyncService) SyncCourse() error {
    var pending []course.LocalCourse
    if err := s.localDB.Model(&course.LocalCourse{}).Where("sync_status = ?", coursee.SyncStatusPending).Find(&pending).Error; err != nil {
        return err
    }

    return s.remoteDB.Transaction(func(tx *gorm.DB) error {
        for _, local := range pending {
            remote := assignment.Course{
                Title:       local.Title,
                Todo:        local.Todo,
                Deadline:    local.Deadline,
                Link:        local.Link,
                CourseCode:  local.CourseCode,
                TypeName:    local.TypeName,
                StatusName:  local.StatusName,
                UserID:      local.UserID,
            }

            if err := tx.Create(&remote).Error; err != nil {
                return err
            }

            // Update local record
            if err := s.localDB.Model(&local).Updates(map[string]interface{}{
                "notion_id":   remote.NotionID,
                "sync_status": assignment.yncStatusSynced,
                "updated_at":  time.Now(),
            }).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
