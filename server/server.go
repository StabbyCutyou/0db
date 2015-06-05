package server

import (
	"github.com/StabbyCutyou/0db/config"
	"github.com/StabbyCutyou/0db/consensus"
)

type ZeroDBServer struct {
	consensus *consensus.Slaxos
}

func New(cfg *config.Config) *ZeroDBServer {
	return &ZeroDBServer{consensus: consensus.NewSlaxos(cfg.Membership)}
}

func (z *ZeroDBServer) Write(key string, data string, ack bool) error {
	return z.consensus.Write(key, data, ack)
}

func (z *ZeroDBServer) Read(key string) (string, error) {
	return z.consensus.Read(key)
}
