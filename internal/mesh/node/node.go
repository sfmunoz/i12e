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

func GetNodeInterface() string {
	return nodeInterface
}

type Node struct {
	id         uint32
	meshNet    *netip.Prefix
	wgKey      *wgkey.WgKey
	wgEndpoint *netip.AddrPort
	tsFirst    *time.Time // first record of the series
	tsCurr     *time.Time // current record
}

func (n *Node) String() string {
	where := "R"
	if n.GetLocal() {
		where = "L"
	}
	tsCurrStr := "<undefined-tsCurr>"
	tsCurr := n.GetTsCurr()
	if tsCurr != nil {
		tsCurrStr = tsCurr.Format(time.RFC3339Nano)
	}
	ageStr := "<undefined-age>"
	age := n.GetAge()
	if age != nil {
		ageStr = age.String()
	}
	return fmt.Sprintf(
		"%s|%s|%s/%d|%s|%s|%s|%s",
		where,
		n.GetNodeName(),
		n.GetMeshIP(),
		n.GetMeshNet().Bits(),
		n.GetWgKey(),
		n.GetWgEndpoint(),
		tsCurrStr,
		ageStr,
	)
}

func (n *Node) GetNodeName() string {
	return getNodeNameFromNodeId(n.id)
}

func (n *Node) GetMeshIP() *netip.Addr {
	x, _ := nodeIdToIp(n.GetMeshNet(), n.id) // err ignored: already validated
	return x
}

func (n *Node) GetMeshNet() *netip.Prefix {
	return n.meshNet
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

func (n *Node) GetTsFirst() *time.Time {
	return n.tsFirst
}

func (n *Node) SetTsFirst(ts *time.Time) {
	n.tsFirst = ts
}

func (n *Node) GetTsCurr() *time.Time {
	return n.tsCurr
}

func (n *Node) SetTsCurr(ts *time.Time) {
	n.tsCurr = ts
}

func (n *Node) GetAge() *time.Duration {
	tsFirst := n.GetTsFirst()
	if tsFirst == nil {
		return nil
	}
	tsCurr := n.GetTsCurr()
	if tsCurr == nil {
		return nil
	}
	d := tsCurr.Sub(*tsFirst)
	return &d
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
	cmd = exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%d", n.GetMeshIP(), n.GetMeshNet().Bits()), "dev", nodeInt)
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

func (n *Node) PurgeFromRemote(remMeshBase string) error {
	if !n.GetLocal() {
		return fmt.Errorf("cannot purge node: it's not local")
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
	lastEntry := len(entries) - 1
	for i, entry := range entries {
		if i == 0 || i == lastEntry {
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
