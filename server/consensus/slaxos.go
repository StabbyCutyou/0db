package consensus

import (
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/message"
	"github.com/StabbyCutyou/0db/server/config"
	"github.com/StabbyCutyou/buffstreams"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/memberlist"
	"hash/fnv"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Slaxos struct {
	consensusPort int
	adminPort     int
	members       *memberlist.Memberlist
	buffManager   *buffstreams.BuffManager
}

type ServerEntry struct {
	Id      uint64
	Address string
}

func NewSlaxos(cfg *config.MembershipConfig) *Slaxos {
	buffManager := buffstreams.New()

	s := &Slaxos{
		buffManager:   buffManager,
		consensusPort: cfg.ReceivePort,
		adminPort:     cfg.AdminPort,
	}
	// TODO clean this up - this is ugly and I don't like it
	memList := newMemberlist(cfg, s)
	s.members = memList
	// Listen for admin messages
	logrus.Info("Starting listening for Admin connections...")
	buffManager.StartListening(strconv.Itoa(cfg.AdminPort), s.adminListenerCallback)
	// Listen for cross-node messages
	logrus.Info("Started listening for Cross-Node connections...")
	buffManager.StartListening(strconv.Itoa(cfg.ReceivePort), s.crossNodeListenerCallback)
	return s
}

func (s *Slaxos) crossNodeListenerCallback(data []byte) error {
	// We now have a message in the dataButter, we should handle it
	msg := &message.DistributedWrite{}
	err := proto.Unmarshal(data, msg)
	if err != nil || msg == nil {
		// Error decoding
		logrus.Error("Error trying to unmarshall")
		logrus.Error(err)
	} else {
		logrus.Info("This is where i'd write ", msg)
	}
	return err
}

func (s *Slaxos) adminListenerCallback(data []byte) error {
	logrus.Info("Bytes are")
	logrus.Info(data)
	msg := &message.AdminMessage{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		logrus.Error("Error unmsarshalling data")
		logrus.Error(err)
	}
	err = s.runCommand(msg)
	logrus.Debug("Ending Admin Connection")
	return err
}

func newMemberlist(cfg *config.MembershipConfig, s *Slaxos) *memberlist.Memberlist {
	mListCfg := memberlist.DefaultLANConfig()
	hostname, err := os.Hostname()
	if err != nil {
		// Could happen, i guess?
		logrus.Error("Could not obtain hostname")
		logrus.Error(err)
		hostname = "localhostfail"
	}
	mListCfg.Name = hostname
	mListCfg.BindPort = cfg.MemberPort
	// TODO replace the below with the new callback system, need new interface to handle invoking callbacks
	//connList := NewConnectionList(cfg)
	mListCfg.Events = s

	logrus.Info("Creating Membership Listener")
	list, err := memberlist.Create(mListCfg)

	if err != nil {
		logrus.Error("Unable to create a Membership Listener. This node is considered Partitioned until the network heals")
		logrus.Error(err)
	}

	return list
}

func (s *Slaxos) runCommand(msg *message.AdminMessage) error {
	var err error = nil
	switch *msg.Command {
	case message.COMMAND_JOIN_CLUSTER:
		err = s.joinCluster(*msg.Message)
	case message.COMMAND_LEAVE_CLUSTER:
		err = s.leaveCluster()
	}
	return err
}

func (s *Slaxos) joinCluster(address string) error {
	logrus.Infof("About to join to Cluster at address %s", address)
	nodesFound, err := s.members.Join([]string{address})
	if err != nil {
		logrus.Errorf("Unable to join Cluster at address %s", address)
		logrus.Error(err)
	} else {
		logrus.Infof("Joined cluster at %s, found %d other members", address, nodesFound)
	}
	return err
}

func (s *Slaxos) leaveCluster() error {
	logrus.Info("Leaving Cluster")
	err := s.members.Leave(50 * time.Millisecond)
	if err != nil {
		logrus.Error("Error leaving Cluster")
		logrus.Error(err)
	}
	return err
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
	return keyHash % uint64(s.members.NumMembers())
}

func isLocalInterface(address string) bool {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			if ipv4.String() == address {
				return true
			}
		}
	}
	return false
}

func (s *Slaxos) writeToNode(nodeId uint64, key string, data string, ack bool) error {
	// TODO this functionality should live in a dispacher...
	// Get the address of the node
	chosenNode := s.members.Members()[nodeId]
	logrus.Debug("ADDR IS")
	logrus.Debug(chosenNode.Addr.String())
	if isLocalInterface(chosenNode.Addr.String()) {
		// The local node owns it
		logrus.Debug("Writing to local storage")
		cmd := exec.Command("echo", data, "> /dev/null")
		if ack == true {
			return cmd.Run()
		} else {
			cmd.Start()
			return nil
		}
	} else {
		logrus.Debug("Writing to remote storage")
		// Network send the write
		// TODO support clustering
		msg := &message.DistributedWrite{Key: &key, Data: &data}

		msgBytes, err := proto.Marshal(msg)
		if err != nil {
			return err
		}

		logrus.Info("ABOUT TO WRITE")
		logrus.Info(msgBytes)
		// Keep the connection open on each write
		_, err = s.buffManager.WriteTo(chosenNode.Addr.String(), strconv.FormatInt(int64(s.consensusPort), 10), msgBytes, true)
		return err
	}
}

func (s *Slaxos) readFromNode(nodeId uint64, key string) (string, error) {
	// TODO this functionality should live in a dispacher...
	// Get the address of the node
	chosenNode := s.members.Members()[nodeId]
	logrus.Debug("ADDR IS")
	logrus.Debug(chosenNode.Addr.String())
	if isLocalInterface(chosenNode.Addr.String()) {
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

// TODO Need a better place to keep these I feel... but where
// When a new node joins the list, open a connection to/from it
func (s *Slaxos) NotifyJoin(n *memberlist.Node) {
	// Do nothing if this is our own node, which does join the cluster
	ourName, err := os.Hostname()
	if err != nil {
		logrus.Error("Failed to resolve our own hostname for comparison on incoming node joining cluster")
		//????
	}
	if n.Name != ourName {
		logrus.Info("NODE JOINED ", n)
		logrus.Infof("Adding write connection %s", n.Addr.String())
		s.buffManager.DialOut(n.Addr.String(), strconv.FormatInt(int64(n.Port), 10))
	}
}

func (s *Slaxos) NotifyLeave(n *memberlist.Node) {
	logrus.Infof("Removing connection %s", n.Addr.String())
	// TODO need to add a way to close out connections in buffstreams
}

func (s *Slaxos) NotifyUpdate(n *memberlist.Node) {

}
