package node

import (
	"fmt"
	"net/netip"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/mesh/netutil"
	"github.com/sfmunoz/i12e/internal/mesh/wgkey"
	"github.com/vishvananda/netlink"
)

type NodeLocal struct {
	Node
	wgKeyPriv *wgkey.WgKeyPriv
}

func (n *NodeLocal) String() string {
	return fmt.Sprintf(
		"L|%s|%s",
		n.Node.String(),
		n.GetWgKeyPriv().Pub(),
	)
}

func (n *NodeLocal) GetWgKeyPriv() *wgkey.WgKeyPriv {
	return n.wgKeyPriv
}

func (n *NodeLocal) HostnameConfig() error {
	nodeName := n.GetNodeName()
	cmd := exec.Command("hostname", nodeName)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("SetHostname(): 'hostname %s' failed': %s (stdout=%s, stderr=%s)", nodeName, err, bo.String(), be.String())
	}
	return nil
}

func (n *NodeLocal) IfaceLocalConfig() error {
	nodeInt := GetNodeInterface()
	link, err := netutil.IfaceCreate(nodeInt)
	if err != nil {
		return err
	}
	if err := netutil.IfaceSyncAddresses(link, n.GetMeshIP(), n.GetMeshNet().Bits()); err != nil {
		return err
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}
	cmd := exec.Command("wg", "set", nodeInt, "listen-port", fmt.Sprintf("%d", n.GetWgEndpoint().Port()), "private-key", nodePrivKeyFname)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'wg set' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return nil
}

func (n *NodeLocal) PushToRemote(remMeshBase string) error {
	ts := time.Now().UTC()
	wgEndpoint := n.GetWgEndpoint()
	touchPath := fmt.Sprintf(
		"%s/%s/%s_%09d/%s/%s/%d",
		remMeshBase,
		n.GetNodeName(),
		ts.Format("20060102_150405"),
		ts.Nanosecond(),
		n.GetWgKeyPriv().Pub().K32().Hex(),
		wgEndpoint.Addr(),
		wgEndpoint.Port(),
	)
	mode, err := getEtcI12eMode()
	if err == nil && mode == "main" { // TODO unhardcode
		touchPath = fmt.Sprintf("%s/kmain", touchPath)
	}
	cmd := exec.Command("rclone", "touch", touchPath)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("'rclone touch' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return nil
}

func (n *NodeLocal) PurgeFromRemote(remMeshBase string) error {
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

func NewNodeLocal(meshNet *netip.Prefix, reset bool) (*NodeLocal, error) {
	if reset {
		if err := deleteEtcHostname(); err != nil {
			return nil, err
		}
	}
	nodeId, err := readEtcHostname(meshNet)
	if err != nil {
		if err := writeRandomEtcHostname(meshNet); err != nil {
			return nil, err
		}
	}
	nodeId, err = readEtcHostname(meshNet)
	if err != nil {
		return nil, err
	}
	wgKeyPriv, err := wgkey.NewWgKeyPriv(nodePrivKeyFname)
	if err != nil {
		return nil, err
	}
	ii, err := netutil.IfaceIP()
	if err != nil {
		return nil, err
	}
	if ii == nil {
		return nil, fmt.Errorf("IfaceIP() returned empty value")
	}
	wgEndpoint := netip.AddrPortFrom(*ii.IP, nodeEndpointPort)
	return &NodeLocal{
		Node: Node{
			id:         nodeId,
			meshNet:    meshNet,
			wgEndpoint: &wgEndpoint,
		},
		wgKeyPriv: wgKeyPriv,
	}, nil
}
