package config

import (
	"fmt"
)

type Bout string
type Mode string

const (
	BoutBashB64  Bout = "bash_b64"
	BoutBashRaw  Bout = "bash_raw"
	BoutIgnition Bout = "ignition"
	BoutDebug    Bout = "debug"
)

const (
	ModeMain   Mode = "main"
	ModeServer Mode = "server"
	ModeAgent  Mode = "agent"
)

var bouts = []Bout{BoutBashB64, BoutBashRaw, BoutIgnition, BoutDebug}
var modes = []Mode{ModeMain, ModeServer, ModeAgent}

func (b Bout) String() string {
	return string(b)
}

func ValidBouts() []string {
	ret := make([]string, len(bouts))
	for i, v := range bouts {
		ret[i] = string(v)
	}
	return ret
}

func GetBout(b string) (Bout, error) {
	for _, v := range bouts {
		if string(v) == b {
			return v, nil
		}
	}
	return BoutBashB64, fmt.Errorf("unknown butane output '%s' (valid: %q)", b, ValidBouts())
}

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
