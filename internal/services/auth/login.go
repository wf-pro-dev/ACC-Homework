package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/williamfotso/acc/internal/services/client"

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
		"http://localhost:3000/acc-homework/login",
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

	return client.SaveCookies(new_client.Jar.Cookies(nil))
}
