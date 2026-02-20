package cmd

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server to run on target hosts",
	Long: `Run server on target hosts:

  - pulls artifacts from rclone server
  - configures network`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if ctx == nil {
			return fmt.Errorf("undefined command context")
		}
		cfg, ok := ctx.Value(cfgKey).(*config.Config)
		if !ok {
			return fmt.Errorf("cannot get config from context")
		}
		return server.Run(cfg)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
