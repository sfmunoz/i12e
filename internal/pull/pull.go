package pull

import (
	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/netutil"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

const rcloneScript = `#!/bin/sh
[ -f /etc/i12e/flags/artifact-pulled ] && exit 0
set -x -e -o pipefail
rclone cat rem:artifact.tar.gz | tar -C / -xvz
rm -f /etc/i12e/flags/artifact-tuned
`

type Pull struct{}

func newPull() *Pull {
	return &Pull{}
}

func (p *Pull) run() error {
	err := cmdutil.RunCmd("/bin/sh", "-c", rcloneScript)
	if err != nil {
		log.Error("Pull.run(): rcloneScript failed", "err", err)
		return err
	}
	ii, err := netutil.IfaceIP()
	if err != nil {
		log.Error("Pull.run(): 'netutil.IfaceIP()' failed", "err", err)
		return err
	}
	iface, ip := "", ""
	if ii != nil {
		iface, ip = ii.Iface, ii.IP
	}
	log.Info("Pull.run()", "iface", iface, "ip", ip)
	err = cmdutil.RunCmd("/opt/libexec/i12e/artifact-tune.sh", iface, ip)
	if err != nil {
		log.Error("/opt/libexec/i12e/artifact-tune.sh failed", "err", err, "iface", iface, "ip", ip)
		return err
	}
	return nil
}

func Run() error {
	return newPull().run()
}
