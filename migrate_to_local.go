package main

import (
	"fmt"
	"log"
	"strconv"
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

	// 2. Seed initial data
	if err := local.SeedInitialData(tx); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to seed initial data: %v", err)
	}
	fmt.Println("✅ Initial data seeded")

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

	count := 0

	remoteCourses, err := client.GetCourses()
	if err != nil {
		fmt.Printf("ERROR : %s", err)
		return err
	}

	count := 0

	for _, rc := range remoteCourses {
		remote_id, err := strconv.Atoi(rc["id"])
		if err != nil {
			return fmt.Errorf("Error formating remote_id : %s", err)
		}
		localCourse := course.LocalCourse{
			RemoteID:   uint(remote_id),
			Code:       rc["code"],
			Name:       rc["name"],
			NotionID:   rc["notion_id"],
			Duration:   rc["duration"],
			RoomNumber: rc["room_number"],
			SyncStatus: course.SyncStatusSynced,
		}
		fmt.Printf("Course remote_id: %d\n", remote_id)
		if err := localDB.First(&localCourse, "remote_id = ?", remote_id).Error; err == nil {
			continue
		}

		if err := localDB.First(&localCourse, "remote_id = ?", remote_id).Error; err == nil {
			continue
		}

		if err := localDB.Create(&localCourse).Error; err != nil {
			count++
			return err
		}
		count++
	}

	fmt.Printf("✅ Migrated %d courses\n", count)

	return nil
}

func migrateAssignments(localDB *gorm.DB) error {

	count := 0

	remoteAssignments, err := client.GetAssignments()
	if err != nil {
		fmt.Printf("ERROR : %s", err)
		return err
	}

	count := 0

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

		if err := localDB.First(&localAssignment, "remote_id = ?", remote_id).Error; err == nil {
			continue
		}

		if err := localDB.Create(&localAssignment).Error; err != nil {
			count++
			return err
		}
		count++
	}

	fmt.Printf("✅ Migrated %d assignments\n", count)
	return nil
}
