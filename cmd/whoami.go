package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/services/auth"
)

var getUserCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Get current user information",
	Run: func(cmd *cobra.Command, args []string) {
		user, err := auth.GetUser()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Println("User information:")
		for k, v := range user {
			fmt.Printf("%-15s: %v\n", k, v)
		}
	},
}

func init() {
	rootCmd.AddCommand(getUserCmd)
}
