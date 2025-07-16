package cmd

import (
	"fmt"
	"log"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/services/auth"
)

var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Login to the system",
	Long:  `Authenticate with the server and save session credentials.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		fmt.Printf("Password for %s: ", username)
		password, err := gopass.GetPasswdMasked()
		if err != nil {
			log.Fatalf("\nError reading password: %v", err)
		}

		err = auth.Login(username, string(password))
		if err != nil {
			log.Fatalf("\nLogin failed: %v", err)
		}

		fmt.Println("\nLogin successful!")
		fmt.Println("Run 'acc listen' in a separate terminal to start receiving notifications.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
