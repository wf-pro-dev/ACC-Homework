package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	"fmt"
	"io"
	"net/http"

	"github.com/williamfotso/acc/internal/core/models/course"
	"github.com/williamfotso/acc/internal/services/network"
	"github.com/williamfotso/acc/internal/storage/local"
	"gorm.io/gorm"
)

func GetCourses() ([]map[string]string, error) {

	var response struct {
		Message string              `json:"message"`
		Courses []map[string]string `json:"courses"`
		Error   string              `json:"error,omitempty"`
	}

	isOnline := network.IsOnline()

	if isOnline {

		new_client, err := NewClient()
		if err != nil {
			return nil, err
		}

		resp, err := new_client.Get("https://newsroom.dedyn.io/acc-homework/course/get")

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

			return nil, errors.New(response.Error)

		}

		if response.Courses == nil {
			return nil, errors.New("no assignment data in response")
		}

	}

	return response.Courses, nil
}

func CreateCourse(courseData map[string]string) (map[string]string, error) {

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

	// Create local assignment
	c := course.LocalCourse{
		Code:       courseData["code"],
		Name:       courseData["name"],
		Duration:   courseData["duration"],
		RoomNumber: courseData["room_number"],
	}

	isOnline := network.IsOnline()

	if isOnline {

		new_client, err := NewClient()
		if err != nil {
			return nil, err
		}

		jsonData, _ := json.Marshal(courseData)

		resp, err := new_client.Post(
			"https://newsroom.dedyn.io/acc-homework/course",
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
			Message string            `json:"message"`
			Course  map[string]string `json:"course"`
			Error   string            `json:"error,omitempty"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		if response.Error != "" {
			return nil, errors.New(response.Error)
		}

		if response.Course == nil {
			return nil, errors.New("no course data in response")
		}

		c.NotionID = response.Course["notion_id"]
		c.SyncStatus = course.SyncStatusSynced

	} else {
		c.SyncStatus = course.SyncStatusPending
	}

	if err := tx.Create(&c).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("local create failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return c.ToMap(), nil
}

func MigrateCourses(db *gorm.DB) error {

	count := 0

	remoteCourses, err := GetCourses()
	if err != nil {
		fmt.Printf("ERROR : %s", err)
		return err
	}

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

		if err := db.First(&localCourse, "remote_id = ?", remote_id).Error; err == nil {
			continue
		}

		if err := db.Create(&localCourse).Error; err != nil {
			count++
			return err
		}
		count++
	}

	return nil
}
