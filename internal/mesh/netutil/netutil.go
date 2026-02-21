package netutil

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"strings"

	"github.com/vishvananda/netlink"
)

func ifaceLoad() (*net.Interface, error) {
	path := "/etc/i12e/iface.txt" // TODO: unhardcode
	_, err := os.Stat(path)
	if err != nil {
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
	return net.InterfaceByName(iname)
}

func ifaceGuess() (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		return &iface, nil
	}
	return nil, fmt.Errorf("ifaceGuess(): cannot figure out the interface")
}

func ifaceGet() (*net.Interface, error) {
	iface, err := ifaceLoad()
	if err != nil {
		return ifaceGuess()
	}
	return iface, nil
}

func MeshEndpointAddr() (*netip.Addr, error) {
	iface, err := ifaceGet()
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
		return &addr, nil
	}
	return nil, fmt.Errorf("cannot get mesh endpoint address")
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

func IfaceSyncAddresses(link netlink.Link, ipTarget *netip.Addr, bitsTarget int) error {
	addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		return err
	}
	found := false
	for _, addr := range addrs {
		ipCurrent, ok := netip.AddrFromSlice(addr.IP)
		if !ok {
			return fmt.Errorf("cannot process addr.IP=%v", addr.IP)
		}
		bitsCurrent, _ := addr.Mask.Size()
		if ipTarget.Compare(ipCurrent) == 0 && bitsCurrent == bitsTarget {
			found = true
			continue
		}
		if err := netlink.AddrDel(link, &addr); err != nil {
			return err
		}
	}
	if found {
		return nil
	}
	addr, err := netlink.ParseAddr(fmt.Sprintf("%s/%d", ipTarget, bitsTarget))
	if err != nil {
		return err
	}
	if err := netlink.AddrAdd(link, addr); err != nil {
		return err
	}
	return nil
}
