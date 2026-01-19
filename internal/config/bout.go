package config

import (
	"fmt"
)

type Bout string

const (
	BoutBashB64  Bout = "bash_b64"
	BoutBashRaw  Bout = "bash_raw"
	BoutIgnition Bout = "ignition"
	BoutDebug    Bout = "debug"
)

var bouts = []Bout{BoutBashB64, BoutBashRaw, BoutIgnition, BoutDebug}

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
