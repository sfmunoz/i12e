package artifact

import (
	"fmt"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "artifact")

func Run(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("artifact.Run(): undefined config")
	}
	log.Info("artifact.Run()", "cfg", cfg)
	return nil
}
