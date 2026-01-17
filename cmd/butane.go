package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// butaneCmd represents the butane command
var butaneCmd = &cobra.Command{
	Use:   "butane",
	Short: "Run butane to generate ignition code",
	Long: `Run butane to generate ignition code:

  - input: butane
  - output: ignition`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("butane called")
	},
}

func init() {
	rootCmd.AddCommand(butaneCmd)
	// butaneCmd.PersistentFlags().String("foo", "", "A help for foo")
	// butaneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
