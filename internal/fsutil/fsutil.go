package fsutil

import (
	"os"
	"path/filepath"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("app", "i12e").
	With("mod", "install")

func FileContentSet(file, bufIn string) error {
	file = filepath.Clean(file)
	dir, _ := filepath.Split(file)
	dir = filepath.Clean(dir)
	log.Info("FileContentSet()", "dir", dir, "file", file)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Error("FileContentSet(): os.MkdirAll() failed", "dir", dir, "err", err)
		return err
	}
	s, err := os.Stat(dir)
	if err != nil {
		log.Error("FileContentSet(): os.Stat() failed", "dir", dir, "err", err)
		return err
	}
	if !s.Mode().IsDir() {
		log.Error("FileContentSet(): dir doesn't exist but it should", "dir", dir)
		return err
	}
	s, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = os.WriteFile(file, []byte(bufIn), 0644)
		if err != nil {
			log.Error("FileContentSet(): file creation failed", "file", file, "err", err)
			return err
		}
		log.Info("FileContentSet(): file created", "file", file)
		return nil
	}
	if err != nil {
		log.Error("FileContentSet(): os.Stat() failed", "file", file, "err", err)
		return err
	}
	buf, err := os.ReadFile(file)
	if err != nil {
		log.Error("FileContentSet(): os.ReadFile() failed", "file", file, "err", err)
		return err
	}
	buf_s := string(buf)
	if buf_s == bufIn {
		log.Info("FileContentSet(): file is up-to-date", "file", file)
		return nil
	}
	err = os.WriteFile(file, []byte(bufIn), 0644)
	if err != nil {
		log.Error("FileContentSet(): file update failed", "file", file, "err", err)
		return err
	}
	log.Info("FileContentSet(): file updated", "file", file)
	return nil
}
