package node

import (
	"testing"

	"net/netip"
)

func TestBackAndForth(t *testing.T) {
	addrs := [][4]byte{
		{192, 168, 88, 99},
		{10, 20, 18, 49},
	}
	for _, b4 := range addrs {
		addr := netip.AddrFrom4(b4)
		u := addrToU32(&addr)
		addr2 := u32ToAddr(u)
		res := addr.Compare(*addr2)
		if res != 0 {
			t.Errorf("%q != %q", addr, addr2)
		}
	}
}

func TestAdd100(t *testing.T) {
	addr0 := netip.AddrFrom4([4]byte{192, 168, 88, 199})
	addr1 := netip.AddrFrom4([4]byte{192, 168, 89, 43})
	u := addrToU32(&addr0)
	addr2 := u32ToAddr(u + 100)
	res := addr2.Compare(addr1)
	if res != 0 {
		t.Errorf("%q != %q", addr1, addr2)
	}
}

func TestNetMinMax(t *testing.T) {
	netPref := netip.MustParsePrefix("10.119.0.0/16")
	addrMin1 := netip.AddrFrom4([4]byte{10, 119, 0, 0})
	addrMin2 := netMin(&netPref)
	res := addrMin1.Compare(*addrMin2)
	if res != 0 {
		t.Errorf("%q != %q", addrMin2, addrMin1)
	}
	addrMax1 := netip.AddrFrom4([4]byte{10, 119, 255, 255})
	addrMax2 := netMax(&netPref)
	res = addrMax1.Compare(*addrMax2)
	if res != 0 {
		t.Errorf("%q != %q", addrMax2, addrMax1)
	}
}
