package config

import (
	"fmt"
)

type Mode string

const (
	ModeMain   Mode = "main"
	ModeServer Mode = "server"
	ModeAgent  Mode = "agent"
)

var modes = []Mode{ModeMain, ModeServer, ModeAgent}

func (m Mode) String() string {
	return string(m)
}

func ValidModes() []string {
	ret := make([]string, len(modes))
	for i, v := range modes {
		ret[i] = string(v)
	}
	return ret
}

func GetMode(m string) (Mode, error) {
	for _, v := range modes {
		if string(v) == m {
			return v, nil
		}
	}
	return ModeMain, fmt.Errorf("unknown mode '%s' (valid: %q)", m, ValidModes())
}
