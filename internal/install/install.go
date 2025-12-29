package install

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("app", "i12e").
	With("mod", "install")

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
ExecStart=/opt/bin/i12e
`

func i12eServiceInstall() {
	fname := "/etc/systemd/system/i12e.service"
	err := os.WriteFile(fname, []byte(i12eService), 0644)
	if err != nil {
		log.Fatal("error: 'os.WriteFile()' failed", "fname", fname, "err", err)
	}
	log.Info("i12e_service_install() complete", "fname", fname)
}

func runCmd(name string, arg ...string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Fatal("error: 'cmd.Run()' failed", "err", err, "stderr", stderr.String())
	}
	for line := range strings.SplitSeq(stdout.String(), "\n") {
		if line == "" {
			continue
		}
		log.Info(line)
	}
	log.Info("i12e_service_enable() complete")
}

func i12eServiceEnable() {
	runCmd("systemctl", "enable", "i12e.service")
	log.Info("i12e_service_enable() complete")
}

func i12eServiceStart() {
	runCmd("systemctl", "start", "i12e.service")
	log.Info("i12e_service_start() complete")
}

func Install() {
	log.Info("installing...")
	i12eServiceInstall()
	i12eServiceEnable()
	i12eServiceStart()
}
