package server

import (
	_ "embed"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

//go:embed static/artifact-pull.sh
var artifactPullSh string

func artifactPull() error {
	cmd := exec.Command("/bin/bash", "-c", artifactPullSh)
	if err := cmdutil.RunCmd(cmd); err != nil {
		log.Error("artifactPull(): artifactPullSh failed", "err", err)
		return err
	}
	return nil
}
