package net

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

const WgPrivKeyFname = "/etc/i12e/wg-priv-key" // FIXME unhardcode this

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

func (w *WgKey) Hex() string {
	return hex.EncodeToString(w.data)
}

func (w *WgKey) Len() int {
	return len(w.data)
}

func getWgKeyFromHex(s string, private bool) (*WgKey, error) {
	s_len := len(s)
	if s_len != 64 {
		return nil, fmt.Errorf("len(hex)=%d (64 expected)", s_len)
	}
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return &WgKey{Private: private, data: data}, nil
}

func getWgPubKey() (*WgKey, error) {
	buf, err := os.ReadFile(WgPrivKeyFname)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("wg", "pubkey")
	cmd.Stdin = bytes.NewBuffer(buf)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("getWgPubKey(): 'wg pubkey' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	data, err := base64.StdEncoding.DecodeString(bo.String())
	if err != nil {
		return nil, err
	}
	data_len := len(data)
	if data_len != 32 {
		return nil, fmt.Errorf("getWgPubKey(): len(data)=%d (32 expected)", data_len)
	}
	return &WgKey{Private: false, data: data}, nil
}

func getWgPrivKey() (*WgKey, error) {
	for i := range 2 {
		_, err := os.Stat(WgPrivKeyFname)
		if err == nil {
			break
		}
		if !errors.Is(err, os.ErrNotExist) || i > 0 {
			return nil, fmt.Errorf("getwgPrivKey(): 'os.Stat()' failed: %s", err)
		}
		cmd := exec.Command("wg", "genkey")
		bo, be, err := cmdutil.RunSimple(cmd)
		if err != nil {
			return nil, fmt.Errorf("getWgPrivKey(): 'wg genkey' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
		}
		if err := os.WriteFile(WgPrivKeyFname, bo.Bytes(), 0600); err != nil {
			return nil, fmt.Errorf("getWgPrivKey(): 'os.WriteFile()' failed: %s", err)
		}
	}
	buf, err := os.ReadFile(WgPrivKeyFname)
	if err != nil {
		return nil, err
	}
	data, err := base64.StdEncoding.DecodeString(string(buf))
	if err != nil {
		return nil, err
	}
	data_len := len(data)
	if data_len != 32 {
		return nil, fmt.Errorf("getWgPrivKey(): len(data)=%d (32 expected)", data_len)
	}
	return &WgKey{Private: true, data: data}, nil
}
