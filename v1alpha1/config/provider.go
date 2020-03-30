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

func Get(orchestrator *CreateCluster) error {

	switch orchestrator.Provider {
	case MultipassProvider:
		{

			masters, workers, haproxy, err := orchestrator.Multipass.Provision()
			if err != nil {
				return err
			}

			orchestrator.WorkerNodes = workers
			orchestrator.MasterNodes = masters
			orchestrator.HAProxy = haproxy
			return nil
		}
	case BaremetalProvider:
		{
			masters, workers, haproxy, err := orchestrator.BareMetal.Provision()
			if err != nil {
				return err
			}

			orchestrator.WorkerNodes = workers
			orchestrator.MasterNodes = masters
			orchestrator.HAProxy = haproxy
			return nil
		}
	}

	return errors.New("provisioner not found")
}
