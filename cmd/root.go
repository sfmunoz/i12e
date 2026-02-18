package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "i12e",
	Short: "infrastructure management tool",
	Long: `Usage: i12e [OPTIONS] COMMAND

i12e is an infrastructure management tool for task automation:

  - artifact generation
  - butane to ignition translation`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.i12e.yaml)")
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// rootCmd.PersistentFlags().BoolP("prod", "p", false, "Environment: 'prod' if set (default: 'dev')")
}
