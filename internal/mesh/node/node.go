package node

import (
	"fmt"
	"net/netip"
	"time"

	"github.com/sfmunoz/i12e/internal/mesh/wgkey"
)

const nodeInterface = "wgi"
const nodePrivKeyFname = "/etc/i12e/wg-priv-key"

func GetNodeInterface() string {
	return nodeInterface
}

type Node struct {
	id         uint32
	meshNet    *netip.Prefix
	wgKeyPriv  *wgkey.WgKeyPriv
	wgKeyPub   *wgkey.WgKeyPub
	wgEndpoint *netip.AddrPort
	tsFirst    *time.Time // first record of the series
	tsCurr     *time.Time // current record
}

func (n *Node) String() string {
	where := "R"
	isLocal := n.GetLocal()
	if isLocal {
		where = "L"
	}
	tsDetails := ""
	if !isLocal {
		tsFirstStr := "<undefined-tsFirst>"
		tsFirst := n.GetTsFirst()
		if tsFirst != nil {
			tsFirstStr = tsFirst.Format(time.RFC3339Nano)
		}
		tsCurrStr := "<undefined-tsCurr>"
		tsCurr := n.GetTsCurr()
		if tsCurr != nil {
			tsCurrStr = tsCurr.Format(time.RFC3339Nano)
		}
		tsDeltaStr := "|"
		tsDelta := n.GetDelta()
		if tsDelta != nil {
			tsDeltaStr = fmt.Sprintf("--[%s]->", tsDelta.String())
		}
		ageStr := "<undefined-age>"
		age := n.GetAge(nil)
		if age != nil {
			ageStr = age.String()
		}
		tsDetails = fmt.Sprintf("|%s%s%s|%s", tsFirstStr, tsDeltaStr, tsCurrStr, ageStr)
	}
	return fmt.Sprintf(
		"%s|%s|%s/%d|%s|%s%s",
		where,
		n.GetNodeName(),
		n.GetMeshIP(),
		n.GetMeshNet().Bits(),
		n.GetWgKeyPub(),
		n.GetWgEndpoint(),
		tsDetails,
	)
}

func (n *Node) GetNodeName() string {
	return getNodeNameFromNodeId(n.id)
}

func (n *Node) GetMeshIP() *netip.Addr {
	x, _ := nodeIdToIp(n.GetMeshNet(), n.id) // err ignored: already validated
	return x
}

func (n *Node) GetMeshNet() *netip.Prefix {
	return n.meshNet
}

func (n *Node) GetWgKeyPriv() *wgkey.WgKeyPriv {
	return n.wgKeyPriv
}

func (n *Node) GetWgKeyPub() *wgkey.WgKeyPub {
	return n.wgKeyPub
}

func (n *Node) GetLocal() bool {
	return n.GetWgKeyPriv() != nil
}

func (n *Node) GetWgEndpoint() *netip.AddrPort {
	return n.wgEndpoint
}

func (n *Node) GetTsFirst() *time.Time {
	return n.tsFirst
}

func (n *Node) SetTsFirst(ts *time.Time) {
	n.tsFirst = ts
}

func (n *Node) GetTsCurr() *time.Time {
	return n.tsCurr
}

func (n *Node) SetTsCurr(ts *time.Time) {
	n.tsCurr = ts
}

func (n *Node) GetAge(tsNow *time.Time) *time.Duration {
	if tsNow == nil {
		ts := time.Now().UTC()
		tsNow = &ts
	}
	tsFirst := n.GetTsFirst()
	if tsFirst == nil {
		return nil
	}
	d := tsNow.Sub(*tsFirst)
	return &d
}

func (n *Node) GetDelta() *time.Duration {
	tsFirst := n.GetTsFirst()
	if tsFirst == nil {
		return nil
	}
	tsCurr := n.GetTsCurr()
	if tsCurr == nil {
		return nil
	}
	d := tsCurr.Sub(*tsFirst)
	return &d
}
