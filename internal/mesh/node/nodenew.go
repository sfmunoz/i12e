package node

import (
	"errors"
	"fmt"
	"net/netip"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sfmunoz/i12e/internal/mesh/ifaceip"
	"github.com/sfmunoz/i12e/internal/mesh/wgkey"
)

const nodeEndpointPort = 51830 // default '51820'
const nodeEtcHostname = "/etc/hostname"

// ny7a1/20260128_152841_153793688/54a2cc8d5e78755ff1debc4a4e6b2fa657ccf86a868b53f9f1b5140487377cc8/192.168.56.53/51820
var nodeRegex = regexp.MustCompile(
	"^(n[0-9a-z]{4})" +
		"/([0-9]{8}_[0-9]{6}_[0-9]{9})" +
		"/([0-9a-f]{64})" +
		`/([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)` +
		"/([0-9]+)" +
		"$",
)

func getNodeNameFromNodeId(nodeId uint32) string {
	s := strconv.FormatUint(uint64(nodeId), 36)
	for len(s) < 4 {
		s = "0" + s
	}
	return "n" + s
}

func getNodeIdFromNodeName(meshNet *netip.Prefix, nodeName string) (uint32, error) {
	nodeNameLen := len(nodeName)
	if nodeNameLen != 5 {
		return 0, fmt.Errorf("len(nodename=%s)=%d (5 expected)", nodeName, nodeNameLen)
	}
	n0 := nodeName[0]
	if n0 != 'n' {
		return 0, fmt.Errorf("nodename='%s' starts with '%c' ('n' expected)", nodeName, n0)
	}
	nodeIdInt64, err := strconv.ParseInt(nodeName[1:], 36, 32)
	if err != nil {
		return 0, err
	}
	nodeId := uint32(nodeIdInt64)
	if err := nodeIdValid(meshNet, nodeId); err != nil {
		return 0, err
	}
	return nodeId, nil
}

func deleteEtcHostname() error {
	_, err := os.Stat(nodeEtcHostname)
	if err == nil {
		return os.Remove(nodeEtcHostname)
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func writeRandomEtcHostname(meshNet *netip.Prefix) error {
	nodeId, err := getRandomNodeId(meshNet)
	if err != nil {
		return err
	}
	nodeName := getNodeNameFromNodeId(nodeId)
	if err := os.WriteFile(nodeEtcHostname, fmt.Appendf(make([]byte, 0), "%s\n", nodeName), 0644); err != nil {
		return err
	}
	return nil
}

func getNodeLocal(meshNet *netip.Prefix, reset bool) (*Node, error) {
	if reset {
		if err := deleteEtcHostname(); err != nil {
			return nil, err
		}
	}
	buf, err := os.ReadFile(nodeEtcHostname)
	if err != nil {
		if err := writeRandomEtcHostname(meshNet); err != nil {
			return nil, err
		}
	}
	buf, err = os.ReadFile(nodeEtcHostname)
	if err != nil {
		return nil, err
	}
	nodeId, err := getNodeIdFromNodeName(meshNet, strings.TrimSpace(string(buf)))
	if err != nil {
		return nil, err
	}
	wgKey, err := wgkey.NewWgKey(wgkey.WithLocal(nodePrivKeyFname))
	if err != nil {
		return nil, err
	}
	ii, err := ifaceip.IfaceIP()
	if err != nil {
		return nil, err
	}
	if ii == nil {
		return nil, fmt.Errorf("IfaceIP() returned empty value")
	}
	wgEndpoint := netip.AddrPortFrom(*ii.IP, nodeEndpointPort)
	return &Node{
		id:         nodeId,
		meshNet:    meshNet,
		wgKey:      wgKey,
		wgEndpoint: &wgEndpoint,
	}, nil
}

func getNodeRemote(meshNet *netip.Prefix, entry string) (*Node, error) {
	arr := nodeRegex.FindStringSubmatch(entry)
	if arr == nil {
		return nil, fmt.Errorf("'nodeRegex.FindStringSubmatch(%s)' returned nil", entry)
	}
	nodeId, err := getNodeIdFromNodeName(meshNet, arr[1])
	if err != nil {
		return nil, err
	}
	wgKey, err := wgkey.NewWgKey(wgkey.WithRemote((arr[3])))
	if err != nil {
		return nil, fmt.Errorf("'getWgKeyRemote(%s)' failed: %s", arr[3], err)
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
	return &Node{
		id:         nodeId,
		meshNet:    meshNet,
		wgKey:      wgKey,
		wgEndpoint: &wgEndpoint,
	}, nil
}

type NodeOption func() (*Node, error)

func WithLocal(meshNet *netip.Prefix, reset bool) NodeOption {
	return func() (*Node, error) {
		return getNodeLocal(meshNet, reset)
	}
}

func WithRemote(meshNet *netip.Prefix, entry string) NodeOption {
	return func() (*Node, error) {
		return getNodeRemote(meshNet, entry)
	}
}

func NewNode(opt NodeOption) (*Node, error) {
	return opt()
}
