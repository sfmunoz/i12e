package node

import (
	"errors"
	"fmt"
	"net/netip"
	"os"
	"strings"
)

const etcHostname = "/etc/hostname"

func deleteEtcHostname() error {
	_, err := os.Stat(etcHostname)
	if err == nil {
		return os.Remove(etcHostname)
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func writeRandomEtcHostname(p *netip.Prefix) error {
	nodeId, err := getRandomNodeId(p)
	if err != nil {
		return err
	}
	nodeName := getNodeNameFromNodeId(nodeId)
	if err := os.WriteFile(etcHostname, fmt.Appendf(make([]byte, 0), "%s\n", nodeName), 0644); err != nil {
		return err
	}
	return nil
}

func readEtcHostname(p *netip.Prefix) (uint32, error) {
	buf, err := os.ReadFile(etcHostname)
	if err != nil {
		return 0, err
	}
	nodeId, err := getNodeIdFromNodeName(p, strings.TrimSpace(string(buf)))
	if err != nil {
		return 0, err
	}
	return nodeId, nil
}
