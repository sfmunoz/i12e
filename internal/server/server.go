package server

import (
	"math/rand/v2"
	"time"

	"github.com/sfmunoz/i12e/internal/config"
	"github.com/sfmunoz/i12e/internal/mesh"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Server struct {
	cfg                               *config.ServerConfig
	steps, bInit, bStep, jInit, jStep int64
}

func newServer(cfg *config.ServerConfig) *Server {
	steps := int64(100)
	bInit := 16 * time.Second
	jInit := 4 * time.Second
	return &Server{
		cfg:   cfg,
		steps: steps,
		bInit: int64(bInit),
		bStep: int64(cfg.Server.SlumberBase-bInit) / steps,
		jInit: int64(jInit),
		jStep: int64(cfg.Server.SlumberJitter-jInit) / steps,
	}
}

func (s *Server) runOne() error {
	if err := rclonePull(); err != nil {
		log.Error("rclonePull() failed", "err", err)
		return err
	}
	if err := artifactPull(); err != nil {
		log.Error("artifactPull() failed", "err", err)
		return err
	}
	if err := artifactTune(); err != nil {
		log.Error("artifactTune() failed", "err", err)
		return err
	}
	if err := mesh.Run(s.cfg); err != nil {
		log.Error("mesh.Run() failed", "err", err)
		return err
	}
	if err := k3sInstall(s.cfg); err != nil {
		log.Error("k3sInstall() failed", "err", err)
		return err
	}
	if err := pluginsRun(); err != nil {
		log.Error("pluginsRun() failed", "err", err)
		return err
	}
	return nil
}

func (s *Server) slumber(i int64, err error) time.Duration {
	x := min(i, s.steps)
	if err != nil {
		x = 0
	}
	return time.Duration(s.bInit + x*s.bStep + rand.Int64N(s.jInit+x*s.jStep))
}

func (s *Server) run() error {
	var i int64 = 0
	for {
		err := s.runOne()
		slumber := s.slumber(i, err)
		log.Info("i12e sleeping...", "slumber", slumber)
		time.Sleep(slumber)
		i = min(i+1, s.steps)
	}
}

func Run(cfg *config.ServerConfig) error {
	return newServer(cfg).run()
}
