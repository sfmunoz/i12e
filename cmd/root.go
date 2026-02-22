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

func viperNew() *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	setDefaults(v)
	return v
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
			v := viperNew()
			if cfg.Env == config.EnvNone {
				if err := v.Unmarshal(cfg, decodeHook); err != nil {
					return err
				}
				//if err := cfg.Validate(cmd.Name()); err != nil {
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
			return cfg.Validate(cmd.Name())
		},
	}
	cmd.AddCommand(serverCmd())
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
