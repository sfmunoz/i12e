package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Artifact generation and push with rclone",
	Long: `Artifact management:

  - generation: tar+gz artifact
  - push to remote using rclone`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("artifact called")
	},
}

func init() {
	rootCmd.AddCommand(artifactCmd)
	// artifactCmd.PersistentFlags().String("foo", "", "A help for foo")
	// artifactCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
