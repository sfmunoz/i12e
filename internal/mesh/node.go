package mesh

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

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

// 252_158/20260128_152841_153793688/54a2cc8d5e78755ff1debc4a4e6b2fa657ccf86a868b53f9f1b5140487377cc8/192.168.56.53/51820
var nodeRegex = regexp.MustCompile(
	"^([0-9]{3}_[0-9]{3})" +
		"/([0-9]{8}_[0-9]{6}_[0-9]{9})" +
		"/([0-9a-f]{64})" +
		`/([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)` +
		"/([0-9]+)" +
		"$",
)

type node struct {
	id             uint16 // uint16 instead of uint8 to avoid pervasive type conversion
	wgPrivKey      *WgKey
	wgPubKey       *WgKey
	wgEndpointIp   string
	wgEndpointPort uint16
}

func (n *node) tuple() [2]uint8 {
	return [2]uint8{uint8(n.id / 256), uint8(n.id % 256)}
}

func (n *node) String() string {
	return fmt.Sprintf("id=%d|name=%s|ip=%s|path=%s", n.getNodeId(), n.getNodeName(), n.getNodeIP(), n.getNodePath())
}

func (n *node) getNodeId() uint16 {
	return n.id
}

func (n *node) getNodeName() string {
	t := n.tuple()
	return fmt.Sprintf("n-%03d-%03d", t[0], t[1])
}

func (n *node) getNodePath() string {
	t := n.tuple()
	return fmt.Sprintf("%03d_%03d", t[0], t[1])
}

func (n *node) getNodeIP() string {
	t := n.tuple()
	return fmt.Sprintf("%s.%d.%d", nodeNet, t[0], t[1])
}

func (n *node) getWgPrivKey() *WgKey {
	return n.wgPrivKey
}

func (n *node) getWgPubKey() *WgKey {
	return n.wgPubKey
}

func (n *node) getLocal() bool {
	return n.wgPrivKey != nil
}

func (n *node) getWgEndpoint() string {
	return fmt.Sprintf("%s:%d", n.getWgEndpointIp(), n.getWgEndpointPort())
}

func (n *node) getWgEndpointIp() string {
	return n.wgEndpointIp
}

func (n *node) getWgEndpointPort() uint16 {
	return n.wgEndpointPort
}

func (n *node) hostnameConfig() error {
	if err := os.WriteFile(nodeEtcHostname, fmt.Appendf(make([]byte, 0), "%s\n", n.getNodeName()), 0644); err != nil {
		return err
	}
	cmd := exec.Command("hostname", "-F", nodeEtcHostname)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("SetHostname(): 'hostname -F' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return nil
}

func (n *node) ifaceLocalConfig() error {
	wgIpInt := n.getNodeIP()
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
	cmd = exec.Command("wg", "set", nodeInterface, "listen-port", fmt.Sprintf("%d", n.getWgEndpointPort()), "private-key", nodePrivKeyFname)
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

func (n *node) pushToRemote() error {
	if !n.getLocal() {
		return fmt.Errorf("cannot push node: it's not local")
	}
	ts := time.Now().UTC()
	nodePath := n.getNodePath()
	touchPath := fmt.Sprintf(
		"%s/%s/%s_%09d/%s/%s/%d",
		remMeshBase,
		nodePath,
		ts.Format("20060102_150405"),
		ts.Nanosecond(),
		n.getWgPubKey().Hex(),
		n.getWgEndpointIp(),
		n.getWgEndpointPort(),
	)
	cmd := exec.Command("rclone", "touch", touchPath)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'rclone touch' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	nodePrefix := fmt.Sprintf("%s/%s", remMeshBase, nodePath)
	cmd = exec.Command("rclone", "lsf", nodePrefix)
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'rclone lsf' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	entries := strings.Split(strings.TrimSpace(bo.String()), "\n")
	slices.Sort(entries)
	slices.Reverse(entries)
	for i, entry := range entries {
		if i < 1 {
			continue
		}
		deletePath := fmt.Sprintf("%s/%s/%s", remMeshBase, nodePath, entry)
		cmd := exec.Command("rclone", "delete", deletePath)
		bo, be, err := cmdutil.RunSimple(cmd)
		if err != nil {
			return fmt.Errorf("'rclone delete' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
		}
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

func getNodeIdFromPath(path string) (uint16, error) {
	parts := strings.Split(path, "_")
	parts_len := len(parts)
	if parts_len != 2 {
		return 0, fmt.Errorf("path='%s' -> len(parts=%s)=%d (2 expected)", path, parts, parts_len)
	}
	ids := [2]uint16{0, 0}
	for i, part := range parts {
		plen := len(part)
		if plen < 1 {
			return 0, fmt.Errorf("len(parts[%d])=%d (>0 expected)", i, plen)
		}
		pint, err := strconv.Atoi(part)
		if err != nil {
			return 0, err
		}
		if pint < 0 {
			return 0, fmt.Errorf("wrong path: '%s[%d]' = '%s' < 0", path, i, part)
		}
		if pint > 255 {
			return 0, fmt.Errorf("wrong path: '%s[%d]' = '%s' > 255", path, i, part)
		}
		ids[i] = uint16(pint)
	}
	nodeId := 256*int(ids[0]) + int(ids[1])
	if err := validNodeId(nodeId); err != nil {
		return 0, err
	}
	return uint16(nodeId), nil
}

func getNodeRemote(entry string) (*node, error) {
	arr := nodeRegex.FindStringSubmatch(entry)
	if arr == nil {
		return nil, fmt.Errorf("'FindStringSubmatch()' returned nil (entry=%s)", entry)
	}
	nodeId, err := getNodeIdFromPath(arr[1])
	if err != nil {
		return nil, err
	}
	wgPubKey, err := getWgKeyFromHex(arr[3], false)
	if err != nil {
		return nil, fmt.Errorf("'getWgKeyFromHex(%s)' failed: %s", arr[3], err)
	}
	wgEndpointPort, err := strconv.Atoi(arr[5])
	if err != nil {
		return nil, fmt.Errorf("'strconv.Atoi(%s)' failed: %s", arr[5], err)
	}
	return &node{
		id:             nodeId,
		wgPrivKey:      nil,
		wgPubKey:       wgPubKey,
		wgEndpointIp:   arr[4],
		wgEndpointPort: uint16(wgEndpointPort),
	}, nil
}
