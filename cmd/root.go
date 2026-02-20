package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cfgKey = "config"

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "i12e",
		Short: "infrastructure management tool",
		Long: `Usage: i12e [OPTIONS] COMMAND

i12e is an infrastructure management tool for task automation:

  - artifact generation
  - butane to ignition translation`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default locations: ., $HOME/.i12e/)")
	// cobra.OnInitialize(func() {
	// 	fmt.Println("cobraOnInitialize()")
	// })
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetEnvPrefix("i12e")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	v.AutomaticEnv()
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		v.AddConfigPath(".")
		v.AddConfigPath("config/dev")
		v.AddConfigPath(home + "/.i12e")
		v.SetConfigName("i12e")
		v.SetConfigType("yaml")
	}
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}
	err := v.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}
	fmt.Println("Configuration initialized. Using config file:", v.ConfigFileUsed())
	fmt.Println("i12e.version .....", v.Get("i12e.version"))
	fmt.Println("i12e.sha256sum ...", v.Get("i12e.sha256sum"))
	for i, j := range v.AllSettings() {
		fmt.Println(i, j)
	}
	// temporary gateway
	v2 := viper.New()
	cfg2 := &config.Config{}
	prod := false
	if err := config.LoadConfig(v2, cfg2, prod); err != nil {
		return err
	}
	ctx := context.WithValue(context.Background(), cfgKey, cfg2)
	cmd.SetContext(ctx)
	return nil
}
