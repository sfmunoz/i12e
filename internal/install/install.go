package install

import (
	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "install")

func Install() {
	log.Info("Install()...")
	cmdutil.RunCmd("/opt/libexec/i12e/k3s-install.sh")
}
