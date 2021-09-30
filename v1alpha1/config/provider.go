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

func (orchestrator *DeleteCluster) GetDeleteCluster() error {

	switch orchestrator.Base.Provider {
	case MultipassProvider:
		{
			resp, err := orchestrator.Multipass.DeleteInstances()
			if err != nil {
				return err
			}

			orchestrator.Master = resp.Master
			orchestrator.Worker = resp.Worker

			return nil

		}
	case BaremetalProvider:
		{
			resp, err := orchestrator.BareMetal.DeleteInstance()
			if err != nil {
				return err
			}

			orchestrator.Master = resp.Master
			orchestrator.Worker = resp.Worker

			return nil

		}
	}

	return errors.New("provisioner not found")
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
