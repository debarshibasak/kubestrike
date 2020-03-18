package config

import (
	"errors"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"

	"github.com/debarshibasak/kubestrike/v1alpha1"
	"github.com/ghodss/yaml"
)

type DeleteNode struct {
	Base
	Multipass         *v1alpha1.MultiPassDeleteNode `yaml:"multipass" json:"multipass"`
	BareMetal         *v1alpha1.BaremetalAddNode    `yaml:"baremetal" json:"baremetal"`
	SkipWorkerFailure bool                          `yaml:"skipWorkerFailure" json:"skipWorkerFailure"`
}

func (d *DeleteNode) Parse(config []byte) (ClusterOperation, error) {
	var orchestration DeleteNode

	err := yaml.Unmarshal(config, &orchestration)
	if err != nil {
		return nil, errors.New("error while parsing second configuration - " + err.Error())
	}

	return &orchestration, nil
}

func (d *DeleteNode) Validate() error {

	if d.Kind != RemoveNodeKind {
		return errKind
	}

	if d.Provider == MultipassProvider && d.Multipass == nil {
		return errMultipass
	}

	if d.Provider == MultipassProvider && d.Multipass.Master == nil {
		return errMultipass
	}

	if d.Provider == MultipassProvider && len(d.Multipass.Master) == 0 {
		return errMultipass
	}

	if d.Provider == BaremetalProvider && d.BareMetal == nil {
		return errBaremetal
	}

	return nil
}

func (d *DeleteNode) Run(verbose bool) error {

	master, workers, err := getNodeForDeletion(d)
	if err != nil {
		return err
	}

	kadmClient := kubeadmclient.Kubeadm{
		MasterNodes: []*kubeadmclient.MasterNode{
			master,
		},
		WorkerNodes:       workers,
		VerboseMode:       verbose,
		SkipWorkerFailure: d.SkipWorkerFailure,
	}

	return kadmClient.RemoveNode()
}

func getNodeForDeletion(d *DeleteNode) (*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, error) {
	switch d.Provider {
	case MultipassProvider:
		return d.Multipass.GetNodesForDeletion()
	case BaremetalProvider:
		return d.BareMetal.GetNodesForDeletion()
	default:
		return nil, nil, errors.New("no provider found")
	}
}
