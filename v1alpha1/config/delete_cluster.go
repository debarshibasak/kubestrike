package config

import (
	"log"

	"github.com/pkg/errors"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"

	"github.com/debarshibasak/kubestrike/v1alpha1"
)

type DeleteCluster struct {
	Base      Base
	Multipass *v1alpha1.MultiPassDeleteCluster `yaml:"multipass" json:"multipass"`
	BareMetal *v1alpha1.BaremetalDeleteCluster `yaml:"baremetal" json:"baremetal"`
}

func (d *DeleteCluster) Run(verbose bool) error {

	log.Println("[kubestrike] provider found - " + d.Base.Provider)

	master, worker, err := GetDeleteCluster(d)
	if err != nil {
		return err
	}

	kadmClient := kubeadmclient.Kubeadm{
		ClusterName:          d.Base.ClusterName,
		MasterNodes:          master,
		WorkerNodes:          worker,
		ResetOnDeleteCluster: false,
	}

	if err := kadmClient.DeleteCluster(); err != nil {
		return err
	}

	return nil
}

func (d *DeleteCluster) Validate() error {

	if d.Multipass.OnlyKube && len(d.Multipass.MasterIP) == 0 {
		return errors.New("no master specified")
	}

	return nil
}

func (d *DeleteCluster) Parse(config []byte) (ClusterOperation, error) {
	return &DeleteCluster{}, nil
}
