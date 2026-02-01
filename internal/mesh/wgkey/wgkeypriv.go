package wgkey

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
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
