package net

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

type MachineId struct {
	data string
}

func (m *MachineId) String() string {
	return m.data
}

func (m *MachineId) NodeId() (int, error) {
	buf, err := hex.DecodeString(m.data)
	if err != nil {
		return 0, err
	}
	for i := 0; i < len(buf)-1; i++ {
		x := 256*int(buf[i]) + int(buf[i+1])
		if x > 0 && x < 0xffff {
			// '/16' network -> skip network and broadcast addresses
			return x, nil
		}
	}
	return 0, fmt.Errorf("MachineId.NodeId(): cannot get proper id")
}

func getMachineId() (*MachineId, error) {
	buf, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return nil, err
	}
	mid := strings.TrimSpace(string(buf))
	mid_len := len(mid)
	if mid_len != 32 {
		return nil, fmt.Errorf("getMachineId(): len(%s)=%d (32 expected)", mid, mid_len)
	}
	return &MachineId{
		data: mid,
	}, nil
}
