package pull

import (
	"errors"
	"net"
	"os"
	"strings"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "pull")

const rcloneScript = `#!/bin/sh
[ -f /etc/i12e/pull-done ] && exit 0
set -x -e -o pipefail
rclone cat rem:artifact.tar.gz | tar -C / -xvz
rm -f /etc/i12e/artifact-tuned
`

type II struct {
	Iface string
	IP    string
}

func ifaceIP() *II {
	path := "/etc/i12e/iface.txt"
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Info("ifaceIP(): file does not exist", "path", path)
		} else {
			log.Error("ifaceIP(): os.Stat() failed", "path", path, "err", err)
		}
		return nil
	}
	buf, err := os.ReadFile(path)
	if err != nil {
		log.Error("ifaceIP(): os.ReadFile() failed", "path", path, "err", err)
		return nil
	}
	iname := strings.TrimSpace(string(buf))
	if len(iname) < 1 {
		log.Error("ifaceIP(): file is empty", "path", path)
		return nil
	}
	iface, err := net.InterfaceByName(iname)
	if err != nil {
		log.Error("ifaceIP(): net.InterfaceByName() failed", "iname", iname, "err", err)
		return nil
	}
	addrs, err := iface.Addrs()
	if err != nil {
		log.Error("ifaceIP(): iface.Addrs() failed", "iname", iname, "err", err)
		return nil
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			log.Info("ifaceIP(): cannot get ipNet", "iname", iname, "addr", addr)
			continue
		}
		ip := ipNet.IP
		if ip.To4() == nil {
			log.Info("ifaceIP(): ip.To4() returned nil", "iname", iname, "addr", addr)
			continue
		}
		ipStr := ip.String()
		if len(ipStr) < 1 {
			log.Info("ifaceIP(): empty ipStr", "iname", iname, "addr", addr)
			continue
		}
		ret := &II{Iface: iname, IP: ip.String()}
		log.Info("ifaceIP() ok", "Iface", ret.Iface, "IP", ret.IP)
		return ret
	}
	log.Warn("ifaceIP(): no IPv4 address found")
	return nil
}

func Pull() {
	log.Info("Pull()...")
	cmdutil.RunCmd("/bin/sh", "-c", rcloneScript)
	iface := ""
	ip := ""
	ii := ifaceIP()
	if ii != nil {
		iface = ii.Iface
		ip = ii.IP
	}
	cmdutil.RunCmd("/opt/libexec/i12e/artifact-tune.sh", iface, ip)
}
