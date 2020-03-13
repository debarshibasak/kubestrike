package config

import (
	"log"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/debarshibasak/kubestrike/providers"
)

type Machine struct {
	Username        string `yaml:"username" json:"username"`
	IP              string `yaml:"ip" json:"ip"`
	PrivateKey      string `yaml:"privateKey" json:"privateKey"`
	PrivateLocation string `yaml:"privateKeyLocation" json:"privateKeyLocation"`
}
type Kind string

const (
	ClusterOrchestration Kind = "ClusterOrchestration"
)

type ClusterOrchestrator struct {
	APIVersion  string             `yaml:"apiVersion" json:"apiVersion"`
	Kind        Kind               `yaml:"kind" json:"kind"`
	Provider    providers.Provider `yaml:"provider" json:"provider"`
	ClusterName string             `yaml:"clusterName" json:"clusterName"`
	Multipass   *struct {
		MasterCount int `yaml:"masterCount" json:"masterCount"`
		WorkerCount int `yaml:"workerCount" json:"workerCount"`
	} `yaml:"multipass" json:"multipass"`
	BareMetal *struct {
		Master                    []Machine `yaml:"master" json:"master"`
		Worker                    []Machine `yaml:"worker" json:"worker"`
		HAProxy                   Machine   `yaml:"haproxy" json:"haproxy"`
		DefaultPrivateKey         string    `yaml:"defaultPrivateKey" json:"defaultPrivateKey"`
		DefaultPrivateKeyLocation string    `yaml:"defaultPrivateKeyLocation" json:"defaultPrivateKeyLocation"`
		DefaultUsername           string    `yaml:"defaultUsername" json:"defaultUsername"`
	} `yaml:"baremetal" json:"baremetal"`
	Networking *struct {
		Plugin  string `yaml:"plugin" json:"plugin"`
		PodCidr string `yaml:"podCidr" json:"podCidr"`
	} `yaml:"networking" json:"networking"`
}

func (clusterOrchestrator *ClusterOrchestrator) Install() error {

	log.Println("[kubestrike] provider found - " + clusterOrchestrator.Provider)
	log.Println("[kubestrike] creating vm...")

	if clusterOrchestrator.Provider == providers.MultipassProvider {

		masterNodes, workerNodes, haproxy, err := providers.Get(
			clusterOrchestrator.Provider,
			clusterOrchestrator.Multipass.MasterCount,
			clusterOrchestrator.Multipass.WorkerCount,
		)
		if err != nil {
			log.Fatal(err)
		}

		var networking *kubeadmclient.Networking

		cni := clusterOrchestrator.Networking.Plugin
		if cni == "" {
			networking = kubeadmclient.Flannel
		} else {
			networking := kubeadmclient.LookupNetworking(cni)
			if networking == nil {
				log.Fatal("network plugin in empty")
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
			log.Fatal(err)
		}
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
