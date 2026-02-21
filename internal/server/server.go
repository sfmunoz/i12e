package server

import (
	"math/rand/v2"
	"time"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/mesh"
	"github.com/sfmunoz/i12e/internal/pull"
	"github.com/sfmunoz/logit"
)

const remBase = "rem:mesh" // TODO: unhardcode

const serverSlumberBase = 8 * time.Second   // TODO: unhardcode
const serverSlumberJitter = 4 * time.Second // TODO: unhardcode

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Server struct {
	cfg           *config.Config
	slumberBase   time.Duration
	slumberJitter time.Duration
}

func newServer(cfg *config.Config, slumberBase time.Duration, slumberJitter time.Duration) *Server {
	return &Server{cfg, slumberBase, slumberJitter}
}

func (s *Server) run() error {
	for {
		log.Info("i12e running...")
		if err := mesh.Run(s.cfg, remBase); err != nil {
			log.Error("mesh.Run() failed", "err", err)
		}
		if err := pull.Run(s.cfg); err != nil {
			log.Error("pull.Run() failed", "err", err)
		}
		slumber := s.slumberBase + time.Duration(rand.Int64N(int64(s.slumberJitter)))
		log.Info("i12e sleeping...", "slumber", slumber)
		time.Sleep(slumber)
	}
}

func Run(cfg *config.Config) error {
	return newServer(cfg, serverSlumberBase, serverSlumberJitter).run()
}
