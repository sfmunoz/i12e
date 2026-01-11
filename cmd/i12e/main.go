package main

import (
	"time"

	"github.com/sfmunoz/i12e/internal/pull"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "i12e").
	With("pkg", "main")

func main() {
	slumber := 3 * time.Second
	for {
		log.Info("i12e running...")
		pull.Pull()
		log.Info("i12e sleeping...", "slumber", slumber)
		time.Sleep(slumber)
	}
}
