package wgkey

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type WgKeyPub struct {
	data []byte
}

func (w *WgKeyPub) String() string {
	return w.B64()
}

func (w *WgKeyPub) Raw() []byte {
	return w.data
}

func (w *WgKeyPub) B64() string {
	return base64.StdEncoding.EncodeToString(w.data)
}

func (w *WgKeyPub) Hex() string {
	return hex.EncodeToString(w.data)
}

func (w *WgKeyPub) Len() int {
	return len(w.data)
}

func NewWgKeyPub(s string) (*WgKeyPub, error) {
	s_len := len(s)
	if s_len != 64 {
		return nil, fmt.Errorf("len(hex)=%d (64 expected)", s_len)
	}
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return &WgKeyPub{data}, nil
}
