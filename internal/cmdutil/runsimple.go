package cmdutil

import (
	"bytes"
	"os/exec"
)

func RunSimple(name string, arg ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return &stdout, &stderr, err
	}
	return &stdout, &stderr, nil
}
