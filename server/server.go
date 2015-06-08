package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/server/config"
	httpendpoint "github.com/StabbyCutyou/0db/server/endpoints/http/v1"
	"github.com/StabbyCutyou/0db/server/node"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	bootNode()
}

// TODO Refactor this to live under the node package?
func bootNode() {
	logrus.Info("Booting 0DB - The database that is 0% a database!")
	cfg, err := config.GetConfig()

	if err != nil {
		logrus.Error(err)
	}

	zdb := node.New(cfg)
	// This will block the main thread
	httpendpoint.Listen(cfg.Rest.HttpPort, zdb)
}
