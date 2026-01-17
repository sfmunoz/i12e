package cmd

import (
	"github.com/sfmunoz/i12e/internal/artifact"
	"github.com/spf13/cobra"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Artifact generation and push with rclone",
	Long: `Artifact management:

  - generation: tar+gz artifact
  - push to remote using rclone`,
	Run: artifact.Run,
}

func init() {
	rootCmd.AddCommand(artifactCmd)
	// artifactCmd.PersistentFlags().String("foo", "", "A help for foo")
	// artifactCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
