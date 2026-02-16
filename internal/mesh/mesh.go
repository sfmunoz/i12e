package mesh

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/mesh/node"
	"github.com/sfmunoz/logit"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

const settleTime = 15 * time.Second

type Mesh struct {
	meshNet *netip.Prefix
	remBase string
}

func (m *Mesh) setNodeListTimestamps(nodeListRaw []*node.NodeRemote) {
	tsMap := make(map[string]*time.Time)
	for _, n := range nodeListRaw {
		k := n.GetNodeName() + "_" + n.GetWgKeyPub().K32().Hex()
		if v, ok := tsMap[k]; ok {
			n.SetTsFirst(v)
			continue
		}
		ts := n.GetTsCurr()
		n.SetTsFirst(ts)
		tsMap[k] = ts
	}
}

func (m *Mesh) getRemoteNodeList() ([]*node.NodeRemote, error) {
	cmd := exec.Command("rclone", "lsf", "-R", "--files-only", m.remBase)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("'rclone lsf' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	entries := strings.Split(strings.TrimSpace(bo.String()), "\n")
	slices.Sort(entries)
	nodeListRaw := make([]*node.NodeRemote, 0)
	for _, entry := range entries {
		entryTrimmed := strings.TrimSpace(entry)
		if len(entryTrimmed) < 1 { // when entries == ""
			continue
		}
		n, err := node.NewNodeRemote(m.meshNet, entryTrimmed)
		if err != nil {
			log.Error("'node.NewNode()' failed", "err", err, "entry", entry)
			continue
		}
		nodeListRaw = append(nodeListRaw, n)
	}
	m.setNodeListTimestamps(nodeListRaw)
	return nodeListRaw, nil
}

func (m *Mesh) appendNodeToBlock(nodeBlock []*node.NodeRemote, n *node.NodeRemote) ([]*node.NodeRemote, bool) {
	if n == nil {
		return nodeBlock, true
	}
	nbLen := len(nodeBlock)
	if nbLen > 0 && n.GetNodeName() != nodeBlock[nbLen-1].GetNodeName() {
		return nodeBlock, true
	}
	return append(nodeBlock, n), false
}

func (m *Mesh) squeezeBlock(nodeList []*node.NodeRemote, nodeBlock []*node.NodeRemote) []*node.NodeRemote {
	if len(nodeBlock) < 1 {
		return nodeList
	}
	nb0 := nodeBlock[0]
	hex0 := nb0.GetWgKeyPub().K32().Hex()
	for i := len(nodeBlock) - 1; i > 0; i-- {
		n := nodeBlock[i]
		if n.GetWgKeyPub().K32().Hex() == hex0 {
			return append(nodeList, n)
		}
	}
	return append(nodeList, nb0)
}

func (m *Mesh) squeezeNodeList(nodeListRaw []*node.NodeRemote) []*node.NodeRemote {
	nodeList := make([]*node.NodeRemote, 0)
	nodeBlock := make([]*node.NodeRemote, 0)
	for _, n := range append(nodeListRaw, nil) {
		var blockDone bool
		nodeBlock, blockDone = m.appendNodeToBlock(nodeBlock, n)
		if !blockDone {
			continue
		}
		nodeList = m.squeezeBlock(nodeList, nodeBlock)
		nodeBlock = make([]*node.NodeRemote, 1)
		nodeBlock[0] = n
	}
	return nodeList
}

func (m *Mesh) getHomonymFromNodeList(nodeList []*node.NodeRemote, nodeLocal *node.NodeLocal) *node.NodeRemote {
	nodeNameLocal := nodeLocal.GetNodeName()
	for _, n := range nodeList {
		if n.GetNodeName() == nodeNameLocal {
			return n
		}
	}
	return nil
}

func (m *Mesh) nodeGiveUp(nodeLocalOld *node.NodeLocal) error {
	nodeLocalNew, err := node.NewNodeLocal(m.meshNet, true)
	if err != nil {
		return fmt.Errorf("node-reset failed (nodeLocal=%s): %s", nodeLocalOld, err)
	}
	log.Info("node-reset OK", "nodeLocalNew", nodeLocalNew, "nodeLocalOld", nodeLocalOld)
	return nil
}

func (m *Mesh) ifacePeersAdd(nodeList []*node.NodeRemote, nodeLocal *node.NodeLocal) error {
	nodeNameLocal := nodeLocal.GetNodeName()
	for _, n := range nodeList {
		if n.GetNodeName() == nodeNameLocal {
			log.Info("peers: skipping local node", "node", n)
			continue
		}
		nAge := *n.GetAge(nil)
		if nAge < settleTime {
			log.Warn("peers: skipping node: nodeAge < settleTime", "nodeAge", nAge, "settleTime", settleTime, "node", n)
			continue
		}
		cmd := exec.Command(
			"wg", "set", node.GetNodeInterface(),
			"peer", n.GetWgKeyPub().K32().B64(),
			"endpoint", n.GetWgEndpoint().String(),
			"allowed-ips", fmt.Sprintf("%s/32", n.GetMeshIP()),
		)
		bo, be, err := cmdutil.RunSimple(cmd)
		if err != nil {
			return fmt.Errorf("'wg set' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
		}
	}
	return nil
}

func (m *Mesh) ifacePeersPurge(nodeList []*node.NodeRemote, wgCli *wgctrl.Client) error {
	iface := node.GetNodeInterface()
	wgDev, err := wgCli.Device(iface)
	if err != nil {
		return err
	}
	peerConfigs := make([]wgtypes.PeerConfig, 0)
peerFor:
	for _, p := range wgDev.Peers {
		pubKey := p.PublicKey
		pubKeyStr := pubKey.String()
		for _, n := range nodeList {
			if n.GetWgKeyPub().K32().B64() == pubKeyStr {
				continue peerFor
			}
		}
		peerConfig := wgtypes.PeerConfig{PublicKey: pubKey, Remove: true}
		peerConfigs = append(peerConfigs, peerConfig)
	}
	config := wgtypes.Config{Peers: peerConfigs}
	return wgCli.ConfigureDevice(iface, config)
}

func (m *Mesh) ifacePeersConfig(nodeList []*node.NodeRemote, nodeLocal *node.NodeLocal, wgCli *wgctrl.Client) error {
	if err := m.ifacePeersAdd(nodeList, nodeLocal); err != nil {
		return err
	}
	return m.ifacePeersPurge(nodeList, wgCli)
}

func (m *Mesh) wireguardDebug(wgCli *wgctrl.Client) error {
	devices, err := wgCli.Devices()
	if err != nil {
		return err
	}
	for _, dev := range devices {
		fmt.Printf("Name ........... %s\n", dev.Name)
		fmt.Printf("  PublicKey .... %s\n", dev.PublicKey.String())
		fmt.Printf("  ListenPort ... %d\n", dev.ListenPort)
		for _, p := range dev.Peers {
			fmt.Printf("  Peer ............ %s\n", p.PublicKey.String())
			fmt.Printf("    Endpoint ...... %v\n", p.Endpoint)
			fmt.Printf("    Allowed IPs ... %v\n", p.AllowedIPs)
		}
	}
	return nil
}

func (m *Mesh) ifacePeersConfigNew(nodeList []*node.NodeRemote, nodeLocal *node.NodeLocal, wgCli *wgctrl.Client) error {
	if err := m.wireguardDebug(wgCli); err != nil {
		return err
	}
	nodeNameLocal := nodeLocal.GetNodeName()
	peerConfigs := make([]wgtypes.PeerConfig, 0)
	for _, n := range nodeList {
		if n.GetNodeName() == nodeNameLocal {
			continue
		}
		nAge := *n.GetAge(nil)
		if nAge < settleTime {
			log.Notice("peers: skipping node: nodeAge < settleTime", "nodeAge", nAge, "settleTime", settleTime, "node", n)
			continue
		}
		// "wg", "set", node.GetNodeInterface(),
		// "peer", n.GetWgKeyPub().K32().B64(),
		// "endpoint", n.GetWgEndpoint().String(),
		// "allowed-ips", fmt.Sprintf("%s/32", n.GetMeshIP()),
		peerKey, err := wgtypes.ParseKey(n.GetWgKeyPub().K32().B64())
		if err != nil {
			return nil
		}
		ep := n.GetWgEndpoint()
		endpoint := &net.UDPAddr{IP: net.ParseIP(ep.Addr().String()), Port: int(ep.Port())}
		keepAlive := 25 * time.Second
		peerConfig := wgtypes.PeerConfig{
			PublicKey:                   peerKey,
			Remove:                      false,
			UpdateOnly:                  false,
			Endpoint:                    endpoint,
			PersistentKeepaliveInterval: &keepAlive,
			ReplaceAllowedIPs:           true,
			AllowedIPs: []net.IPNet{
				{
					IP:   net.ParseIP(n.GetMeshIP().String()),
					Mask: net.CIDRMask(32, 32),
				},
			},
		}
		peerConfigs = append(peerConfigs, peerConfig)
	}
	// wg set <nodeInt> listen-port fmt.Sprintf("%d", n.GetWgEndpoint().Port()) private-key <nodePrivKeyFname>
	epLocalPort := int(nodeLocal.GetWgEndpoint().Port())
	ifaceName := node.GetNodeInterface()
	// privateKey := wgtypes.Key([32]byte{})
	config := wgtypes.Config{
		// PrivateKey:   &privateKey,
		ListenPort:   &epLocalPort,
		ReplacePeers: false,
		Peers:        peerConfigs,
	}
	if err := wgCli.ConfigureDevice(ifaceName, config); err != nil {
		return err
	}
	return m.wireguardDebug(wgCli)
}

func (m *Mesh) etcHostsUpdate(nodeList []*node.NodeRemote) error {
	// Flatcar's default
	lines := []string{
		"#",
		"# autogenerated by i12e - do not edit",
		"#",
		"",
		"127.0.0.1 localhost",
		"::1 localhost",
	}
	for _, n := range nodeList {
		nAge := *n.GetAge(nil)
		if nAge < settleTime {
			log.Warn("/etc/hosts: skipping node: nodeAge < settleTime", "nodeAge", nAge, "settleTime", settleTime, "node", n)
			continue
		}
		nodeAlias := n.GetNodeAlias()
		if nodeAlias != "" {
			nodeAlias = " " + nodeAlias
		}
		lines = append(lines, fmt.Sprintf("%s %s%s", n.GetMeshIP(), n.GetNodeName(), nodeAlias))
	}
	lines = append(lines, "") // NL to the end
	buf := strings.Join(lines, "\n")
	return os.WriteFile("/etc/hosts", []byte(buf), 0644)
}

func (m *Mesh) run() error {
	nodeLocal, err := node.NewNodeLocal(m.meshNet, false)
	if err != nil {
		return err
	}
	if err := nodeLocal.PushToRemote(m.remBase); err != nil {
		return err
	}
	nodeListRaw, err := m.getRemoteNodeList()
	if err != nil {
		return err
	}
	nodeList := m.squeezeNodeList(nodeListRaw)
	nodeRemote := m.getHomonymFromNodeList(nodeList, nodeLocal)
	if nodeRemote == nil {
		return fmt.Errorf("cannot find nodeLocal='%s' homonym in node list", nodeLocal)
	}
	if nodeRemote.GetWgKeyPub().K32().Hex() != nodeLocal.GetWgKeyPriv().Pub().K32().Hex() {
		log.Warn("battle lost, giving up", "nodeLocal", nodeLocal, "nodeRemote", nodeRemote)
		return m.nodeGiveUp(nodeLocal)
	}
	nodeAge := *nodeRemote.GetAge(nil)
	if nodeAge < settleTime {
		log.Warn("nodeAge < settleTime; waiting...", "nodeAge", nodeAge, "settleTime", settleTime, "nodeLocal", nodeLocal, "nodeAge", nodeAge)
		return nil
	}
	if err := nodeLocal.PurgeFromRemote(m.remBase); err != nil {
		return err
	}
	if err := nodeLocal.HostnameConfig(); err != nil {
		return err
	}
	wgCli, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer wgCli.Close()
	if err := nodeLocal.IfaceLocalConfig(wgCli); err != nil {
		return err
	}
	if err := m.ifacePeersConfig(nodeList, nodeLocal, wgCli); err != nil {
		return err
	}
	if err := m.ifacePeersConfigNew(nodeList, nodeLocal, wgCli); err != nil {
		return err
	}
	if err := m.etcHostsUpdate(nodeList); err != nil {
		return err
	}
	return nil
}

func newMesh(meshNet *netip.Prefix, remBase string) *Mesh {
	return &Mesh{meshNet, remBase}
}

func Run(meshNet *netip.Prefix, remBase string) error {
	return newMesh(meshNet, remBase).run()
}
