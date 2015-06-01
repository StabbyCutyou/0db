package server

import "github.com/StabbyCutyou/0db/consensus"

type ZeroDBServer struct {
	consensus *consensus.Slaxos
}

func New() *ZeroDBServer {
	return &ZeroDBServer{consensus: consensus.NewSlaxos()}
}

func (z *ZeroDBServer) Write(key string, data string, ack bool) error {
	return z.consensus.Write(key, data, ack)
}

func (z *ZeroDBServer) Read(key string) (string, error) {
	return z.consensus.Read(key)
}
