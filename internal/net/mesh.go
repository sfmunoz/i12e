package net

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

type Mesh struct {
	base string
	data []string
}

func (m *Mesh) String() string {
	return strings.Join(m.data, "\n")
}

func (m *Mesh) NodePush(nodePath string, ts time.Time, wgPubKey *WgKey, wgEndpointIp string, wgEndpointPort uint16) error {
	cmd := exec.Command(
		"/bin/sh",
		"-s",
		"-",
		"node-push",
		nodePath,
		fmt.Sprintf("%s_%09d", ts.Format("20060102_150405"), ts.Nanosecond()),
		wgPubKey.Hex(),
		wgEndpointIp,
		fmt.Sprintf("%d", wgEndpointPort),
	)
	cmd.Stdin = bytes.NewBuffer(netSh)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("Net.push(): 'net.sh' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return nil
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
