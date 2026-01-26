package net

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sfmunoz/i12e/internal/netutil"
	"github.com/sfmunoz/logit"
)

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
	WgIface   string
	WgPort    uint16
	Iface     string
	WgFname   string
	Tnow      time.Time
	MachineId string
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
	return &Net{
		WgIface:   "wg0",
		WgPort:    51820,
		Iface:     ii.Iface,
		WgFname:   "/etc/i12e/wg-privkey",
		Tnow:      time.Now().UTC(),
		MachineId: machineId,
	}, nil
}

func (n *Net) run() error {
	log.Info("Net.run()", "WgIface", n.WgIface)
	log.Info("Net.run()", "WgPort", n.WgPort)
	log.Info("Net.run()", "Iface", n.Iface)
	log.Info("Net.run()", "WgFname", n.WgFname)
	log.Info("Net.run()", "Tnow", n.Tnow)
	log.Info("Net.run()", "MachineId", n.MachineId)
	return nil
}

func Run() error {
	n, err := newNet()
	if err != nil {
		return err
	}
	return n.run()
}
