package config

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/spf13/viper"
)

type I12e struct {
	Version   string `mapstructure:"version"`
	Sha256sum string `mapstructure:"sha256sum"`
}

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
	Mode              Mode
	Bout              Bout
	I12e              *I12e     `mapstructure:"i12e"`
	RcloneRemote      string    `mapstructure:"rclone_remote"`
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

func validateI12e(i12e *I12e) error {
	if i12e == nil {
		return fmt.Errorf("config: undefined 'i12e'")
	}
	s256_len := len(i12e.Sha256sum)
	if s256_len < 1 {
		return fmt.Errorf("config: undefined 'i12e.sha256sum'")
	}
	if s256_len != 64 {
		return fmt.Errorf("config: len(i12e.sha256sum)=%d (64 expected)", s256_len)
	}
	if len(i12e.Version) < 1 {
		return fmt.Errorf("config: undefined 'i12e.version'")
	}
	return nil
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
		return nil // kube_vip is optional
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
		return nil // flannel is optional
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
	if err := validateI12e(cfg.I12e); err != nil {
		return err
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
	fname := "config/dev/i12e.yaml"
	if prod {
		fname = "config/prod/i12e.yaml"
	}
	fp, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	fnameEnc := "config/dev/i12e.enc.yaml"
	if prod {
		fnameEnc = "config/prod/i12e.enc.yaml"
	}
	bufOut, bufErr, err := cmdutil.RunSimple(exec.Command("sops", "decrypt", fnameEnc))
	if err != nil {
		return nil, fmt.Errorf("'sops decrypt' failed: err=%s; buf_err=%s; prod=%t", err, bufErr, prod)
	}
	v := viper.New()
	//v.SetEnvPrefix("I12E")
	//v.AutomaticEnv()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(fp); err != nil {
		return nil, err
	}
	if err := v.MergeConfig(bufOut); err != nil {
		return nil, err
	}
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
