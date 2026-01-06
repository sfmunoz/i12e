package cmdutil

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "cmdutil")

func RunCmd(name string, arg ...string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Fatal("error: 'cmd.Run()' failed", "err", err, "stderr", stderr.String())
	}
	for line := range strings.SplitSeq(stdout.String(), "\n") {
		if line == "" {
			continue
		}
		log.Info(line)
	}
}
