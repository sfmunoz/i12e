package pull

import (
	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "pull")

func Pull() {
	cmd := "rclone cat rem:artifact.tar.gz | tar -C / -xvz"
	log.Info("$ " + cmd)
	cmdutil.RunCmd("/bin/sh", "-c", cmd)
	log.Info("systemctlDaemonReload() complete")
}
