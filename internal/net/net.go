package net

import (
	"time"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

const WgPrivKeyFname = "/etc/i12e/wg-priv-key" // FIXME unhardcode this
const remMeshBase = "rem:mesh"                 // FIXME unhardcode this

type Net struct {
	Tnow           time.Time
	nodeLocal      *node
	WgInterface    string
	WgEndpointPort uint16
	WgEndpointIp   string
	WgPrivKey      *WgKey
	WgPubKey       *WgKey
}

func newNet() (*Net, error) {
	ii, err := IfaceIP()
	if err != nil {
		log.Error("newNet(): 'IfaceIP()' failed", "err", err)
		return nil, err
	}
	if ii == nil {
		log.Error("newNet(): 'IfaceIP()' returned 'nil'")
		return nil, err
	}
	nodeLocal, err := getNodeLocal()
	if err != nil {
		log.Error("newNet(): 'getNodeLocal()' failed", "err", err)
		return nil, err
	}
	if err := nodeLocal.SetHostname(); err != nil {
		return nil, err
	}
	wgPrivKey, err := getWgPrivKey(WgPrivKeyFname)
	if err != nil {
		log.Error("newNet(): 'getWgPrivKey()' failed", "err", err)
		return nil, err
	}
	wgPubKey, err := getWgPubKey(WgPrivKeyFname)
	if err != nil {
		log.Error("newNet(): 'getWgPubKey()' failed", "err", err)
		return nil, err
	}
	return &Net{
		Tnow:           time.Now().UTC(),
		nodeLocal:      nodeLocal,
		WgInterface:    "wgi",
		WgEndpointIp:   ii.IP,
		WgEndpointPort: 51830, // default '51820'
		WgPrivKey:      wgPrivKey,
		WgPubKey:       wgPubKey,
	}, nil
}

func (n *Net) run() error {
	log.Info("Net.run()", "nodeLocal", n.nodeLocal)
	log.Info("Net.run()", "Tnow", n.Tnow)
	log.Info("Net.run()", "WgInterface", n.WgInterface)
	log.Info("Net.run()", "WgEndpointIp", n.WgEndpointIp)
	log.Info("Net.run()", "WgEndpointPort", n.WgEndpointPort)
	log.Info("Net.run()", "WgPrivKey", n.WgPrivKey, "WgPrivKeyLen", n.WgPrivKey.Len())
	log.Info("Net.run()", "WgPubKey", n.WgPubKey, "WgPubKeyHex", n.WgPubKey.Hex(), "WgPubKeyLen", n.WgPubKey.Len())
	mesh, err := getMesh(remMeshBase)
	if err != nil {
		return err
	}
	if err := mesh.NodePush(n.nodeLocal.NodePath(), n.Tnow, n.WgPubKey, n.WgEndpointIp, n.WgEndpointPort); err != nil {
		return err
	}
	if err := mesh.NodeConfig(n.WgInterface, n.nodeLocal.NodeIP(), n.WgEndpointPort, WgPrivKeyFname, n.nodeLocal.NodePath()); err != nil {
		return err
	}
	log.Info("Net.run()", "mesh", mesh)
	return nil
}

func Run() error {
	n, err := newNet()
	if err != nil {
		return err
	}
	return n.run()
}
