package cmdutil

import (
	"bytes"
	"os/exec"
)

func RunSimple(cmd *exec.Cmd) (*bytes.Buffer, *bytes.Buffer, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return &stdout, &stderr, err
	}
	return &stdout, &stderr, nil
}
