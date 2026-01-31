package pull

import (
	_ "embed"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/mesh"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

//go:embed static/artifact-pull.sh
var rcloneScript string

type Pull struct{}

func newPull() *Pull {
	return &Pull{}
}

func (p *Pull) run() error {
	if err := cmdutil.RunCmd("/bin/sh", "-c", rcloneScript); err != nil {
		log.Error("Pull.run(): rcloneScript failed", "err", err)
		return err
	}
	ii, err := mesh.IfaceIP()
	if err != nil {
		log.Error("Pull.run(): 'net.IfaceIP()' failed", "err", err)
		return err
	}
	iface, ip := "", ""
	if ii != nil {
		iface, ip = ii.Iface, ii.IP
	}
	log.Info("Pull.run()", "iface", iface, "ip", ip)
	if err := cmdutil.RunCmd("/opt/libexec/i12e/artifact-tune.sh", iface, ip); err != nil {
		log.Error("/opt/libexec/i12e/artifact-tune.sh failed", "err", err, "iface", iface, "ip", ip)
		return err
	}
	return nil
}

func Run() error {
	return newPull().run()
}
