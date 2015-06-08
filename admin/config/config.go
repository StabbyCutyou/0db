package config

import (
	"flag"
	"github.com/StabbyCutyou/0db/message"
)

// TODO too much re-use between how I like to do argument and config parsing.
// I should really be able to turn this into a generic approach that I can reuse.

const FLAG_JOIN_CLUSTER = "j"
const FLAG_LEAVE_CLUSTER = "l"
const FLAG_ADMIN_PORT = "p"

const FLAG_JOIN_CLUSTER_DEFAULT = ""
const FLAG_LEAVE_CLUSTER_DEFAULT = false
const FLAG_ADMIN_PORT_DEFAULT = 6063

type Config struct {
	Command string
	Message string
	Port    int
}

func GetConfig() *Config {
	joinCluster := flag.String(FLAG_JOIN_CLUSTER, FLAG_JOIN_CLUSTER_DEFAULT, "location of cluster to join")
	leaveCluster := flag.Bool(FLAG_LEAVE_CLUSTER, FLAG_LEAVE_CLUSTER_DEFAULT, "leave the current cluster")
	port := flag.Int(FLAG_ADMIN_PORT, FLAG_ADMIN_PORT_DEFAULT, "local port to issue command on")
	flag.Parse()

	if *leaveCluster != FLAG_LEAVE_CLUSTER_DEFAULT {
		return &Config{Command: message.COMMAND_LEAVE_CLUSTER, Message: "", Port: *port}
	} else if *joinCluster != FLAG_JOIN_CLUSTER_DEFAULT {
		return &Config{Command: message.COMMAND_JOIN_CLUSTER, Message: *joinCluster, Port: *port}
	}

	return &Config{}
}
