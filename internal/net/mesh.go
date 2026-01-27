package net

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

type Mesh struct {
	base string
	data []string
}

func (m *Mesh) String() string {
	return strings.Join(m.data, "\n")
}

func getMesh(base string) (*Mesh, error) {
	cmd := exec.Command("rclone", "lsf", "-R", "--files-only", base)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("getMesh() failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	buf := strings.TrimSpace(bo.String())
	return &Mesh{
		base: base,
		data: strings.Split(buf, "\n"),
	}, nil
}
