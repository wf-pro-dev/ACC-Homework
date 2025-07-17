package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/services/sync"
	"github.com/williamfotso/acc/internal/storage/local"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync local changes with remote server",
	Run: func(cmd *cobra.Command, args []string) {
		userID, err := local.GetCurrentUserID()
		if err != nil {
			log.Fatal("Failed to get current user ID:", err)
		}

		db, err := local.GetLocalDB(userID)
		if err != nil {
			log.Fatal("Failed to get local DB:", err)
		}

		if err := sync.Sync(db); err != nil {
			log.Fatal("Sync failed:", err)
		}

		log.Println("Sync completed successfully")
	},
}
