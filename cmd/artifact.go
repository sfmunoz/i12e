package cmd

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/artifact"
	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
)

func artifactCmd() *cobra.Command {
	cfg := &config.ArtifactConfig{}
	cmd := &cobra.Command{
		Use:   "artifact",
		Short: "Artifact generation and push with rclone",
		Long: `Artifact management:

  - generation: tar+gz artifact
  - push to remote using rclone`,
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
			bufOut, bufErr, err := cmdutil.RunSimple(cfg.Env.SopsCmd("i12e-conf"))
			if err != nil {
				return fmt.Errorf("'SopsCmd(i12e-conf)' failed: err=%s; buf_err=%s", err, bufErr)
			}
			v := viperNew()
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
			return artifact.Run(cfg)
		},
	}
	cmd.Flags().BoolP("prod", "p", false, "Environment: 'prod' if set (default: 'dev')")
	return cmd
}
