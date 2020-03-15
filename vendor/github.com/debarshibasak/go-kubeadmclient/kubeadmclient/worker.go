package kubeadmclient

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type WorkerNode struct {
	*Node
}

func NewWorkerNode(username string,
	ipOrHost string,
	privateKeyLocation string) *WorkerNode {

	return &WorkerNode{
		Node: &Node{
			username:           username,
			ipOrHost:           ipOrHost,
			privateKeyLocation: privateKeyLocation,
			clientID:           uuid.New().String(),
		},
	}
}

func (n *WorkerNode) ctlCommand(cmd string) error {
	return n.sshClientWithTimeout(1 * time.Minute).Run([]string{
		"sudo KUBECONFIG=/etc/kubernetes/kubelet.conf " + cmd,
	})
}

func (n *WorkerNode) getHostName() (string, error) {

	hostname, err := n.sshClient().Collect("hostname")
	if err != nil {
		return "", err
	}

	return hostname, nil
}

func (n *WorkerNode) setWorkerLabel() error {
	return n.ctlCommand("kubectl label nodes --selector=kubernetes.io/hostname=`hostname` node-type=worker")
}

func (n *WorkerNode) drainAndReset() (string, error) {

	osType, err := n.determineOS()
	if err != nil {
		return "", err
	}

	if osType == nil {
		return "", errors.New("could not determine ostype, may be it could not ssh into it, or does not support the os")
	}

	hostname, err := n.getHostName()
	if err != nil {
		return "", err
	}

	if err := n.ctlCommand("kubectl drain `hostname`"); err != nil {
		return hostname, err
	}

	if err := n.sshClientWithTimeout(30 * time.Minute).Run([]string{"sudo kubeadm reset -f"}); err != nil {
		return hostname, err
	}

	return hostname, nil
}

func (n *WorkerNode) install(joinCommand string) error {

	osType, err := n.determineOS()
	if err != nil {
		return err
	}

	if err := n.sshClientWithTimeout(30 * time.Minute).Run(osType.Commands()); err != nil {
		return err
	}

	if err := n.sshClientWithTimeout(30 * time.Minute).Run([]string{
		"sudo " + joinCommand,
	}); err != nil {
		return err
	}

	return n.setWorkerLabel()
}
