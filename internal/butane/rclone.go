package butane

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/sfmunoz/i12e/internal/cmdutil"
	"github.com/sfmunoz/i12e/internal/config"
)

type RcloneBlock map[string]string
type RcloneConf map[string]RcloneBlock

func RcloneDump(cfg *config.Config) error {
	bo, be, err := cmdutil.RunSimple(exec.Command("rclone", "config", "dump"))
	if err != nil {
		return fmt.Errorf("'sops decrypt' failed: err=%s; buf_err=%s; prod=%t", err, be, cfg.Prod)
	}
	var rcloneConf RcloneConf
	if err := json.Unmarshal(bo.Bytes(), &rcloneConf); err != nil {
		return err
	}
	rem := cfg.RcloneRemote
	var out bytes.Buffer
	for k1, v1 := range rcloneConf {
		if k1 != rem {
			continue
		}
		fmt.Fprintf(&out, "[%s]\n", rem)
		for k2, v2 := range v1 {
			fmt.Fprintf(&out, "%s = %s\n", k2, v2)
		}
	}
	log.Info("RcloneDump", "out", out.String())
	return nil
}
