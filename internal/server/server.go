package server

import (
	"fmt"
	"time"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/pull"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Server struct {
	cfg *config.Config
}

func newServer(cfg *config.Config) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("newServer(): undefined config")
	}
	return &Server{
		cfg: cfg,
	}, nil
}

func (s *Server) run() error {
	slumber := 3 * time.Second
	for {
		log.Info("i12e running...")
		pull.Pull()
		log.Info("i12e sleeping...", "slumber", slumber)
		time.Sleep(slumber)
	}
}

func Run(cfg *config.Config) error {
	a, err := newServer(cfg)
	if err != nil {
		return err
	}
	return a.run()
}
