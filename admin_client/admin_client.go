package admin_client

import (
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/message"
)

type AdminClient struct {
	address string
}

func New(address string) (*AdminClient, error) {
	return &AdminClient{address: address}, nil
}

func (ac *AdminClient) openConnection() (net.Conn, err) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		logrus.Error("Error connecting to 0db node")
		logrus.Error(err)
		return nil, err
	}
	return conn, nil
}

func (ac *AdminClient) JoinCluster(address string) {
	msg := message.AdminMessage{Command: "JoinCluster", Message: address}
	msgBytes, err := proto.Marshal(msg)
	conn := ac.openConnection()
	conn.Write(msgBytes)
	conn.Close()
}

func (ac *AdminClient) LeaveCluster() {
	msg := message.AdminMessage{Command: "LeaveCluster", Message: ""}
	msgBytes, err := proto.Marshal(msg)
	conn := ac.openConnection()
	conn.Write(msgBytes)
	conn.Close()
}
