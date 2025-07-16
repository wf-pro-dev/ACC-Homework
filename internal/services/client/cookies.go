package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/net/publicsuffix"
)

const (
	appName = "acc-homework" // Change this to your application name
)

// getCookieFilePath returns the canonical path for the cookie file.
func getCookieFilePath() (string, error) {
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

// SaveCookies serializes the cookies from the client's jar and saves them to a file.
func SaveCookies(client *http.Client) error {
	cookieFile, err := getCookieFilePath()
	if err != nil {
		return err
	}

	cookies := client.Jar.Cookies(nil) // Pass nil to get all cookies
	data, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal cookies: %w", err)
	}

	return ioutil.WriteFile(cookieFile, data, 0600)
}

// NewClientWithCookies creates a new http.Client and loads cookies from the file.
func NewClientWithCookies() (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	cookieFile, err := getCookieFilePath()
	if err != nil {
		return nil, err
	}

	// It's okay if the cookie file doesn't exist yet.
	if _, err := os.Stat(cookieFile); os.IsNotExist(err) {
		log.Println("No cookie file found, creating new client.")
		return &http.Client{Jar: jar}, nil
	}

	data, err := ioutil.ReadFile(cookieFile)
	if err != nil {
		return nil, fmt.Errorf("could not read cookie file: %w", err)
	}

	var cookies []*http.Cookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return nil, fmt.Errorf("could not unmarshal cookies: %w", err)
	}

	// Assuming the cookies are for the correct domain.
	if len(cookies) > 0 {
		// We need a URL to set cookies in the jar.
		// This should be the base URL of your service.
		url, _ := url.Parse("https://newsroom.dedyn.io")
		jar.SetCookies(url, cookies)
	}

	return &http.Client{Jar: jar}, nil
}

func LoadCookies() ([]*http.Cookie, error) {
	path, err := getCookieFilePath()
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

// ClearCookies removes the cookie file from disk.
func ClearCookies() error {
	cookieFile, err := getCookieFilePath()
	if err != nil {
		return err
	}
	// It's not an error if the file doesn't exist.
	if err := os.Remove(cookieFile); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
