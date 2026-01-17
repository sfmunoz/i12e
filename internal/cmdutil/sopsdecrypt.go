package cmdutil

import (
	"bytes"
	"os/exec"
)

func SopsDecrypt(fname string) (*bytes.Buffer, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("sops", "decrypt", fname)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return &stderr, err
	}
	return &stdout, nil
}
