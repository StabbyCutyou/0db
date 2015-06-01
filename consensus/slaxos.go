package consensus

import (
	"hash/fnv"
	"os/exec"
)

type Slaxos struct {
	Cluster map[uint64]ServerEntry
}

type ServerEntry struct {
	Id      uint64
	Address string
}

func NewSlaxos() *Slaxos {
	s := &Slaxos{}
	s.initCluster()
	return s
}

func (s *Slaxos) initCluster() {
	cluster := make(map[uint64]ServerEntry)
	cluster[0] = ServerEntry{Id: 0, Address: "localhost"}
	s.Cluster = cluster
}

func (s *Slaxos) Write(key string, data string, ack bool) error {
	// First, hash the key
	keyHash := s.hashKey(key)
	// Now, modulus by the number of servers in the cluster to determine where the write goes
	nodeId := keyHash % uint64(len(s.Cluster))
	// Now we know the node, send the write command to that node
	return s.writeToNode(nodeId, key, data, ack)
}

func (s *Slaxos) Read(key string) (string, error) {
	// First, hash the key
	keyHash := s.hashKey(key)
	// Now, modulus by the number of servers in the cluster to determine where the write goes
	nodeId := keyHash % uint64(len(s.Cluster))
	// Now we know the node, send the read command to that node and return the result
	return s.readFromNode(nodeId, key)
}

func (s *Slaxos) writeToNode(node uint64, key string, data string, ack bool) error {
	// TODO this functionality should live in a dispacher...
	// Get the address of the node
	address := s.Cluster[node].Address
	if address == "localhost" {
		// The local node owns it
		cmd := exec.Command("echo", data, "> /dev/null")
		if ack == true {
			return cmd.Run()
		} else {
			cmd.Start()
			return nil
		}
	} else {
		// Network send the write
		// TODO support clustering
		return nil
	}
}

func (s *Slaxos) readFromNode(node uint64, key string) (string, error) {
	// TODO this functionality should live in a dispacher...
	// Get the address of the node
	address := s.Cluster[node].Address
	if address == "localhost" {
		// The local node owns it
		cmd := exec.Command("cat", "/dev/null")
		bytes, err := cmd.Output()
		return string(bytes[:]), err
	} else {
		// Network send the write
		// TODO support clustering
		return "", nil
	}
}

func (s *Slaxos) hashKey(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return h.Sum64()
}
