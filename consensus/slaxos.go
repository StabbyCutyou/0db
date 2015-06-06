package consensus

import (
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/config"
	"github.com/StabbyCutyou/0db/message"
	"hash/fnv"
	"os/exec"
)

type Slaxos struct {
	receivePort  int
	dispatchPort int
	membership   *Membership
	admin        *Admin
}

type ServerEntry struct {
	Id      uint64
	Address string
}

func NewSlaxos(cfg config.MembershipConfig) *Slaxos {
	return &Slaxos{
		membership:   NewMembershipListener(cfg.MemberPort),
		admin:        NewAdminListener(cfg.AdminPort),
		receivePort:  cfg.ReceivePort,
		dispatchPort: cfg.DispatchPort,
	}
}

func (s *Slaxos) ProcessAdminCommands() {
	for {
		data := s.admin.CommandQueue.Dequeue()
		if data != nil {
			msg := data.(message.AdminMessage)
			switch *msg.Command {
			case "JoinCluster":
				s.JoinCluster(*msg.Message)
			case "LeaveCluster":
				s.LeaveCluster()
			}
		}
	}
}

func (s *Slaxos) JoinCluster(address string) {
	logrus.Infof("About to join to Cluster at address %s", address)
	nodesFound, err := s.membership.JoinCluster(address)
	if err != nil {
		logrus.Errorf("Unable to join Cluster at address %s", address)
	} else {
		logrus.Infof("Joined cluster at %s, found %d other members", address, nodesFound)
	}
}

func (s *Slaxos) LeaveCluster() {
	logrus.Info("Leaving Cluster")
	err := s.membership
	if err != nil {
		logrus.Error("Error leaving Cluster")
		logrus.Error(err)
	}
}

func (s *Slaxos) Write(key string, data string, ack bool) error {
	// First, hash the key
	keyHash := s.hashKey(key)
	// Now, modulus by the number of servers in the cluster to determine where the write goes
	nodeId := s.calculateNodeIndex(keyHash)
	// Now we know the node, send the write command to that node
	return s.writeToNode(nodeId, key, data, ack)
}

func (s *Slaxos) Read(key string) (string, error) {
	// First, hash the key
	keyHash := s.hashKey(key)
	// Now, modulus by the number of servers in the cluster to determine where the write goes
	nodeId := s.calculateNodeIndex(keyHash)
	// Now we know the node, send the read command to that node and return the result
	return s.readFromNode(nodeId, key)
}

func (s *Slaxos) calculateNodeIndex(keyHash uint64) uint64 {
	return keyHash % uint64(s.membership.ClusterSize())
}

func (s *Slaxos) writeToNode(nodeId uint64, key string, data string, ack bool) error {
	// TODO this functionality should live in a dispacher...
	// Get the address of the node
	address := s.membership.MemberAt(nodeId)
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

func (s *Slaxos) readFromNode(nodeId uint64, key string) (string, error) {
	// TODO this functionality should live in a dispacher...
	// Get the address of the node
	address := s.membership.MemberAt(nodeId)
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
