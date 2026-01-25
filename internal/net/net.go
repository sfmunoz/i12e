package net

import (
	"github.com/sfmunoz/i12e/internal/netutil"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Net struct{}

func newNet() *Net {
	return &Net{}
}

func (n *Net) run() error {
	log.Info("Net()...")
	iface := ""
	ip := ""
	ii := netutil.IfaceIP()
	if ii != nil {
		iface = ii.Iface
		ip = ii.IP
	}
	log.Info("Net.run()", "iface", iface, "ip", ip)
	return nil
}

func Run() error {
	return newNet().run()
}
