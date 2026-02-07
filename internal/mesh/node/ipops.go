package node

import (
	"fmt"
	"math/rand/v2"
	"net/netip"
)

func b4ToU32(ip [4]byte) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func u32ToB4(u uint32) [4]byte {
	return [4]byte{byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u)}
}

func addrToU32(addr *netip.Addr) uint32 {
	return b4ToU32(addr.As4()) // panic("As4 called on IPv6 address")
}

func u32ToAddr(u uint32) *netip.Addr {
	addr := netip.AddrFrom4(u32ToB4(u))
	return &addr
}

func netAddr(p *netip.Prefix) *netip.Addr {
	addr := p.Addr()
	u := addrToU32(&addr)
	b := 32 - p.Bits()
	return u32ToAddr((u >> b) << b)
}

func netBcast(p *netip.Prefix) *netip.Addr {
	u := addrToU32(netAddr(p))
	x := uint32(0xffffffff) >> p.Bits()
	return u32ToAddr(u | x)
}

func validBits(p *netip.Prefix) error {
	b := p.Bits()
	if b < 12 {
		return fmt.Errorf("invalid '%s': bits=%b (min=12)", p, b)
	}
	if b > 29 {
		return fmt.Errorf("invalid '%s': bits=%b (max=29)", p, b)
	}
	return nil
}

func nodeIdMin(p *netip.Prefix) (uint32, error) {
	if err := validBits(p); err != nil {
		return 0, err
	}
	m := map[int]uint32{
		29: 4,  // net=0,[1,2,3],[4,5,6],broadcast=7
		28: 6,  // net=0,[1,...,5],[6,...,14],broadcast=15
		27: 10, // net=0,[1,...,9],[10,...,30],broadcast=31
		26: 16, // net=0,[1,...,15],[16,...,62],broadcast=63
		25: 30, // net=0,[1,...,29],[30,...,126],broadcast=127
		24: 50, // net=0,[1,...,49],[50,...,254],broadcast=255
	}
	b := p.Bits()
	ret, ok := m[b]
	if ok {
		return ret, nil
	}
	return 100, nil // net=0,[1,...,99],[100,...,N-1],broadcast=N (N>=511)
}

func nodeIdMax(p *netip.Prefix) (uint32, error) {
	if err := validBits(p); err != nil {
		return 0, err
	}
	x := uint32(0xffffffff) >> p.Bits()
	return x - 1, nil // '-1' -> leave out broadcast address
}

func nodeIdValid(p *netip.Prefix, id uint32) error {
	nmin, err := nodeIdMin(p)
	if err != nil {
		return err
	}
	nmax, err := nodeIdMax(p)
	if err != nil {
		return err
	}
	if id < nmin {
		return fmt.Errorf("nodeIdValid(): id=%d (min=%d)", id, nmin)
	}
	if id > nmax {
		return fmt.Errorf("nodeIdValid(): id=%d (max=%d)", id, nmax)
	}
	return nil
}

func nodeIdToIp(p *netip.Prefix, id uint32) (*netip.Addr, error) {
	if err := nodeIdValid(p, id); err != nil {
		return nil, err
	}
	addr := netAddr(p)
	u := addrToU32(addr)
	return u32ToAddr(u | id), nil
}

func getRandomNodeId(p *netip.Prefix) (uint32, error) {
	nmin, err := nodeIdMin(p)
	if err != nil {
		return 0, err
	}
	nmax, err := nodeIdMax(p)
	if err != nil {
		return 0, err
	}
	return nmin + rand.Uint32N(nmax-nmin+1), nil
}
