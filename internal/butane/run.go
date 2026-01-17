package butane

import (
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "butane")

func Run(prod bool) {
	log.Info("butane.Run()", "prod", prod)
}
