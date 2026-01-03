package install

import (
	"os"
)

const k3sOverrideConfBuf = `[Service]
ExecStartPre=-/usr/bin/sh -c 'rm -f /var/lib/rancher/k3s/server/db/state.db*'
ExecStart=
ExecStart=/opt/bin/k3s server
`

func k3sOverrideConf() {
	dir := "/etc/systemd/system/k3s.service.d"
	file := dir + "/override.conf"
	os.MkdirAll(dir, 0755)
	s, err := os.Stat(dir)
	if err != nil {
		log.Error("k3sOverrideConf(): os.Stat() failed", "dir", dir, "err", err)
		return
	}
	if !s.Mode().IsDir() {
		log.Error("k3sOverrideConf(): dir doesn't exist but it should", "dir", dir)
		return
	}
	s, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = os.WriteFile(file, []byte(k3sOverrideConfBuf), 0644)
		if err != nil {
			log.Error("k3sOverrideConf(): file creation failed", "file", file, "err", err)
			return
		}
		log.Info("k3sOverrideConf(): file created", "file", file)
		return
	}
	if err != nil {
		log.Error("k3sOverrideConf(): os.Stat() failed", "file", file, "err", err)
		return
	}
	buf, err := os.ReadFile(file)
	if err != nil {
		log.Error("k3sOverrideConf(): os.ReadFile() failed", "file", file, "err", err)
		return
	}
	buf_s := string(buf)
	if buf_s == k3sOverrideConfBuf {
		log.Info("k3sOverrideConf(): file is up-to-date", "file", file)
		return
	}
	err = os.WriteFile(file, []byte(k3sOverrideConfBuf), 0644)
	if err != nil {
		log.Error("k3sOverrideConf(): file update failed", "file", file, "err", err)
		return
	}
	log.Info("k3sOverrideConf(): file updated", "file", file)
}

func K3sInstall() {
	k3sOverrideConf()
}
