package config

import (
	"fmt"
)

type Mode int

const (
	ModeMain Mode = iota
	ModeServer
	ModeAgent
)

var modeMap = map[Mode]string{
	ModeMain:   "main",
	ModeServer: "server",
	ModeAgent:  "agent",
}

func (m Mode) String() string {
	return modeMap[m]
}

func ValidModes() []string {
	return []string{ModeMain.String(), ModeServer.String(), ModeAgent.String()}
}

func GetMode(m string) (Mode, error) {
	for k, v := range modeMap {
		if v == m {
			return k, nil
		}
	}
	return ModeMain, fmt.Errorf("unknown mode '%s' (valid: %q)", m, ValidModes())
}
