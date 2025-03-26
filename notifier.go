package main

import (
	"ACC-HOMEWORK/crud"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// sendNotification triggers a macOS notification
func sendNotification(assign map[string]string) {

	deadline_raw, err_date := time.Parse(time.DateOnly, assign["deadline"])
	deadline := strings.Join(strings.Split(deadline_raw.Format(time.RFC1123), " ")[:3], " ")

	subtitle := fmt.Sprintf("%v / %v", deadline, assign["course_code"])
	if err_date != nil {
		panic(err_date)
	}

	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%v" with title "%v" subtitle "%v" sound name "Frog"`, assign["todo"], assign["title"], subtitle))
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error sending notification:", err)
	}
}

// scheduleNotifications runs a loop that checks the time and triggers notifications
func scheduleNotifications() {
	notificationTimes := []string{"06:00", "13:00", "20:00"} // 6:00 AM, 1:00 PM, 8:00 PM

	db, err_db := crud.GetDB()

	if err_db != nil {
		panic(err_db)
	}

	for {
		now := time.Now()
		currentTime := now.Format("15:04") // Format time as HH:MM
		today := now.Format("2006-01-02")
		days_7_later := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

		assignements, err := crud.GetHandler(fmt.Sprintf("SELECT * FROM assignements WHERE deadline > '%v' AND deadline < '%v' ", today, days_7_later), db)

		if err != nil {
			panic(err)
		}

		for _, t := range notificationTimes {
			if currentTime == t {
				for _, assign := range assignements {
					sendNotification(assign)
					time.Sleep(10 * time.Second)
				}
				// Avoid multiple notifications in one minute
				time.Sleep(time.Minute)
			}
		}
		time.Sleep(time.Second) // Check every second
	}
}

func main() {
	scheduleNotifications()
}
