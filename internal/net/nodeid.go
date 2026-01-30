package net

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

const nodeIdMin = 100
const nodeIdMax = 65_534 // 65_535 = 255.255 = broadcast address
const nodeIdFname = "/etc/i12e/node-id.txt"
const nodeNet = "10.56"

type nodeId struct {
	nid uint16 // uint16 instead of uint8 to avoid pervasive type conversion
}

func (n *nodeId) tuple() [2]uint8 {
	return [2]uint8{uint8(n.nid / 256), uint8(n.nid % 256)}
}

func (n *nodeId) String() string {
	return fmt.Sprintf("[%s|%s|%s]", n.NodeName(), n.NodeIP(), n.NodePath())
}

func (n *nodeId) NodeId() uint16 {
	return n.nid
}

func (n *nodeId) NodeName() string {
	t := n.tuple()
	return fmt.Sprintf("n-%03d-%03d", t[0], t[1])
}

func (n *nodeId) NodePath() string {
	t := n.tuple()
	return fmt.Sprintf("%03d_%03d", t[0], t[1])
}

func (n *nodeId) NodeIP() string {
	t := n.tuple()
	return fmt.Sprintf("%s.%d.%d", nodeNet, t[0], t[1])
}

func loadNodeId() (*nodeId, error) {
	buf, err := os.ReadFile(nodeIdFname)
	if err != nil {
		return nil, err
	}
	nid, err := strconv.Atoi(strings.TrimSpace(string(buf)))
	if err != nil {
		return nil, err
	}
	if nid < nodeIdMin {
		return nil, fmt.Errorf("'%s' node-id is '%d' (min=%d)", nodeIdFname, nid, nodeIdMin)
	}
	if nid > nodeIdMax {
		return nil, fmt.Errorf("'%s' node-id is '%d' (max=%d)", nodeIdFname, nid, nodeIdMax)
	}
	return &nodeId{nid: uint16(nid)}, nil
}

func writeNodeId() error {
	x := rand.Int32N(nodeIdMax-nodeIdMin+1) + nodeIdMin
	if err := os.WriteFile(nodeIdFname, []byte(fmt.Sprintf("%d\n", x)), 0600); err != nil {
		return err
	}
	return nil
}

func getNodeIdLocal() (*nodeId, error) {
	nodeId, err := loadNodeId()
	if err == nil {
		return nodeId, nil
	}
	if err := writeNodeId(); err != nil {
		return nil, err
	}
	return loadNodeId()
}
