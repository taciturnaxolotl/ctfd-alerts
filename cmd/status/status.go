package status

import (
	"github.com/spf13/cobra"
)

// StatusCmd represents the status command
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show CTFd status information",
	Long:  "Shows the current CTFd scoreboard and list of challenges in a tabular format",
	Run:   runDashboard,
}
