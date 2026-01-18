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
	cfg, err := config.LoadConfig(prod)
	if err != nil {
		return err
	}
	return artifact.Run(cfg)
}

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Artifact generation and push with rclone",
	Long: `Artifact management:

  - generation: tar+gz artifact
  - push to remote using rclone`,
	RunE: artifactRun,
}

func init() {
	rootCmd.AddCommand(artifactCmd)
	// artifactCmd.PersistentFlags().String("foo", "", "A help for foo")
	// artifactCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
