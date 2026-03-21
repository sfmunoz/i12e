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
		log.Info("pluginsRun(): running plugin", "plugin", plugin)
		cmd := exec.Command("/bin/bash", plugin)
		if err := cmdutil.RunCmd(cmd); err != nil {
			log.Error("plugin failed", "err", err, "plugin", plugin)
			return err
		}
	}
	return nil
}
