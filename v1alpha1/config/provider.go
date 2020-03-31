package config

import (
	"errors"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
)

type Provider string

const (
	MultipassProvider Provider = "Multipass"
	BaremetalProvider Provider = "Baremetal"
)

type Providers interface {
	Provision() ([]*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, *kubeadmclient.HaProxyNode, error)
}

func GetDeleteCluster(orchestrator *DeleteCluster) ([]*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, error) {

	switch orchestrator.Base.Provider {
	case MultipassProvider:
		{
			return orchestrator.Multipass.DeleteInstances()
		}
	case BaremetalProvider:
		{
			return orchestrator.BareMetal.DeleteInstance()
		}
	}

	return nil, nil, errors.New("provisioner not found")
}

func Get(createCluster *CreateCluster) error {

	if createCluster.Multipass != nil {

		masters, workers, haproxy, err := createCluster.Multipass.Provision()
		if err != nil {
			return err
		}

		createCluster.WorkerNodes = workers
		createCluster.MasterNodes = masters
		createCluster.HAProxy = haproxy
		return nil
	}

	if createCluster.BareMetal != nil {
		masters, workers, haproxy, err := createCluster.BareMetal.Provision()
		if err != nil {
			return err
		}

		createCluster.WorkerNodes = workers
		createCluster.MasterNodes = masters
		createCluster.HAProxy = haproxy
		return nil
	}

	return errors.New("provisioner not found")
}
