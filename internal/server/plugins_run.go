package server

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/sfmunoz/i12e/internal/cmdutil"
)

const pluginsFolder = "/opt/libexec/i12e/plugins"

func pluginsRun() error {
	if _, err := os.Stat(pluginsFolder); os.IsNotExist(err) {
		return nil
	}
	matches, err := filepath.Glob(filepath.Join(pluginsFolder, "*.sh"))
	if err != nil {
		return err
	}
	sort.Strings(matches)
	for _, plugin := range matches {
		if _, err := os.Stat(plugin); err != nil {
			if os.IsNotExist(err) {
				log.Warn("plugin vanished: it surely was purged", "plugin", plugin)
				continue
			}
			log.Error("os.Stat() failed", "err", err, "plugin", plugin)
			return err
		}
		log.Info("running", "plugin", plugin)
		cmd := exec.Command("/bin/bash", plugin)
		if err := cmdutil.RunCmd(cmd); err != nil {
			log.Error("cmdutil.RunCmd() failed", "err", err, "plugin", plugin)
			return err
		}
	}
	return nil
}
