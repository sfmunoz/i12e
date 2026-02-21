package node

import (
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"

	"github.com/sfmunoz/i12e/internal/config"
)

const nodePrivKeyFname = "/etc/i12e/wg-priv-key" // TODO: unhardcode

type Node struct {
	cfg        *config.Config
	id         uint32
	wgEndpoint *netip.AddrPort
}

func (n *Node) String() string {
	return fmt.Sprintf(
		"%s|%s/%d|%s",
		n.GetNodeName(),
		n.GetMeshIP(),
		n.cfg.Mesh.NetworkAddress.Bits(),
		n.GetWgEndpoint(),
	)
}

func (n *Node) GetNodeName() string {
	return getNodeNameFromNodeId(n.id)
}

func (n *Node) GetMeshIP() *netip.Addr {
	x, _ := nodeIdToIp(n.cfg.Mesh.NetworkAddress, n.id) // err ignored: already validated
	return x
}

func (n *Node) GetWgEndpoint() *netip.AddrPort {
	return n.wgEndpoint
}

func getNodeNameFromNodeId(nodeId uint32) string {
	s := strconv.FormatUint(uint64(nodeId), 36)
	for len(s) < 4 {
		s = "0" + s
	}
	return "n" + s
}

func getNodeIdFromNodeName(p *netip.Prefix, nodeName string) (uint32, error) {
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
	if err := nodeIdValid(p, nodeId); err != nil {
		return 0, err
	}
	return nodeId, nil
}

func getEtcI12eMode() (string, error) {
	fname := "/etc/i12e/mode" // TODO: unhardcode
	buf, err := os.ReadFile(fname)
	if err != nil {
		return "", err
	}
	mode := strings.TrimSpace(string(buf))
	if len(mode) < 1 {
		return "", fmt.Errorf("'%s' file is empty", fname)
	}
	return mode, nil
}
