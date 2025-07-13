package main

import (
	"fmt"

	"github.com/williamfotso/acc/internal/services/network"
)

func main() {
	if network.IsOnline() {
		fmt.Println("Online")
	} else {
		fmt.Println("Offline")
	}
}
