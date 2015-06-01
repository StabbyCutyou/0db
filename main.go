package main

import (
	"github.com/Sirupsen/logrus"
	httpendpoint "github.com/StabbyCutyou/0db/endpoints/http/v1"
	"github.com/StabbyCutyou/0db/server"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("Booting 0DB - The database that is 0% a database!")
	zdb := server.New()
	// This will block the main thread
	httpendpoint.Listen(5050, zdb)
}
