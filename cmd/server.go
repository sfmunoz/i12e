package cmd

import (
	"fmt"
	"time"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func setDefaults(v *viper.Viper) {
	// sub-object definition implies struct definition
	v.SetDefault("k3s.tls_san", "kmain")
	v.SetDefault("mesh.endpoint_port", 51823)             // wireguard default: 51820
	v.SetDefault("mesh.network_address", "10.119.0.0/28") // from /12 (20 bits for host) to /29 (3 bits for host)
	v.SetDefault("mesh.wireguard_interface", "wgi")
	v.SetDefault("mesh.wireguard_priv_key_fname", "/etc/i12e/wg-priv-key")
	v.SetDefault("mesh.remote_base", "rem:mesh")
	v.SetDefault("server.slumber_base", 10*time.Second)
	v.SetDefault("server.slumber_jitter", 5*time.Second)
}

func serverCmd() *cobra.Command {
	cfg := &config.ServerConfig{}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server to run on target hosts",
		Long: `Run server on target hosts:

  - pulls artifacts from rclone server
  - configures network`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			decodeHook := viper.DecodeHook(config.PrefixDecodeHook())
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
