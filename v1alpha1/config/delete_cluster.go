package config

import (
	"errors"
	"log"

	"github.com/debarshibasak/kubestrike/v1alpha1/provider"

	"github.com/ghodss/yaml"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
)

type DeleteCluster struct {
	Base
	Multipass *provider.MultiPassDeleteCluster `yaml:"multipass" json:"multipass"`
	BareMetal *provider.BaremetalDeleteCluster `yaml:"baremetal" json:"baremetal"`
}

func (d *DeleteCluster) Run(verbose bool) error {

	log.Println("[kubestrike] provider found - " + d.Base.Provider)

	master, worker, err := GetDeleteCluster(d)
	if err != nil {
		return err
	}

	kadmClient := kubeadmclient.Kubeadm{
		ClusterName:          d.ClusterName,
		MasterNodes:          master,
		WorkerNodes:          worker,
		ResetOnDeleteCluster: true,
		VerboseMode:          verbose,
	}

	return kadmClient.DeleteCluster()
}

func (d *DeleteCluster) Validate() error {

	if d.Multipass != nil && d.Multipass.OnlyKube && len(d.Multipass.MasterIP) == 0 {
		return errors.New("no master specified")
	}

	return nil
}

func (d *DeleteCluster) Parse(config []byte) (ClusterOperation, error) {

	var orchestration DeleteCluster

	err := yaml.Unmarshal(config, &orchestration)
	if err != nil {
		return nil, errors.New("error while parsing configuration")
	}

	return &orchestration, nil
}
