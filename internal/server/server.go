package server

import (
	"time"

	"github.com/sfmunoz/i12e/internal/net"
	"github.com/sfmunoz/i12e/internal/pull"
	"github.com/sfmunoz/logit"
)

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
	for {
		log.Info("i12e running...")
		pull.Pull()
		net.Run()
		log.Info("i12e sleeping...", "slumber", s.slumber)
		time.Sleep(s.slumber)
	}
}

func Run() error {
	return newServer().run()
}
