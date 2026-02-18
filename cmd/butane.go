package cmd

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/butane"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func buildConfig(cmd *cobra.Command) (*config.Config, error) {
	prod, err := cmd.Flags().GetBool("prod")
	if err != nil {
		return nil, err
	}
	e := "dev"
	if prod {
		e = "prod"
	}
	cfg := &config.Config{}
	cfg.Butane = &config.Butane{
		EncYaml: fmt.Sprintf("config/%s/butane.enc.yaml", e),
	}
	v := viper.New()
	v.BindPFlag("butane.mode", cmd.Flags().Lookup("mode"))
	v.BindPFlag("butane.bout", cmd.Flags().Lookup("output"))
	if err := config.LoadConfig(v, cfg, prod); err != nil {
		return nil, err
	}
	return cfg, nil
}

func butaneCmd() *cobra.Command {
	var cfg *config.Config = nil
	cmd := &cobra.Command{
		Use:   "butane",
		Short: "Run butane to generate ignition code",
		Long: `Run butane to generate ignition code

Examples:
  Reset flatcar host over ssh (default: '-o bash_b64'):
    $ i12e butane | ssh core@192.168.56.51 bash

  Generate ignition file:
    $ i12e butane -o ignition`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cfg, err = buildConfig(cmd)
			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return butane.Run(cfg)
		},
	}
	cmd.Flags().BoolP("prod", "p", false, "Environment: 'prod' if set (default: 'dev')")
	cmd.Flags().StringP("mode", "m", config.ModeMain.String(), fmt.Sprintf("Set target mode: %q", config.ValidModes()))
	cmd.Flags().StringP("output", "o", config.BoutBashB64.String(), fmt.Sprintf("Set output format: %q", config.ValidBouts()))
	return cmd
}

func init() {
	rootCmd.AddCommand(butaneCmd())
}
