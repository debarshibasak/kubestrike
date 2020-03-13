package config

import (
	"log"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/debarshibasak/kubestrike/providers"
	"github.com/pkg/errors"
)

type Kind string

const (
	ClusterOrchestration Kind = "ClusterOrchestration"
)

type ClusterOrchestrator struct {
	APIVersion  string               `yaml:"apiVersion" json:"apiVersion"`
	Kind        Kind                 `yaml:"kind" json:"kind"`
	Provider    providers.Provider   `yaml:"provider" json:"provider"`
	ClusterName string               `yaml:"clusterName" json:"clusterName"`
	Multipass   *providers.Multipass `yaml:"multipass" json:"multipass"`
	BareMetal   *providers.Baremetal `yaml:"baremetal" json:"baremetal"`
	Networking  *struct {
		Plugin  string `yaml:"plugin" json:"plugin"`
		PodCidr string `yaml:"podCidr" json:"podCidr"`
	} `yaml:"networking" json:"networking"`
}

func (clusterOrchestrator *ClusterOrchestrator) Install() error {

	log.Println("[kubestrike] provider found - " + clusterOrchestrator.Provider)

	masterNodes, workerNodes, haproxy, err := providers.Get(clusterOrchestrator)
	if err != nil {
		return err
	}

	var networking *kubeadmclient.Networking

	cni := clusterOrchestrator.Networking.Plugin
	if cni == "" {
		networking = kubeadmclient.Flannel
	} else {
		networking := kubeadmclient.LookupNetworking(cni)
		if networking == nil {
			return errors.New("network plugin in empty")
		}
	}

	log.Println("[kubestrike] creating cluster...")

	kubeadmClient := kubeadmclient.Kubeadm{
		ClusterName: clusterOrchestrator.ClusterName,
		HaProxyNode: haproxy,
		MasterNodes: masterNodes,
		WorkerNodes: workerNodes,
		VerboseMode: false,
		Netorking:   networking,
	}

	err = kubeadmClient.CreateCluster()
	if err != nil {
		return err
	}

	return nil
}

func (clusterOrchestrator *ClusterOrchestrator) Validate() error {

	if clusterOrchestrator.ClusterName == "" {
		return errClusterNameIsEmpty
	}
	if clusterOrchestrator.Kind != ClusterOrchestration {
		return errKind
	}

	if clusterOrchestrator.Provider == providers.MultipassProvider && clusterOrchestrator.Multipass == nil {
		return errMultipass
	}

	if clusterOrchestrator.Provider == providers.BaremetalProvider && clusterOrchestrator.BareMetal == nil {
		return errBaremetal
	}

	if clusterOrchestrator.Networking == nil {
		return errNetworking
	}

	return nil
}
