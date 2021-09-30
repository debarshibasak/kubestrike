package config

import (
	"errors"
	"log"

	"github.com/debarshibasak/go-k3s/k3sclient"
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"

	"github.com/debarshibasak/kubestrike/v1alpha1/engine"
	"github.com/debarshibasak/kubestrike/v1alpha1/provider"
	"github.com/debarshibasak/machina"

	"github.com/ghodss/yaml"
)

type DeleteCluster struct {
	Base
	Multipass *provider.MultiPassDeleteCluster `yaml:"multipass" json:"multipass"`
	BareMetal *provider.BaremetalDeleteCluster `yaml:"baremetal" json:"baremetal"`

	Master              []*machina.Node     `yaml:"-" json:"-"`
	Worker              []*machina.Node     `yaml:"-" json:"-"`
	OrchestrationEngine engine.Orchestrator `yaml:"-" json:"-"`

	KubeadmEngine *engine.KubeadmEngine `yaml:"kubeadm" json:"kubeadm"`
	K3sEngine     *engine.K3SEngine     `yaml:"k3s" json:"k3s"`
}

func (orchestrator *DeleteCluster) Run(verbose bool) error {

	log.Println("[kubestrike] provider found - " + orchestrator.Base.Provider)

	err := orchestrator.GetDeleteCluster()
	if err != nil {
		return err
	}

	return orchestrator.getOrchestrator().DeleteCluster()
}

func (orchestrator *DeleteCluster) Validate() error {

	if orchestrator.Multipass != nil && orchestrator.Multipass.OnlyKube && len(orchestrator.Multipass.MasterIP) == 0 {
		return errors.New("no master specified")
	}

	return nil
}

func (orchestrator *DeleteCluster) Parse(config []byte) (ClusterOperation, error) {

	var orchestration DeleteCluster

	err := yaml.Unmarshal(config, &orchestration)
	if err != nil {
		return nil, errors.New("error while parsing configuration")
	}

	return &orchestration, nil
}

func (orchestrator *DeleteCluster) getOrchestrator() engine.Orchestrator {
	switch orchestrator.OrchestrationEngine.(type) {

	case *engine.KubeadmEngine:
		{
			var orch *engine.KubeadmEngine

			orch = orchestrator.OrchestrationEngine.(*engine.KubeadmEngine)

			var masterNodes []*kubeadmclient.MasterNode
			var workerNodes []*kubeadmclient.WorkerNode

			for _, master := range orchestrator.Master {
				masterNodes = append(masterNodes, kubeadmclient.NewMasterNode(master.GetUsername(), master.GetIP(), master.GetPrivateKey()))
			}

			for _, worker := range orchestrator.Worker {
				workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode(worker.GetUsername(), worker.GetIP(), worker.GetPrivateKey()))
			}

			orch.ClusterName = orchestrator.ClusterName
			orch.Masters = masterNodes
			orch.Workers = workerNodes

			return orch
		}

	case *engine.K3SEngine:
		{
			var orch *engine.K3SEngine

			orch = orchestrator.OrchestrationEngine.(*engine.K3SEngine)

			var masterNodes []*k3sclient.Master
			var workerNodes []*k3sclient.Worker

			for _, master := range orchestrator.Master {
				masterNodes = append(masterNodes, k3sclient.NewMaster(master.GetUsername(), master.GetIP(), master.GetPrivateKey()))
			}

			for _, worker := range orchestrator.Worker {
				workerNodes = append(workerNodes, k3sclient.NewWorker(worker.GetUsername(), worker.GetIP(), worker.GetPrivateKey()))
			}

			orch.ClusterName = orchestrator.ClusterName
			orch.Masters = masterNodes
			orch.Workers = workerNodes

			return orch
		}
	default:
		return nil
	}
}
