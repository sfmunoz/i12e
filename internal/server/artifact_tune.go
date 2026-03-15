package server

import (
	_ "embed"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

func artifactTune() error {
	log.Info("artifactTune()")
	cmd := exec.Command("/opt/libexec/i12e/artifact-tune.sh")
	if err := cmdutil.RunCmd(cmd); err != nil {
		log.Error("/opt/libexec/i12e/artifact-tune.sh failed", "err", err)
		return err
	}
	return nil
}
