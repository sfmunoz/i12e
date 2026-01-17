package artifact

import (
	"github.com/sfmunoz/logit"
	"github.com/spf13/cobra"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "artifact")

func Run(cmd *cobra.Command, args []string) {
	prod, err := cmd.Flags().GetBool("prod")
	if err != nil {
		log.Fatal("'cmd.Flags().GetBool()' failed", "err", err)
	}
	log.Info("artifact.Run()", "prod", prod)
}
