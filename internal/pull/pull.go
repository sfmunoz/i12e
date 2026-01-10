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

const script = `#!/bin/sh
FLAG_FILE="/etc/i12e/z.flag"
REBOOT_FILE="/etc/i12e/reboot-required"
set -e -o pipefail
function reboot_if_required {
  [ -f "$REBOOT_FILE" ] || return 0
  set -x
  systemctl daemon-reload
  [ -f /etc/systemd/system/k3s.service ] || touch /etc/i12e/k3s-install-required
  rm -f "$REBOOT_FILE"
  { set +x; } 2> /dev/null
  if [ -f "$REBOOT_FILE" ]
  then
    echo "error: reboot aborted: cannot delete '$REBOOT_FILE' before 'systemctl reboot' execution"
    return 1
  fi
  set -x
  systemctl reboot
}
function pull_if_needed {
  if [ -f "$FLAG_FILE" ]
  then
    echo "pull not needed: '${FLAG_FILE}' already exists"
  else
    echo "pulling artifact.tar.gz provided that '${FLAG_FILE}' doesn't exist..."
    set -x
    rclone cat rem:artifact.tar.gz | tar -C / -xvz
    touch "$REBOOT_FILE"
    { set +x; } 2>/dev/null
  fi
}
pull_if_needed
reboot_if_required
`

type II struct {
	Iface string
	IP    string
}

func iface_and_ip() *II {
	path := "/etc/i12e/iface.txt"
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Info("file does not exist", "path", path)
		} else {
			log.Error("os.Stat() failed", "path", path, "err", err)
		}
		return nil
	}
	buf, err := os.ReadFile(path)
	if err != nil {
		log.Error("os.ReadFile() failed", "path", path, "err", err)
		return nil
	}
	iname := strings.TrimSpace(string(buf))
	if len(iname) < 1 {
		log.Error("file is empty", "path", path)
		return nil
	}
	iface, err := net.InterfaceByName(iname)
	if err != nil {
		log.Error("net.InterfaceByName() failed", "iname", iname, "err", err)
		return nil
	}
	addrs, err := iface.Addrs()
	if err != nil {
		log.Error("iface.Addrs() failed", "iname", iname, "err", err)
		return nil
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipNet.IP
		if ip.To4() == nil {
			continue
		}
		ipStr := ip.String()
		if len(ipStr) < 1 {
			log.Error("empty ipStr", "iname", iname)
			continue
		}
		return &II{
			Iface: iname,
			IP:    ip.String(),
		}
	}
	log.Warn("no IPv4 address found")
	return nil
}

func Pull() {
	log.Info("Pull()...")
	ii := iface_and_ip()
	if ii != nil {
		log.Info("iface_and_ip() ok", "Iface", ii.Iface, "IP", ii.IP)
	}
	cmdutil.RunCmd("/bin/sh", "-c", script)
}
