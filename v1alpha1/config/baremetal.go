package config

import (
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/pkg/errors"
)

type Machine struct {
	Username        string `yaml:"username" json:"username"`
	IP              string `yaml:"ip" json:"ip"`
	PrivateKey      string `yaml:"privateKey" json:"privateKey"`
	PrivateLocation string `yaml:"privateKeyLocation" json:"privateKeyLocation"`
}

type Baremetal struct {
	Master                    []Machine `yaml:"master" json:"master"`
	Worker                    []Machine `yaml:"worker" json:"worker"`
	HAProxy                   Machine   `yaml:"haproxy" json:"haproxy"`
	DefaultPrivateKey         string    `yaml:"defaultPrivateKey" json:"defaultPrivateKey"`
	DefaultPrivateKeyLocation string    `yaml:"defaultPrivateKeyLocation" json:"defaultPrivateKeyLocation"`
	DefaultUsername           string    `yaml:"defaultUsername" json:"defaultUsername"`
}

func (m *Baremetal) Provision() ([]*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, *kubeadmclient.HaProxyNode, error) {
	var (
		masterNodes []*kubeadmclient.MasterNode
		workerNodes []*kubeadmclient.WorkerNode
		haproxy     *kubeadmclient.HaProxyNode
	)

	if len(m.Master) > 1 {

		var username string

		if m.HAProxy.Username == "" {
			if m.DefaultUsername == "" {
				return masterNodes, workerNodes, haproxy, errors.New("username is empty")
			}
			username = m.DefaultUsername
		} else {
			username = m.HAProxy.Username
		}

		if m.HAProxy.IP == "" {
			return masterNodes, workerNodes, haproxy, errors.New("ip is not set for haproxy machine")
		}

		//TODO change it to correct key location

		haproxy = kubeadmclient.NewHaProxyNode(username, m.HAProxy.IP, m.DefaultPrivateKeyLocation)
	}

	for _, workerMachine := range m.Worker {
		workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode(workerMachine.IP, m.DefaultUsername, m.DefaultPrivateKeyLocation))
	}

	for _, masterMachine := range m.Master {
		masterNodes = append(masterNodes, kubeadmclient.NewMasterNode(masterMachine.IP, m.DefaultUsername, m.DefaultPrivateKeyLocation))
	}

	return masterNodes, workerNodes, haproxy, nil
}
