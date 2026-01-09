package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sfmunoz/i12e/internal/install"
	"github.com/sfmunoz/i12e/internal/pull"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "main")

func ip4addrs(iname string) ([]string, error) {
	iface, err := net.InterfaceByName(iname)
	if err != nil {
		return nil, fmt.Errorf("net.InterfaceByName() failed: %s", err.Error())
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("iface.Addrs() failed: %s", err.Error())
	}
	ret := make([]string, 0)
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipNet.IP
		if ip.To4() == nil {
			continue
		}
		ret = append(ret, ip.String())
	}
	if len(ret) < 1 {
		return nil, fmt.Errorf("no IPv4 address found")
	}
	return ret, nil
}

func main() {
	if os.Getenv("I12E_INSTALL") == "1" {
		install.I12eInstall()
		return
	}
	slumber := 3 * time.Second
	for {
		log.Info("i12e running...")
		install.K3sInstall()
		pull.Pull()
		log.Info("i12e sleeping...", "slumber", slumber)
		addr, err := ip4addrs("enp0s8")
		if err != nil {
			log.Error("ip4addrs() failed", "err", err)
		} else {
			log.Info("ip4addrs()", "addr", addr[0])
		}
		time.Sleep(slumber)
	}
}
