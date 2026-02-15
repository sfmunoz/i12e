package netutil

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os"
	"strings"

	"github.com/vishvananda/netlink"
)

type II struct {
	Iface string
	IP    *netip.Addr
}

func ifaceTxtLoad() (string, error) {
	path := "/etc/i12e/iface.txt"
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	buf, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	iname := strings.TrimSpace(string(buf))
	if len(iname) < 1 {
		return "", fmt.Errorf("file '%s' is empty", path)
	}
	return iname, nil
}

func ifaceNameGuess() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		return iface.Name, nil
	}
	return "", fmt.Errorf("cannot guess interface name")
}

func IfaceIP() (*II, error) {
	iname, err := ifaceTxtLoad()
	if err != nil {
		return nil, err
	}
	if iname == "" {
		var err error
		iname, err = ifaceNameGuess()
		if err != nil {
			return nil, err
		}
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

func IfaceCreate(ifaceName string) (netlink.Link, error) {
	if link, err := netlink.LinkByName(ifaceName); err == nil {
		return link, nil
	}
	la := netlink.NewLinkAttrs()
	la.Name = ifaceName
	wgi := &netlink.Wireguard{LinkAttrs: la}
	if err := netlink.LinkAdd(wgi); err != nil {
		return nil, err
	}
	link, err := netlink.LinkByName(ifaceName)
	if err != nil {
		return nil, err
	}
	return link, nil
}
