package main

import (
	"context"
	"log"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
	"github.com/taciturnaxolotl/ctfd-alerts/clients"
	"github.com/taciturnaxolotl/ctfd-alerts/cmd/serve"
	"github.com/taciturnaxolotl/ctfd-alerts/cmd/status"
)

var (
	debugLog *log.Logger
)

// rootCmd represents the base command
var cmd = &cobra.Command{
	Use:   "ctfd-alerts",
	Short: "A tool for monitoring CTFd competitions",
	Long: `      _    __     _             _           _
  ___| |_ / _| __| |       __ _| | ___ _ __| |_ ___
 / __| __| |_ / _, |_____ / _  | |/ _ \ '__| __/ __|
| (__| |_|  _| (_| |_____| (_| | |  __/ |  | |_\__ \
 \___|\__|_|  \__,_|      \__,_|_|\___|_|   \__|___/

ctfd-alerts is a command-line tool that helps you monitor CTFd-based
competitions by sending you ntfy notifications when someone bypasses
you or a new challenge is announced. You can also use the fancy status command :)`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")
		var err error
		config, err = loadConfig(configFile)
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		setupLogging(config.Debug)

		// Create a new CTFd client and add it to context
		ctfdClient := clients.NewCTFdClient(config.CTFdConfig.ApiBase, config.CTFdConfig.ApiKey)
		ctx := context.WithValue(cmd.Context(), "ctfd_client", ctfdClient)
		ctx = context.WithValue(ctx, "config", config)
		cmd.SetContext(ctx)
	},
}

func init() {
	// Add persistent flags that work across all commands
	cmd.PersistentFlags().StringP("config", "c", "config.toml", "config file path")

	// Add commands
	cmd.AddCommand(status.StatusCmd)
	cmd.AddCommand(serve.ServeCmd)
}

func main() {
	if err := fang.Execute(
		context.Background(),
		cmd,
	); err != nil {
		os.Exit(1)
	}
}
