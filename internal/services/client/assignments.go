package client
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CreateAssignment( assignmentData map[string]string) (map[string]interface{}  ,error) {
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
                Message string                          `json:"message"`
                Assignment map[string]interface{}       `json:"assignment"`
                Error   string                          `json:"error,omitempty"`
        }

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
                return nil, fmt.Errorf("failed to decode response: %w", err)
        }

        if response.Error != "" {
                return nil, fmt.Errorf(response.Error)
        }

        if response.Assignment == nil {
                return nil, fmt.Errorf("no assignment data in response")
        }

	return response.Assignment, nil 
}


func UpdateAssignment(id, column, value string) error {
	new_client, err := NewClient()
        if err != nil {
                return err
        }

	updateData := map[string]interface{}{
		"id" : id,
                "value": value,
                "column": column,
        }
        
	jsonData, _ := json.Marshal(updateData)

        resp, err := new_client.Post(
                "http://localhost:3000/acc-homework/assignment/update",
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
