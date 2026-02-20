package cmd

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/artifact"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Artifact generation and push with rclone",
	Long: `Artifact management:

  - generation: tar+gz artifact
  - push to remote using rclone`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if ctx == nil {
			return fmt.Errorf("undefined command context")
		}
		if cfg, ok := ctx.Value(cfgKey).(*config.Config); ok {
			return artifact.Run(cfg)
		}
		return fmt.Errorf("cannot get config from context")
	},
}

func init() {
	artifactCmd.Flags().BoolP("prod", "p", false, "Environment: 'prod' if set (default: 'dev')")
	rootCmd.AddCommand(artifactCmd)
}
