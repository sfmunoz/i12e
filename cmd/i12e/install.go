package main

import "os"

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

func install() {
	log.Info("installing...")
	fname := "/etc/systemd/system/i12e.service"
	err := os.WriteFile(fname, []byte(i12e_service), 0644)
	if err != nil {
		log.Fatal("os.WriteFile() failed", "fname", fname, "err", err)
	}
}
