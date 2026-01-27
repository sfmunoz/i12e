package net

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

type WgKey struct {
	Private bool
	data    []byte
}

func (w *WgKey) String() string {
	if w.Private {
		return strings.Repeat("*", len(w.data))
	}
	return w.B64()
}

func (w *WgKey) B64() string {
	return base64.StdEncoding.EncodeToString(w.data)
}

func getWgKey(private bool) (*WgKey, error) {
	c := "pub-key"
	if private {
		c = "priv-key"
	}
	cmd := exec.Command("/bin/sh", "-s", "-", c, "/etc/i12e/wg-priv-key")
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
	return &WgKey{
		Private: private,
		data:    data,
	}, nil
}
