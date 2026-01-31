package mesh

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

const WgPrivKeyFname = "/etc/i12e/wg-priv-key" // FIXME unhardcode this
const remMeshBase = "rem:mesh"                 // FIXME unhardcode this

func getRegex() (*regexp.Regexp, error) {
	// 252_158/20260128_152841_153793688/54a2cc8d5e78755ff1debc4a4e6b2fa657ccf86a868b53f9f1b5140487377cc8/192.168.56.53/51820
	expr := "^([0-9]{3}_[0-9]{3})" +
		"/([0-9]{8}_[0-9]{6}_[0-9]{9})" +
		"/([0-9a-f]{64})" +
		`/([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)` +
		"/([0-9]+)" +
		"$"
	return regexp.Compile(expr)
}

type Mesh struct {
	base           string
	re             *regexp.Regexp
	Tnow           time.Time
	nodeLocal      *node
	WgInterface    string
	WgEndpointPort uint16
	WgEndpointIp   string
	WgPrivKey      *WgKey
	WgPubKey       *WgKey
}

func (m *Mesh) NodePush(nodeLocal *node, ts time.Time, wgPubKey *WgKey, wgEndpointIp string, wgEndpointPort uint16) error {
	nodePath := nodeLocal.NodePath()
	touchPath := fmt.Sprintf(
		"%s/%s/%s_%09d/%s/%s/%d",
		m.base,
		nodePath,
		ts.Format("20060102_150405"),
		ts.Nanosecond(),
		wgPubKey.Hex(),
		wgEndpointIp,
		wgEndpointPort,
	)
	cmd := exec.Command("rclone", "touch", touchPath)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("Mesh.NodePush(): 'rclone touch' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	nodePrefix := fmt.Sprintf("%s/%s", m.base, nodePath)
	cmd = exec.Command("rclone", "lsf", nodePrefix)
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("Mesh.NodePush(): 'rclone lsf' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	entries := strings.Split(strings.TrimSpace(bo.String()), "\n")
	slices.Sort(entries)
	slices.Reverse(entries)
	for i, entry := range entries {
		if i < 1 {
			continue
		}
		deletePath := fmt.Sprintf("%s/%s/%s", m.base, nodePath, entry)
		cmd := exec.Command("rclone", "delete", deletePath)
		bo, be, err := cmdutil.RunSimple(cmd)
		if err != nil {
			return fmt.Errorf("Mesh.NodePush(): 'rclone delete' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
		}
	}
	return nil
}

func (m *Mesh) ifaceLocalConfig(nodeLocal *node, wgInterface string, wgEndpointPort uint16, wgPrivKeyFname string) error {
	wgIpInt := nodeLocal.NodeIP()
	cmd := exec.Command("ip", "link", "set", wgInterface, "down")
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		log.Notice("Mesh.ifaceLocalConfig(): 'ip link set' failed (it's OK)", "err", err)
	}
	cmd = exec.Command("ip", "link", "del", wgInterface)
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		log.Notice("Mesh.ifaceLocalConfig(): 'ip link del' failed (it's OK)", "err", err)
	}
	cmd = exec.Command("ip", "link", "add", wgInterface, "type", "wireguard")
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("Mesh.ifaceLocalConfig(): 'ip link add' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	cmd = exec.Command("ip", "addr", "add", fmt.Sprintf("%s/16", wgIpInt), "dev", wgInterface)
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("Mesh.ifaceLocalConfig(): 'ip link add' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	cmd = exec.Command("wg", "set", wgInterface, "listen-port", fmt.Sprintf("%d", wgEndpointPort), "private-key", wgPrivKeyFname)
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("Mesh.ifaceLocalConfig(): 'wg set' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	cmd = exec.Command("ip", "link", "set", wgInterface, "up")
	bo, be, err = cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("Mesh.ifaceLocalConfig(): 'ip link set' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return nil
}

func (m *Mesh) getNodeList(nodeLocal *node) ([]*node, error) {
	cmd := exec.Command("rclone", "lsf", "-R", "--files-only", m.base)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("'rclone lsf' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	data := strings.Split(strings.TrimSpace(bo.String()), "\n")
	slices.Sort(data)
	nodeList := make([]*node, 0)
	wgPathNameLocal := nodeLocal.NodePath()
	for _, entry := range data {
		x := m.re.FindStringSubmatch(entry)
		if x == nil {
			log.Error("'FindStringSubmatch()' failed", "entry", entry)
			continue
		}
		node, err := getNodeFromPath(x[1])
		if err != nil {
			log.Error("'getNodeFromPath()' failed", "err", err, "nodepath", x[1])
			continue
		}
		k, err := getWgKeyFromHex(x[3], false)
		if err != nil {
			log.Error("'getWgKeyFromHex()' failed", "err", err, "hex", x[3])
			continue
		}
		node.SetWgKey(k)
		node.SetLocal(node.NodePath() == wgPathNameLocal)
		node.SetEndPoint(fmt.Sprintf("%s:%s", x[4], x[5]))
		nodeList = append(nodeList, node)
	}
	return nodeList, nil
}

func (m *Mesh) ifacePeersConfig(nodeList []*node, wgInterface string) error {
	for _, node := range nodeList {
		if node.GetLocal() {
			log.Info("skipping local node", "node", node)
			continue
		}
		k := node.GetWgKey()
		if k == nil {
			return fmt.Errorf("'node.WgKey()' returned nil")
		}
		endPoint := node.GetEndPoint()
		if len(endPoint) < 1 {
			return fmt.Errorf("'node.GetEndPoint()' returned empty data")
		}
		cmd := exec.Command(
			"wg", "set", wgInterface,
			"peer", k.B64(),
			"allowed-ips", fmt.Sprintf("%s/32", node.NodeIP()),
			"endpoint", endPoint,
		)
		bo, be, err := cmdutil.RunSimple(cmd)
		if err != nil {
			return fmt.Errorf("'wg set' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
		}
	}
	return nil
}

func (m *Mesh) etcHostsUpdate(nodeList []*node) error {
	// Flatcar's default
	lines := []string{
		"#",
		"# autogenerated by i12e - do not edit",
		"#",
		"",
		"127.0.0.1 localhost",
		"::1 localhost",
	}
	for _, node := range nodeList {
		lines = append(lines, fmt.Sprintf("%s %s", node.NodeIP(), node.NodeName()))
	}
	lines = append(lines, "") // NL to the end
	buf := strings.Join(lines, "\n")
	return os.WriteFile("/etc/hosts", []byte(buf), 0644)
}

func (m *Mesh) NodeConfig(nodeLocal *node, wgInterface string, wgEndpointPort uint16, wgPrivKeyFname string) error {
	if err := m.ifaceLocalConfig(nodeLocal, wgInterface, wgEndpointPort, wgPrivKeyFname); err != nil {
		return err
	}
	nodeList, err := m.getNodeList(nodeLocal)
	if err != nil {
		return err
	}
	if err := m.ifacePeersConfig(nodeList, wgInterface); err != nil {
		return err
	}
	if err := m.etcHostsUpdate(nodeList); err != nil {
		return err
	}
	return nil
}

func newMesh(base string) (*Mesh, error) {
	ii, err := IfaceIP()
	if err != nil {
		log.Error("newNet(): 'IfaceIP()' failed", "err", err)
		return nil, err
	}
	if ii == nil {
		log.Error("newNet(): 'IfaceIP()' returned 'nil'")
		return nil, err
	}
	nodeLocal, err := getNodeLocal()
	if err != nil {
		log.Error("newNet(): 'getNodeLocal()' failed", "err", err)
		return nil, err
	}
	if err := nodeLocal.SetHostname(); err != nil {
		return nil, err
	}
	wgPrivKey, err := getWgPrivKey(WgPrivKeyFname)
	if err != nil {
		log.Error("newNet(): 'getWgPrivKey()' failed", "err", err)
		return nil, err
	}
	wgPubKey, err := getWgPubKey(WgPrivKeyFname)
	if err != nil {
		log.Error("newNet(): 'getWgPubKey()' failed", "err", err)
		return nil, err
	}
	re, err := getRegex()
	if err != nil {
		return nil, err
	}
	return &Mesh{
		base:           base,
		re:             re,
		Tnow:           time.Now().UTC(),
		nodeLocal:      nodeLocal,
		WgInterface:    "wgi",
		WgEndpointIp:   ii.IP,
		WgEndpointPort: 51830, // default '51820'
		WgPrivKey:      wgPrivKey,
		WgPubKey:       wgPubKey,
	}, nil
}

func (m *Mesh) run() error {
	log.Info("Net.run()", "nodeLocal", m.nodeLocal)
	log.Info("Net.run()", "Tnow", m.Tnow)
	log.Info("Net.run()", "WgInterface", m.WgInterface)
	log.Info("Net.run()", "WgEndpointIp", m.WgEndpointIp)
	log.Info("Net.run()", "WgEndpointPort", m.WgEndpointPort)
	log.Info("Net.run()", "WgPrivKey", m.WgPrivKey, "WgPrivKeyLen", m.WgPrivKey.Len())
	log.Info("Net.run()", "WgPubKey", m.WgPubKey, "WgPubKeyHex", m.WgPubKey.Hex(), "WgPubKeyLen", m.WgPubKey.Len())
	if err := m.NodePush(m.nodeLocal, m.Tnow, m.WgPubKey, m.WgEndpointIp, m.WgEndpointPort); err != nil {
		return err
	}
	if err := m.NodeConfig(m.nodeLocal, m.WgInterface, m.WgEndpointPort, WgPrivKeyFname); err != nil {
		return err
	}
	log.Info("Mesh.run()", "mesh", m)
	return nil
}

func Run() error {
	mesh, err := newMesh(remMeshBase)
	if err != nil {
		return err
	}
	return mesh.run()
}
