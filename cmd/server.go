package cmd

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/server"
	"github.com/spf13/cobra"
)

func serverCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server to run on target hosts",
		Long: `Run server on target hosts:

  - pulls artifacts from rclone server
  - configures network`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg == nil {
				return fmt.Errorf("undefined config")
			}
			return server.Run(cfg)
		},
	}
	return cmd
}
