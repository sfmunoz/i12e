package wgkey

type WgKeyPub struct {
	k32 *K32
}

func (w *WgKeyPub) String() string {
	return w.K32().B64()
}

func (w *WgKeyPub) K32() *K32 {
	return w.k32
}

func NewWgKeyPub(k32 *K32) *WgKeyPub {
	return &WgKeyPub{k32}
}
