package client

import (
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/message"
	"github.com/StabbyCutyou/buffstreams"
	"github.com/golang/protobuf/proto"
	"strconv"
)

type AdminClient struct {
	port        int
	buffManager *buffstreams.BuffManager
}

func New(port int) (*AdminClient, error) {
	return &AdminClient{
		port:        port,
		buffManager: buffstreams.New(),
	}, nil
}

func (ac *AdminClient) JoinCluster(address string) error {
	logrus.Info("Joining Cluster")
	cmd := message.COMMAND_JOIN_CLUSTER
	msg := &message.AdminMessage{Command: &cmd, Message: &address}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Debug("About to write")
	logrus.Debug(msgBytes)
	_, err = ac.buffManager.WriteTo("127.0.0.1", strconv.Itoa(ac.port), msgBytes, false)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

// Don't use this yet, it does dumb shit and breaks things
func (ac *AdminClient) LeaveCluster() error {
	cmd := message.COMMAND_LEAVE_CLUSTER
	msg := &message.AdminMessage{Command: &cmd, Message: nil}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = ac.buffManager.WriteTo("127.0.0.1", strconv.Itoa(ac.port), msgBytes, false)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
