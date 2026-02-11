package node

import (
	"fmt"
	"net/netip"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/mesh/wgkey"
)

// ny7a1/20260128_152841_153793688/54a2cc8d5e78755ff1debc4a4e6b2fa657ccf86a868b53f9f1b5140487377cc8/192.168.56.53/51820
var nodeRegex = regexp.MustCompile(
	"^(n[0-9a-z]{4})" +
		"/([0-9]{8}_[0-9]{6}_[0-9]{9})" +
		"/([0-9a-f]{64})" +
		`/([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)` +
		"/([0-9]+)" +
		"$",
)

type NodeRemote struct {
	id         uint32
	meshNet    *netip.Prefix
	wgKeyPriv  *wgkey.WgKeyPriv
	wgKeyPub   *wgkey.WgKeyPub
	wgEndpoint *netip.AddrPort
	tsFirst    *time.Time // first record of the series
	tsCurr     *time.Time // current record
}

func (n *NodeRemote) String() string {
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
	return fmt.Sprintf(
		"R|%s|%s/%d|%s|%s|%s%s%s|%s",
		n.GetNodeName(),
		n.GetMeshIP(),
		n.GetMeshNet().Bits(),
		n.GetWgKeyPub(),
		n.GetWgEndpoint(),
		tsFirstStr,
		tsDeltaStr,
		tsCurrStr,
		ageStr,
	)
}

func (n *NodeRemote) GetNodeName() string {
	return getNodeNameFromNodeId(n.id)
}

func (n *NodeRemote) GetMeshIP() *netip.Addr {
	x, _ := nodeIdToIp(n.GetMeshNet(), n.id) // err ignored: already validated
	return x
}

func (n *NodeRemote) GetMeshNet() *netip.Prefix {
	return n.meshNet
}

func (n *NodeRemote) GetWgEndpoint() *netip.AddrPort {
	return n.wgEndpoint
}

func (n *NodeRemote) GetWgKeyPub() *wgkey.WgKeyPub {
	return n.wgKeyPub
}

func (n *NodeRemote) GetTsFirst() *time.Time {
	return n.tsFirst
}

func (n *NodeRemote) SetTsFirst(ts *time.Time) {
	n.tsFirst = ts
}

func (n *NodeRemote) GetTsCurr() *time.Time {
	return n.tsCurr
}

func (n *NodeRemote) SetTsCurr(ts *time.Time) {
	n.tsCurr = ts
}

func (n *NodeRemote) GetAge(tsNow *time.Time) *time.Duration {
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

func (n *NodeRemote) GetDelta() *time.Duration {
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

func getTimestamp(tsStrIn string) (*time.Time, error) {
	ts, err := time.Parse(
		"20060102.150405.999999999",
		strings.ReplaceAll(tsStrIn, "_", "."), // to properly parse ns
	)
	if err != nil {
		return nil, err
	}
	return &ts, nil
}

func NewNodeRemote(meshNet *netip.Prefix, entry string) (*NodeRemote, error) {
	arr := nodeRegex.FindStringSubmatch(entry)
	if arr == nil {
		return nil, fmt.Errorf("'nodeRegex.FindStringSubmatch(%s)' returned nil", entry)
	}
	nodeId, err := getNodeIdFromNodeName(meshNet, arr[1])
	if err != nil {
		return nil, err
	}
	tsCurr, err := getTimestamp(arr[2])
	if err != nil {
		return nil, err
	}
	k32, err := wgkey.NewK32(wgkey.WithHex(arr[3]))
	if err != nil {
		return nil, fmt.Errorf("'wgkey.NewK32(%s)' failed: %s", arr[3], err)
	}
	addr, err := netip.ParseAddr(arr[4])
	if err != nil {
		return nil, fmt.Errorf("'netip.ParseAddr(%s)' failed: %s", arr[4], err)
	}
	wgEndpointPort, err := strconv.Atoi(arr[5])
	if err != nil {
		return nil, fmt.Errorf("'strconv.Atoi(%s)' failed: %s", arr[5], err)
	}
	wgEndpoint := netip.AddrPortFrom(addr, uint16(wgEndpointPort))
	return &NodeRemote{
		id:         nodeId,
		meshNet:    meshNet,
		wgKeyPriv:  nil,
		wgKeyPub:   wgkey.NewWgKeyPub(k32),
		wgEndpoint: &wgEndpoint,
		tsFirst:    nil,
		tsCurr:     tsCurr,
	}, nil
}
