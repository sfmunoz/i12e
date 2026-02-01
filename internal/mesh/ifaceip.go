package mesh

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os"
	"strings"
)

type II struct {
	Iface string
	IP    *netip.Addr
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
		ip4 := ipNet.IP.To4()
		if ip4 == nil {
			continue
		}
		addr := netip.AddrFrom4([4]byte(ip4))
		return &II{Iface: iname, IP: &addr}, nil
	}
	return nil, nil
}
