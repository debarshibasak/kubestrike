package config

import (
	"errors"

	"github.com/debarshibasak/machina"

	"github.com/debarshibasak/go-k3s/k3sclient"
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"

	"github.com/debarshibasak/kubestrike/v1alpha1/engine"

	"github.com/debarshibasak/kubestrike/v1alpha1/provider"

	"github.com/ghodss/yaml"
)

type RemoveNode struct {
	Base
	Multipass           *provider.MultiPassDeleteNode `yaml:"multipass" json:"multipass"`
	BareMetal           *provider.BaremetalAddNode    `yaml:"baremetal" json:"baremetal"`
	SkipWorkerFailure   bool                          `yaml:"skipWorkerFailure" json:"skipWorkerFailure"`
	Master              *machina.Node                 `yaml:"-" json:"-"`
	Worker              []*machina.Node               `yaml:"-" json:"-"`
	OrchestrationEngine engine.Orchestrator           `yaml:"-" json:"-"`

	KubeadmEngine *engine.KubeadmEngine `yaml:"kubeadm" json:"kubeadm"`
	K3sEngine     *engine.K3SEngine     `yaml:"k3s" json:"k3s"`
}

func (d *RemoveNode) Parse(config []byte) (ClusterOperation, error) {
	var orchestration RemoveNode

	err := yaml.Unmarshal(config, &orchestration)
	if err != nil {
		return nil, errors.New("error while parsing second configuration - " + err.Error())
	}

	return &orchestration, nil
}

func (d *RemoveNode) Validate() error {

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

func (d *RemoveNode) Run(verbose bool) error {

	if err := d.getNodeForRemoving(); err != nil {
		return err
	}

	return d.getOrchestrator().RemoveNode()
}

func (d *RemoveNode) getNodeForRemoving() error {
	switch d.Provider {
	case MultipassProvider:
		resp, err := d.Multipass.GetNodesForDeletion()
		if err != nil {
			return err
		}

		d.Master = resp.Master
		d.Worker = resp.Worker
		return nil

	case BaremetalProvider:
		resp, err := d.BareMetal.GetNodesForDeletion()
		if err != nil {
			return err
		}

		d.Master = resp.Master
		d.Worker = resp.Worker
		return nil

	default:
		return errors.New("no provider found")
	}
}

func (d *RemoveNode) getOrchestrator() engine.Orchestrator {

	switch d.OrchestrationEngine.(type) {

	case *engine.KubeadmEngine:
		{
			var orch *engine.KubeadmEngine

			orch = d.OrchestrationEngine.(*engine.KubeadmEngine)

			var masterNodes []*kubeadmclient.MasterNode
			var workerNodes []*kubeadmclient.WorkerNode

			masterNodes = append(masterNodes, kubeadmclient.NewMasterNode(d.Master.GetUsername(), d.Master.GetIP(), d.Master.GetPrivateKey()))

			for _, worker := range d.Worker {
				workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode(worker.GetUsername(), worker.GetIP(), worker.GetPrivateKey()))
			}

			orch.ClusterName = d.ClusterName
			orch.Masters = masterNodes
			orch.Workers = workerNodes

			return orch
		}

	case *engine.K3SEngine:
		{
			var orch *engine.K3SEngine

			orch = d.OrchestrationEngine.(*engine.K3SEngine)

			var masterNodes []*k3sclient.Master
			var workerNodes []*k3sclient.Worker

			masterNodes = append(masterNodes, k3sclient.NewMaster(d.Master.GetUsername(), d.Master.GetIP(), d.Master.GetPrivateKey()))

			for _, worker := range d.Worker {
				workerNodes = append(workerNodes, k3sclient.NewWorker(worker.GetUsername(), worker.GetIP(), worker.GetPrivateKey()))
			}

			orch.ClusterName = d.ClusterName
			orch.Masters = masterNodes
			orch.Workers = workerNodes

			return orch
		}
	default:
		return nil
	}
}
