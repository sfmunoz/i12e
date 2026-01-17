package artifact

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "artifact")

func Run(prod bool) error {
	log.Info("artifact.Run()", "prod", prod)
	fname := "secrets-dev.yaml"
	if prod {
		fname = "secrets-prod.yaml"
	}
	buf, err := cmdutil.SopsDecrypt(fname)
	if err != nil {
		return fmt.Errorf("cmdutil.SopsDecrypt() failed: err=%s; prod=%t", err, prod)
	}
	log.Info("cmdutil.SopsDecrypt() OK", "buf", buf)
	return nil
}
