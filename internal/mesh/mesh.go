package mesh

import (
	"cmp"
	"fmt"
	"net/netip"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/mesh/node"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Mesh struct {
	meshNet *netip.Prefix
	remBase string
}

func (m *Mesh) setNodeListTimestamps(nodeListRaw []*node.Node) {
	tsMap := make(map[string]*time.Time)
	nLast := len(nodeListRaw) - 1
	for i := nLast; i >= 0; i-- {
		n := nodeListRaw[i]
		k := n.GetNodeName() + "_" + n.GetWgKey().GetPubKey().Hex()
		if v, ok := tsMap[k]; ok {
			n.SetTsFirst(v)
			continue
		}
		ts := n.GetTsCurr()
		n.SetTsFirst(ts)
		tsMap[k] = ts
	}
}

func (m *Mesh) getRemoteNodeList() ([]*node.Node, error) {
	cmd := exec.Command("rclone", "lsf", "-R", "--files-only", m.remBase)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("'rclone lsf' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	entries := strings.Split(strings.TrimSpace(bo.String()), "\n")
	slices.Sort(entries)
	slices.Reverse(entries)
	nodeListRaw := make([]*node.Node, 0)
	for _, entry := range entries {
		entryTrimmed := strings.TrimSpace(entry)
		if len(entryTrimmed) < 1 { // when entries == ""
			continue
		}
		n, err := node.NewNode(node.WithRemote(m.meshNet, entryTrimmed))
		if err != nil {
			log.Error("'node.NewNode()' failed", "err", err, "entry", entry)
			continue
		}
		nodeListRaw = append(nodeListRaw, n)
	}
	m.setNodeListTimestamps(nodeListRaw)
	return nodeListRaw, nil
}

func (m *Mesh) appendNodeToBlock(nodeBlock []*node.Node, n *node.Node) ([]*node.Node, bool) {
	if n == nil {
		return nodeBlock, true
	}
	nbLen := len(nodeBlock)
	if nbLen > 0 && n.GetNodeName() != nodeBlock[nbLen-1].GetNodeName() {
		return nodeBlock, true
	}
	return append(nodeBlock, n), false
}

func (m *Mesh) appendBlockToNodeList(nodeList []*node.Node, nodeBlock []*node.Node) []*node.Node {
	nodeCmp := func(n1, n2 *node.Node) int {
		// -1: n1 < n2 | 0: n1 == n2 | +1: n1 > n2
		return cmp.Compare(*n2.GetAge(), *n1.GetAge())
	}
	if len(nodeBlock) < 1 {
		return nodeList
	}
	slices.SortFunc(nodeBlock, nodeCmp)
	nodeSeen := make(map[string]bool, len(nodeBlock))
	nodeBlockTrimmed := make([]*node.Node, 0, len(nodeBlock))
	for _, n := range nodeBlock {
		k := n.GetWgKey().GetPubKey().Hex()
		if _, ok := nodeSeen[k]; !ok {
			nodeSeen[k] = true
			nodeBlockTrimmed = append(nodeBlockTrimmed, n)
		}
	}
	if len(nodeBlockTrimmed) < 1 {
		return nodeList
	}
	n0 := nodeBlockTrimmed[0]
	if *n0.GetAge() < 3*time.Second {
		return nodeList
	}
	return append(nodeList, nodeBlock[0])
}

func (m *Mesh) getConfirmedNodeList(nodeListRaw []*node.Node) []*node.Node {
	nodeList := make([]*node.Node, 0)
	nodeBlock := make([]*node.Node, 0)
	for _, n := range append(nodeListRaw, nil) {
		var blockDone bool
		nodeBlock, blockDone = m.appendNodeToBlock(nodeBlock, n)
		if !blockDone {
			continue
		}
		nodeList = m.appendBlockToNodeList(nodeList, nodeBlock)
		nodeBlock = make([]*node.Node, 1)
		nodeBlock[0] = n
	}
	return nodeList
}

func (m *Mesh) nodeLocalInNodeList(nodeList []*node.Node, nodeLocal *node.Node, idOnly bool) *node.Node {
	nodeNameLocal := nodeLocal.GetNodeName()
	nodePubKeyLocal := nodeLocal.GetWgKey().GetPubKey().Hex()
	for _, n := range nodeList {
		if n.GetNodeName() != nodeNameLocal {
			continue
		}
		if idOnly || n.GetWgKey().GetPubKey().Hex() == nodePubKeyLocal {
			return n
		}
	}
	return nil
}

func (m *Mesh) getContenderNodes(nodeListRaw []*node.Node, nodeLocal *node.Node) []*node.Node {
	nodeNameLocal := nodeLocal.GetNodeName()
	nodePubKeyLocal := nodeLocal.GetWgKey().GetPubKey().Hex()
	pubKeys := make([]string, 0)
	nodeListRet := make([]*node.Node, 0)
	for _, n := range nodeListRaw {
		if n.GetNodeName() != nodeNameLocal {
			continue
		}
		pubKey := n.GetWgKey().GetPubKey().Hex()
		if pubKey == nodePubKeyLocal {
			continue
		}
		if slices.Contains(pubKeys, pubKey) {
			continue
		}
		pubKeys = append(pubKeys, pubKey)
		nodeListRet = append(nodeListRet, n)
	}
	return nodeListRet
}

func (m *Mesh) nodeGiveUp(nodeLocalOld *node.Node) error {
	nodeLocalNew, err := node.NewNode(node.WithLocal(m.meshNet, true))
	if err != nil {
		return fmt.Errorf("node-reset failed (nodeLocal=%s): %s", nodeLocalOld, err)
	}
	log.Info("node-reset OK", "nodeLocalNew", nodeLocalNew, "nodeLocalOld", nodeLocalOld)
	return nil
}

func (m *Mesh) nodeGiveUpOrPush(nodeListRaw []*node.Node, nodeLocal *node.Node) error {
	ncList := m.getContenderNodes(nodeListRaw, nodeLocal)
	ncListLen := len(ncList)
	if ncListLen > 0 {
		for i, n := range ncList {
			log.Warn("contender", "i", i+1, "tot", ncListLen, "node", n)
		}
		log.Warn("contention detected: giving up", "nodeLocal", nodeLocal)
		return m.nodeGiveUp(nodeLocal)
	}
	log.Info("nodeLocal.PushToRemote()...")
	if err := nodeLocal.PushToRemote(m.remBase); err != nil {
		return err
	}
	return nil
}

func (m *Mesh) ifacePeersConfig(nodeList []*node.Node, nodeLocal *node.Node) error {
	nodeNameLocal := nodeLocal.GetNodeName()
	for _, n := range nodeList {
		if n.GetNodeName() == nodeNameLocal {
			log.Info("skipping local node", "node", n)
			continue
		}
		cmd := exec.Command(
			"wg", "set", node.GetNodeInterface(),
			"peer", n.GetWgKey().GetPubKey().B64(),
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

func (m *Mesh) etcHostsUpdate(nodeList []*node.Node) error {
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
		lines = append(lines, fmt.Sprintf("%s %s", n.GetMeshIP(), n.GetNodeName()))
	}
	lines = append(lines, "") // NL to the end
	buf := strings.Join(lines, "\n")
	return os.WriteFile("/etc/hosts", []byte(buf), 0644)
}

func (m *Mesh) run() error {
	nodeListRaw, err := m.getRemoteNodeList()
	if err != nil {
		return err
	}
	nodeList := m.getConfirmedNodeList(nodeListRaw)
	nodeLocal, err := node.NewNode(node.WithLocal(m.meshNet, false))
	if err != nil {
		return err
	}
	nodeInList := m.nodeLocalInNodeList(nodeList, nodeLocal, true)
	if nodeInList == nil {
		return m.nodeGiveUpOrPush(nodeListRaw, nodeLocal)
	}
	if m.nodeLocalInNodeList(nodeList, nodeLocal, false) == nil {
		log.Warn("conflict detected: giving up", "nodeLocal", nodeLocal, "nodeInList", nodeInList)
		return m.nodeGiveUp(nodeLocal)
	}
	if err := nodeLocal.PushToRemote(m.remBase); err != nil {
		return err
	}
	if err := nodeLocal.PurgeFromRemote(m.remBase); err != nil {
		return err
	}
	if err := nodeLocal.HostnameConfig(); err != nil {
		return err
	}
	if err := nodeLocal.IfaceLocalConfig(); err != nil {
		return err
	}
	if err := m.ifacePeersConfig(nodeList, nodeLocal); err != nil {
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
