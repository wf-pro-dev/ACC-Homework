package main

import (
	"fmt"
	"os"
	"github.com/williamfotso/acc/internal/services/auth"

)

func main () {

	/*err := auth.Login("William Fotso", "securepassword")
	//err := auth.Logout()
	if err != nil {
		fmt.Println("\nLogin failed:", err)
		os.Exit(1)
	}

	fmt.Println("\nLogin successful!")*/

	user, err := auth.GetUser()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("User information:")
	for k, v := range user {
		fmt.Printf("%-15s: %v\n", k, v)
	}

}
