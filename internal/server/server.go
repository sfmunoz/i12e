package server

import (
	"time"

	"github.com/sfmunoz/i12e/internal/mesh"
	"github.com/sfmunoz/i12e/internal/pull"
	"github.com/sfmunoz/logit"
)

const remBase = "rem:mesh" // FIXME unhardcode this

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Server struct {
	slumber time.Duration
}

func newServer() *Server {
	return &Server{
		slumber: 3 * time.Second,
	}
}

func (s *Server) run() error {
	firstTime := true
	for {
		log.Info("i12e running...")
		if err := pull.Run(); err != nil {
			log.Error("pull.Run() failed", "err", err)
		}
		if firstTime {
			firstTime = false
			if err := mesh.Run(remBase); err != nil {
				log.Error("net.Run() failed", "err", err)
			}
		}
		log.Info("i12e sleeping...", "slumber", s.slumber)
		time.Sleep(s.slumber)
	}
}

func Run() error {
	return newServer().run()
}
