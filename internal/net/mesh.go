package net

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

type Mesh struct {
	base string
	re   *regexp.Regexp
	data []string
}

func (m *Mesh) String() string {
	return strings.Join(m.data, "\n")
}

func (m *Mesh) DumpToLog() {
	for _, entry := range m.data {
		x := m.re.FindStringSubmatch(entry)
		for i, v := range x {
			log.Info(">>", "entry", entry, "i", i, "v", v)
		}
	}
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
	re, err := regexp.Compile("^([0-9]{3})/([0-9]{3})/([0-9]{8})_([0-9]{6})_([0-9]{9})/.*")
	cmd := exec.Command("rclone", "lsf", "-R", "--files-only", base)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("getMesh() failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	data := strings.Split(strings.TrimSpace(bo.String()), "\n")
	slices.Sort(data)
	return &Mesh{base, re, data}, nil
}
