package config

import (
	"fmt"
)

type Output string
type Mode string

const (
	OutputBashB64  Output = "bash_b64"
	OutputBashRaw  Output = "bash_raw"
	OutputIgnition Output = "ignition"
	OutputDebug    Output = "debug"
)

const (
	ModeMain   Mode = "main"
	ModeServer Mode = "server"
	ModeAgent  Mode = "agent"
)

var outputs = []Output{OutputBashB64, OutputBashRaw, OutputIgnition, OutputDebug}
var modes = []Mode{ModeMain, ModeServer, ModeAgent}

func (o Output) String() string {
	return string(o)
}

func ValidOutputs() []string {
	ret := make([]string, len(outputs))
	for i, v := range outputs {
		ret[i] = string(v)
	}
	return ret
}

func GetOutput(o string) (Output, error) {
	for _, v := range outputs {
		if string(v) == o {
			return v, nil
		}
	}
	return OutputBashB64, fmt.Errorf("unknown butane output '%s' (valid: %q)", o, ValidOutputs())
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
