package wgkey

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
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
	kPriv, err := wgtypes.NewKey(k32.Raw())
	if err != nil {
		return nil, err
	}
	kPub := kPriv.PublicKey()
	k32_out, err := NewK32(WithBytes(kPub[:]))
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
		key, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(wgPrivKeyFname, []byte(key.String()+"\n"), 0600); err != nil {
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
