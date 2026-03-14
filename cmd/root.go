package cmd

import (
	"fmt"
	"net/netip"
	"os"
	"reflect"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func rootCmd(version, commitSHA, buildTime string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "i12e",
		Short: "infrastructure management tool",
		Long: `Usage: i12e [OPTIONS] COMMAND

i12e is an infrastructure management tool for task automation:

  - artifact generation
  - butane to ignition translation`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commitSHA, buildTime),
	}
	cmd.AddCommand(serverCmd(), artifactCmd(), butaneCmd(version))
	return cmd
}

func Execute(version, commitSHA, buildTime string) {
	err := rootCmd(version, commitSHA, buildTime).Execute()
	if err != nil {
		os.Exit(1)
	}
}

func PrefixDecodeHook() viper.DecoderConfigOption {
	f := func(from reflect.Type, to reflect.Type, data any) (any, error) {
		if from.Kind() != reflect.String {
			return data, nil
		}
		if to != reflect.TypeFor[*netip.Prefix]() {
			return data, nil
		}
		s := data.(string)
		p, err := netip.ParsePrefix(s)
		if err != nil {
			return nil, fmt.Errorf("invalid prefix %q: %w", s, err)
		}
		return &p, nil
	}
	return viper.DecodeHook(f)
}

func viperNew() *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetDefault("k3s.tls_san", "kmain")
	v.SetDefault("mesh.endpoint_port", 51823)             // wireguard default: 51820
	v.SetDefault("mesh.network_address", "10.119.0.0/28") // from /12 (20 bits for host) to /29 (3 bits for host)
	v.SetDefault("mesh.wireguard_interface", "wgi")
	v.SetDefault("mesh.wireguard_priv_key_fname", "/etc/i12e/wg-priv-key")
	v.SetDefault("mesh.remote_base", "rem:mesh")
	v.SetDefault("server.slumber_base", 19*time.Minute)
	v.SetDefault("server.slumber_jitter", 2*time.Minute)
	return v
}
