package config

import (
	"fmt"
	"net/netip"
)

type Config struct {
	I12e *struct {
		Version   string `mapstructure:"version"`
		Sha256sum string `mapstructure:"sha256sum"`
	} `mapstructure:"i12e"`
	K3s *struct {
		Token      string `mapstructure:"token"`
		AgentToken string `mapstructure:"agent_token"`
	} `mapstructure:"k3s"`
	RcloneRemote string `mapstructure:"rclone_remote"`
	PortKnocking []int  `mapstructure:"port_knocking"`
	KubeVip      *struct {
		Vip       string `mapstructure:"vip"`
		Interface string `mapstructure:"interface"`
		Kvversion string `mapstructure:"kvversion"`
	} `mapstructure:"kube_vip"`
	Mesh *struct {
		EndpointInterface  string        `mapstructure:"endpoint_interface"`
		EndpointPort       int           `mapstructure:"endpoint_port"`
		NetworkAddress     *netip.Prefix `mapstructure:"network_address"`
		WireGuardInterface string        `mapstructure:"wireguard_interface"`
	} `mapstructure:"mesh"`
	Pushover *struct {
		UserKey string `mapstructure:"user_key"`
		Token   string `mapstructure:"token"`
	} `mapstructure:"pushover"`
	SshAuthorizedKeys []string `mapstructure:"ssh_authorized_keys"`
	Butane            *struct {
		Mode    Mode   `mapstructure:"mode"`
		Bout    Bout   `mapstructure:"bout"`
		EncYaml string `mapstructure:"enc_yaml"`
	} `mapstructure:"butane"`
}

func (c *Config) validateI12e() error {
	i12e := c.I12e
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

func (c *Config) validateK3s() error {
	k3s := c.K3s
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

func (c *Config) validatePortKnocking() error {
	portKnocking := c.PortKnocking
	if len(portKnocking) < 1 {
		return fmt.Errorf("config: undefined 'port_knocking'")
	}
	l := len(portKnocking)
	if l != 4 {
		return fmt.Errorf("config: 'len(port_knocking)=%d' (valid: 4, see https://github.com/sfmunoz/i12e/issues/138)", l)
	}
	return nil
}

func (c *Config) validateKubeVip() error {
	kubeVip := c.KubeVip
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

func (c *Config) validateMesh() error {
	mesh := c.Mesh
	if mesh == nil {
		return fmt.Errorf("config: undefined 'mesh'")
	}
	if mesh.EndpointPort < 1024 {
		return fmt.Errorf("config: 'mesh.endpoint_port=%d' is too low (min=1024)", mesh.EndpointPort)
	}
	if mesh.EndpointPort > 65_535 {
		return fmt.Errorf("config: 'mesh.endpoint_port=%d' is too high (max=65535)", mesh.EndpointPort)
	}
	if mesh.NetworkAddress == nil {
		return fmt.Errorf("config: undefined 'mesh.network_address")
	}
	meshAddr := mesh.NetworkAddress.Addr()
	if !meshAddr.Is4() {
		return fmt.Errorf("config: 'mesh.network_address=%s' is not IPv4", mesh.NetworkAddress)
	}
	if !meshAddr.IsPrivate() {
		return fmt.Errorf("config: 'mesh.network_address=%s' is not private", mesh.NetworkAddress)
	}
	b := mesh.NetworkAddress.Bits() // from /12 (20 bits for host) to /29 (3 bits for host)
	if b < 12 {
		return fmt.Errorf("config: wrong 'mesh.network_address=%s' (bits=%d, min=12)", mesh.NetworkAddress, b)
	}
	if b > 29 {
		return fmt.Errorf("config: wrong 'mesh.network_address=%s' (bits=%d, max=29)", mesh.NetworkAddress, b)
	}
	if len(mesh.WireGuardInterface) < 1 {
		return fmt.Errorf("config: undefined 'mesh.wireguard_interface'")
	}
	return nil
}

func (c *Config) validatePushover() error {
	pushover := c.Pushover
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

func (c *Config) validateSshAuthorizedKeys() error {
	sshAuthorizedKeys := c.SshAuthorizedKeys
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

func (cfg *Config) Validate() error {
	if len(cfg.RcloneRemote) < 1 {
		return fmt.Errorf("config: undefined 'rclone_remote'")
	}
	if err := cfg.validateI12e(); err != nil {
		return err
	}
	if err := cfg.validateK3s(); err != nil {
		return err
	}
	if err := cfg.validatePortKnocking(); err != nil {
		return err
	}
	if err := cfg.validateKubeVip(); err != nil {
		return err
	}
	if err := cfg.validateMesh(); err != nil {
		return err
	}
	if err := cfg.validatePushover(); err != nil {
		return err
	}
	if err := cfg.validateSshAuthorizedKeys(); err != nil {
		return err
	}
	return nil
}
