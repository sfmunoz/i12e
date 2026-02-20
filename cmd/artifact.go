package cmd

import (
	"github.com/sfmunoz/i12e/internal/artifact"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
)

func artifactCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifact",
		Short: "Artifact generation and push with rclone",
		Long: `Artifact management:

  - generation: tar+gz artifact
  - push to remote using rclone`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return artifact.Run(cfg)
		},
	}
	cmd.Flags().BoolP("prod", "p", false, "Environment: 'prod' if set (default: 'dev')")
	return cmd
}
