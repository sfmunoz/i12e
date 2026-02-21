package config

type Env int

const (
	EnvNone Env = iota
	EnvDev
	EnvProd
)

func (e Env) String() string {
	if e == EnvDev {
		return "dev"
	}
	if e == EnvProd {
		return "prod"
	}
	return "none"
}
