package cmd

import (
	"fmt"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/butane"
	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
)

func butaneCmd(versionDefault string) *cobra.Command {
	cfg := &config.ButaneConfig{}
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
			cfg.Env = config.EnvNone // when the command doesn't define -p/--prod use default config
			if prod, err := cmd.Flags().GetBool("prod"); err == nil {
				if prod {
					cfg.Env = config.EnvProd
				} else {
					cfg.Env = config.EnvDev
				}
			}
			if cfg.Env != config.EnvDev && cfg.Env != config.EnvProd {
				return fmt.Errorf("wrong env '%s': it must be '%s' or '%s'", cfg.Env, config.EnvDev, config.EnvProd)
			}
			bufOut, bufErr, err := cmdutil.RunSimple(exec.Command("sops", "decrypt", cfg.Env.I12eEncYaml()))
			if err != nil {
				return fmt.Errorf("'sops decrypt' failed: err=%s; buf_err=%s", err, bufErr)
			}
			v := viperNew()
			v.BindPFlag("butane.version", cmd.Flags().Lookup("version"))
			v.BindPFlag("butane.mode", cmd.Flags().Lookup("mode"))
			v.BindPFlag("butane.output", cmd.Flags().Lookup("output"))
			if err := v.ReadConfig(bufOut); err != nil {
				return err
			}
			if err := v.Unmarshal(cfg, PrefixDecodeHook()); err != nil {
				return err
			}
			return cfg.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg == nil {
				return fmt.Errorf("undefined config")
			}
			return butane.Run(cfg)
		},
	}
	cmd.Flags().BoolP("prod", "p", false, "Environment: 'prod' if set (default: 'dev')")
	cmd.Flags().StringP("version", "v", versionDefault, "Set version to deploy")
	cmd.Flags().StringP("mode", "m", config.ModeMain.String(), fmt.Sprintf("Set target mode: %q", config.ValidModes()))
	cmd.Flags().StringP("output", "o", config.OutputBashB64.String(), fmt.Sprintf("Set output format: %q", config.ValidOutputs()))
	return cmd
}
