package net

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func getTuple(s string) ([2]int, error) {
	buf, err := hex.DecodeString(s)
	if err != nil {
		return [2]int{0, 0}, err
	}
	for i := 0; i < len(buf)-1; i++ {
		a, b := int(buf[i]), int(buf[i+1])
		if a == 0 && b == 0 || a == 255 && b == 255 {
			// '/16' network -> skip network and broadcast addresses
			continue
		}
		return [2]int{a, b}, nil
	}
	return [2]int{0, 0}, fmt.Errorf("getTuple(): cannot get proper value")
}

type MachineId struct {
	data  string
	tuple [2]int
}

func (m *MachineId) String() string {
	return m.data
}

func (m *MachineId) NodeId() int {
	return 256*m.tuple[0] + m.tuple[1]
}

func (m *MachineId) NodeName() string {
	return fmt.Sprintf("n-%03d-%03d", m.tuple[0], m.tuple[1])
}

func (m *MachineId) PathName() string {
	return fmt.Sprintf("%03d/%03d", m.tuple[0], m.tuple[1])
}

func (m *MachineId) IP() string {
	return fmt.Sprintf("10.56.%d.%d", m.tuple[0], m.tuple[1])
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
	tuple, err := getTuple(mid)
	if err != nil {
		return nil, err
	}
	return &MachineId{
		data:  mid,
		tuple: tuple,
	}, nil
}
