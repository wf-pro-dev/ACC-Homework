package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/services/auth"
)

func init() {
	rootCmd.AddCommand(logoutCmd)
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out and stop the background listener",
	Run: func(cmd *cobra.Command, args []string) {
		// Stop the running listener daemon.
		pidFile, err := auth.GetDaemonPIDFilePath()
		if err != nil {
			log.Printf("Warning: could not get PID file path: %v", err)
		} else {
			if _, err := os.Stat(pidFile); err == nil {
				pidBytes, err := ioutil.ReadFile(pidFile)
				if err != nil {
					log.Printf("Warning: could not read PID file: %v", err)
				} else {
					pid, err := strconv.Atoi(string(pidBytes))
					if err != nil {
						log.Printf("Warning: invalid PID in file: %v", err)
					} else {
						// Find and signal the process to terminate.
						process, err := os.FindProcess(pid)
						if err != nil {
							log.Printf("Warning: could not find process with PID %d: %v", pid, err)
						} else {
							// Send SIGTERM for graceful shutdown.
							if err := process.Signal(syscall.SIGTERM); err != nil {
								log.Printf("Warning: failed to send signal to listener: %v", err)
							}
						}
					}
				}
				// Clean up the PID file regardless of success.
				os.Remove(pidFile)
			}
		}

		// Perform the cleanup of credentials and cookies.
		if err := auth.Logout(); err != nil {
			log.Fatalf("Logout failed: %v", err)
		}
		fmt.Println("Logout successful.")
	},
}
