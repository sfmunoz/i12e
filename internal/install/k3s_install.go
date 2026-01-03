package install

import (
	"fmt"
	"strings"

	"github.com/sfmunoz/i12e/internal/fsutil"
)

const k3sOverrideConfBuf = `[Service]
ExecStartPre=-/usr/bin/sh -c 'rm -f /var/lib/rancher/k3s/server/db/state.db*'
ExecStart=
ExecStart=/opt/bin/k3s server
`

func k3sConfigYaml() error {
	// https://docs.k3s.io/installation/configuration
	file := "/etc/rancher/k3s/config.yaml"
	position := 1
	k3s_cmd := "server"
	token := "main-token"
	agentToken := "agent-token"
	tlsSan := "192.168.56.50"
	nodeIp := "192.168.56.51"
	flannelIface := "enp0s8"
	lines := make([]string, 0, 20)
	if k3s_cmd == "server" {
		lines = append(lines, fmt.Sprintf("token: \"%s\"", token))
		lines = append(lines, fmt.Sprintf("agent-token: \"%s\"", agentToken))
		lines = append(lines, "secrets-encryption: true")
		lines = append(lines, "secrets-encryption-provider: secretbox")
		lines = append(lines, "flannel-backend: \"wireguard-native\"")
		if tlsSan != "" {
			lines = append(lines, fmt.Sprintf("tls-san: \"%s\"", tlsSan))
		}
	} else {
		lines = append(lines, fmt.Sprintf("token: \"%s\"", agentToken))
	}
	if position == 1 {
		lines = append(lines, "cluster-init: true")
	} else {
		lines = append(lines, fmt.Sprintf("server: \"https://%s:6443\"", tlsSan))
	}
	if nodeIp != "" {
		lines = append(lines, fmt.Sprintf("node-ip: \"%s\"", nodeIp))
	}
	if flannelIface != "" {
		lines = append(lines, fmt.Sprintf("flannel-iface: \"%s\"", flannelIface))
	}
	lines = append(lines, "")
	buf := strings.Join(lines, "\n")
	updated, err := fsutil.FileContentSet(file, buf)
	if err != nil {
		return err
	}
	if updated {
		log.Info("k3sConfigYaml(): file updated", "file", file)
	}
	return nil
}

func k3sOverrideConf() error {
	file := "/etc/systemd/system/k3s.service.d/override.conf"
	updated, err := fsutil.FileContentSet(file, k3sOverrideConfBuf)
	if err != nil {
		return err
	}
	if updated {
		log.Info("k3sOverrideConf(): file updated", "file", file)
	}
	return nil
}

func K3sInstall() {
	k3sConfigYaml()
	k3sOverrideConf()
}
