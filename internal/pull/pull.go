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

const script1 = `#!/bin/sh
[ -f /etc/i12e/pull-done ] && exit 0
set -x -e -o pipefail
rclone cat rem:artifact.tar.gz | tar -C / -xvz
rm -f /etc/i12e/config-patched
`

const scriptBuf2 = `#!/bin/sh
FLAG_FILE="/etc/i12e/config-patched"
[ -f "$FLAG_FILE" ] && exit 0
IFACE="{{ .Iface }}"
IP="{{ .Ip }}"
[ "$IFACE" = "" -o "$IP" = "" ] && exit 0
set -x -e -o pipefail
awk \
  -v IFACE="$IFACE" \
  -v IP="$IP" \
  '!/^(node-ip|flannel-iface):/ {
    print
  }
  END {
    printf("flannel-iface: \"%s\"\nnode-ip: \"%s\"\n",IFACE,IP)
  }' \
  /etc/rancher/k3s/config.yaml > /etc/rancher/k3s/config.yaml.new
cat /etc/rancher/k3s/config.yaml.new > /etc/rancher/k3s/config.yaml
rm -f /etc/rancher/k3s/config.yaml.new
touch "$FLAG_FILE"
`

var scriptTpl2 = template.Must(template.New("script").Parse(scriptBuf2))

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
			continue
		}
		ip := ipNet.IP
		if ip.To4() == nil {
			continue
		}
		ipStr := ip.String()
		if len(ipStr) < 1 {
			log.Error("ifaceIP(): empty ipStr", "iname", iname)
			continue
		}
		return &II{
			Iface: iname,
			IP:    ip.String(),
		}
	}
	log.Warn("ifaceIP(): no IPv4 address found")
	return nil
}

func script2() string {
	var buf bytes.Buffer
	iface := ""
	ip := ""
	ii := ifaceIP()
	if ii != nil {
		log.Info("ifaceIP() ok", "Iface", ii.Iface, "IP", ii.IP)
		iface = ii.Iface
		ip = ii.IP
	}
	err := scriptTpl2.Execute(&buf, map[string]string{"Iface": iface, "Ip": ip})
	if err != nil {
		log.Fatal("scriptTpl2.Execute() failed", "err", err, "iface", iface, "Ip", ip)
	}
	return buf.String()
}

func Pull() {
	log.Info("Pull()...")
	cmdutil.RunCmd("/bin/sh", "-c", script1)
	cmdutil.RunCmd("/bin/sh", "-c", script2())
}
