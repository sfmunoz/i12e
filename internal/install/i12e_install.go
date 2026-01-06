package install

import (
	"os"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

const i12eService = `[Unit]
Description=i12e
Wants=network-online.target
After=network-online.target

[Install]
WantedBy=multi-user.target

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStartPre=-/usr/bin/rm -f /etc/i12e/reboot-required
ExecStart=/bin/sh -c 'if [ -x /opt/bin/i12e ]; then exec /opt/bin/i12e; else exec /usr/bin/i12e; fi'
`

func i12eServiceInstall() {
	fname := "/etc/systemd/system/i12e.service"
	err := os.WriteFile(fname, []byte(i12eService), 0644)
	if err != nil {
		log.Fatal("error: 'os.WriteFile()' failed", "fname", fname, "err", err)
	}
	log.Info("i12eServiceInstall() complete", "fname", fname)
}

func systemctlDaemonReload() {
	log.Info("$ systemctl daemon-reload")
	cmdutil.RunCmd("systemctl", "daemon-reload")
	log.Info("systemctlDaemonReload() complete")
}

func i12eServiceEnable() {
	log.Info("$ systemctl enable i12e.service")
	cmdutil.RunCmd("systemctl", "enable", "i12e.service")
	log.Info("i12eServiceEnable() complete")
}

func i12eServiceStart() {
	log.Info("$ systemctl start i12e.service")
	cmdutil.RunCmd("systemctl", "start", "i12e.service")
	log.Info("i12eServiceStart() complete")
}

func Install() {
	log.Info("installing...")
	i12eServiceInstall()
	systemctlDaemonReload()
	i12eServiceEnable()
	i12eServiceStart()
}
