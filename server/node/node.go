package node

import (
	"github.com/StabbyCutyou/0db/server/config"
	"github.com/StabbyCutyou/0db/server/consensus"
)

type ZeroDBNode struct {
	consensus *consensus.Slaxos
}

func New(cfg *config.Config) *ZeroDBNode {
	return &ZeroDBNode{consensus: consensus.NewSlaxos(cfg.Membership)}
}

func (z *ZeroDBNode) Write(key string, data string, ack bool) error {
	return z.consensus.Write(key, data, ack)
}

func (z *ZeroDBNode) Read(key string) (string, error) {
	return z.consensus.Read(key)
}
