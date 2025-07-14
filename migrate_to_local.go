package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"


	"github.com/williamfotso/acc/internal/core/models/assignment"
	"github.com/williamfotso/acc/internal/core/models/course"
	"github.com/williamfotso/acc/internal/services/client"
	"github.com/williamfotso/acc/internal/storage/local"
)

func main() {

	userID, err := local.GetCurrentUserID()
	if err != nil {
		log.Fatalf("Failed to get current user ID: %v", err)
	}

	fmt.Printf("Current user ID: %d\n", userID)

	// 1. Connect to SQLite (local)
	localDB, err := local.GetLocalDB(userID)
	if err != nil {
		log.Fatalf("Failed to connect to local DB: %v", err)
	}

	tx := localDB.Begin().Debug()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 3. Migrate courses
	if err := migrateCourses(tx); err != nil {
		tx.Rollback()
		log.Fatalf("Course migration failed: %v", err)
	}

	// 4. Migrate assignments
	if err := migrateAssignments(tx); err != nil {
		tx.Rollback()
		log.Fatalf("Assignment migration failed: %v", err)
	}

	tx.Commit()

	fmt.Println("✅ Migration completed successfully")
}


func migrateCourses(localDB *gorm.DB) error {

	remoteCourses, err := client.GetCourses()
	if err != nil {
		fmt.Printf("ERROR : %s", err)
		return err
	}

	for _, rc := range remoteCourses {
		localCourse := course.LocalCourse{
			Code:       rc["code"],
			Name:       rc["name"],
			NotionID:   rc["notion_id"],
			Duration:   rc["duration"],
			RoomNumber: rc["room_number"],
			SyncStatus: course.SyncStatusSynced,
		}

		if err := localDB.Create(&localCourse).Error; err != nil {
			return err
		}
	}

	fmt.Printf("✅ Migrated %d courses\n", len(remoteCourses))
	return nil
}

func migrateAssignments(localDB *gorm.DB) error {

	remoteAssignments, err := client.GetAssignments()
	if err != nil {
		fmt.Printf("ERROR : %s", err)
		return err
	}

	for _, ra := range remoteAssignments {

		deadline, err := time.Parse(time.DateOnly, ra["deadline"])
		if err != nil {

			return fmt.Errorf("Error formating deadline : %s", err)
		}

		localAssignment := assignment.LocalAssignment{
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

		if err := localDB.Create(&localAssignment).Error; err != nil {
			return err
		}
	}

	fmt.Printf("✅ Migrated %d assignments\n", len(remoteAssignments))
	return nil
}
