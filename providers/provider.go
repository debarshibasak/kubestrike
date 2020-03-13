package providers

import (
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/pkg/errors"
)

type Provider string

const (
	MultipassProvider Provider = "Multipass"
	BaremetalProvider Provider = "Baremetal"
)

type Providers interface {
	Provision() ([]*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, *kubeadmclient.HaProxyNode, error)
}

func Get(provider string, mastercount int, workercount int) ([]*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, *kubeadmclient.HaProxyNode, error) {

	switch provider {
	case "multipass":
		{
			p := &Multipass{
				Worker: workercount,
				Master: mastercount,
			}
			return p.Provision()
		}
	}

	return nil, nil, nil, errors.New("provisioner not found")
}
