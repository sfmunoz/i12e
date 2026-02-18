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

type K3s struct {
	Token      string `mapstructure:"token"`
	AgentToken string `mapstructure:"agent_token"`
}

type KubeVip struct {
	Vip       string `mapstructure:"vip"`
	Interface string `mapstructure:"interface"`
	Kvversion string `mapstructure:"kvversion"`
}

type Mesh struct {
	EndpointInterface  string `mapstructure:"endpoint_interface"`
	WireGuardInterface string `mapstructure:"wireguard_interface"`
}

type Pushover struct {
	UserKey string `mapstructure:"user_key"`
	Token   string `mapstructure:"token"`
}

type Butane struct {
	Mode    Mode
	Bout    Bout
	EncYaml string
}

type Config struct {
	I12e              *I12e     `mapstructure:"i12e"`
	K3s               *K3s      `mapstructure:"k3s"`
	RcloneRemote      string    `mapstructure:"rclone_remote"`
	PortKnocking      []int     `mapstructure:"port_knocking"`
	KubeVip           *KubeVip  `mapstructure:"kube_vip"`
	Mesh              *Mesh     `mapstructure:"mesh"`
	Pushover          *Pushover `mapstructure:"pushover"`
	SshAuthorizedKeys []string  `mapstructure:"ssh_authorized_keys"`
	Butane            *Butane
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

func validateK3s(k3s *K3s) error {
	if k3s == nil {
		return fmt.Errorf("config: undefined 'k3s'")
	}
	if len(k3s.Token) < 1 {
		return fmt.Errorf("config: undefined 'k3s.token'")
	}
	if len(k3s.AgentToken) < 1 {
		return fmt.Errorf("config: undefined 'k3s.agent_token'")
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

func validateMesh(mesh *Mesh) error {
	if mesh == nil {
		return fmt.Errorf("config: undefined 'mesh'")
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
	if err := validateI12e(cfg.I12e); err != nil {
		return err
	}
	if err := validateK3s(cfg.K3s); err != nil {
		return err
	}
	if err := validatePortKnocking(cfg.PortKnocking); err != nil {
		return err
	}
	if err := validateKubeVip(cfg.KubeVip); err != nil {
		return err
	}
	if err := validateMesh(cfg.Mesh); err != nil {
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

func setDefaults(v *viper.Viper) {
	// implies definition of 'Mesh' structure
	v.SetDefault("mesh.wireguard_interface", "wgi")
}

func LoadConfig(cfg *Config, prod bool) error {
	e := "dev"
	if prod {
		e = "prod"
	}
	i12eYaml := fmt.Sprintf("config/%s/i12e.yaml", e)
	i12eEncYaml := fmt.Sprintf("config/%s/i12e.enc.yaml", e)
	fp, err := os.Open(i12eYaml)
	if err != nil {
		return err
	}
	defer fp.Close()
	bufOut, bufErr, err := cmdutil.RunSimple(exec.Command("sops", "decrypt", i12eEncYaml))
	if err != nil {
		return fmt.Errorf("'sops decrypt' failed: err=%s; buf_err=%s", err, bufErr)
	}
	v := viper.New()
	setDefaults(v)
	//v.SetEnvPrefix("I12E")
	//v.AutomaticEnv()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(fp); err != nil {
		return err
	}
	if err := v.MergeConfig(bufOut); err != nil {
		return err
	}
	if err := v.Unmarshal(cfg); err != nil {
		return err
	}
	if err := validateConfig(cfg); err != nil {
		return err
	}
	return nil
}
