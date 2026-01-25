package netutil

import (
	"errors"
	"net"
	"os"
	"strings"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type II struct {
	Iface string
	IP    string
}

func IfaceIP() *II {
	path := "/etc/i12e/iface.txt"
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Info("IfaceIP(): file does not exist", "path", path)
		} else {
			log.Error("IfaceIP(): os.Stat() failed", "path", path, "err", err)
		}
		return nil
	}
	buf, err := os.ReadFile(path)
	if err != nil {
		log.Error("IfaceIP(): os.ReadFile() failed", "path", path, "err", err)
		return nil
	}
	iname := strings.TrimSpace(string(buf))
	if len(iname) < 1 {
		log.Error("IfaceIP(): file is empty", "path", path)
		return nil
	}
	iface, err := net.InterfaceByName(iname)
	if err != nil {
		log.Error("IfaceIP(): net.InterfaceByName() failed", "iname", iname, "err", err)
		return nil
	}
	addrs, err := iface.Addrs()
	if err != nil {
		log.Error("IfaceIP(): iface.Addrs() failed", "iname", iname, "err", err)
		return nil
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			log.Info("IfaceIP(): cannot get ipNet", "iname", iname, "addr", addr)
			continue
		}
		ip := ipNet.IP
		if ip.To4() == nil {
			log.Info("IfaceIP(): ip.To4() returned nil", "iname", iname, "addr", addr)
			continue
		}
		ipStr := ip.String()
		if len(ipStr) < 1 {
			log.Info("IfaceIP(): empty ipStr", "iname", iname, "addr", addr)
			continue
		}
		ret := &II{Iface: iname, IP: ip.String()}
		log.Info("IfaceIP() ok", "Iface", ret.Iface, "IP", ret.IP)
		return ret
	}
	log.Warn("IfaceIP(): no IPv4 address found")
	return nil
}
