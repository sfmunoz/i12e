package node

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/netip"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/mesh/netutil"
	"github.com/vishvananda/netlink"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type NodeLocal struct {
	Node
	wgKeyPriv wgtypes.Key
}

func (n *NodeLocal) String() string {
	return fmt.Sprintf(
		"L|%s|%s",
		n.Node.String(),
		n.GetWgKeyPriv().PublicKey(),
	)
}

func (n *NodeLocal) GetWgKeyPriv() wgtypes.Key {
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

func (n *NodeLocal) IfaceLocalConfig(wgCli *wgctrl.Client) error {
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
	privateKey := n.GetWgKeyPriv()
	listenPort := int(n.GetWgEndpoint().Port())
	config := wgtypes.Config{
		PrivateKey: &privateKey,
		ListenPort: &listenPort,
	}
	return wgCli.ConfigureDevice(nodeInt, config)
}

func (n *NodeLocal) PushToRemote(remMeshBase string) error {
	ts := time.Now().UTC()
	wgEndpoint := n.GetWgEndpoint()
	pubKey := n.GetWgKeyPriv().PublicKey()
	touchPath := fmt.Sprintf(
		"%s/%s/%s_%09d/%s/%s/%d",
		remMeshBase,
		n.GetNodeName(),
		ts.Format("20060102_150405"),
		ts.Nanosecond(),
		hex.EncodeToString(pubKey[:]),
		wgEndpoint.Addr(),
		wgEndpoint.Port(),
	)
	mode, err := getEtcI12eMode()
	if err == nil && mode == "main" { // TODO: unhardcode
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
	wgKeyPriv, err := getWgKeyPriv(nodePrivKeyFname)
	if err != nil {
		return nil, err
	}
	addr, err := netutil.MeshEndpointAddr()
	if err != nil {
		return nil, err
	}
	wgEndpoint := netip.AddrPortFrom(*addr, nodeEndpointPort)
	return &NodeLocal{
		Node: Node{
			id:         nodeId,
			meshNet:    meshNet,
			wgEndpoint: &wgEndpoint,
		},
		wgKeyPriv: wgKeyPriv,
	}, nil
}

func getWgKeyPriv(wgPrivKeyFname string) (wgtypes.Key, error) {
	for i := range 2 {
		_, err := os.Stat(wgPrivKeyFname)
		if err == nil {
			break
		}
		if !errors.Is(err, os.ErrNotExist) || i > 0 {
			return wgtypes.Key{}, fmt.Errorf("'os.Stat()' failed: %s", err)
		}
		key, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return wgtypes.Key{}, err
		}
		if err := os.WriteFile(wgPrivKeyFname, []byte(key.String()+"\n"), 0600); err != nil {
			return wgtypes.Key{}, fmt.Errorf("'os.WriteFile()' failed: %s", err)
		}
	}
	buf, err := os.ReadFile(wgPrivKeyFname)
	if err != nil {
		return wgtypes.Key{}, err
	}
	return wgtypes.ParseKey(string(buf))
}
