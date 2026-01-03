package install

import (
	"github.com/sfmunoz/i12e/internal/fsutil"
)

const k3sOverrideConfBuf = `[Service]
ExecStartPre=-/usr/bin/sh -c 'rm -f /var/lib/rancher/k3s/server/db/state.db*'
ExecStart=
ExecStart=/opt/bin/k3s server
`

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
	k3sOverrideConf()
}
