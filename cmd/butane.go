package cmd

import (
	"github.com/sfmunoz/i12e/internal/butane"
	"github.com/spf13/cobra"
)

// butaneCmd represents the butane command
var butaneCmd = &cobra.Command{
	Use:   "butane",
	Short: "Run butane to generate ignition code",
	Long: `Run butane to generate ignition code:

  - input: butane
  - output: ignition`,
	Run: butane.Run,
}

func init() {
	rootCmd.AddCommand(butaneCmd)
	// butaneCmd.PersistentFlags().String("foo", "", "A help for foo")
	// butaneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
