package node

import (
	"fmt"
	"net/netip"
	"strconv"
)

const nodeEndpointPort = 51830 // default '51820'
const nodeInterface = "wgi"
const nodePrivKeyFname = "/etc/i12e/wg-priv-key"

type Node struct {
	id         uint32
	meshNet    *netip.Prefix
	wgEndpoint *netip.AddrPort
}

func (n *Node) GetNodeName() string {
	return getNodeNameFromNodeId(n.id)
}

func (n *Node) GetMeshNet() *netip.Prefix {
	return n.meshNet
}

func (n *Node) GetMeshIP() *netip.Addr {
	x, _ := nodeIdToIp(n.GetMeshNet(), n.id) // err ignored: already validated
	return x
}

func (n *Node) GetWgEndpoint() *netip.AddrPort {
	return n.wgEndpoint
}

func GetNodeInterface() string {
	return nodeInterface
}

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
