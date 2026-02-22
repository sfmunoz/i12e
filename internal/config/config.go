package config

import (
	"net/netip"
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
	Env  Env `validate:"i12e_env"`
	I12e *struct {
		Version   string `mapstructure:"version" validate:"i12e_semver_v"` // "semver" doesn't accept the leading v
		Sha256sum string `mapstructure:"sha256sum" validate:"sha256"`
	} `mapstructure:"i12e" validate:"required"`
	K3s *struct {
		Token      string `mapstructure:"token" validate:"gte=20"`
		AgentToken string `mapstructure:"agent_token" validate:"gte=20"`
		TlsSan     string `mapstructure:"tls_san" validate:"gte=1"`
	} `mapstructure:"k3s" validate:"required"`
	Rclone *struct {
		Remote string `mapstructure:"remote" validate:"gte=1"`
	} `mapstructure:"rclone" validate:"required"`
	Artifact *struct {
		PortKnocking []int `mapstructure:"port_knocking" validate:"len=4"` // valid: 4, see https://github.com/sfmunoz/i12e/issues/138
	} `mapstructure:"artifact" validate:"required"`
	KubeVip *struct {
		Vip       string `mapstructure:"vip" validate:"ip4_addr"`
		Interface string `mapstructure:"interface" validate:"gte=2"`
		Kvversion string `mapstructure:"kvversion" validate:"required,i12e_semver_v"` // "semver" doesn't accept the leading v
	} `mapstructure:"kube_vip"` // not required
	Mesh *struct {
		EndpointInterface     string        `mapstructure:"endpoint_interface"` // not required
		EndpointPort          int           `mapstructure:"endpoint_port" validate:"gte=1024,lt=65535"`
		NetworkAddress        *netip.Prefix `mapstructure:"network_address" validate:"required,i12e_mesh_network"`
		WireGuardInterface    string        `mapstructure:"wireguard_interface" validate:"gte=1"`
		WireGuardPrivKeyFname string        `mapstructure:"wireguard_priv_key_fname" validate:"gte=1"`
		RemoteBase            string        `mapstructure:"remote_base" validate:"startsnotwith=:,contains=:,endsnotwith=:"`
	} `mapstructure:"mesh" validate:"required"`
	Pushover *struct {
		UserKey string `mapstructure:"user_key" validate:"gte=1"`
		Token   string `mapstructure:"token" validate:"gte=1"`
	} `mapstructure:"pushover" validate:"required"`
	Butane *struct {
		SshAuthorizedKeys []string `mapstructure:"ssh_authorized_keys" validate:"gte=1,dive,gte=1"`
		Mode              Mode     `mapstructure:"mode" validate:"i12e_butane_mode"`
		Output            Output   `mapstructure:"output" validate:"i12e_butane_output"`
	} `mapstructure:"butane"` // not required
	Server *struct {
		SlumberBase   time.Duration `mapstructure:"slumber_base" validate:"gte=10s"`
		SlumberJitter time.Duration `mapstructure:"slumber_jitter" validate:"gte=1s"`
	} `mapstructure:"server" validate:"required"`
}

func (cfg *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("i12e_semver_v", validSemverV)
	validate.RegisterValidation("i12e_env", validEnv)
	validate.RegisterValidation("i12e_butane_mode", validButaneMode)
	validate.RegisterValidation("i12e_butane_output", validButaneOutput)
	validate.RegisterValidation("i12e_mesh_network", validMeshNetwork)
	if err := validate.Struct(cfg); err != nil {
		return err
	}
	return nil
}
