package main

import (
	"fmt"
	"log"

	"github.com/williamfotso/acc/internal/storage/local"
)

func main() {
	userID := uint(1)

	// 1. Initialize database
	db, err := local.GetLocalDB(userID)
	if err != nil {
		log.Fatalf("Failed to create local DB: %v", err)
	}
	fmt.Println("✅ Database connection established")

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := local.SeedInitialData(tx); err != nil {
		tx.Rollback()
		log.Fatalf("Failed to seed initial data: %v", err)
	}
	fmt.Println("✅ Initial data seeded")

	tx.Commit()

	/*
		// 2. Create test course FIRST
		testCourse := course.LocalCourse{
			Code:   "TEST-101",
			Name:   "Test Course",
			UserID: userID,
		}

		if err := tx.Create(&testCourse).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to create course: %v", err)
		}
		fmt.Println("✅ Test course created")

		// 3. Now create assignment that references the course
		testAssignment := assignment.LocalAssignment{
			Title:      "Test Assignment",
			Todo:       "This is a test assignment",
			Deadline:   time.Now().Add(7 * 24 * time.Hour),
			CourseCode: testCourse.Code, // Must match existing course
			TypeName:   "HW",            // Must exist in local_assignment_types
			StatusName: "Not started",   // Must exist in local_assignment_statuses
			UserID:     userID,
		}

		if err := tx.Create(&testAssignment).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to create assignment: %v", err)
		}
		fmt.Println("✅ Test assignment created")

		// 4. Verify data
		var foundAssignment assignment.LocalAssignment
		if err := tx.Preload("Course").
			Where("title = ?", "Test Assignment").
			First(&foundAssignment).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to read assignment: %v", err)
		}

		tx.Commit()

		fmt.Printf("✅ Retrieved assignment with course:\n%+v\n", foundAssignment)
		fmt.Printf("Course details: %+v\n", foundAssignment.Course)

		// 5. Cleanup
		//db.Exec("DELETE FROM local_assignments WHERE title = ?", "Test Assignment")
		//db.Exec("DELETE FROM local_courses WHERE code = ?", "TEST-101")
		fmt.Println("✅ Test data cleaned up")*/
}
