package node

import (
	"fmt"
	"math/rand/v2"
	"net/netip"
	"os"
	"strconv"
	"strings"

	"github.com/sfmunoz/i12e/internal/mesh/ifaceip"
	"github.com/sfmunoz/i12e/internal/mesh/wgkey"
)

func validNodeId(nodeId int) error {
	if nodeId < nodeIdMin {
		return fmt.Errorf("'%s' node-id is '%d' (min=%d)", nodeIdFname, nodeId, nodeIdMin)
	}
	if nodeId > nodeIdMax {
		return fmt.Errorf("'%s' node-id is '%d' (max=%d)", nodeIdFname, nodeId, nodeIdMax)
	}
	return nil
}

func loadNode() (*Node, error) {
	buf, err := os.ReadFile(nodeIdFname)
	if err != nil {
		return nil, err
	}
	nodeId, err := strconv.Atoi(strings.TrimSpace(string(buf)))
	if err != nil {
		return nil, err
	}
	if err := validNodeId(nodeId); err != nil {
		return nil, err
	}
	wgKey, err := wgkey.GetWgKeyLocal(nodePrivKeyFname)
	if err != nil {
		return nil, err
	}
	ii, err := ifaceip.IfaceIP()
	if err != nil {
		return nil, err
	}
	if ii == nil {
		return nil, fmt.Errorf("IfaceIP() returned empty value")
	}
	wgEndpoint := netip.AddrPortFrom(*ii.IP, nodeEndpointPort)
	return &Node{
		id:         uint16(nodeId),
		wgKey:      wgKey,
		wgEndpoint: &wgEndpoint,
	}, nil
}

func writeNode() error {
	x := rand.Int32N(nodeIdMax-nodeIdMin+1) + nodeIdMin
	if err := os.WriteFile(nodeIdFname, fmt.Appendf(make([]byte, 0), "%d\n", x), 0600); err != nil {
		return err
	}
	return nil
}

func GetNodeLocal() (*Node, error) {
	node, err := loadNode()
	if err == nil {
		return node, nil
	}
	if err := writeNode(); err != nil {
		return nil, err
	}
	return loadNode()
}

func getNodeIdFromPath(path string) (uint16, error) {
	parts := strings.Split(path, "_")
	parts_len := len(parts)
	if parts_len != 2 {
		return 0, fmt.Errorf("path='%s' -> len(parts=%s)=%d (2 expected)", path, parts, parts_len)
	}
	ids := [2]uint16{0, 0}
	for i, part := range parts {
		plen := len(part)
		if plen < 1 {
			return 0, fmt.Errorf("len(parts[%d])=%d (>0 expected)", i, plen)
		}
		pint, err := strconv.Atoi(part)
		if err != nil {
			return 0, err
		}
		if pint < 0 {
			return 0, fmt.Errorf("wrong path: '%s[%d]' = '%s' < 0", path, i, part)
		}
		if pint > 255 {
			return 0, fmt.Errorf("wrong path: '%s[%d]' = '%s' > 255", path, i, part)
		}
		ids[i] = uint16(pint)
	}
	nodeId := 256*int(ids[0]) + int(ids[1])
	if err := validNodeId(nodeId); err != nil {
		return 0, err
	}
	return uint16(nodeId), nil
}

func GetNodeRemote(entry string) (*Node, error) {
	arr := nodeRegex.FindStringSubmatch(entry)
	if arr == nil {
		return nil, fmt.Errorf("'nodeRegex.FindStringSubmatch(%s)' returned nil", entry)
	}
	nodeId, err := getNodeIdFromPath(arr[1])
	if err != nil {
		return nil, err
	}
	wgKey, err := wgkey.GetWgKeyRemote(arr[3])
	if err != nil {
		return nil, fmt.Errorf("'getWgKeyRemote(%s)' failed: %s", arr[3], err)
	}
	addr, err := netip.ParseAddr(arr[4])
	if err != nil {
		return nil, fmt.Errorf("'netip.ParseAddr(%s)' failed: %s", arr[4], err)
	}
	wgEndpointPort, err := strconv.Atoi(arr[5])
	if err != nil {
		return nil, fmt.Errorf("'strconv.Atoi(%s)' failed: %s", arr[5], err)
	}
	wgEndpoint := netip.AddrPortFrom(addr, uint16(wgEndpointPort))
	return &Node{
		id:         nodeId,
		wgKey:      wgKey,
		wgEndpoint: &wgEndpoint,
	}, nil
}
