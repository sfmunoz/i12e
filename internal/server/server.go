package server

import (
	"math/rand/v2"
	"time"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/mesh"
	"github.com/sfmunoz/i12e/internal/pull"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Server struct {
	cfg *config.ServerConfig
}

func newServer(cfg *config.ServerConfig) *Server {
	return &Server{cfg}
}

func (s *Server) run() error {
	i := 0
	for {
		log.Info("i12e running...")
		if err := mesh.Run(s.cfg); err != nil {
			log.Error("mesh.Run() failed", "err", err)
		}
		if err := pull.Run(s.cfg); err != nil {
			log.Error("pull.Run() failed", "err", err)
		}
		slumber := s.cfg.Server.SlumberBase + time.Duration(rand.Int64N(int64(s.cfg.Server.SlumberJitter)))
		if i < 4 {
			i += 1
			slumber = 16*time.Second + time.Duration(rand.Int64N(int64(4*time.Second)))
		}
		log.Info("i12e sleeping...", "slumber", slumber)
		time.Sleep(slumber)
	}
}

func Run(cfg *config.ServerConfig) error {
	return newServer(cfg).run()
}
