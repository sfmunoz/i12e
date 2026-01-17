package cmd

import (
	"github.com/sfmunoz/i12e/internal/artifact"
	"github.com/spf13/cobra"
)

func artifactRun(cmd *cobra.Command, args []string) {
	prod, err := cmd.Flags().GetBool("prod")
	if err != nil {
		log.Fatal("'cmd.Flags().GetBool()' failed", "err", err)
	}
	artifact.Run(prod)
}

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Artifact generation and push with rclone",
	Long: `Artifact management:

  - generation: tar+gz artifact
  - push to remote using rclone`,
	Run: artifactRun,
}

func init() {
	rootCmd.AddCommand(artifactCmd)
	// artifactCmd.PersistentFlags().String("foo", "", "A help for foo")
	// artifactCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
