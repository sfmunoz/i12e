package server

import (
	"bytes"
	_ "embed"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

//go:embed static/rclone-pull.sh
var rclonePullSh []byte

func rclonePull() error {
	cmd := exec.Command("bash")
	cmd.Stdin = bytes.NewBuffer(rclonePullSh)
	return cmdutil.RunCmd(cmd)
}
