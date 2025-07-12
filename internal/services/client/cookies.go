package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"net/http"
)

const (
	appName = "acc-homework"  // Change this to your application name
)

// getCookiePath returns the path for cookies in config_dir/your-cli-app/cookies/cookies.txt
func getCookiePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	
	// Create app directory and cookies subdirectory
	appDir := filepath.Join(configDir, appName)
	cookiesDir := filepath.Join(appDir, "cookies")
	
	// Create directories with 0755 permissions (rwx for owner, rx for others)
	if err := os.MkdirAll(cookiesDir, 0755); err != nil {
		return "", err
	}
	
	return filepath.Join(cookiesDir, "cookies.txt"), nil
}

// SaveCookies stores cookies to the cookies.txt file
func SaveCookies(cookies []*http.Cookie) error {
	path, err := getCookiePath()
	if err != nil {
		return err
	}

	// Create or truncate the file with 0600 permissions (rw for owner only)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(cookies)
}

// LoadCookies reads cookies from cookies.txt
func LoadCookies() ([]*http.Cookie, error) {
	path, err := getCookiePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No cookies yet
		}
		return nil, err
	}
	defer file.Close()

	var cookies []*http.Cookie
	err = json.NewDecoder(file).Decode(&cookies)
	return cookies, err
}

// ClearCookies removes the cookies.txt file
func ClearCookies() error {
	path, err := getCookiePath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}
