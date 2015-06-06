package consensus

import (
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/message"
	"github.com/golang/protobuf/proto"
	"github.com/oleiade/lane"
	"net"
	"strconv"
)

type Admin struct {
	adminPort    int
	CommandQueue *lane.Queue
}

func NewAdminListener(adminPort int) *Admin {
	logrus.Info("Creating Admin Listener")
	admin := &Admin{
		adminPort:    adminPort,
		CommandQueue: lane.NewQueue(),
	}
	go admin.startAdminListener()
	return admin
}

func (a *Admin) startAdminListener() {
	socket, err := net.Listen("tcp", ":"+strconv.Itoa(a.adminPort))
	if err != nil {
		logrus.Errorf("Error binding to TCP Port %d while attempting to Create Admin Listener", a.adminPort)
		logrus.Error(err)
	}

	// Begin listen loop
	for {
		logrus.Debug("Awaiting Admin Connection...")

		conn, err := socket.Accept()
		if err != nil {
			logrus.Error("Error accepting remote Admin connection")
			logrus.Error(err)
		} else {
			logrus.Debug("Accepted Admin Connection")
			go a.handleAdminConnection(conn)
		}
	}
}

func (a *Admin) handleAdminConnection(conn net.Conn) {
	buffer := make([]byte, 256)    // Read 256 bytes at time
	bytesRead := make([]byte, 256) // Hold the data read outside the buffer
	// Do an initial read
	bytesLen, err := conn.Read(buffer)
	// While we're reading bytes, and there is no error or the error is a natural EOF
	for bytesLen >= 0 && (err == nil || err.Error() == "EOF") {
		bytesRead = append(bytesRead, buffer...)
	}

	if err != nil && err.Error() == "EOF" {
		// The connection has reaced a natural conclusion - parse the command and run it
		msg := &message.AdminMessage{}
		proto.Unmarshal(bytesRead, msg)
		a.CommandQueue.Enqueue(msg)
		logrus.Debug("Ending Admin Connection")
		return
	} else if err != nil {
		logrus.Error("Underlying network failure during Admin Connection")
	}

}
