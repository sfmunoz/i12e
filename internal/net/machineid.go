package net

import (
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
