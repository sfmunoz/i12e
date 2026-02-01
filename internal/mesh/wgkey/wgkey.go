package wgkey

type WgKey struct {
	privKey *wgKeyPriv
	pubKey  *wgKeyPub
}

func (w *WgKey) GetPrivKey() *wgKeyPriv {
	return w.privKey
}

func (w *WgKey) GetPubKey() *wgKeyPub {
	return w.pubKey
}

func (w *WgKey) GetLocal() bool {
	return w.GetPrivKey() != nil
}
