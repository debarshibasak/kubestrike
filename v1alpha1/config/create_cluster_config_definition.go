package config

import (
	"log"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/networking"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"

	"errors"
)

type Kind string

const (
	CreateClusterKind Kind = "CreateCluster"
	AddNodeKind       Kind = "AddNode"
	RemoveNodeKind    Kind = "RemoveNode"
	DeleteClusterKind Kind = "DeleteCluster"
)

type Base struct {
	APIVersion  string   `yaml:"apiVersion" json:"apiVersion"`
	Kind        Kind     `yaml:"kind" json:"kind"`
	Provider    Provider `yaml:"provider" json:"provider"`
	ClusterName string   `yaml:"clusterName" json:"clusterName"`
}

type CreateCluster struct {
	Base
	Multipass  *Multipass `yaml:"multipass" json:"multipass"`
	BareMetal  *Baremetal `yaml:"baremetal" json:"baremetal"`
	Networking *struct {
		Plugin  string `yaml:"plugin" json:"plugin"`
		PodCidr string `yaml:"podCidr" json:"podCidr"`
	} `yaml:"networking" json:"networking"`
}

func (createCluster *CreateCluster) Install() error {

	log.Println("[kubestrike] provider found - " + createCluster.Provider)

	masterNodes, workerNodes, haproxy, err := Get(createCluster)
	if err != nil {
		return err
	}

	var networkingPlugin *networking.Networking

	cni := createCluster.Networking.Plugin
	if cni == "" {
		networkingPlugin = networking.Flannel
	} else {
		networkingPlugin := networking.LookupNetworking(cni)
		if networkingPlugin == nil {
			return errors.New("network plugin in empty")
		}
	}

	log.Println("\n[kubestrike] creating cluster...")

	kubeadmClient := kubeadmclient.Kubeadm{
		ClusterName: createCluster.ClusterName,
		HaProxyNode: haproxy,
		MasterNodes: masterNodes,
		WorkerNodes: workerNodes,
		VerboseMode: false,
		Networking:  networkingPlugin,
	}

	err = kubeadmClient.CreateCluster()
	if err != nil {
		return err
	}

	return nil
}

func (createCluster *CreateCluster) Validate() error {

	if createCluster.ClusterName == "" {
		return errClusterNameIsEmpty
	}
	if createCluster.Kind != CreateClusterKind {
		return errKind
	}

	if createCluster.Provider == MultipassProvider && createCluster.Multipass == nil {
		return errMultipass
	}

	if createCluster.Provider == BaremetalProvider && createCluster.BareMetal == nil {
		return errBaremetal
	}

	if createCluster.Networking == nil {
		return errNetworking
	}

	return nil
}
