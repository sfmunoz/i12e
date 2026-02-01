package node

import (
	"fmt"
	"net/netip"
)

func addrToU32(addr *netip.Addr) (uint32, error) {
	if addr.Is6() {
		return 0, fmt.Errorf("cannot convert IPv6 address %s to uint32", addr)
	}
	return b4ToU32(addr.As4()), nil
}

func u32ToAddr(u uint32) *netip.Addr {
	addr := netip.AddrFrom4(u32ToB4(u))
	return &addr
}

func b4ToU32(ip [4]byte) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func u32ToB4(u uint32) [4]byte {
	return [4]byte{
		byte(u >> 24),
		byte(u >> 16),
		byte(u >> 8),
		byte(u),
	}
}
