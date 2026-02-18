package cmd

import (
	"github.com/sfmunoz/i12e/internal/artifact"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
)

func artifactRun(cmd *cobra.Command, args []string) error {
	prod, err := cmd.Flags().GetBool("prod")
	if err != nil {
		return err
	}
	cfg := &config.Config{}
	if err := config.LoadConfig(cfg, prod); err != nil {
		return err
	}
	return artifact.Run(cfg)
}

func artifactCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifact",
		Short: "Artifact generation and push with rclone",
		Long: `Artifact management:

  - generation: tar+gz artifact
  - push to remote using rclone`,
		RunE: artifactRun,
	}
	cmd.Flags().BoolP("prod", "p", false, "Environment: 'prod' if set (default: 'dev')")
	return cmd
}

func init() {
	rootCmd.AddCommand(artifactCmd())
}
