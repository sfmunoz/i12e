package cmd

import (
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/server"
	"github.com/spf13/cobra"
)

func serverRun(cmd *cobra.Command, args []string) error {
	prod, err := cmd.Flags().GetBool("prod")
	if err != nil {
		return err
	}
	cfg, err := config.LoadConfig(prod)
	if err != nil {
		return err
	}
	return server.Run(cfg)
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server to run on target hosts",
	Long: `Run server on target hosts:

  - pulls artifacts from rclone server
  - configures network`,
	RunE: serverRun,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
