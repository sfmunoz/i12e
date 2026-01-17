package butane

import (
	"github.com/sfmunoz/logit"
	"github.com/spf13/cobra"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "butane")

func Run(cmd *cobra.Command, args []string) {
	log.Info("butane.Run()")
}
