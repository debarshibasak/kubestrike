package config

import (
	"fmt"
	"log"

	"github.com/debarshibasak/go-k3s/k3sclient"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"

	"github.com/debarshibasak/machina"

	"github.com/debarshibasak/kubestrike/v1alpha1/engine"

	"github.com/debarshibasak/kubestrike/v1alpha1/provider"

	"github.com/ghodss/yaml"

	"errors"
)

type CreateCluster struct {
	Base
	Multipass           *provider.MultipassCreateCluster `yaml:"multipass" json:"multipass"`
	BareMetal           *provider.Baremetal              `yaml:"baremetal" json:"baremetal"`
	OrchestrationEngine engine.Orchestrator              `yaml:"-" json:"-"`
	KubeadmEngine       *engine.KubeadmEngine            `yaml:"kubeadm" json:"kubeadm"`
	K3sEngine           *engine.K3SEngine                `yaml:"k3s" json:"k3s"`
	WorkerNodes         []*machina.Node                  `yaml:"-" json:"-"`
	MasterNodes         []*machina.Node                  `yaml:"-" json:"-"`
	HAProxy             *machina.Node                    `yaml:"-" json:"-"`
}

func (c *CreateCluster) Parse(config []byte) (ClusterOperation, error) {

	var createClusterConfiguration CreateCluster

	err := yaml.Unmarshal(config, &createClusterConfiguration)
	if err != nil {
		return nil, errors.New("error while parsing inner configuration")
	}

	if c.Multipass != nil && c.BareMetal != nil {
		return nil, errors.New("only 1 provider is allowed (options are multipass and baremetal)")
	}

	if createClusterConfiguration.KubeadmEngine != nil && createClusterConfiguration.K3sEngine != nil {
		return nil, errors.New("only 1 orchestration engine is allowed")
	}

	if createClusterConfiguration.KubeadmEngine != nil {
		createClusterConfiguration.OrchestrationEngine = createClusterConfiguration.KubeadmEngine
	}

	if createClusterConfiguration.K3sEngine != nil {
		createClusterConfiguration.OrchestrationEngine = createClusterConfiguration.K3sEngine
	}

	return &createClusterConfiguration, nil
}

func (c *CreateCluster) getOrchestrator() engine.Orchestrator {

	switch c.OrchestrationEngine.(type) {

	case *engine.KubeadmEngine:
		{

			var orch *engine.KubeadmEngine

			orch = c.OrchestrationEngine.(*engine.KubeadmEngine)

			var masterNodes []*kubeadmclient.MasterNode
			var workerNodes []*kubeadmclient.WorkerNode
			var haproxy *kubeadmclient.HaProxyNode

			for _, master := range c.MasterNodes {
				masterNodes = append(masterNodes, kubeadmclient.NewMasterNode("ubuntu", master.GetIP(), master.GetPrivateKey()))
			}

			if c.HAProxy != nil {
				haproxy = kubeadmclient.NewHaProxyNode("ubuntu", c.HAProxy.GetIP(), c.HAProxy.GetPrivateKey())
			}

			for _, worker := range c.WorkerNodes {
				workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode("ubuntu", worker.GetIP(), worker.GetPrivateKey()))
			}

			orch.HAProxy = haproxy
			orch.ClusterName = c.ClusterName
			orch.Masters = masterNodes
			orch.Workers = workerNodes

			return orch
		}

	case *engine.K3SEngine:
		{
			var orch *engine.K3SEngine

			orch = c.OrchestrationEngine.(*engine.K3SEngine)

			var masterNodes []*k3sclient.Master
			var workerNodes []*k3sclient.Worker
			var haproxy *k3sclient.HAProxy

			for _, master := range c.MasterNodes {
				masterNodes = append(masterNodes, k3sclient.NewMaster("ubuntu", master.GetIP(), master.GetPrivateKey()))
			}

			if c.HAProxy != nil {
				haproxy = k3sclient.NewHAProxy("ubuntu", c.HAProxy.GetIP(), c.HAProxy.GetPrivateKey())
			}

			for _, worker := range c.WorkerNodes {
				workerNodes = append(workerNodes, k3sclient.NewWorker("ubuntu", worker.GetIP(), worker.GetPrivateKey()))
			}

			orch.HAProxy = haproxy
			orch.ClusterName = c.ClusterName
			orch.Masters = masterNodes
			orch.Workers = workerNodes

			return orch
		}
	default:
		return nil
	}
}

func (c *CreateCluster) Run(verbose bool) error {

	log.Println("[kubestrike] provider found - " + c.Provider)

	err := Get(c)
	if err != nil {
		return err
	}

	log.Println("\n[kubestrike] creating cluster...")

	orchestrator := c.getOrchestrator()

	if orchestrator == nil {
		return errors.New("could not determine the orchestration engine")
	}

	err = orchestrator.CreateCluster()
	if err != nil {
		return err
	}

	fmt.Println("")
	log.Println("[kubestrike] You can access the cluster now")
	fmt.Println("")

	return nil
}

func (c *CreateCluster) Validate() error {

	if c.ClusterName == "" {
		return errClusterNameIsEmpty
	}
	if c.Kind != CreateClusterKind {
		return errKind
	}

	if c.Multipass != nil && c.BareMetal != nil {
		return errors.New("only one provider is allowed")
	}
	return nil
}
