package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func setDefaults(v *viper.Viper) {
	// sub-object definition implies struct definition
	v.SetDefault("mesh.endpoint_port", 51823)             // wireguard default: 51820
	v.SetDefault("mesh.network_address", "10.119.0.0/28") // from /12 (20 bits for host) to /29 (3 bits for host)
	v.SetDefault("mesh.wireguard_interface", "wgi")
	v.SetDefault("mesh.wireguard_priv_key_fname", "/etc/i12e/wg-priv-key")
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
			prod, err := cmd.Flags().GetBool("prod")
			decodeHook := viper.DecodeHook(config.PrefixDecodeHook())
			if err != nil {
				// when the command doesn't define -p/--prod use default config
				v := viper.New()
				setDefaults(v)
				if err := v.Unmarshal(cfg, decodeHook); err != nil {
					return err
				}
				//if err := cfg.Validate(); err != nil {
				//	return err
				//}
				return nil
			}
			e := "dev"
			if prod {
				e = "prod"
			}
			v := viper.New()
			v.BindPFlag("butane.mode", cmd.Flags().Lookup("mode"))
			v.BindPFlag("butane.bout", cmd.Flags().Lookup("output"))
			setDefaults(v)
			v.SetDefault("butane.enc_yaml", fmt.Sprintf("config/%s/butane.enc.yaml", e)) // implies 'Butane' structure definition
			v.SetConfigType("yaml")
			fp, err := os.Open(fmt.Sprintf("config/%s/i12e.yaml", e))
			if err != nil {
				return err
			}
			defer fp.Close()
			bufOut, bufErr, err := cmdutil.RunSimple(exec.Command("sops", "decrypt", fmt.Sprintf("config/%s/i12e.enc.yaml", e)))
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
