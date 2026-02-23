package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/artifact"
	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			fp, err := os.Open(cfg.Env.I12eYaml())
			if err != nil {
				return err
			}
			defer fp.Close()
			bufOut, bufErr, err := cmdutil.RunSimple(exec.Command("sops", "decrypt", cfg.Env.I12eEncYaml()))
			if err != nil {
				return fmt.Errorf("'sops decrypt' failed: err=%s; buf_err=%s", err, bufErr)
			}
			v := viperNew()
			if err := v.ReadConfig(fp); err != nil {
				return err
			}
			if err := v.MergeConfig(bufOut); err != nil {
				return err
			}
			decodeHook := viper.DecodeHook(PrefixDecodeHook())
			if err := v.Unmarshal(cfg, decodeHook); err != nil {
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
