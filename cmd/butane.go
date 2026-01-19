package cmd

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/butane"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
)

func getModeFlag(cmd *cobra.Command, cfg *config.Config) error {
	mode, err := cmd.Flags().GetString("mode")
	if err != nil {
		return err
	}
	m, err := config.GetMode(mode)
	if err != nil {
		return err
	}
	cfg.Mode = m
	return nil
}

func getOutputFlag(cmd *cobra.Command, cfg *config.Config) error {
	bout, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	b, err := config.GetBout(bout)
	if err != nil {
		return err
	}
	cfg.Bout = b
	return nil
}

func buildConfig(cmd *cobra.Command) (*config.Config, error) {
	prod, err := cmd.Flags().GetBool("prod")
	if err != nil {
		return nil, err
	}
	cfg, err := config.LoadConfig(prod)
	if err != nil {
		return nil, err
	}
	cfg.Prod = prod
	if err := getModeFlag(cmd, cfg); err != nil {
		return nil, err
	}
	if err := getOutputFlag(cmd, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func butaneCmd() *cobra.Command {
	var cfg *config.Config = nil
	cmd := &cobra.Command{
		Use:   "butane",
		Short: "Run butane to generate ignition code",
		Long: `Run butane to generate ignition code:

  - input: butane
  - output: ignition`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			c, err := buildConfig(cmd)
			if err != nil {
				return err
			}
			cfg = c
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return butane.Run(cfg)
		},
	}
	cmd.Flags().StringP("mode", "m", config.ModeMain.String(), fmt.Sprintf("Set target mode: %q", config.ValidModes()))
	cmd.Flags().StringP("output", "o", config.BoutBashB64.String(), fmt.Sprintf("Set output format: %q", config.ValidBouts()))
	return cmd
}

func init() {
	rootCmd.AddCommand(butaneCmd())
}
