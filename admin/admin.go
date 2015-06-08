package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/admin/client"
	"github.com/StabbyCutyou/0db/admin/config"
	"github.com/StabbyCutyou/0db/message"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	cfg := config.GetConfig()
	if cfg.Command == "" {
		logrus.Error("Please provide a valid command.")
		return
	}
	adminClient, err := client.New(cfg.Port)
	if err != nil {
		logrus.Error("Error initiating client")
		logrus.Fatal(err)
	}

	switch cfg.Command {
	case message.COMMAND_JOIN_CLUSTER:
		err = adminClient.JoinCluster(cfg.Message)
	case message.COMMAND_LEAVE_CLUSTER:
		err = adminClient.LeaveCluster()
	default:
		logrus.Error("Please specify a valid command")
		return
	}

	if err != nil {
		logrus.Error("Error running command")
		logrus.Fatal(err)
	}
}
