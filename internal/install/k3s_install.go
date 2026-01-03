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
	return fsutil.FileContentSet("/etc/systemd/system/k3s.service.d/override.conf", k3sOverrideConfBuf)
}

func K3sInstall() {
	k3sOverrideConf()
}
