package config

import (
	"code.google.com/p/gcfg"
	"flag"
	"github.com/Sirupsen/logrus"
)

const FLAG_CONFIG = "c"

type Config struct {
	Rest       RestConfig
	Membership MembershipConfig
}

type RestConfig struct {
	HttpPort int
}

type MembershipConfig struct {
	MemberPort   int
	ReceivePort  int
	DispatchPort int
}

func GetConfig() (*Config, error) {
	flagData := parseFlags()
	filename := flagData[FLAG_CONFIG]

	logrus.Info("0db configuring with file located at ", *filename)
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, *filename)

	if err != nil {
		logrus.Error("0db encountered an error when reading the file located at ", *filename)
		logrus.Fatal(err)
	}

	logrus.Info("0db configuration loaded")
	return &cfg, err
}

func parseFlags() map[string]*string {
	flagData := make(map[string]*string)
	flagData[FLAG_CONFIG] = flag.String(FLAG_CONFIG, "./config/0db.cfg", "location of config file")

	flag.Parse()

	if *flagData[FLAG_CONFIG] == "" {
		logrus.Warn("Empty value provided for config file location from flag -c : Falling back to default location './config/0db.gcfg'")
		*flagData[FLAG_CONFIG] = "./config/0db.cfg"
	}

	return flagData
}
