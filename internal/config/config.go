package config

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/logit"
	"github.com/spf13/viper"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "config")

type Config struct {
	RcloneRemote      string   `mapstructure:"rclone_remote"`
	SshAuthorizedKeys []string `mapstructure:"ssh_authorized_keys"`
}

func LoadConfig(prod bool) (*Config, error) {
	fname := "secrets-dev.yaml"
	if prod {
		fname = "secrets-prod.yaml"
	}
	buf, err := cmdutil.SopsDecrypt(fname)
	if err != nil {
		return nil, fmt.Errorf("cmdutil.SopsDecrypt() failed: err=%s; prod=%t", err, prod)
	}
	v := viper.New()
	//v.SetEnvPrefix("I12E")
	//v.AutomaticEnv()
	v.SetConfigType("yaml")
	v.ReadConfig(buf)
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if len(cfg.SshAuthorizedKeys) < 1 {
		return nil, fmt.Errorf("undefined 'ssh_authorized_keys'")
	}
	if len(cfg.RcloneRemote) < 1 {
		return nil, fmt.Errorf("undefined 'rclone_remote'")
	}
	return &cfg, nil
}
