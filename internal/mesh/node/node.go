package node

import (
	"fmt"
	"net/netip"
	"time"

	"github.com/sfmunoz/i12e/internal/mesh/wgkey"
)

const nodeEndpointPort = 51830 // default '51820'
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

func (n *Node) GetLocal() bool {
	return n.GetWgKeyPriv() != nil
}

func (n *Node) GetWgEndpoint() *netip.AddrPort {
	return n.wgEndpoint
}
