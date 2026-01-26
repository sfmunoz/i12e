package net

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/cmdutil"
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

func getWgKey(c string) (string, error) {
	cmd := exec.Command("/bin/sh", "-s", "-", c, "/etc/i12e/wg-priv-key")
	cmd.Stdin = bytes.NewBuffer(netSh)
	bo, be, err := cmdutil.RunSimple(cmd)
	if err != nil {
		return "", fmt.Errorf("'net.sh' failed': %s (stdout=%s, stderr=%s)", err, bo.String(), be.String())
	}
	return bo.String(), nil
}

func getWgPrivKey() (string, error) {
	return getWgKey("priv-key")
}

func getWgPubKey() (string, error) {
	return getWgKey("pub-key")
}

type Net struct {
	Tnow           time.Time
	MachineId      string
	WgInterface    string
	WgEndpointPort uint16
	WgEndpointIp   string
	WgPrivKey      string
	WgPubKey       string
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
	wgPrivKey, err := getWgPrivKey()
	if err != nil {
		log.Error("newNet(): 'getWgPrivKey()' failed", "err", err)
		return nil, err
	}
	wgPubKey, err := getWgPubKey()
	if err != nil {
		log.Error("newNet(): 'getWgPubKey(false)' failed", "err", err)
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
	log.Info("Net.run()", "WgPrivKey", strings.Repeat("*", len(n.WgPrivKey)))
	log.Info("Net.run()", "WgPubKey", n.WgPubKey)
	return nil
}

func Run() error {
	n, err := newNet()
	if err != nil {
		return err
	}
	return n.run()
}
