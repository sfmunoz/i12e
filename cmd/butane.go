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
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return butane.Run(cfg)
		},
	}
	cmd.Flags().StringP("mode", "m", config.ModeMain.String(), fmt.Sprintf("Set output mode (valid: %q)", config.ValidModes()))
	return cmd
}

func init() {
	rootCmd.AddCommand(butaneCmd())
}
