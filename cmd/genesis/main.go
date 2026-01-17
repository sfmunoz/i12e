package main

import (
	"time"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("mod", "genesis").
	With("pkg", "main")

func main() {
	slumber := time.Second
	for {
		log.Info("genesis running...")
		time.Sleep(slumber)
	}
}
