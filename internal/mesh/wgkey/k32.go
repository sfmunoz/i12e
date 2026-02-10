package wgkey

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type K32 struct {
	data [32]byte
}

func (k *K32) Raw() []byte {
	return k.data[:]
}

func (k *K32) B64() string {
	return base64.StdEncoding.EncodeToString(k.data[:])
}

func (k *K32) Hex() string {
	return hex.EncodeToString(k.data[:])
}

type K32Option func() (*K32, error)

func newK32(b []byte) (*K32, error) {
	bLen := len(b)
	if bLen != 32 {
		return nil, fmt.Errorf("len(b)=%d (32 expected)", bLen)
	}
	return &K32{[32]byte(b)}, nil
}

func WithBytes(b []byte) K32Option {
	return func() (*K32, error) {
		return newK32(b)
	}
}

func WithHex(s string) K32Option {
	return func() (*K32, error) {
		b, err := hex.DecodeString(s)
		if err != nil {
			return nil, err
		}
		return newK32(b)
	}
}

func WithB64(s string) K32Option {
	return func() (*K32, error) {
		b, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return nil, err
		}
		return newK32(b)
	}
}

func NewK32(opt K32Option) (*K32, error) {
	return opt()
}
