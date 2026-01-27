package net

import (
	_ "embed"
	"time"

	"github.com/sfmunoz/i12e/internal/netutil"
	"github.com/sfmunoz/logit"
)

//go:embed static/net.sh
var netSh []byte

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Net struct {
	Tnow           time.Time
	MachineId      *MachineId
	WgInterface    string
	WgEndpointPort uint16
	WgEndpointIp   string
	WgPrivKey      *WgKey
	WgPubKey       *WgKey
}

func newNet() (*Net, error) {
	ii, err := netutil.IfaceIP()
	if err != nil {
		log.Error("newNet(): 'netutil.IfaceIP()' failed", "err", err)
		return nil, err
	}
	if ii == nil {
		log.Error("newNet(): 'netutil.IfaceIP()' returned 'nil'")
		return nil, err
	}
	machineId, err := getMachineId()
	if err != nil {
		log.Error("newNet(): 'getMachineId()' failed", "err", err)
		return nil, err
	}
	wgPrivKey, err := getWgKey(true)
	if err != nil {
		log.Error("newNet(): 'getWgKey(true)' failed", "err", err)
		return nil, err
	}
	wgPubKey, err := getWgKey(false)
	if err != nil {
		log.Error("newNet(): 'getWgKey(false)' failed", "err", err)
		return nil, err
	}
	return &Net{
		Tnow:           time.Now().UTC(),
		MachineId:      machineId,
		WgInterface:    "wgi",
		WgEndpointIp:   ii.IP,
		WgEndpointPort: 51830, // default '51820'
		WgPrivKey:      wgPrivKey,
		WgPubKey:       wgPubKey,
	}, nil
}

func (n *Net) run() error {
	log.Info("Net.run()", "Tnow", n.Tnow)
	log.Info("Net.run()", "MachineId", n.MachineId, "NodeId", n.MachineId.NodeId(), "NodeName", n.MachineId.NodeName(), "PathName", n.MachineId.PathName(), "IP", n.MachineId.IP())
	log.Info("Net.run()", "WgInterface", n.WgInterface)
	log.Info("Net.run()", "WgEndpointIp", n.WgEndpointIp)
	log.Info("Net.run()", "WgEndpointPort", n.WgEndpointPort)
	log.Info("Net.run()", "WgPrivKey", n.WgPrivKey, "WgPrivKeyLen", n.WgPrivKey.Len())
	log.Info("Net.run()", "WgPubKey", n.WgPubKey, "WgPubKeyHex", n.WgPubKey.Hex(), "WgPubKeyLen", n.WgPubKey.Len())
	return nil
}

func Run() error {
	n, err := newNet()
	if err != nil {
		return err
	}
	return n.run()
}
