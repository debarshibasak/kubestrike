package kubeadmclient

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	osType "github.com/debarshibasak/go-kubeadmclient/kubeadmclient/ostype"

	"github.com/debarshibasak/go-kubeadmclient/sshclient"
)

type Node struct {
	username           string
	ipOrHost           string
	osType             string
	privateKeyLocation string
	verboseMode        bool
	clientID           string
}

func (n *Node) String() string {
	return fmt.Sprintf("ip=%v username=%v key=%v", n.ipOrHost, n.username, n.privateKeyLocation)
}

func (n *Node) determineOS() (osType.OsType, error) {

	client := n.sshClient()
	out, err := client.Collect("uname -a")
	if err != nil {
		return nil, err
	}

	if strings.Contains(out, "Ubuntu") {
		return &osType.Ubuntu{}, err
	}

	if err := client.Run([]string{"ls /etc/centos-release"}); err == nil {
		return &osType.Centos{}, err
	}

	if err := client.Run([]string{"ls /etc/redhat-release"}); err == nil {
		return &osType.Centos{}, err
	}

	return &osType.Unknown{}, errors.New("unknown os type")
}

func (n *Node) sshClient() *sshclient.SSHConnection {
	return &sshclient.SSHConnection{
		Username:    n.username,
		IP:          n.ipOrHost,
		KeyLocation: n.privateKeyLocation,
		VerboseMode: n.verboseMode,
		ClientID:    n.clientID,
	}
}

func (n *Node) sshClientWithTimeout(duration time.Duration) *sshclient.SSHConnection {
	return &sshclient.SSHConnection{
		Username:    n.username,
		IP:          n.ipOrHost,
		KeyLocation: n.privateKeyLocation,
		VerboseMode: n.verboseMode,
		Timeout:     duration,
		ClientID:    n.clientID,
	}
}
