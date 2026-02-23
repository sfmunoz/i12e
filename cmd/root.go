package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func viperNew() *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetDefault("k3s.tls_san", "kmain")
	v.SetDefault("mesh.endpoint_port", 51823)             // wireguard default: 51820
	v.SetDefault("mesh.network_address", "10.119.0.0/28") // from /12 (20 bits for host) to /29 (3 bits for host)
	v.SetDefault("mesh.wireguard_interface", "wgi")
	v.SetDefault("mesh.wireguard_priv_key_fname", "/etc/i12e/wg-priv-key")
	v.SetDefault("mesh.remote_base", "rem:mesh")
	v.SetDefault("server.slumber_base", 10*time.Second)
	v.SetDefault("server.slumber_jitter", 5*time.Second)
	return v
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "i12e",
		Short: "infrastructure management tool",
		Long: `Usage: i12e [OPTIONS] COMMAND

i12e is an infrastructure management tool for task automation:

  - artifact generation
  - butane to ignition translation`,
	}
	flist := []func() *cobra.Command{serverCmd, artifactCmd, butaneCmd}
	for _, f := range flist {
		cmd.AddCommand(f())
	}
	return cmd
}

func Execute() {
	err := rootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
