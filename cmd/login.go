package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/williamfotso/acc/internal/services/auth"
)

var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Login to the system",
	Long:  `Authenticate with the server and start background sync`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		fmt.Printf("Password for %s: ", username)
		password, err := gopass.GetPasswdMasked()
		if err != nil {
			log.Fatalf("Error reading password: %v", err)
		}

		err = auth.Login(username, string(password))
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		}

		// Start listener as daemon
		daemonCmd := exec.Command(os.Args[0], "listen")
		daemonCmd.Stdout = nil
		daemonCmd.Stderr = nil
		daemonCmd.Stdin = nil
		daemonCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

		if err := daemonCmd.Start(); err != nil {
			log.Fatalf("Failed to start listener: %v", err)
		}

		// Write PID file
		pidFile, err := auth.GetDaemonPIDFilePath()
		if err != nil {
			log.Fatalf("Could not get PID file path: %v", err)
		}

		if err := os.MkdirAll(filepath.Dir(pidFile), 0755); err != nil {
			log.Fatalf("Could not create config directory: %v", err)
		}

		if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", daemonCmd.Process.Pid)), 0644); err != nil {
			log.Fatalf("Could not write PID file: %v", err)
		}

		fmt.Printf("Welcome %s! Background sync started.\n", username)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
