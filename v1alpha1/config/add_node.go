package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/debarshibasak/go-k3s/k3sclient"
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"

	"github.com/debarshibasak/machina"

	"github.com/debarshibasak/kubestrike/v1alpha1/engine"

	"github.com/debarshibasak/kubestrike/v1alpha1/provider"

	"github.com/ghodss/yaml"
)

type AddNode struct {
	Base
	Multipass           *provider.MultiPassAddNode `yaml:"multipass" json:"multipass"`
	BareMetal           *provider.BaremetalAddNode `yaml:"baremetal" json:"baremetal"`
	SkipWorkerFailure   bool                       `yaml:"skip_worker_failure" json:"skip_worker_failure"`
	OrchestrationEngine engine.Orchestrator        `yaml:"-" json:"-"`
	KubeadmEngine       *engine.KubeadmEngine      `yaml:"kubeadm" json:"kubeadm"`
	K3sEngine           *engine.K3SEngine          `yaml:"k3s" json:"k3s"`
	WorkerNodes         []*machina.Node            `yaml:"-" json:"-"`
	MasterNodes         *machina.Node              `yaml:"-" json:"-"`
	HAProxy             *machina.Node              `yaml:"-" json:"-"`
}

func (a *AddNode) Parse(config []byte) (ClusterOperation, error) {
	var orchestration AddNode

	err := yaml.Unmarshal(config, &orchestration)
	if err != nil {
		return nil, errors.New("error while parsing configuration - " + err.Error())
	}

	if a.Multipass != nil && a.BareMetal != nil {
		return nil, errors.New("only 1 provider is allowed (options are multipass and baremetal)")
	}

	if orchestration.KubeadmEngine != nil && orchestration.K3sEngine != nil {
		return nil, errors.New("only 1 orchestration engine is allowed")
	}

	if orchestration.KubeadmEngine != nil {
		orchestration.OrchestrationEngine = orchestration.KubeadmEngine
	}

	if orchestration.K3sEngine != nil {
		orchestration.OrchestrationEngine = orchestration.K3sEngine
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

//Nice feature would be to Delete the acquired VM on failure
func (a *AddNode) Run(verbose bool) error {

	err := a.getNodes()
	if err != nil {
		return err
	}

	orch := a.getOrchestrator()

	if orch == nil {
		return errors.New("no orchestrator found")
	}

	if err := orch.AddNode(); err != nil {
		log.Fatal("[kubestring] err=" + err.Error())
	}

	log.Println("[kubestrke] nodes added")
	fmt.Println("")
	fmt.Println("------")
	for _, node := range a.WorkerNodes {
		fmt.Println(node.GetIP())
	}
	fmt.Println("------")
	fmt.Println("")

	return nil
}

func (a *AddNode) getOrchestrator() engine.Orchestrator {

	switch a.OrchestrationEngine.(type) {

	case *engine.KubeadmEngine:
		{
			var orch *engine.KubeadmEngine

			orch = a.OrchestrationEngine.(*engine.KubeadmEngine)

			var masterNodes []*kubeadmclient.MasterNode
			var workerNodes []*kubeadmclient.WorkerNode

			masterNodes = append(masterNodes, kubeadmclient.NewMasterNode("ubuntu", a.MasterNodes.GetIP(), a.MasterNodes.GetPrivateKey()))

			for _, worker := range a.WorkerNodes {
				workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode("ubuntu", worker.GetIP(), worker.GetPrivateKey()))
			}

			orch.ClusterName = a.ClusterName
			orch.Masters = masterNodes
			orch.Workers = workerNodes

			return orch
		}

	case *engine.K3SEngine:
		{
			var orch *engine.K3SEngine

			orch = a.OrchestrationEngine.(*engine.K3SEngine)

			var masterNodes []*k3sclient.Master
			var workerNodes []*k3sclient.Worker

			masterNodes = append(masterNodes, k3sclient.NewMaster("ubuntu", a.MasterNodes.GetIP(), a.MasterNodes.GetPrivateKey()))

			for _, worker := range a.WorkerNodes {
				workerNodes = append(workerNodes, k3sclient.NewWorker("ubuntu", worker.GetIP(), worker.GetPrivateKey()))
			}

			orch.ClusterName = a.ClusterName
			orch.Masters = masterNodes
			orch.Workers = workerNodes

			return orch
		}
	default:
		return nil
	}
}

func (a *AddNode) getNodes() error {

	if a.Multipass != nil {

		addResponse, err := a.Multipass.GetNodes()
		if err != nil {
			return err
		}
		a.MasterNodes = addResponse.Master
		a.WorkerNodes = addResponse.Worker

		return nil
	} else if a.BareMetal != nil {
		addResponse, err := a.BareMetal.GetNodes()
		if err != nil {
			return err
		}
		a.MasterNodes = addResponse.Master
		a.WorkerNodes = addResponse.Worker

		return nil
	} else {
		return errors.New("no provider found")
	}
}
