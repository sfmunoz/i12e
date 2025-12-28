package main

import (
	"os"
	"time"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().
	WithLevel(logit.LevelInfo).
	With("app", "i12e")

func main() {
	if os.Getenv("I12E_INSTALL") == "1" {
		install()
		return
	}
	slumber := 3 * time.Second
	for {
		log.Info("i12e running...", "slumber", slumber)
		time.Sleep(slumber)
	}
}
