package net

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/netutil"
	"github.com/sfmunoz/logit"
)

//go:embed static/net.sh
var netSh []byte

var log = logit.Logit().WithLevel(logit.LevelInfo)

func getMachineId() (string, error) {
	buf, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "", err
	}
	mid := strings.TrimSpace(string(buf))
	mid_len := len(mid)
	if mid_len != 32 {
		return "", fmt.Errorf("getMachineId(): len(%s)=%d (32 expected)", mid, mid_len)
	}
	return mid, nil
}

type Net struct {
	Tnow           time.Time
	MachineId      string
	WgInterface    string
	WgEndpointPort uint16
	WgEndpointIp   string
	WgPrivKey      WgKey
	WgPubKey       WgKey
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
	wgPrivKey, err := getWgKey(WgCmdPrivKey)
	if err != nil {
		log.Error("newNet(): 'getWgKey(WgcmdPrivKey)' failed", "err", err)
		return nil, err
	}
	wgPubKey, err := getWgKey(WgCmdPubKey)
	if err != nil {
		log.Error("newNet(): 'getWgKey(WgCmdPubKey)' failed", "err", err)
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
	log.Info("Net.run()", "MachineId", n.MachineId)
	log.Info("Net.run()", "WgInterface", n.WgInterface)
	log.Info("Net.run()", "WgEndpointIp", n.WgEndpointIp)
	log.Info("Net.run()", "WgEndpointPort", n.WgEndpointPort)
	log.Info("Net.run()", "WgPrivKey", strings.Repeat("*", len(n.WgPrivKey)), "WgPrivKeyLen", len(n.WgPrivKey))
	log.Info("Net.run()", "WgPubKey", n.WgPubKey, "WgPubKeyLen", len(n.WgPubKey))
	return nil
}

func Run() error {
	n, err := newNet()
	if err != nil {
		return err
	}
	return n.run()
}
