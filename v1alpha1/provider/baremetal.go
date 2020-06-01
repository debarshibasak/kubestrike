package provider

import (
	"errors"

	"github.com/debarshibasak/machina"
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

func (m *BaremetalAddNode) GetNodes() (*AddNodeResponse, error) {

	var addNodeResponse AddNodeResponse
	var workerNodes []*machina.Node
	for _, workerMachine := range m.Worker {
		workerNodes = append(workerNodes, machina.NewNode(m.DefaultUsername, workerMachine.IP, m.DefaultPrivateKeyLocation))
	}

	addNodeResponse.Master = machina.NewNode(m.DefaultUsername, m.Master.IP, m.DefaultPrivateKeyLocation)
	return &addNodeResponse, nil
}

func (m *BaremetalAddNode) GetNodesForDeletion() (*RemoveNodeResponse, error) {

	var resp RemoveNodeResponse
	var workerNodes []*machina.Node
	for _, workerMachine := range m.Worker {
		workerNodes = append(workerNodes, machina.NewNode(m.DefaultUsername, workerMachine.IP, m.DefaultPrivateKeyLocation))
	}

	resp.Worker = workerNodes
	resp.Master = machina.NewNode(m.DefaultUsername, m.Master.IP, m.DefaultPrivateKeyLocation)

	return &resp, nil
}

func (m *BaremetalDeleteCluster) DeleteInstance() (*DeleteClusterResponse, error) {

	var masterNodes []*machina.Node
	var workerNodes []*machina.Node

	var resp DeleteClusterResponse

	//TODO Do alternative possibilities check here
	for _, node := range m.Master {
		masterNodes = append(masterNodes, machina.NewNode(m.DefaultUsername, node.IP, m.DefaultPrivateKeyLocation))
	}

	for _, node := range m.Worker {
		workerNodes = append(workerNodes, machina.NewNode(m.DefaultUsername, node.IP, m.DefaultPrivateKeyLocation))
	}

	resp.Master = masterNodes
	resp.Worker = workerNodes

	return &resp, nil
}

func (m *Baremetal) Provision() ([]*machina.Node, []*machina.Node, *machina.Node, error) {
	var (
		masterNodes []*machina.Node
		workerNodes []*machina.Node
		haproxy     *machina.Node
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
		haproxy = machina.NewNode(username, m.HAProxy.IP, m.DefaultPrivateKeyLocation)
	}

	for _, workerMachine := range m.Worker {
		workerNodes = append(workerNodes, machina.NewNode(m.DefaultUsername, workerMachine.IP, m.DefaultPrivateKeyLocation))
	}

	for _, masterMachine := range m.Master {
		masterNodes = append(masterNodes, machina.NewNode(m.DefaultUsername, masterMachine.IP, m.DefaultPrivateKeyLocation))
	}

	return masterNodes, workerNodes, haproxy, nil
}
