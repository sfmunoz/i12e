package server

import (
	_ "embed"
	"fmt"
	"net"
	"os"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/mesh/node"
)

func meshCfg(cfg *config.ServerConfig) error {
	fname := "/etc/i12e/k3s/mesh.cfg"
	_, err := os.Stat(fname)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		log.Error("meshCfg(): os.Stat() failed", "err", err, "fname", fname)
		return err
	}
	iface := cfg.Mesh.WireGuardInterface
	if _, err := net.InterfaceByName(iface); err != nil {
		log.Notice("meshCfg(): network interface doesn't exist yet", "iface", iface)
		return nil
	}
	nodeLocal, err := node.NewNodeLocal(cfg, false)
	if err != nil {
		return err
	}
	ip := nodeLocal.GetMeshIP().String()
	log.Info("meshCfg(): creating config", "fname", fname, "iface", iface, "ip", ip)
	buf := fmt.Sprintf("IFACE=\"%s\"\nIP=\"%s\"\n", iface, ip)
	if err := os.WriteFile(fname, []byte(buf), 0600); err != nil {
		log.Error("meshCfg(): os.WriteFile() write failed", "err", err, "fname", fname, "buf", buf)
		return err
	}
	return nil
}
