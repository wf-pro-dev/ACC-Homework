package auth

import (
	"fmt"
	"net/http"
	"github.com/williamfotso/acc/internal/services/client" // Update with your correct import path
)

func Logout() error {
	// Create new client that will load existing cookies
	c, err := client.NewClient()
	if err != nil {
		return err
	}

	// Make POST request to logout endpoint (empty body)
	resp, err := c.Post(
		"http://localhost:3000/acc-homework/logout", // Note: changed from /login to /logout
		"application/json",
		nil, // No body needed for logout
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Only consider status 200 OK as successful logout
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logout failed with status: %d", resp.StatusCode)
	}

	// Clear local cookies regardless of server response
	if err := client.ClearCookies(); err != nil {
		return fmt.Errorf("failed to clear local cookies: %w", err)
	}

	return nil
}
