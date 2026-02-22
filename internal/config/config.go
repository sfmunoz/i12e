package config

import (
	"fmt"
	"net/netip"
	"slices"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Docs refs:
//   https://pkg.go.dev/github.com/go-playground/validator/v10
//   ~/go/pkg/mod/github.com/go-playground/validator/v10@v10.30.1/baked_in.go
//   ~/go/pkg/mod/github.com/go-playground/validator/v10@v10.30.1/_examples/simple/main.go
// Home:
//   https://github.com/go-playground/validator

type Config struct {
	Env  Env
	I12e *struct {
		Version   string `mapstructure:"version" validate:"gte=6,lowercase,startswith=v"` // gte=6 -> v0.0.0
		Sha256sum string `mapstructure:"sha256sum" validate:"len=64,lowercase"`
	} `mapstructure:"i12e" validate:"required"`
	K3s *struct {
		Token      string `mapstructure:"token" validate:"gte=20"`
		AgentToken string `mapstructure:"agent_token" validate:"gte=20"`
		TlsSan     string `mapstructure:"tls_san" validate:"gte=1"`
	} `mapstructure:"k3s" validate:"required"`
	RcloneRemote string `mapstructure:"rclone_remote"`
	PortKnocking []int  `mapstructure:"port_knocking"`
	KubeVip      *struct {
		Vip       string `mapstructure:"vip"`
		Interface string `mapstructure:"interface"`
		Kvversion string `mapstructure:"kvversion"`
	} `mapstructure:"kube_vip"`
	Mesh *struct {
		EndpointInterface     string        `mapstructure:"endpoint_interface"`
		EndpointPort          int           `mapstructure:"endpoint_port"`
		NetworkAddress        *netip.Prefix `mapstructure:"network_address"`
		WireGuardInterface    string        `mapstructure:"wireguard_interface"`
		WireGuardPrivKeyFname string        `mapstructure:"wireguard_priv_key_fname"`
		RemoteBase            string        `mapstructure:"remote_base"`
	} `mapstructure:"mesh"`
	Pushover *struct {
		UserKey string `mapstructure:"user_key"`
		Token   string `mapstructure:"token"`
	} `mapstructure:"pushover"`
	SshAuthorizedKeys []string `mapstructure:"ssh_authorized_keys"`
	Butane            *struct {
		Mode   Mode   `mapstructure:"mode"`
		Output Output `mapstructure:"output"`
	} `mapstructure:"butane"`
	Server *struct {
		SlumberBase   time.Duration `mapstructure:"slumber_base"`
		SlumberJitter time.Duration `mapstructure:"slumber_jitter"`
	} `mapstructure:"server"`
}

func (c *Config) fname(bname string) string {
	if c.Env == EnvNone {
		return fmt.Sprintf("/etc/i12e/%s", bname)
	}
	return fmt.Sprintf("config/%s/%s", c.Env.String(), bname)
}

func (c *Config) ButaneEncYaml() string {
	return c.fname("butane.enc.yaml")
}

func (c *Config) I12eYaml() string {
	return c.fname("i12e.yaml")
}

func (c *Config) I12eEncYaml() string {
	return c.fname("i12e.enc.yaml")
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
	if len(mesh.WireGuardPrivKeyFname) < 1 {
		return fmt.Errorf("config: undefined 'mesh.wireguard_priv_key_fname'")
	}
	if len(mesh.RemoteBase) < 1 {
		return fmt.Errorf("config: undefined 'mesh.remote_base'")
	}
	parts := strings.Split(mesh.RemoteBase, ":")
	partsLen := len(parts)
	if partsLen != 2 {
		return fmt.Errorf("config: wrong 'mesh.remote_base=%s (parts=%d, required=2)'", mesh.RemoteBase, partsLen)
	}
	for k, v := range parts {
		if len(v) < 1 {
			return fmt.Errorf("config: wrong 'mesh.remote_base=%s (part[%d] is undefined2)'", mesh.RemoteBase, k)
		}
	}
	strings.Contains(mesh.RemoteBase, ":")
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

func (c *Config) validateButane() error {
	butane := c.Butane
	if butane == nil {
		return nil
	}
	if len(butane.Mode) < 1 {
		return fmt.Errorf("config: undefined 'butane.mode'")
	}
	validModes := ValidModes()
	if !slices.Contains(validModes, butane.Mode.String()) {
		return fmt.Errorf("config: invalid 'butane.mode=%s' (valid: %q)", butane.Mode.String(), validModes)
	}
	if len(butane.Output) < 1 {
		return fmt.Errorf("config: undefined 'butane.output'")
	}
	validOutputs := ValidOutputs()
	if !slices.Contains(validOutputs, butane.Output.String()) {
		return fmt.Errorf("config: invalid 'butane.output=%s' (valid: %q)", butane.Output.String(), validOutputs)
	}
	return nil
}

func (c *Config) validateServer() error {
	server := c.Server
	if server == nil {
		return fmt.Errorf("config: undefined 'server'")
	}
	if server.SlumberBase < 10*time.Second {
		return fmt.Errorf("config: wrong 'mesh.slumber_base=%s' (min=10s)", server.SlumberBase)
	}
	if server.SlumberJitter < time.Second {
		return fmt.Errorf("config: wrong 'mesh.slumber_jitter=%s' (min=1s)", server.SlumberJitter)
	}
	return nil
}

func (cfg *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(cfg); err != nil {
		return err
	}
	if len(cfg.RcloneRemote) < 1 {
		return fmt.Errorf("config: undefined 'rclone_remote'")
	}
	flist := []func() error{
		cfg.validatePortKnocking,
		cfg.validateKubeVip,
		cfg.validateMesh,
		cfg.validatePushover,
		cfg.validateSshAuthorizedKeys,
		cfg.validateButane,
		cfg.validateServer,
	}
	for _, f := range flist {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
