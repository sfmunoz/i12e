package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
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

func rootCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "i12e",
		Short: "infrastructure management tool",
		Long: `Usage: i12e [OPTIONS] COMMAND

i12e is an infrastructure management tool for task automation:

  - artifact generation
  - butane to ignition translation`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg.Env = config.EnvNone // when the command doesn't define -p/--prod use default config
			if prod, err := cmd.Flags().GetBool("prod"); err == nil {
				if prod {
					cfg.Env = config.EnvProd
				} else {
					cfg.Env = config.EnvDev
				}
			}
			decodeHook := viper.DecodeHook(config.PrefixDecodeHook())
			v := viper.New()
			v.SetConfigType("yaml")
			setDefaults(v)
			if cfg.Env == config.EnvNone {
				if err := v.Unmarshal(cfg, decodeHook); err != nil {
					return err
				}
				//if err := cfg.Validate(); err != nil {
				//	return err
				//}
				return nil
			}
			if cmd.Name() == "butane" {
				v.BindPFlag("butane.mode", cmd.Flags().Lookup("mode"))
				v.BindPFlag("butane.output", cmd.Flags().Lookup("output"))
			}
			fp, err := os.Open(cfg.Env.I12eYaml())
			if err != nil {
				return err
			}
			defer fp.Close()
			bufOut, bufErr, err := cmdutil.RunSimple(exec.Command("sops", "decrypt", cfg.Env.I12eEncYaml()))
			if err != nil {
				return fmt.Errorf("'sops decrypt' failed: err=%s; buf_err=%s", err, bufErr)
			}
			if err := v.ReadConfig(fp); err != nil {
				return err
			}
			if err := v.MergeConfig(bufOut); err != nil {
				return err
			}
			if err := v.Unmarshal(cfg, decodeHook); err != nil {
				return err
			}
			return cfg.Validate()
		},
	}
	cmd.AddCommand(serverCmd(cfg))
	cmd.AddCommand(artifactCmd(cfg))
	cmd.AddCommand(butaneCmd(cfg))
	return cmd
}

func Execute() {
	cfg := &config.Config{}
	err := rootCmd(cfg).Execute()
	if err != nil {
		os.Exit(1)
	}
}
