package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func viperNew() *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	return v
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "i12e",
		Short: "infrastructure management tool",
		Long: `Usage: i12e [OPTIONS] COMMAND

i12e is an infrastructure management tool for task automation:

  - artifact generation
  - butane to ignition translation`,
	}
	flist := []func() *cobra.Command{serverCmd, artifactCmd, butaneCmd}
	for _, f := range flist {
		cmd.AddCommand(f())
	}
	return cmd
}

func Execute() {
	err := rootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
