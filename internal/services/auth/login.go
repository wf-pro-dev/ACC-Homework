package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/williamfotso/acc/internal/services/client"
	"github.com/williamfotso/acc/internal/storage/local"
)

func Login(username, password string) error {
	new_client, err := client.NewClient()
	if err != nil {
		return err
	}

	loginData := map[string]string{
		"username": username,
		"password": password,
	}
	jsonData, _ := json.Marshal(loginData)

	resp, err := new_client.Post(
		"https://newsroom.dedyn.io/acc-homework/login",
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

	var response struct {
		Message  string `json:"message"`
		Username string `json:"username"`
		UserID   string `json:"user_id"`
		Error    string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Error != "" {
		return errors.New(response.Error)
	}

	id, err := strconv.ParseUint(response.UserID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	if err := local.StoreCredentials(
		uint(id),
		response.Username,
	); err != nil {
		log.Printf("Failed to store local credentials: %v", err)
		// Continue anyway - this is non-fatal
	}

	db, err := local.GetLocalDB(uint(id))
	if err != nil {
		return err
	}

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

	tx.Commit()

	return client.SaveCookies(new_client.Jar.Cookies(nil))
}
