package wgkey

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

type WgKeyPriv struct {
	k32      *K32
	wgKeyPub *WgKeyPub
}

func (w *WgKeyPriv) String() string {
	return strings.Repeat("*", 32)
}

func (w *WgKeyPriv) K32() *K32 {
	return w.k32
}

func (w *WgKeyPriv) Pub() *WgKeyPub {
	return w.wgKeyPub
}

func privToPub(k32 *K32) (*WgKeyPub, error) {
	cmd := exec.Command("wg", "pubkey")
	cmd.Stdin = bytes.NewBuffer([]byte(k32.B64()))
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return nil, fmt.Errorf("'wg pubkey' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	k32_out, err := NewK32(WithB64(bo.String()))
	if err != nil {
		return nil, err
	}
	return NewWgKeyPub(k32_out), nil
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
	k32, err := NewK32(WithB64(string(buf)))
	if err != nil {
		return nil, err
	}
	wgKeyPub, err := privToPub(k32)
	if err != nil {
		return nil, err
	}
	return &WgKeyPriv{k32, wgKeyPub}, nil
}
