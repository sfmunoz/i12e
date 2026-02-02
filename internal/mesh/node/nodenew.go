package node

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/netip"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sfmunoz/i12e/internal/mesh/ifaceip"
	"github.com/sfmunoz/i12e/internal/mesh/wgkey"
)

const nodeEndpointPort = 51830 // default '51820'
const nodeIdFname = "/etc/i12e/node-id.txt"

const nodeIdMin uint16 = 100

var nodeIdMax = uint16(addrToU32(netMax(nodeNet)) - addrToU32(netMin(nodeNet)) - 1) // '-1' -> avoid broadcast address

// ny7a1/20260128_152841_153793688/54a2cc8d5e78755ff1debc4a4e6b2fa657ccf86a868b53f9f1b5140487377cc8/192.168.56.53/51820
var nodeRegex = regexp.MustCompile(
	"^(n[0-9a-z]{4})" +
		"/([0-9]{8}_[0-9]{6}_[0-9]{9})" +
		"/([0-9a-f]{64})" +
		`/([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)` +
		"/([0-9]+)" +
		"$",
)

func validNodeId(nodeId uint16) error {
	if nodeId < nodeIdMin {
		return fmt.Errorf("'%s' node-id is '%d' (min=%d)", nodeIdFname, nodeId, nodeIdMin)
	}
	if nodeId > nodeIdMax {
		return fmt.Errorf("'%s' node-id is '%d' (max=%d)", nodeIdFname, nodeId, nodeIdMax)
	}
	return nil
}

func resetNodeLocal() error {
	_, err := os.Stat(nodeIdFname)
	if err == nil {
		return os.Remove(nodeIdFname)
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func procNodeList(nodeList []*Node, nodeLocal *Node, remMeshBase string) error {
	nodeIdLocal := nodeLocal.GetNodeId()
	for i, n := range nodeList {
		if n.GetNodeId() != nodeIdLocal {
			continue
		}
		nodeList[i] = nodeLocal
		return nil
	}
	if err := nodeLocal.pushToRemote(remMeshBase); err != nil {
		return fmt.Errorf("procNodeList(): couldn't find nodeLocal in nodeList but pushToRemote() failed: %q", err)
	}
	// FIXME add nodeLocal to nodeList
	return fmt.Errorf("procNodeList(): couldn't find nodeLocal in nodeList; pushToRemote() performed")
}

func loadNode(nodeList []*Node, remMeshBase string) (*Node, error) {
	buf, err := os.ReadFile(nodeIdFname)
	if err != nil {
		return nil, err
	}
	nodeIdInt, err := strconv.Atoi(strings.TrimSpace(string(buf)))
	if err != nil {
		return nil, err
	}
	nodeId := uint16(nodeIdInt)
	if err := validNodeId(nodeId); err != nil {
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
	nodeLocal := &Node{
		id:         uint16(nodeId),
		wgKey:      wgKey,
		wgEndpoint: &wgEndpoint,
	}
	if err := procNodeList(nodeList, nodeLocal, remMeshBase); err != nil {
		return nil, err
	}
	return nodeLocal, nil
}

func writeNode() error {
	x := uint16(rand.Int32N(int32(nodeIdMax-nodeIdMin+1))) + nodeIdMin
	if err := os.WriteFile(nodeIdFname, fmt.Appendf(make([]byte, 0), "%d\n", x), 0600); err != nil {
		return err
	}
	return nil
}

func getNodeLocal(nodeList []*Node, remMeshBase string) (*Node, error) {
	node, err := loadNode(nodeList, remMeshBase)
	if err == nil {
		return node, nil
	}
	if err := writeNode(); err != nil {
		return nil, err
	}
	return loadNode(nodeList, remMeshBase)
}

func getNodeIdFromNodeName(nodeName string) (uint16, error) {
	nodeNameLen := len(nodeName)
	if nodeNameLen != 5 {
		return 0, fmt.Errorf("len(nodename=%s)=%d (5 expected)", nodeName, nodeNameLen)
	}
	n0 := nodeName[0]
	if n0 != 'n' {
		return 0, fmt.Errorf("nodename='%s' starts with '%c' ('n' expected)", nodeName, n0)
	}
	nodeIdInt64, err := strconv.ParseInt(nodeName[1:], 36, 64)
	if err != nil {
		return 0, err
	}
	nodeId := uint16(nodeIdInt64)
	if err := validNodeId(nodeId); err != nil {
		return 0, err
	}
	return nodeId, nil
}

func getNodeRemote(entry string) (*Node, error) {
	arr := nodeRegex.FindStringSubmatch(entry)
	if arr == nil {
		return nil, fmt.Errorf("'nodeRegex.FindStringSubmatch(%s)' returned nil", entry)
	}
	nodeId, err := getNodeIdFromNodeName(arr[1])
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
		wgKey:      wgKey,
		wgEndpoint: &wgEndpoint,
	}, nil
}

type NodeOption func() (*Node, error)

func WithLocal(nodeList []*Node, remMeshBase string) NodeOption {
	return func() (*Node, error) {
		return getNodeLocal(nodeList, remMeshBase)
	}
}

func WithRemote(entry string) NodeOption {
	return func() (*Node, error) {
		return getNodeRemote(entry)
	}
}

func NewNode(opt NodeOption) (*Node, error) {
	return opt()
}
