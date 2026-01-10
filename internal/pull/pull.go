package pull

import (
	"bytes"
	"errors"
	"net"
	"os"
	"strings"
	"text/template"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "pull")

const scriptTplBuf = `#!/bin/sh
FLAG_FILE="/etc/i12e/z.flag"
REBOOT_FILE="/etc/i12e/reboot-required"
IFACE="{{ .Iface }}"
IP="{{ .Ip }}"
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
    if [ "$IFACE" != "" -a "$IP" != "" ]
    then
      echo "node-ip: \"${IP}\"" >> /etc/rancher/k3s/config.yaml
      echo "flannel-iface: \"${IFACE}\"" >> /etc/rancher/k3s/config.yaml
    fi
    touch "$REBOOT_FILE"
    { set +x; } 2>/dev/null
  fi
}
pull_if_needed
reboot_if_required
`

var scriptTpl = template.Must(template.New("script").Parse(scriptTplBuf))

type II struct {
	Iface string
	IP    string
}

func ifaceIP() *II {
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

func script() string {
	var buf bytes.Buffer
	iface := ""
	ip := ""
	ii := ifaceIP()
	if ii != nil {
		log.Info("ifaceIP() ok", "Iface", ii.Iface, "IP", ii.IP)
		iface = ii.Iface
		ip = ii.IP
	}
	err := scriptTpl.Execute(&buf, map[string]string{"Iface": iface, "Ip": ip})
	if err != nil {
		log.Fatal("scriptTpl.Execute() failed", "err", err, "iface", iface, "Ip", ip)
	}
	return buf.String()
}

func Pull() {
	log.Info("Pull()...")
	cmdutil.RunCmd("/bin/sh", "-c", script())
}
