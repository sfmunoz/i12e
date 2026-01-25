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
	ii, err := netutil.IfaceIP()
	if err != nil {
		log.Error("Net.run(): 'netutil.IfaceIP()' failed", "err", err)
		return err
	}
	iface, ip := "", ""
	if ii != nil {
		iface, ip = ii.Iface, ii.IP
	}
	log.Info("Net.run()", "iface", iface, "ip", ip)
	return nil
}

func Run() error {
	return newNet().run()
}
