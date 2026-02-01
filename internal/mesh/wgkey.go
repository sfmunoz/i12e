package mesh

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

type wgKeyPriv struct {
	data []byte
}

func (w *wgKeyPriv) String() string {
	return strings.Repeat("*", len(w.data))
}

func (w *wgKeyPriv) Raw() []byte {
	return w.data
}

func (w *wgKeyPriv) B64() string {
	return base64.StdEncoding.EncodeToString(w.data)
}

func (w *wgKeyPriv) Hex() string {
	return hex.EncodeToString(w.data)
}

func (w *wgKeyPriv) Len() int {
	return len(w.data)
}

type wgKeyPub struct {
	data []byte
}

func (w *wgKeyPub) String() string {
	return w.B64()
}

func (w *wgKeyPub) Raw() []byte {
	return w.data
}

func (w *wgKeyPub) B64() string {
	return base64.StdEncoding.EncodeToString(w.data)
}

func (w *wgKeyPub) Hex() string {
	return hex.EncodeToString(w.data)
}

func (w *wgKeyPub) Len() int {
	return len(w.data)
}

type wgKey struct {
	privKey *wgKeyPriv
	pubKey  *wgKeyPub
}

func (w *wgKey) getPrivKey() *wgKeyPriv {
	return w.privKey
}

func (w *wgKey) getPubKey() *wgKeyPub {
	return w.pubKey
}

func (w *wgKey) getLocal() bool {
	return w.getPrivKey() != nil
}

func getWgKeyRemote(s string) (*wgKey, error) {
	s_len := len(s)
	if s_len != 64 {
		return nil, fmt.Errorf("len(hex)=%d (64 expected)", s_len)
	}
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return &wgKey{
		privKey: nil,
		pubKey:  &wgKeyPub{data},
	}, nil
}

func privToPub(kPriv *wgKeyPriv) (*wgKeyPub, error) {
	cmd := exec.Command("wg", "pubkey")
	cmd.Stdin = bytes.NewBuffer([]byte(kPriv.B64()))
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
	return &wgKeyPub{data}, nil
}

func getWgKeyLocal(wgPrivKeyFname string) (*wgKey, error) {
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
	privKey := &wgKeyPriv{data}
	pubKey, err := privToPub(privKey)
	if err != nil {
		return nil, err
	}
	return &wgKey{privKey, pubKey}, nil
}
