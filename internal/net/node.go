package net

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

const nodeIdMin = 100
const nodeIdMax = 65_534 // 65_535 = 255.255 = broadcast address
const nodeIdFname = "/etc/i12e/node-id.txt"
const nodeNet = "10.56"
const etcHostname = "/etc/hostname"

type node struct {
	id uint16 // uint16 instead of uint8 to avoid pervasive type conversion
}

func (n *node) tuple() [2]uint8 {
	return [2]uint8{uint8(n.id / 256), uint8(n.id % 256)}
}

func (n *node) String() string {
	return fmt.Sprintf("id=%d|name=%s|ip=%s|path=%s", n.NodeId(), n.NodeName(), n.NodeIP(), n.NodePath())
}

func (n *node) NodeId() uint16 {
	return n.id
}

func (n *node) NodeName() string {
	t := n.tuple()
	return fmt.Sprintf("n-%03d-%03d", t[0], t[1])
}

func (n *node) NodePath() string {
	t := n.tuple()
	return fmt.Sprintf("%03d_%03d", t[0], t[1])
}

func (n *node) NodeIP() string {
	t := n.tuple()
	return fmt.Sprintf("%s.%d.%d", nodeNet, t[0], t[1])
}

func (n *node) SetHostname() error {
	if err := os.WriteFile(etcHostname, fmt.Appendf(make([]byte, 0), "%s\n", n.NodeName()), 0644); err != nil {
		return err
	}
	cmd := exec.Command("hostname", "-F", etcHostname)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("SetHostname(): 'hostname -F' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return nil
}

func loadNode() (*node, error) {
	buf, err := os.ReadFile(nodeIdFname)
	if err != nil {
		return nil, err
	}
	nodeId, err := strconv.Atoi(strings.TrimSpace(string(buf)))
	if err != nil {
		return nil, err
	}
	if nodeId < nodeIdMin {
		return nil, fmt.Errorf("'%s' node-id is '%d' (min=%d)", nodeIdFname, nodeId, nodeIdMin)
	}
	if nodeId > nodeIdMax {
		return nil, fmt.Errorf("'%s' node-id is '%d' (max=%d)", nodeIdFname, nodeId, nodeIdMax)
	}
	return &node{id: uint16(nodeId)}, nil
}

func writeNode() error {
	x := rand.Int32N(nodeIdMax-nodeIdMin+1) + nodeIdMin
	if err := os.WriteFile(nodeIdFname, fmt.Appendf(make([]byte, 0), "%d\n", x), 0600); err != nil {
		return err
	}
	return nil
}

func getNodeLocal() (*node, error) {
	node, err := loadNode()
	if err == nil {
		return node, nil
	}
	if err := writeNode(); err != nil {
		return nil, err
	}
	return loadNode()
}
