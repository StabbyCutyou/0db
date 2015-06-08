package client

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/message"
	"github.com/golang/protobuf/proto"
	"net"
)

type AdminClient struct {
	Port int
}

func New(port int) (*AdminClient, error) {
	return &AdminClient{Port: port}, nil
}

func (ac *AdminClient) openConnection() (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", ac.Port))
	if err != nil {
		logrus.Error("Error connecting to 0db node")
		logrus.Error(err)
		return nil, err
	}
	return conn, nil
}

func (ac *AdminClient) JoinCluster(address string) error {
	cmd := message.COMMAND_JOIN_CLUSTER
	msg := &message.AdminMessage{Command: &cmd, Message: &address}
	msgBytes, err := proto.Marshal(msg)
	conn, err := ac.openConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func (ac *AdminClient) LeaveCluster() error {
	cmd := message.COMMAND_LEAVE_CLUSTER
	msg := &message.AdminMessage{Command: &cmd, Message: nil}
	msgBytes, err := proto.Marshal(msg)
	conn, err := ac.openConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(msgBytes)
	if err != nil {
		return err
	}
	return nil
}
