package cmd

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/services/auth"
	"github.com/williamfotso/acc/internal/services/events"
	"syscall"
)

var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Login to the system",
	Long:  `Authenticate with the server using your username and password`,
	Args:  cobra.ExactArgs(1), // Requires exactly 1 argument (username)
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		
		// Prompt for password with masking
		fmt.Printf("Password for %s: ", username)
		password, err := gopass.GetPasswdMasked()
		if err != nil {
			fmt.Println("\nError reading password:", err)
			return
		}
		
		// Perform login
		err = auth.Login(username, string(password))
		if err != nil {
			fmt.Println("\nLogin failed:", err)
			syscall.Exit(1)
		}

		eventHandler.Start()
		
		fmt.Println("\nLogin successful!")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
