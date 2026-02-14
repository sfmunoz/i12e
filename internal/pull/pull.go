package pull

import (
	_ "embed"
	"net"
	"net/netip"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/mesh/node"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

//go:embed static/artifact-pull.sh
var rcloneScript string

type Pull struct {
	meshNet *netip.Prefix
}

func (p *Pull) run() error {
	iface := node.GetNodeInterface()
	if _, err := net.InterfaceByName(iface); err != nil {
		log.Notice("Pull.run(): network interface doesn't exist yet", "iface", iface)
		return nil
	}
	if err := cmdutil.RunCmd("/bin/sh", "-c", rcloneScript); err != nil {
		log.Error("Pull.run(): rcloneScript failed", "err", err)
		return err
	}
	nodeLocal, err := node.NewNodeLocal(p.meshNet, false)
	if err != nil {
		return err
	}
	ip := nodeLocal.GetMeshIP().String()
	log.Info("Pull.run()", "iface", iface, "ip", ip)
	if err := cmdutil.RunCmd("/opt/libexec/i12e/artifact-tune.sh", iface, ip); err != nil {
		log.Error("/opt/libexec/i12e/artifact-tune.sh failed", "err", err, "iface", iface, "ip", ip)
		return err
	}
	return nil
}

func newPull(meshNet *netip.Prefix) *Pull {
	return &Pull{meshNet}
}

func Run(meshNet *netip.Prefix) error {
	return newPull(meshNet).run()
}
