package mesh

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
const nodeIdMax = 65_534       // 65_535 = 255.255 = broadcast address
const nodeEndpointPort = 51830 // default '51820'
const nodeNet = "10.56"
const nodeInterface = "wgi"
const nodePrivKeyFname = "/etc/i12e/wg-priv-key"
const nodeIdFname = "/etc/i12e/node-id.txt"
const nodeEtcHostname = "/etc/hostname"

type node struct {
	id             uint16 // uint16 instead of uint8 to avoid pervasive type conversion
	local          bool
	wgPrivKey      *WgKey
	wgPubKey       *WgKey
	wgEndpointIp   string
	wgEndpointPort uint16
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

func (n *node) GetWgPrivKey() *WgKey {
	return n.wgPrivKey
}

func (n *node) SetWgPubKey(wgPubKey *WgKey) {
	n.wgPubKey = wgPubKey
}

func (n *node) GetWgPubKey() *WgKey {
	return n.wgPubKey
}

func (n *node) SetLocal(local bool) {
	n.local = local
}

func (n *node) GetLocal() bool {
	return n.local
}

func (n *node) GetWgEndpoint() string {
	return fmt.Sprintf("%s:%d", n.GetWgEndpointIp(), n.GetWgEndpointPort())
}

func (n *node) SetWgEndpointIp(endPointIp string) {
	n.wgEndpointIp = endPointIp
}

func (n *node) GetWgEndpointIp() string {
	return n.wgEndpointIp
}

func (n *node) SetWgEndpointPort(wgEndpointPort uint16) {
	n.wgEndpointPort = wgEndpointPort
}

func (n *node) GetWgEndpointPort() uint16 {
	return n.wgEndpointPort
}

func (n *node) SetHostname() error {
	if err := os.WriteFile(nodeEtcHostname, fmt.Appendf(make([]byte, 0), "%s\n", n.NodeName()), 0644); err != nil {
		return err
	}
	cmd := exec.Command("hostname", "-F", nodeEtcHostname)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("SetHostname(): 'hostname -F' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return nil
}

func validNodeId(nodeId int) error {
	if nodeId < nodeIdMin {
		return fmt.Errorf("'%s' node-id is '%d' (min=%d)", nodeIdFname, nodeId, nodeIdMin)
	}
	if nodeId > nodeIdMax {
		return fmt.Errorf("'%s' node-id is '%d' (max=%d)", nodeIdFname, nodeId, nodeIdMax)
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
	if err := validNodeId(nodeId); err != nil {
		return nil, err
	}
	wgPrivKey, err := getWgPrivKey(nodePrivKeyFname)
	if err != nil {
		return nil, err
	}
	wgPubKey, err := getWgPubKey(nodePrivKeyFname)
	if err != nil {
		return nil, err
	}
	ii, err := IfaceIP()
	if err != nil {
		return nil, err
	}
	if ii == nil {
		return nil, fmt.Errorf("IfaceIP() returned empty value")
	}
	return &node{
		id:             uint16(nodeId),
		local:          true,
		wgPrivKey:      wgPrivKey,
		wgPubKey:       wgPubKey,
		wgEndpointIp:   ii.IP,
		wgEndpointPort: nodeEndpointPort,
	}, nil
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

func getNodeFromPath(path string) (*node, error) {
	parts := strings.Split(path, "_")
	parts_len := len(parts)
	if parts_len != 2 {
		return nil, fmt.Errorf("len(parts)=%d (2 expected)", parts_len)
	}
	ids := [2]uint16{0, 0}
	for i := range 2 {
		plen := len(parts[i])
		if plen < 1 {
			return nil, fmt.Errorf("len(parts[%d])=%d (>0 expected)", i, plen)
		}
		pint, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, err
		}
		if pint < 0 {
			return nil, fmt.Errorf("wrong path: '%s[%d]' = '%s' < 0", path, i, parts[i])
		}
		if pint > 255 {
			return nil, fmt.Errorf("wrong path: '%s[%d]' = '%s' > 255", path, i, parts[i])
		}
		ids[i] = uint16(pint)
	}
	nodeId := 256*int(ids[0]) + int(ids[1])
	if err := validNodeId(nodeId); err != nil {
		return nil, err
	}
	return &node{
		id:             uint16(nodeId),
		local:          false,
		wgPrivKey:      nil,
		wgPubKey:       nil,
		wgEndpointIp:   "",
		wgEndpointPort: nodeEndpointPort,
	}, nil
}

func (n *node) ifaceLocalConfig() error {
	wgIpInt := n.NodeIP()
	cmd := exec.Command("ip", "link", "set", nodeInterface, "down")
	bo, be, err := cmdutil.RunSimple(cmd)
	// ignore error: it's OK
	cmd = exec.Command("ip", "link", "del", nodeInterface)
	bo, be, err = cmdutil.RunSimple(cmd)
	// ignore error: it's OK
	cmd = exec.Command("ip", "link", "add", nodeInterface, "type", "wireguard")
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'ip link add' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	cmd = exec.Command("ip", "addr", "add", fmt.Sprintf("%s/16", wgIpInt), "dev", nodeInterface)
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'ip link add' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	cmd = exec.Command("wg", "set", nodeInterface, "listen-port", fmt.Sprintf("%d", n.GetWgEndpointPort()), "private-key", nodePrivKeyFname)
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'wg set' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	cmd = exec.Command("ip", "link", "set", nodeInterface, "up")
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'ip link set' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return nil
}
