package pull

import (
	"fmt"
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

func iface_and_ip() (*II, error) {
	path := "/etc/i12e/iface.txt"
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	iname := strings.TrimSpace(string(buf))
	iface, err := net.InterfaceByName(iname)
	if err != nil {
		return nil, fmt.Errorf("net.InterfaceByName() failed: %s", err.Error())
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("iface.Addrs() failed: %s", err.Error())
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
		return &II{
			Iface: iname,
			IP:    ip.String(),
		}, nil
	}
	return nil, fmt.Errorf("no IPv4 address found")
}

func Pull() {
	log.Info("Pull()...")
	ii, err := iface_and_ip()
	if err != nil {
		log.Error("iface_and_ip() failed", "err", err)
	} else {
		log.Info("iface_and_ip() ok", "Iface", ii.Iface, "IP", ii.IP)
	}
	cmdutil.RunCmd("/bin/sh", "-c", script)
}
