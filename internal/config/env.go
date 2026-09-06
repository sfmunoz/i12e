package config

import (
	"fmt"
	"os"
	"os/exec"
)

type Env string

const (
	EnvNone Env = "none"
	EnvDev  Env = "dev"
	EnvProd Env = "prod"
)

var envs = []Env{EnvNone, EnvDev, EnvProd}

func ValidEnvs() []string {
	ret := make([]string, len(envs))
	for i, v := range envs {
		ret[i] = string(v)
	}
	return ret
}

func GetEnv(e string) (Env, error) {
	for _, v := range envs {
		if string(v) == e {
			return v, nil
		}
	}
	return EnvNone, fmt.Errorf("unknown env '%s' (valid: %q)", e, ValidEnvs())
}

func (e Env) String() string {
	return string(e)
}

func (e Env) SopsCmd(arg ...string) *exec.Cmd {
	return e.innerCmd("./scripts/sops.sh", arg...)
}

func (e Env) RcloneCmd(arg ...string) *exec.Cmd {
	return e.innerCmd("./scripts/rclone.sh", arg...)
}

func (e Env) innerCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	if e == EnvProd {
		cmd.Env = append(os.Environ(), "I12E_ENV=prod")
	} else {
		cmd.Env = os.Environ()
	}
	return cmd
}
