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

func (m *MachineId) tuple() ([2]int, error) {
	buf, err := hex.DecodeString(m.data)
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
	return [2]int{0, 0}, fmt.Errorf("MachineId.tuple(): cannot get properly value")
}

func (m *MachineId) NodeId() (int, error) {
	t, err := m.tuple()
	if err != nil {
		return 0, err
	}
	return 256*t[0] + t[1], nil
}

func (m *MachineId) NodeName() (string, error) {
	t, err := m.tuple()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("n-%03d-%03d", t[0], t[1]), nil
}

func (m *MachineId) IP() (string, error) {
	t, err := m.tuple()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("10.56.%d.%d", t[0], t[1]), nil
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
