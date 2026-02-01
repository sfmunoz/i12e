package node

import (
	"net/netip"
	"testing"
)

func TestBackAndForth(t *testing.T) {
	addrs := [][4]byte{
		{192, 168, 88, 99},
		{10, 20, 18, 49},
	}
	for _, b4 := range addrs {
		addr := netip.AddrFrom4(b4)
		u, err := addrToU32(&addr)
		if err != nil {
			t.Errorf("addrToU32(%q) failed: %v", addr, err)
		}
		addr2 := u32ToAddr(u)
		res := addr.Compare(*addr2)
		if res != 0 {
			t.Errorf("%q != %q", addr, addr2)
		}
	}
}
