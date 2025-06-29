package main

import (
	"ACC-HOMEWORK/crud"
	"fmt"
	"log"
	"math"
	"os/exec"
	"time"
)

func getUntilDeadline(deadline time.Time) string {

	today, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))

	if err != nil {
		return ""
	}

	if deadline.Equal(today) {
		return "Today"
	}

	if deadline.Equal(today.AddDate(0, 0, 1)) {
		return "Tomorrow"
	}

	daysRemaining := int(math.Ceil(time.Until(deadline).Hours() / 24))
	return fmt.Sprintf("%d days remaining", daysRemaining)

}

// sendNotification sends a clickable notification that opens a URL when clicked
func sendNotification(assign map[string]string) error {
	title := fmt.Sprintf("%s: %s", assign["course_code"], assign["title"])

	// Parse deadline and calculate days remaining
	deadline, err := time.Parse(time.DateOnly, assign["deadline"][:10])
	if err != nil {
		return fmt.Errorf("failed to parse deadline: %w", err)
	}

	subtitle := fmt.Sprintf("Due %s (%s)",
		deadline.Format("Jan 2, 2006"),
		getUntilDeadline(deadline))

	// Build the terminal-notifier command
	args := []string{
		"-title", title,
		"-subtitle", subtitle,
		"-message", assign["todo"],
		"-sound", "Frog",
		"-timeout", "30", // Notification stays for 30 seconds
	}

	// Add click action if URL exists
	if link, exists := assign["link"]; exists && link != "" {
		args = append(args, "-open", link)
	}

	// Find terminal-notifier in common locations
	locations := []string{
		"/usr/local/bin/terminal-notifier",                         // Homebrew default
		"/opt/homebrew/bin/terminal-notifier",                      // Apple Silicon Homebrew
		"./terminal-notifier.app/Contents/MacOS/terminal-notifier", // Local copy
	}

	var cmd *exec.Cmd
	var lastError error

	// Try different locations until we find the binary
	for _, path := range locations {
		cmd = exec.Command(path, args...)
		if err := cmd.Run(); err == nil {
			return nil // Success!
		} else {
			lastError = err
		}
	}

	return fmt.Errorf("failed to send notification (tried paths: %v): %w", locations, lastError)
}

// scheduleNotifications checks for upcoming assignments and notifies
func scheduleNotifications() error {
	db, err := crud.GetDB()
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	now := time.Now()

	query := fmt.Sprintf(
		"SELECT * FROM assignements WHERE deadline BETWEEN '%s' AND '%s' ORDER BY deadline ASC",
		now.Format(time.DateOnly),
		now.AddDate(0, 0, 7).Format(time.DateOnly),
	)

	assignments, err := crud.GetHandler(query, db)

	if err != nil {
		return fmt.Errorf("query error: %w", err)
	}

	for _, assign := range assignments {
		if err := sendNotification(assign); err != nil {
			fmt.Printf("Error notifying for assignment %s: %v\n", assign["title"], err)
		}
		time.Sleep(5 * time.Second) // Space out notifications
	}

	return nil
}

func main() {
	err := scheduleNotifications()
	if err != nil {
		log.Fatal(err)
	}
}
