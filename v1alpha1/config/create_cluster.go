package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"

	"github.com/ghodss/yaml"

	"github.com/debarshibasak/kubestrike/v1alpha1"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/networking"

	"errors"
)

type CreateCluster struct {
	Base
	Multipass  *v1alpha1.Multipass `yaml:"multipass" json:"multipass"`
	BareMetal  *v1alpha1.Baremetal `yaml:"baremetal" json:"baremetal"`
	Networking *struct {
		Plugin      string `yaml:"plugin" json:"plugin"`
		PodCidr     string `yaml:"podCidr" json:"podCidr"`
		ServiceCidr string `yaml:"serviceCidr" json:"serviceCidr"`
	} `yaml:"networking" json:"networking"`
}

func (createCluster *CreateCluster) Parse(config []byte) (ClusterOperation, error) {

	var orchestration CreateCluster

	err := yaml.Unmarshal(config, &orchestration)
	if err != nil {
		if err := json.Unmarshal(config, &orchestration); err != nil {
			return nil, errors.New("error while parsing configuration")
		}
	}

	return &orchestration, nil
}

func (createCluster *CreateCluster) Run(verbose bool) error {

	log.Println("[kubestrike] provider found - " + createCluster.Provider)

	masterNodes, workerNodes, haproxy, err := Get(createCluster)
	if err != nil {
		return err
	}

	var networkingPlugin networking.Networking

	cni := strings.TrimSpace(createCluster.Networking.Plugin)
	if cni == "" {
		networkingPlugin = *networking.Flannel
	} else {
		v := networking.LookupNetworking(cni)
		networkingPlugin = *v
		if networkingPlugin.Name == "" {
			return errors.New("network plugin in empty")
		}
	}

	log.Println("\n[kubestrike] creating cluster...")

	kubeadmClient := kubeadmclient.Kubeadm{
		ClusterName:    createCluster.ClusterName,
		HaProxyNode:    haproxy,
		MasterNodes:    masterNodes,
		WorkerNodes:    workerNodes,
		VerboseMode:    verbose,
		Networking:     &networkingPlugin,
		PodNetwork:     createCluster.Networking.PodCidr,
		ServiceNetwork: createCluster.Networking.ServiceCidr,
	}

	err = kubeadmClient.CreateCluster()
	if err != nil {
		return err
	}

	kubeConfig, err := kubeadmClient.GetKubeConfig()
	if err != nil {
		return err
	}

	u, _ := user.Current()

	return ioutil.WriteFile(u.HomeDir+"/.kubeconfig_"+createCluster.ClusterName, []byte(kubeConfig), os.FileMode(0777))
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

	if createCluster.Networking != nil && createCluster.Networking.PodCidr != "" && createCluster.Networking.ServiceCidr == "" {
		return errNetworking
	}

	if createCluster.Networking != nil && createCluster.Networking.PodCidr == "" && createCluster.Networking.ServiceCidr != "" {
		return errNetworking
	}

	return nil
}
