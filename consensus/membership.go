package consensus

import (
	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/memberlist"
	"os"
	"time"
)

type Membership struct {
	members *memberlist.Memberlist
}

func NewMembershipListener(memberPort int) *Membership {
	mListCfg := memberlist.DefaultLANConfig()
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Error("Could not obtain hostname")
		logrus.Error(err)
		hostname = "localhostfail"
	}
	mListCfg.Name = hostname
	mListCfg.BindPort = memberPort

	logrus.Info("Creating Membership Listener")
	list, err := memberlist.Create(mListCfg)

	if err != nil {
		logrus.Error("Unable to create a Membership Listener. This node is considered Partitioned until the network heals")
		logrus.Error(err)
	}

	return &Membership{members: list}
}

func (m *Membership) LeaveCluster() error {
	return m.members.Leave(500 * time.Millisecond)
}

func (m *Membership) JoinCluster(address string) (int, error) {
	return m.members.Join([]string{address})
}

func (m *Membership) ClusterSize() int {
	return len(m.members.Members())
}

func (m *Membership) MemberAt(nodeId uint64) string {
	return m.members.Members()[nodeId].Addr.String()
}
