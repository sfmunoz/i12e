package wgkey

import (
	"encoding/base64"
	"encoding/hex"
)

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
