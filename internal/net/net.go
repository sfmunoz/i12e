package net

import (
	"github.com/sfmunoz/i12e/internal/netutil"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

func Net() {
	log.Info("Net()...")
	iface := ""
	ip := ""
	ii := netutil.IfaceIP()
	if ii != nil {
		iface = ii.Iface
		ip = ii.IP
	}
	log.Info("Net()", "iface", iface, "ip", ip)
}
