package netutil

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

type II struct {
	Iface string
	IP    string
}

func IfaceIP() (*II, error) {
	path := "/etc/i12e/iface.txt"
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	iname := strings.TrimSpace(string(buf))
	if len(iname) < 1 {
		return nil, fmt.Errorf("file '%s' is empty", path)
	}
	iface, err := net.InterfaceByName(iname)
	if err != nil {
		return nil, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipNet.IP
		if ip.To4() == nil {
			continue
		}
		ipStr := ip.String()
		if len(ipStr) < 1 {
			continue
		}
		return &II{Iface: iname, IP: ip.String()}, nil
	}
	return nil, nil
}
