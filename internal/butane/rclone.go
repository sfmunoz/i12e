package butane

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

type rcloneBlock map[string]string
type rcloneConf map[string]rcloneBlock

func remotes(rc *rcloneConf, rem string) ([]string, error) {
	ret := make([]string, 0)
	for k1, v1 := range *rc {
		if k1 != rem {
			continue
		}
		ret = append(ret, rem)
		for k2, v2 := range v1 {
			if k2 != "remote" {
				continue
			}
			parts := strings.Split(v2, ":")
			if len(parts) < 2 {
				continue
			}
			rem2, err := remotes(rc, parts[0])
			if err != nil {
				return nil, err
			}
			ret = append(ret, rem2...)
		}
	}
	if len(ret) < 1 {
		return nil, fmt.Errorf("cannot find '%s' remote", rem)
	}
	return ret, nil
}

func rcloneConfig(rem string) (*bytes.Buffer, error) {
	bo, be, err := cmdutil.RunSimple(exec.Command("rclone", "config", "dump"))
	if err != nil {
		return nil, fmt.Errorf("'rclone config dump' failed: err=%s; buf_err=%s", err, be)
	}
	var rcloneConf rcloneConf
	if err := json.Unmarshal(bo.Bytes(), &rcloneConf); err != nil {
		return nil, err
	}
	rems, err := remotes(&rcloneConf, rem)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	var firstLine bool = true
	for _, r := range rems {
		for k1, v1 := range rcloneConf {
			if k1 != r {
				continue
			}
			if firstLine {
				firstLine = false
			} else {
				fmt.Fprintln(&out)
			}
			if r == rem {
				fmt.Fprint(&out, "[rem]\n")
			} else {
				fmt.Fprintf(&out, "[%s]\n", r)
			}
			for k2, v2 := range v1 {
				fmt.Fprintf(&out, "%s = %s\n", k2, v2)
			}
		}
	}
	return &out, nil
}
