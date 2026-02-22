package node

import (
	"encoding/hex"
	"fmt"
	"net/netip"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/config"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// ny7a1/20260128_152841_153793688/54a2cc8d5e78755ff1debc4a4e6b2fa657ccf86a868b53f9f1b5140487377cc8/192.168.56.53/51820(/kmain)
var nodeRegex = regexp.MustCompile(
	"^(n[0-9a-z]{4})" +
		"/([0-9]{8}_[0-9]{6}_[0-9]{9})" +
		"/([0-9a-f]{64})" +
		`/([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)` +
		"/([0-9]+)" +
		"(/([^/]+))?" +
		"$",
)

type NodeRemote struct {
	Node
	wgKeyPub  wgtypes.Key
	tsFirst   *time.Time // first record of the series
	tsCurr    *time.Time // current record
	nodeAlias string
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
	nodeAliasStr := ""
	nodeAlias := n.GetNodeAlias()
	if nodeAlias != "" {
		nodeAliasStr = fmt.Sprintf("%s=", nodeAlias)
	}
	return fmt.Sprintf(
		"R|%s%s|%s|%s%s%s|%s",
		nodeAliasStr,
		n.Node.String(),
		n.GetWgKeyPub(),
		tsFirstStr,
		tsDeltaStr,
		tsCurrStr,
		ageStr,
	)
}

func (n *NodeRemote) GetWgKeyPub() wgtypes.Key {
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

func (n *NodeRemote) GetNodeAlias() string {
	return n.nodeAlias
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

func NewNodeRemote(cfg *config.ServerConfig, entry string) (*NodeRemote, error) {
	arr := nodeRegex.FindStringSubmatch(entry)
	if arr == nil {
		return nil, fmt.Errorf("'nodeRegex.FindStringSubmatch(%s)' returned nil", entry)
	}
	nodeId, err := getNodeIdFromNodeName(cfg.Mesh.NetworkAddress, arr[1])
	if err != nil {
		return nil, err
	}
	tsCurr, err := getTimestamp(arr[2])
	if err != nil {
		return nil, err
	}
	buf, err := hex.DecodeString(arr[3])
	if err != nil {
		return nil, err
	}
	wgKeyPub, err := wgtypes.NewKey(buf)
	if err != nil {
		return nil, err
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
		Node: Node{
			cfg:        cfg,
			id:         nodeId,
			wgEndpoint: &wgEndpoint,
		},
		wgKeyPub:  wgKeyPub,
		tsFirst:   nil,
		tsCurr:    tsCurr,
		nodeAlias: arr[7], // ok if empty
	}, nil
}
