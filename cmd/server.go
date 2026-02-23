package cmd

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverCmd() *cobra.Command {
	cfg := &config.ServerConfig{}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server to run on target hosts",
		Long: `Run server on target hosts:

  - pulls artifacts from rclone server
  - configures network`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			decodeHook := viper.DecodeHook(PrefixDecodeHook())
			v := viperNew()
			cfg.Env = config.EnvNone
			if err := v.Unmarshal(cfg, decodeHook); err != nil {
				return err
			}
			return cfg.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg == nil {
				return fmt.Errorf("undefined config")
			}
			return server.Run(cfg)
		},
	}
	return cmd
}
