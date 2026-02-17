package server

import (
	"math/rand/v2"
	"net/netip"
	"time"

	"github.com/sfmunoz/i12e/internal/mesh"
	"github.com/sfmunoz/i12e/internal/pull"
	"github.com/sfmunoz/logit"
)

// from /12 (20 bits for host) to /29 (3 bits for host)
var meshNet = netip.MustParsePrefix("10.119.0.0/28") // TODO: unhardcode

const remBase = "rem:mesh" // TODO: unhardcode

const serverSlumberBase = 8 * time.Second   // TODO: unhardcode
const serverSlumberJitter = 4 * time.Second // TODO: unhardcode

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Server struct {
	slumberBase   time.Duration
	slumberJitter time.Duration
}

func newServer(slumberBase time.Duration, slumberJitter time.Duration) *Server {
	return &Server{slumberBase, slumberJitter}
}

func (s *Server) run() error {
	for {
		log.Info("i12e running...")
		if err := mesh.Run(&meshNet, remBase); err != nil {
			log.Error("net.Run() failed", "err", err)
		}
		if err := pull.Run(&meshNet); err != nil {
			log.Error("pull.Run() failed", "err", err)
		}
		slumber := s.slumberBase + time.Duration(rand.Int64N(int64(s.slumberJitter)))
		log.Info("i12e sleeping...", "slumber", slumber)
		time.Sleep(slumber)
	}
}

func Run() error {
	return newServer(serverSlumberBase, serverSlumberJitter).run()
}
