package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func i12e_service_install() {
	const i12e_service = `[Unit]
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
	fname := "/etc/systemd/system/i12e.service"
	err := os.WriteFile(fname, []byte(i12e_service), 0644)
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

func i12e_service_enable() {
	runCmd("systemctl", "enable", "i12e.service")
	log.Info("i12e_service_enable() complete")
}

func i12e_service_start() {
	runCmd("systemctl", "start", "i12e.service")
	log.Info("i12e_service_start() complete")
}

func install() {
	log.Info("installing...")
	i12e_service_install()
	i12e_service_enable()
	i12e_service_start()
}
