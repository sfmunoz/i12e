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

func writeRandomEtcHostname(meshNet *netip.Prefix) error {
	nodeId, err := getRandomNodeId(meshNet)
	if err != nil {
		return err
	}
	nodeName := getNodeNameFromNodeId(nodeId)
	if err := os.WriteFile(etcHostname, fmt.Appendf(make([]byte, 0), "%s\n", nodeName), 0644); err != nil {
		return err
	}
	return nil
}

func readEtcHostname(meshNet *netip.Prefix) (uint32, error) {
	buf, err := os.ReadFile(etcHostname)
	if err != nil {
		return 0, err
	}
	nodeId, err := getNodeIdFromNodeName(meshNet, strings.TrimSpace(string(buf)))
	if err != nil {
		return 0, err
	}
	return nodeId, nil
}
