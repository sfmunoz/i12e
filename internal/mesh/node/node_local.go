package node

import (
	"fmt"
	"net/netip"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/mesh/ifaceip"
	"github.com/sfmunoz/i12e/internal/mesh/wgkey"
)

func (n *Node) GetWgKeyPriv() *wgkey.WgKeyPriv {
	return n.wgKeyPriv
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
		n.GetWgKeyPub().K32().Hex(),
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

func NewNodeLocal(meshNet *netip.Prefix, reset bool) (*Node, error) {
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
	ii, err := ifaceip.IfaceIP()
	if err != nil {
		return nil, err
	}
	if ii == nil {
		return nil, fmt.Errorf("IfaceIP() returned empty value")
	}
	wgEndpoint := netip.AddrPortFrom(*ii.IP, nodeEndpointPort)
	return &Node{
		id:         nodeId,
		meshNet:    meshNet,
		wgKeyPriv:  wgKeyPriv,
		wgKeyPub:   wgKeyPriv.Pub(),
		wgEndpoint: &wgEndpoint,
		tsFirst:    nil,
		tsCurr:     nil,
	}, nil
}
