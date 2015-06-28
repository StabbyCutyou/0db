package consensus

import (
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/message"
	"github.com/StabbyCutyou/0db/server/config"
	"github.com/StabbyCutyou/0db/util"
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
	receivePort    int
	dispatchPort   int
	adminPort      int
	members        *memberlist.Memberlist
	connectionList *ConnectionList
}

type ServerEntry struct {
	Id      uint64
	Address string
}

func NewSlaxos(cfg *config.MembershipConfig) *Slaxos {
	memList, connectionList := newMemberlist(cfg)
	s := Slaxos{
		members:        memList,
		connectionList: connectionList,
		receivePort:    cfg.ReceivePort,
		dispatchPort:   cfg.DispatchPort,
		adminPort:      cfg.AdminPort,
	}
	s.startAdminListener(cfg.AdminPort)
	return &s
}

func newMemberlist(cfg *config.MembershipConfig) (*memberlist.Memberlist, *ConnectionList) {
	mListCfg := memberlist.DefaultLANConfig()
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Error("Could not obtain hostname")
		logrus.Error(err)
		hostname = "localhostfail"
	}
	mListCfg.Name = hostname
	mListCfg.BindPort = cfg.MemberPort
	connList := NewConnectionList(cfg)
	mListCfg.Events = connList

	logrus.Info("Creating Membership Listener")
	list, err := memberlist.Create(mListCfg)

	if err != nil {
		logrus.Error("Unable to create a Membership Listener. This node is considered Partitioned until the network heals")
		logrus.Error(err)
	}

	return list, connList
}

func (s *Slaxos) startAdminListener(adminPort int) {
	logrus.Debug("Starting admin listener")
	go func() {
		socket, err := net.Listen("tcp", ":"+strconv.Itoa(adminPort))
		if err != nil {
			logrus.Errorf("Error binding to TCP Port %d while attempting to create admin listener", adminPort)
			logrus.Error(err)
		}

		// Begin listen loop
		for {
			logrus.Debug("Awaiting admin connection...")

			conn, err := socket.Accept()
			if err != nil {
				logrus.Error("Error accepting remote admin connection")
				logrus.Error(err)
			} else {
				logrus.Debug("Accepted admin connection")
				s.handleAdminConnection(conn)
			}
		}
	}()
}

func (s *Slaxos) handleAdminConnection(conn net.Conn) {
	go func() {
		defer conn.Close()
		buffer := make([]byte, 128)  // Read 128 bytes at time
		bytesRead := make([]byte, 0) // Hold the data read outside the buffer
		// Do an initial read
		bytesLen, err := conn.Read(buffer)
		// While we're reading bytes, and there is no error or the error is a natural EOF
		for bytesLen > 0 && (err == nil || err.Error() == "EOF") {
			logrus.Debug("Reading Bytes")
			bytesRead = append(bytesRead, buffer[:bytesLen]...)
			bytesLen, err = conn.Read(buffer)
		}

		if err != nil && err.Error() == "EOF" {
			// The connection has reaced a natural conclusion - parse the command and run it
			msg := &message.AdminMessage{}
			proto.Unmarshal(bytesRead, msg)
			s.runCommand(msg)
			logrus.Debug("Ending Admin Connection")
			return
		} else if err != nil {
			logrus.Error("Underlying network failure during Admin Connection")
		}
	}()
}

func (s *Slaxos) runCommand(msg *message.AdminMessage) {
	switch *msg.Command {
	case message.COMMAND_JOIN_CLUSTER:
		s.joinCluster(*msg.Message)
	case message.COMMAND_LEAVE_CLUSTER:
		s.leaveCluster()
	}
}

func (s *Slaxos) joinCluster(address string) {
	logrus.Infof("About to join to Cluster at address %s", address)
	nodesFound, err := s.members.Join([]string{address})
	if err != nil {
		logrus.Errorf("Unable to join Cluster at address %s", address)
		logrus.Error(err)
	} else {
		logrus.Infof("Joined cluster at %s, found %d other members", address, nodesFound)
	}
}

func (s *Slaxos) leaveCluster() {
	logrus.Info("Leaving Cluster")
	err := s.members.Leave(50 * time.Millisecond)
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
		// TODO don't hardcode header size
		toWriteLen := util.UInt16ToByteArray(uint16(len(msgBytes)), util.MessageSizeToBitLength(4096))
		logrus.Info(msgBytes)
		if err != nil {
			return err
		}
		toWrite := append(toWriteLen, msgBytes...)
		logrus.Info("ABOUT TO WRITE")
		logrus.Info(toWrite)
		_, err = s.connectionList.dispatchConnections[chosenNode].Write(toWrite)
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
