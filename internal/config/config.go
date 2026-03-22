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

type ArtifactConfig struct {
	Env      Env            `validate:"i12e_env"`
	Artifact *artifact      `mapstructure:"artifact" validate:"required"`
	K3s      *k3s           `mapstructure:"k3s" validate:"required"`
	Mesh     *mesh          `mapstructure:"mesh" validate:"required"`
	Rclone   *rclone        `mapstructure:"rclone" validate:"required"`
	WikiJs         map[string]any `mapstructure:"wikijs" validate:"required"`
	AnkiSyncServer map[string]any `mapstructure:"anki_sync_server" validate:"required"`
	WgConf         []string       `mapstructure:"wg_conf"`
}

type ButaneConfig struct {
	Env    Env     `validate:"i12e_env"`
	Butane *butane `mapstructure:"butane" validate:"required"`
	// I12e   *i12e   `mapstructure:"i12e" validate:"required"`
	Rclone *rclone `mapstructure:"rclone" validate:"required"`
}

type ServerConfig struct {
	Env    Env     `validate:"i12e_env"`
	Server *server `mapstructure:"server" validate:"required"`
	K3s    *k3s    `mapstructure:"k3s" validate:"required"`
	Mesh   *mesh   `mapstructure:"mesh" validate:"required"`
}

// type i12e struct {
// 	Version string `mapstructure:"version" validate:"i12e_semver_v"` // "semver" doesn't accept the leading v
// }

type artifact struct {
	PortKnocking  []int  `mapstructure:"port_knocking" validate:"len=4"` // valid: 4, see https://github.com/sfmunoz/i12e/issues/138
	K3sToken      string `mapstructure:"k3s_token" validate:"gte=20"`
	K3sAgentToken string `mapstructure:"k3s_agent_token" validate:"gte=20"`
}

type butane struct {
	SshAuthorizedKeys []string `mapstructure:"ssh_authorized_keys" validate:"gte=1,dive,gte=1"`
	Version           string   `validate:"gte=1"`
	Mode              Mode     `mapstructure:"mode" validate:"i12e_butane_mode"`
	Output            Output   `mapstructure:"output" validate:"i12e_butane_output"`
}

type server struct {
	SlumberBase   time.Duration `mapstructure:"slumber_base" validate:"gte=10s"`
	SlumberJitter time.Duration `mapstructure:"slumber_jitter" validate:"gte=1s"`
}

type k3s struct {
	TlsSan string `mapstructure:"tls_san" validate:"gte=1"`
}

type rclone struct {
	Remote string `mapstructure:"remote" validate:"gte=1"`
}

type mesh struct {
	EndpointInterface     string        `mapstructure:"endpoint_interface"` // not required
	EndpointPort          int           `mapstructure:"endpoint_port" validate:"gte=1024,lt=65535"`
	NetworkAddress        *netip.Prefix `mapstructure:"network_address" validate:"required,i12e_mesh_network"`
	WireGuardInterface    string        `mapstructure:"wireguard_interface" validate:"gte=1"`
	WireGuardPrivKeyFname string        `mapstructure:"wireguard_priv_key_fname" validate:"gte=1"`
	RemoteBase            string        `mapstructure:"remote_base" validate:"startsnotwith=:,contains=:,endsnotwith=:"`
}

// type pushover struct {
// 	UserKey string `mapstructure:"user_key" validate:"gte=1"`
// 	Token   string `mapstructure:"token" validate:"gte=1"`
// }
//
// type kubeVip struct {
// 	Vip       string `mapstructure:"vip" validate:"ip4_addr"`
// 	Interface string `mapstructure:"interface" validate:"gte=2"`
// 	Kvversion string `mapstructure:"kvversion" validate:"required,i12e_semver_v"` // "semver" doesn't accept the leading v
// }

func (cfg *ArtifactConfig) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	registerValidations(v)
	return v.Struct(cfg)
}

func (cfg *ButaneConfig) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	registerValidations(v)
	return v.Struct(cfg)
}

func (cfg *ServerConfig) Validate() error {
	v := validator.New(validator.WithRequiredStructEnabled())
	registerValidations(v)
	return v.Struct(cfg)
}
