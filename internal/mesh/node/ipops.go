package node

import (
	"net/netip"
)

func netMin(net *netip.Prefix) *netip.Addr {
	addr := net.Addr()
	return &addr
}

func netMax(net *netip.Prefix) *netip.Addr {
	addr := net.Addr()
	u := addrToU32(&addr)
	x := uint32(0xffffffff) >> net.Bits()
	addr2 := u32ToAddr(u | x)
	return addr2
}

func addrToU32(addr *netip.Addr) uint32 {
	return b4ToU32(addr.As4()) // panic("As4 called on IPv6 address")
}

func u32ToAddr(u uint32) *netip.Addr {
	addr := netip.AddrFrom4(u32ToB4(u))
	return &addr
}

func b4ToU32(ip [4]byte) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func u32ToB4(u uint32) [4]byte {
	return [4]byte{byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u)}
}
