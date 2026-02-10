package wgkey

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

type WgKeyPriv struct {
	data     []byte
	wgKeyPub *WgKeyPub
}

func (w *WgKeyPriv) String() string {
	return strings.Repeat("*", len(w.data))
}

func (w *WgKeyPriv) Raw() []byte {
	return w.data
}

func (w *WgKeyPriv) B64() string {
	return base64.StdEncoding.EncodeToString(w.data)
}

func (w *WgKeyPriv) Hex() string {
	return hex.EncodeToString(w.data)
}

func (w *WgKeyPriv) Len() int {
	return len(w.data)
}

func (w *WgKeyPriv) Pub() *WgKeyPub {
	return w.wgKeyPub
}

func privToPub(b []byte) (*WgKeyPub, error) {
	cmd := exec.Command("wg", "pubkey")
	cmd.Stdin = bytes.NewBuffer([]byte(base64.StdEncoding.EncodeToString(b)))
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("'wg pubkey' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	data, err := base64.StdEncoding.DecodeString(bo.String())
	if err != nil {
		return nil, err
	}
	data_len := len(data)
	if data_len != 32 {
		return nil, fmt.Errorf("len(data)=%d (32 expected)", data_len)
	}
	return &WgKeyPub{data}, nil
}

func NewWgKeyPriv(wgPrivKeyFname string) (*WgKeyPriv, error) {
	for i := range 2 {
		_, err := os.Stat(wgPrivKeyFname)
		if err == nil {
			break
		}
		if !errors.Is(err, os.ErrNotExist) || i > 0 {
			return nil, fmt.Errorf("'os.Stat()' failed: %s", err)
		}
		cmd := exec.Command("wg", "genkey")
		bo, be, err := cmdutil.RunSimple(cmd)
		if err != nil {
			return nil, fmt.Errorf("'wg genkey' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
		}
		if err := os.WriteFile(wgPrivKeyFname, bo.Bytes(), 0600); err != nil {
			return nil, fmt.Errorf("'os.WriteFile()' failed: %s", err)
		}
	}
	buf, err := os.ReadFile(wgPrivKeyFname)
	if err != nil {
		return nil, err
	}
	data, err := base64.StdEncoding.DecodeString(string(buf))
	if err != nil {
		return nil, err
	}
	data_len := len(data)
	if data_len != 32 {
		return nil, fmt.Errorf("len(data)=%d (32 expected)", data_len)
	}
	wgKeyPub, err := privToPub(data)
	if err != nil {
		return nil, err
	}
	return &WgKeyPriv{data, wgKeyPub}, nil
}
