package config

import (
	"fmt"
)

type Env string

const (
	EnvNone Env = "none"
	EnvDev  Env = "dev"
	EnvProd Env = "prod"
)

var envs = []Env{EnvNone, EnvDev, EnvProd}

func (e Env) String() string {
	return string(e)
}

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
