package v1alpha1

import (
	"errors"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
)

type Machine struct {
	Username        string `yaml:"username" json:"username"`
	IP              string `yaml:"ip" json:"ip"`
	PrivateKey      string `yaml:"privateKey" json:"privateKey"`
	PrivateLocation string `yaml:"privateKeyLocation" json:"privateKeyLocation"`
}

type BaremetalDeleteCluster struct {
	Key
	Master []Machine `yaml:"master" json:"master"`
	Worker []Machine `yaml:"workers" json:"workers"`
}

type Baremetal struct {
	Master  []Machine `yaml:"master" json:"master"`
	Worker  []Machine `yaml:"worker" json:"worker"`
	HAProxy Machine   `yaml:"haproxy" json:"haproxy"`
	Key
}

type Key struct {
	DefaultPrivateKey         string `yaml:"key" json:"keys"` //TODO
	DefaultPrivateKeyLocation string `yaml:"keyLocation" json:"keyLocation"`
	DefaultUsername           string `yaml:"username" json:"username"`
}

type BaremetalAddNode struct {
	Key
	Worker []Machine `yaml:"workers" json:"workers"`
	Master Machine   `yaml:"master" json:"master"`
}

type BaremetalDeleteNode struct {
	Key
	Worker []Machine `yaml:"worker" json:"worker"`
	Master Machine   `yaml:"master" json:"master"`
}

func (m *BaremetalAddNode) GetNodes() (*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, error) {
	var workerNodes []*kubeadmclient.WorkerNode
	for _, workerMachine := range m.Worker {
		workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode(m.DefaultUsername, workerMachine.IP, m.DefaultPrivateKeyLocation))
	}

	return kubeadmclient.NewMasterNode(m.DefaultUsername, m.Master.IP, m.DefaultPrivateKeyLocation), workerNodes, nil
}

func (m *BaremetalAddNode) GetNodesForDeletion() (*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, error) {
	var workerNodes []*kubeadmclient.WorkerNode
	for _, workerMachine := range m.Worker {
		workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode(m.DefaultUsername, workerMachine.IP, m.DefaultPrivateKeyLocation))
	}

	return kubeadmclient.NewMasterNode(m.DefaultUsername, m.Master.IP, m.DefaultPrivateKeyLocation), workerNodes, nil
}

func (m *BaremetalDeleteCluster) DeleteInstance() ([]*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, error) {

	var masterNodes []*kubeadmclient.MasterNode
	var workerNodes []*kubeadmclient.WorkerNode

	//TODO Do alternative possibilities check here
	for _, node := range m.Master {
		masterNodes = append(masterNodes, kubeadmclient.NewMasterNode(node.Username, node.IP, m.DefaultPrivateKeyLocation))
	}

	for _, node := range m.Worker {
		workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode(node.Username, node.IP, m.DefaultPrivateKeyLocation))
	}

	return masterNodes, workerNodes, nil
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
		workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode(m.DefaultUsername, workerMachine.IP, m.DefaultPrivateKeyLocation))
	}

	for _, masterMachine := range m.Master {
		masterNodes = append(masterNodes, kubeadmclient.NewMasterNode(m.DefaultUsername, masterMachine.IP, m.DefaultPrivateKeyLocation))
	}

	return masterNodes, workerNodes, haproxy, nil
}
