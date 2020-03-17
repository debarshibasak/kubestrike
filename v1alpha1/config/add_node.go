package config

import (
	"errors"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"

	"github.com/debarshibasak/kubestrike/v1alpha1"
	"github.com/ghodss/yaml"
)

type AddNode struct {
	Base
	Multipass         *v1alpha1.MultiPassAddNode `yaml:"multipass" json:"multipass"`
	BareMetal         *v1alpha1.BaremetalAddNode `yaml:"baremetal" json:"baremetal"`
	SkipWorkerFailure bool                       `yaml:"skip_worker_failure" json:"skip_worker_failure"`
}

func (a *AddNode) Parse(config []byte) (ClusterOperation, error) {
	var orchestration AddNode

	err := yaml.Unmarshal(config, &orchestration)
	if err != nil {
		return nil, errors.New("error while parsing configuration")
	}

	return &orchestration, nil
}

func (a *AddNode) Validate() error {

	if a.Kind != AddNodeKind {
		return errKind
	}

	if a.Provider == MultipassProvider && a.Multipass == nil {
		return errMultipass
	}

	if a.Provider == MultipassProvider && a.Multipass.Master == nil {
		return errMultipass
	}

	if a.Provider == MultipassProvider && len(a.Multipass.Master) == 0 {
		return errMultipass
	}

	if a.Provider == BaremetalProvider && a.BareMetal == nil {
		return errBaremetal
	}

	return nil
}

//TODO Delete the acquireVM on failure
func (a *AddNode) Run(verbose bool) error {

	master, workers, err := getNode(a)
	if err != nil {
		return err
	}

	kadmClient := kubeadmclient.Kubeadm{
		MasterNodes: []*kubeadmclient.MasterNode{
			master,
		},
		WorkerNodes:       workers,
		VerboseMode:       verbose,
		SkipWorkerFailure: false,
	}

	return kadmClient.AddNode()
}

func getNode(d *AddNode) (*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, error) {
	switch d.Provider {
	case MultipassProvider:
		return d.Multipass.GetNodes()
	case BaremetalProvider:
		return d.BareMetal.GetNodes()
	default:
		return nil, nil, errors.New("no provider found")
	}
}
