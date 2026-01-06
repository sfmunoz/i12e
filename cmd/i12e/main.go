package main

import (
	"os"
	"time"

	"github.com/sfmunoz/i12e/internal/install"
	"github.com/sfmunoz/i12e/internal/pull"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "main")

func main() {
	if os.Getenv("I12E_INSTALL") == "1" {
		install.Install()
		return
	}
	slumber := 3 * time.Second
	for {
		log.Info("i12e running...", "slumber", slumber)
		pull.Pull()
		time.Sleep(slumber)
	}
}
