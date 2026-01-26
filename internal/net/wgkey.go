package net

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

type WgCmd string

const (
	WgCmdPrivKey WgCmd = "priv-key"
	WgCmdPubKey  WgCmd = "pub-key"
)

func (w WgCmd) String() string {
	return string(w)
}

type WgKey []byte

func getWgKey(c WgCmd) (WgKey, error) {
	cmd := exec.Command("/bin/sh", "-s", "-", string(c), "/etc/i12e/wg-priv-key")
	cmd.Stdin = bytes.NewBuffer(netSh)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("getWgKey(): 'net.sh' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	data, err := base64.StdEncoding.DecodeString(bo.String())
	if err != nil {
		return nil, err
	}
	data_len := len(data)
	if data_len != 32 {
		return nil, fmt.Errorf("getWgKey(): len(data)=%d (32 expected)", data_len)
	}
	return WgKey(data), nil
}
