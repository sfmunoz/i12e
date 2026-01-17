package artifact

import (
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "artifact")

func Run(prod bool) {
	log.Info("artifact.Run()", "prod", prod)
}
