package main

import (
	"log"

	"github.com/williamfotso/acc/notifier"
)

func main() {
	err := notifier.ScheduleNotifications()
	if err != nil {
		log.Fatal(err)
	}
}
