package cmd

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/butane"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
)

var butaneCmd = &cobra.Command{
	Use:   "butane",
	Short: "Run butane to generate ignition code",
	Long: `Run butane to generate ignition code

Examples:
  Reset flatcar host over ssh (default: '-o bash_b64'):
    $ i12e butane | ssh core@192.168.56.51 bash

  Generate ignition file:
    $ i12e butane -o ignition`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if ctx == nil {
			return fmt.Errorf("undefined command context")
		}
		cfg, ok := ctx.Value(cfgKey).(*config.Config)
		if !ok {
			return fmt.Errorf("cannot get config from context")
		}
		prod, err := cmd.Flags().GetBool("prod")
		if err != nil {
			return err
		}
		e := "dev"
		if prod {
			e = "prod"
		}
		m, err := cmd.Flags().GetString("mode")
		if err != nil {
			return err
		}
		butaneMode, err := config.GetMode(m)
		if err != nil {
			return err
		}
		b, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		butaneBout, err := config.GetBout(b)
		if err != nil {
			return err
		}
		cfg.Butane.EncYaml = fmt.Sprintf("config/%s/butane.enc.yaml", e)
		cfg.Butane.Mode = butaneMode
		cfg.Butane.Bout = butaneBout
		return butane.Run(cfg)
	},
}

func init() {
	butaneCmd.Flags().BoolP("prod", "p", false, "Environment: 'prod' if set (default: 'dev')")
	butaneCmd.Flags().StringP("mode", "m", config.ModeMain.String(), fmt.Sprintf("Set target mode: %q", config.ValidModes()))
	butaneCmd.Flags().StringP("output", "o", config.BoutBashB64.String(), fmt.Sprintf("Set output format: %q", config.ValidBouts()))
	rootCmd.AddCommand(butaneCmd)
}
