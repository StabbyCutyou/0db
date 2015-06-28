package consensus

import (
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/server/config"
	"github.com/hashicorp/memberlist"
	"net"
	"strconv"
	"sync"
)

type ConnectionList struct {
	sync.RWMutex
	dispatchConnections map[*memberlist.Node]net.Conn
	receiveConnections  map[*memberlist.Node]net.Conn
	membershipConfig    *config.MembershipConfig
	receiveSocket       *net.Listener
}

func (cl *ConnectionList) addConnection(n *memberlist.Node) {

}

func (cl *ConnectionList) removeConnectionAt(n *memberlist.Node) {
	cl.Lock()
	defer cl.Unlock()
	logrus.Infof("Removing connection %s", n.Addr.String())
	cl.dispatchConnections[n].Close()
	cl.receiveConnections[n].Close()
	// Remove it
	delete(cl.dispatchConnections, n)
	delete(cl.receiveConnections, n)
}

func NewConnectionList(conf *config.MembershipConfig) *ConnectionList {
	receiveSocket, err := net.Listen("tcp", ":"+strconv.Itoa(conf.ReceivePort))
	if err != nil {
		logrus.Error(err)
		// Error opening connection - bail out?
	}
	return &ConnectionList{
		dispatchConnections: make(map[*memberlist.Node]net.Conn),
		receiveConnections:  make(map[*memberlist.Node]net.Conn),
		membershipConfig:    conf,
		receiveSocket:       &receiveSocket,
	}
}

// TODO ported from splitter, needs work
func readFromConnection(reader net.Conn, buffer []byte) (int, error) {
	bytesLen, err := reader.Read(buffer)
	// Output the content of the bytes to the queue
	if bytesLen == 0 {
		if err != nil && err.Error() == "EOF" {
			logrus.Error("End of individual transmission")
			return bytesLen, err
		}
	}

	if err != nil {
		logrus.Error("Underlying network failure?")
		logrus.Error(err)
	}

	return bytesLen, nil
}

// TODO ported from splitter, needs work
func handleReadConnection(conn net.Conn) {
	//headerBuffer := make([]byte, 4) // 4 is hardcoded for now, should be configurable
	//
	//for {
	//	logrus.Debug("Begining Read")
	//	// First, read the number of bytes required to determine the message length
	//	_, err := readFromConnection(conn, headerBuffer)
	//
	//	if err != nil && err.Error() == "EOF" {
	//		conn.Close()
	//		break
	//	}
	//
	//	msgLength, bytesParsed := binary.Uvarint(headerBuffer)
	//	if bytesParsed == 0 {
	//		logrus.Error("Buffer too small")
	//		break
	//	} else if bytesParsed < 0 {
	//		logrus.Error("Buffer overflow")
	//		break
	//	}
	//	logrus.Info(msgLength)
	//	logrus.Info(bytesParsed)
	//	dataBuffer := make([]byte, msgLength)
	//	bytesLen, err := readFromConnection(conn, dataBuffer)
	//	if err != nil && err.Error() == "EOF" {
	//		conn.Close()
	//		break
	//	}
	//
	//	logrus.Info(dataBuffer)
	//
	//	if bytesLen > 0 && (err == nil || (err != nil && err.Error() == "EOF")) {
	//		readQueue.Enqueue(dataBuffer)
	//	}
	//}
}

// When a new node joins the list, open a connection to/from it
func (cl *ConnectionList) NotifyJoin(n *memberlist.Node) {
	logrus.Info("NODE JOINED ", n)
	//cl.Lock()
	//defer cl.Unlock()
	//logrus.Infof("Adding read connection %s", n.Addr.String())
	//conn, err := net.Dial("tcp", writer)
	//if err != nil {
	//	logrus.Error(err)
	//	// Bail out?
	//} else {
	//	// Open read connect
	//
	//}
}

func (cl *ConnectionList) NotifyLeave(n *memberlist.Node) {

}

func (cl *ConnectionList) NotifyUpdate(n *memberlist.Node) {

}
