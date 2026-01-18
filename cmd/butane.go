package cmd

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/butane"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
)

func butaneCmd() *cobra.Command {
	var cfg *config.Config = nil
	cmd := &cobra.Command{
		Use:   "butane",
		Short: "Run butane to generate ignition code",
		Long: `Run butane to generate ignition code:

  - input: butane
  - output: ignition`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			prod, err := cmd.Flags().GetBool("prod")
			if err != nil {
				return err
			}
			cfg, err = config.LoadConfig(prod)
			if err != nil {
				return err
			}
			// mode
			mode, err := cmd.Flags().GetString("mode")
			if err != nil {
				return err
			}
			m, err := config.GetMode(mode)
			if err != nil {
				return err
			}
			cfg.Mode = m
			// bout (butane output)
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
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return butane.Run(cfg)
		},
	}
	cmd.Flags().StringP("mode", "m", config.ModeMain.String(), fmt.Sprintf("Set target mode: %q", config.ValidModes()))
	cmd.Flags().StringP("output", "o", config.BoutBash64.String(), fmt.Sprintf("Set output format: %q", config.ValidBouts()))
	return cmd
}

func init() {
	rootCmd.AddCommand(butaneCmd())
}
