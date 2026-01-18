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

type KubeVip struct {
	Vip       string `mapstructure:"vip"`
	Interface string `mapstructure:"interface"`
	Kvversion string `mapstructure:"kvversion"`
}

type Flannel struct {
	Interface string `mapstructure:"interface"`
}

type Pushover struct {
	UserKey string `mapstructure:"user_key"`
	Token   string `mapstructure:"token"`
}

type Config struct {
	Prod              bool
	RcloneRemote      string    `mapstructure:"rclone_remote"`
	RcloneConfigPass  string    `mapstructure:"rclone_config_pass"`
	K3sToken          string    `mapstructure:"k3s_token"`
	K3sAgentToken     string    `mapstructure:"k3s_agent_token"`
	K3sUrl            string    `mapstructure:"k3s_url"`
	TlsSan            string    `mapstructure:"tls_san"`
	PortKnocking      []int     `mapstructure:"port_knocking"`
	KubeVip           *KubeVip  `mapstructure:"kube_vip"`
	Flannel           *Flannel  `mapstructure:"flannel"`
	Pushover          *Pushover `mapstructure:"pushover"`
	SshAuthorizedKeys []string  `mapstructure:"ssh_authorized_keys"`
}

func validatePortKnocking(portKnocking []int) error {
	if len(portKnocking) < 1 {
		return fmt.Errorf("config: undefined 'port_knocking'")
	}
	l := len(portKnocking)
	if l != 4 {
		return fmt.Errorf("config: 'len(port_knocking)=%d' (valid: 4, see https://github.com/sfmunoz/i12e/issues/138)", l)
	}
	return nil
}

func validateKubeVip(kubeVip *KubeVip) error {
	if kubeVip == nil {
		return fmt.Errorf("config: undefined 'kube_vip'")
	}
	if len(kubeVip.Vip) < 1 {
		return fmt.Errorf("config: undefined 'kube_vip.vip'")
	}
	if len(kubeVip.Interface) < 1 {
		return fmt.Errorf("config: undefined 'kube_vip.interface'")
	}
	if len(kubeVip.Kvversion) < 1 {
		return fmt.Errorf("config: undefined 'kube_vip.kvversion'")
	}
	return nil
}

func validateFlannel(flannel *Flannel) error {
	if flannel == nil {
		return fmt.Errorf("config: undefined 'flannel'")
	}
	if len(flannel.Interface) < 1 {
		return fmt.Errorf("config: undefined 'flannel.interface'")
	}
	return nil
}

func validatePushover(pushover *Pushover) error {
	if pushover == nil {
		return fmt.Errorf("config: undefined 'pushover'")
	}
	if len(pushover.UserKey) < 1 {
		return fmt.Errorf("config: undefined 'pushover.user_key'")
	}
	if len(pushover.Token) < 1 {
		return fmt.Errorf("config: undefined 'pushover.token'")
	}
	return nil
}

func validateSshAuthorizedKeys(sshAuthorizedKeys []string) error {
	if len(sshAuthorizedKeys) < 1 {
		return fmt.Errorf("config: undefined 'ssh_authorized_keys'")
	}
	for i, s := range sshAuthorizedKeys {
		if len(s) < 1 {
			return fmt.Errorf("config: undefined 'ssh_authorized_keys[%d]'", i)
		}
	}
	return nil
}

func validateConfig(cfg *Config) error {
	if len(cfg.RcloneRemote) < 1 {
		return fmt.Errorf("config: undefined 'rclone_remote'")
	}
	if len(cfg.RcloneConfigPass) < 1 {
		return fmt.Errorf("config: undefined 'rclone_config_pass'")
	}
	if len(cfg.RcloneConfigPass) < 1 {
		return fmt.Errorf("config: undefined 'rclone_remote'")
	}
	if len(cfg.K3sToken) < 1 {
		return fmt.Errorf("config: undefined 'k3s_token'")
	}
	if len(cfg.K3sAgentToken) < 1 {
		return fmt.Errorf("config: undefined 'k3s_agent_token'")
	}
	if len(cfg.K3sUrl) < 1 {
		return fmt.Errorf("config: undefined 'k3s_url'")
	}
	if len(cfg.TlsSan) < 1 {
		return fmt.Errorf("config: undefined 'tls_san'")
	}
	if err := validatePortKnocking(cfg.PortKnocking); err != nil {
		return err
	}
	if err := validateKubeVip(cfg.KubeVip); err != nil {
		return err
	}
	if err := validateFlannel(cfg.Flannel); err != nil {
		return err
	}
	if err := validatePushover(cfg.Pushover); err != nil {
		return err
	}
	if err := validateSshAuthorizedKeys(cfg.SshAuthorizedKeys); err != nil {
		return err
	}
	return nil
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
	cfg.Prod = prod
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
