package artifact

import (
	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "artifact")

func Run(prod bool) {
	log.Info("artifact.Run()", "prod", prod)
	fname := "secrets-dev.yaml"
	if prod {
		fname = "secrets-prod.yaml"
	}
	buf, err := cmdutil.SopsDecrypt(fname)
	if err != nil {
		log.Fatal("loadConf() failed", "err", err, "prod", prod)
	}
	log.Info("loadConf() OK", "buf", buf)
}
