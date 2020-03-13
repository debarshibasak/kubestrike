package config

import "github.com/debarshibasak/kubestrike/providers"

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
	APIVersion string             `yaml:"apiVersion" json:"apiVersion"`
	Kind       Kind               `yaml:"kind" json:"kind"`
	Provider   providers.Provider `yaml:"provider" json:"provider"`
	Multipass  *struct {
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

func (clusterOrchestrator *ClusterOrchestrator) validate() error {
	if p.useStrictAPIVersionCheck {
		if err := validateAPIVersion(clusterOrchestrator.APIVersion); err != nil {
			return err
		}
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
