package cmd

import (
        "fmt"
        "github.com/spf13/cobra"
        "github.com/williamfotso/acc/internal/services/auth"
        "os"
)


func init() {
	rootCmd.AddCommand(logoutCmd)
}

// In your Cobra command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from the system",
	Run: func(cmd *cobra.Command, args []string) {
		if err := auth.Logout(); err != nil {
			fmt.Println("Logout failed:", err)
			os.Exit(1)
		}
		fmt.Println("Logout successful")
	},
}

