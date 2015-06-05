package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/StabbyCutyou/0db/config"
	httpendpoint "github.com/StabbyCutyou/0db/endpoints/http/v1"
	"github.com/StabbyCutyou/0db/server"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("Booting 0DB - The database that is 0% a database!")
	cfg, err := config.GetConfig()

	if err != nil {
		logrus.Error(err)
	}

	zdb := server.New(cfg)
	// This will block the main thread
	httpendpoint.Listen(cfg.Rest.HttpPort, zdb)
}
