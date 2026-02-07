package node

import (
	"fmt"
	"net/netip"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/mesh/wgkey"
)

const nodeInterface = "wgi"
const nodePrivKeyFname = "/etc/i12e/wg-priv-key"

var nodeNet = func() *netip.Prefix {
	// XXX future: from /12 (2^20 hosts < 36^4) to /29 (6 hosts)
	x := netip.MustParsePrefix("10.119.0.0/16") // '/16' is mandatory (hardcoded for now)
	return &x
}()

func GetNodeInterface() string {
	return nodeInterface
}

type Node struct {
	id         uint32
	wgKey      *wgkey.WgKey
	wgEndpoint *netip.AddrPort
}

func (n *Node) tuple() [2]byte {
	return [2]byte{byte(n.id / 256), byte(n.id % 256)}
}

func (n *Node) String() string {
	return fmt.Sprintf(
		"name=%s|ip=%s|local=%t|wgkey=%s|endpoint=%s",
		n.GetNodeName(),
		n.GetNodeIP(),
		n.GetLocal(),
		n.GetWgKey(),
		n.GetWgEndpoint(),
	)
}

func (n *Node) GetNodeId() uint32 {
	return n.id
}

func (n *Node) GetNodeName() string {
	return getNodeNameFromNodeId(n.GetNodeId())
}

func (n *Node) GetNodeIP() *netip.Addr {
	x := nodeNet.Addr().As4()
	t := n.tuple()
	addr := netip.AddrFrom4([4]byte{x[0], x[1], t[0], t[1]})
	return &addr
}

func (n *Node) GetWgKey() *wgkey.WgKey {
	return n.wgKey
}

func (n *Node) GetLocal() bool {
	return n.GetWgKey().GetLocal()
}

func (n *Node) GetWgEndpoint() *netip.AddrPort {
	return n.wgEndpoint
}

func (n *Node) HostnameConfig() error {
	if !n.GetLocal() {
		return fmt.Errorf("cannot set hostname: not local node")
	}
	nodeName := n.GetNodeName()
	cmd := exec.Command("hostname", nodeName)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("SetHostname(): 'hostname %s' failed': %s (stdout=%s, stderr=%s)", nodeName, err, bo.String(), be.String())
	}
	return nil
}

func (n *Node) IfaceLocalConfig() error {
	nodeInt := GetNodeInterface()
	cmd := exec.Command("ip", "link", "set", nodeInt, "down")
	bo, be, err := cmdutil.RunSimple(cmd)
	// ignore error: it's OK
	cmd = exec.Command("ip", "link", "del", nodeInt)
	bo, be, err = cmdutil.RunSimple(cmd)
	// ignore error: it's OK
	cmd = exec.Command("ip", "link", "add", nodeInt, "type", "wireguard")
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'ip link add' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	cmd = exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%d", n.GetNodeIP(), nodeNet.Bits()), "dev", nodeInt)
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'ip link add' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	cmd = exec.Command("wg", "set", nodeInt, "listen-port", fmt.Sprintf("%d", n.GetWgEndpoint().Port()), "private-key", nodePrivKeyFname)
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'wg set' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	cmd = exec.Command("ip", "link", "set", nodeInt, "up")
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'ip link set' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return nil
}

func (n *Node) PushToRemote(remMeshBase string) error {
	if !n.GetLocal() {
		return fmt.Errorf("cannot push node: it's not local")
	}
	ts := time.Now().UTC()
	wgEndpoint := n.GetWgEndpoint()
	touchPath := fmt.Sprintf(
		"%s/%s/%s_%09d/%s/%s/%d",
		remMeshBase,
		n.GetNodeName(),
		ts.Format("20060102_150405"),
		ts.Nanosecond(),
		n.GetWgKey().GetPubKey().Hex(),
		wgEndpoint.Addr(),
		wgEndpoint.Port(),
	)
	cmd := exec.Command("rclone", "touch", touchPath)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'rclone touch' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return nil
}

func (n *Node) PurgeFromRemote(remMeshBase string, keep int) error {
	if !n.GetLocal() {
		return fmt.Errorf("cannot purge node: it's not local")
	}
	if keep < 1 {
		return fmt.Errorf("invalid keep=%d (min=1)", keep)
	}
	nodeName := n.GetNodeName()
	nodePrefix := fmt.Sprintf("%s/%s", remMeshBase, nodeName)
	cmd := exec.Command("rclone", "lsf", nodePrefix)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'rclone lsf' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	entries := strings.Split(strings.TrimSpace(bo.String()), "\n")
	slices.Sort(entries)
	slices.Reverse(entries)
	for i, entry := range entries {
		if i < keep {
			continue
		}
		deletePath := fmt.Sprintf("%s/%s/%s", remMeshBase, nodeName, entry)
		cmd := exec.Command("rclone", "delete", deletePath)
		bo, be, err := cmdutil.RunSimple(cmd)
		if err != nil {
			return fmt.Errorf("'rclone delete' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
		}
	}
	return nil
}
