package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cfgKey = "config"

var rootCmd = &cobra.Command{
	Use:   "i12e",
	Short: "infrastructure management tool",
	Long: `Usage: i12e [OPTIONS] COMMAND

i12e is an infrastructure management tool for task automation:

  - artifact generation
  - butane to ignition translation`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		prod, err := cmd.Flags().GetBool("prod") // don't want a global flag -> each command must define it
		if err != nil {
			return err
		}
		e := "dev"
		if prod {
			e = "prod"
		}
		v := viper.New()
		v.BindPFlag("butane.mode", cmd.Flags().Lookup("mode"))
		v.BindPFlag("butane.bout", cmd.Flags().Lookup("output"))
		v.SetDefault("mesh.wireguard_interface", "wgi")                              // implies of 'Mesh' structure definition
		v.SetDefault("butane.enc_yaml", fmt.Sprintf("config/%s/butane.enc.yaml", e)) // implies of 'Butane' structure definition
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
		cfg := &config.Config{}
		if err := v.Unmarshal(cfg); err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}
		ctx := context.WithValue(context.Background(), cfgKey, cfg)
		cmd.SetContext(ctx)
		return nil

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
