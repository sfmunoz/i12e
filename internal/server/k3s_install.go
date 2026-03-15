package server

import (
	_ "embed"
	"net"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/mesh/node"
)

func k3sInstall(cfg *config.ServerConfig) error {
	iface := cfg.Mesh.WireGuardInterface
	if _, err := net.InterfaceByName(iface); err != nil {
		log.Notice("k3sInstall(): network interface doesn't exist yet", "iface", iface)
		return nil
	}
	nodeLocal, err := node.NewNodeLocal(cfg, false)
	if err != nil {
		return err
	}
	ip := nodeLocal.GetMeshIP().String()
	log.Info("k3sInstall()", "iface", iface, "ip", ip)
	cmd := exec.Command("/opt/libexec/i12e/k3s-install.sh", iface, ip)
	if err := cmdutil.RunCmd(cmd); err != nil {
		log.Error("/opt/libexec/i12e/k3s-install.sh failed", "err", err, "iface", iface, "ip", ip)
		return err
	}
	return nil
}
