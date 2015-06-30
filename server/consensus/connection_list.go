package consensus

import (
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/server/config"
	"github.com/hashicorp/memberlist"
	"net"
	"os"
	"strconv"
	"sync"
)

type ConnectionList struct {
	sync.RWMutex
	dispatchConnections map[*memberlist.Node]net.Conn
	receiveConnections  map[*memberlist.Node]net.Conn
	membershipConfig    *config.MembershipConfig
	receiveSocket       *net.TCPListener
}

// When a new node joins the list, open a connection to/from it
func (cl *ConnectionList) NotifyJoin(n *memberlist.Node) {
	// Do nothing if this is our own node, which does join the cluster
	ourName, err := os.Hostname()
	if err != nil {
		logrus.Error("Failed to resolve our own hostname for comparison on incoming node joining cluster")
		//????
	}
	if n.Name != ourName {
		logrus.Info("NODE JOINED ", n)
		cl.Lock()
		defer cl.Unlock()
		logrus.Infof("Adding write connection %s", n.Addr.String())
		// We want to dial into the remote nodes listener port
		tcpAddr, err := net.ResolveTCPAddr("tcp", n.Addr.String()+":"+strconv.FormatInt(int64(cl.membershipConfig.ReceivePort), 10))
		if err != nil {
			logrus.Error(err)
			// Error resolving addr - bail out?
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			logrus.Error(err)
			// Bail out?
		} else {
			// Open read connect
			cl.dispatchConnections[n] = conn
		}
	}
}

func (cl *ConnectionList) NotifyLeave(n *memberlist.Node) {
	cl.Lock()
	defer cl.Unlock()
	logrus.Infof("Removing connection %s", n.Addr.String())
	cl.dispatchConnections[n].Close()
	cl.receiveConnections[n].Close()
	// Remove it
	delete(cl.dispatchConnections, n)
	delete(cl.receiveConnections, n)
}

func (cl *ConnectionList) NotifyUpdate(n *memberlist.Node) {

}
