package server

import (
	_ "embed"
	"net"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/mesh/node"
)

//go:embed static/artifact-pull.sh
var artifactPullSh string

func artifactPull(cfg *config.ServerConfig) error {
	iface := cfg.Mesh.WireGuardInterface
	if _, err := net.InterfaceByName(iface); err != nil {
		log.Notice("Pull.run(): network interface doesn't exist yet", "iface", iface)
		return nil
	}
	cmd := exec.Command("/bin/sh", "-c", artifactPullSh)
	if err := cmdutil.RunCmd(cmd); err != nil {
		log.Error("Pull.run(): rcloneScript failed", "err", err)
		return err
	}
	nodeLocal, err := node.NewNodeLocal(cfg, false)
	if err != nil {
		return err
	}
	ip := nodeLocal.GetMeshIP().String()
	log.Info("Pull.run()", "iface", iface, "ip", ip)
	cmd = exec.Command("/opt/libexec/i12e/artifact-tune.sh", iface, ip)
	if err := cmdutil.RunCmd(cmd); err != nil {
		log.Error("/opt/libexec/i12e/artifact-tune.sh failed", "err", err, "iface", iface, "ip", ip)
		return err
	}
	return nil
}
