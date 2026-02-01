package node

import (
	"fmt"
	"net/netip"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/mesh/wgkey"
)

const nodeInterface = "wgi"
const nodePrivKeyFname = "/etc/i12e/wg-priv-key"
const nodeEtcHostname = "/etc/hostname"

var nodeNet = func() *netip.Prefix {
	x := netip.MustParsePrefix("10.119.0.0/16") // '/16' is mandatory (hardcoded for now)
	return &x
}()

// 252_158/20260128_152841_153793688/54a2cc8d5e78755ff1debc4a4e6b2fa657ccf86a868b53f9f1b5140487377cc8/192.168.56.53/51820
var nodeRegex = regexp.MustCompile(
	"^([0-9]{3}_[0-9]{3})" +
		"/([0-9]{8}_[0-9]{6}_[0-9]{9})" +
		"/([0-9a-f]{64})" +
		`/([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)` +
		"/([0-9]+)" +
		"$",
)

func GetNodeInterface() string {
	return nodeInterface
}

type Node struct {
	id         uint16 // uint16 instead of byte to avoid pervasive type conversion
	wgKey      *wgkey.WgKey
	wgEndpoint *netip.AddrPort
}

func (n *Node) tuple() [2]byte {
	return [2]byte{byte(n.id / 256), byte(n.id % 256)}
}

func (n *Node) String() string {
	return fmt.Sprintf("id=%d|name=%s|ip=%s|path=%s", n.GetNodeId(), n.GetNodeName(), n.GetNodeIP(), n.GetNodePath())
}

func (n *Node) GetNodeId() uint16 {
	return n.id
}

func (n *Node) GetNodeName() string {
	t := n.tuple()
	return fmt.Sprintf("n-%03d-%03d", t[0], t[1])
}

func (n *Node) GetNodePath() string {
	t := n.tuple()
	return fmt.Sprintf("%03d_%03d", t[0], t[1])
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
	if err := os.WriteFile(nodeEtcHostname, fmt.Appendf(make([]byte, 0), "%s\n", n.GetNodeName()), 0644); err != nil {
		return err
	}
	cmd := exec.Command("hostname", "-F", nodeEtcHostname)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("SetHostname(): 'hostname -F' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
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
	nodePath := n.GetNodePath()
	wgEndpoint := n.GetWgEndpoint()
	touchPath := fmt.Sprintf(
		"%s/%s/%s_%09d/%s/%s/%d",
		remMeshBase,
		nodePath,
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
