package cmd

import (
	"github.com/sfmunoz/i12e/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server to run on target hosts",
	Long: `Run server on target hosts:

  - pulls artifacts from rclone server
  - configures network`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
