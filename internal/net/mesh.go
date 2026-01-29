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

func getRegex() (*regexp.Regexp, error) {
	// 252/158/20260128_152841_153793688/54a2cc8d5e78755ff1debc4a4e6b2fa657ccf86a868b53f9f1b5140487377cc8/192.168.56.53/51820
	expr := "^([0-9]{3})" +
		"/([0-9]{3})" +
		"/([0-9]{8})" +
		"_([0-9]{6})" +
		"_([0-9]{9})" +
		"/([0-9a-f]{64})" +
		`/([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)` +
		"/([0-9]+)" +
		"$"
	return regexp.Compile(expr)
}

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
		log.Info("==", "entry", entry)
		for i, v := range x {
			log.Info(">>", "i", i, "v", v)
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
		return fmt.Errorf("Mesh.NodePush(): 'net.sh' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return nil
}

func (m *Mesh) NodeConfig(wgInterface string, wgIpInt string, wgEndpointPort uint16, wgPrivKeyFname string, wgPathName string) error {
	cmd := exec.Command(
		"/bin/sh",
		"-s",
		"-",
		"node-config",
		wgInterface,
		wgIpInt,
		fmt.Sprintf("%d", wgEndpointPort),
		wgPrivKeyFname,
	)
	cmd.Stdin = bytes.NewBuffer(netSh)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return fmt.Errorf("Mesh.NodeConfig(): 'net.sh' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	for _, entry := range m.data {
		x := m.re.FindStringSubmatch(entry)
		if x == nil {
			continue
		}
		if wgPathName == fmt.Sprintf("%s/%s", x[1], x[2]) {
			continue // it's my entry
		}
		k, err := getWgKeyFromHex(x[6], false)
		if err != nil {
			log.Error("error getting 'wg-key'", "err", err, "hex", x[6])
			continue
		}
		cmd := exec.Command(
			"wg", "set", wgInterface,
			"peer", k.B64(),
			"allowed-ips", fmt.Sprintf("10.56.%s.%s/32", strings.TrimLeft(x[1], "0"), strings.TrimLeft(x[2], "0")),
			"endpoint", fmt.Sprintf("%s:%s", x[7], x[8]),
		)
		bo, be, err := cmdutil.RunSimple(cmd)
		if err != nil {
			return fmt.Errorf("'wg set' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
		}
	}
	return nil
}

func getMesh(base string) (*Mesh, error) {
	re, err := getRegex()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("rclone", "lsf", "-R", "--files-only", base)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("getMesh() failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	data := strings.Split(strings.TrimSpace(bo.String()), "\n")
	slices.Sort(data)
	return &Mesh{base, re, data}, nil
}
