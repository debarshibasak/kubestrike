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

func (createCluster *CreateCluster) Parse(config []byte) (ClusterOperation, error) {

	var createClusterConfiguration CreateCluster

	err := yaml.Unmarshal(config, &createClusterConfiguration)
	if err != nil {
		return nil, errors.New("error while parsing inner configuration")
	}

	if createClusterConfiguration.KubeadmEngine != nil && createClusterConfiguration.K3sEngine != nil {
		return nil, errors.New("only 1 orchestration engine is allowed")
	}

	if createClusterConfiguration.KubeadmEngine != nil {
		createClusterConfiguration.OrchestrationEngine = createClusterConfiguration.KubeadmEngine
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

func (createCluster *CreateCluster) Run(verbose bool) error {

	log.Println("[kubestrike] provider found - " + createCluster.Provider)

	err := Get(createCluster)
	if err != nil {
		return err
	}

	log.Println("\n[kubestrike] creating cluster...")

	orchestrator := createCluster.getOrchestrator()

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

	//if createCluster.Networking == nil {
	//	return errNetworking
	//}
	//
	//if createCluster.Networking != nil && createCluster.Networking.PodCidr != "" && createCluster.Networking.ServiceCidr == "" {
	//	return errNetworking
	//}
	//
	//if createCluster.Networking != nil && createCluster.Networking.PodCidr == "" && createCluster.Networking.ServiceCidr != "" {
	//	return errNetworking
	//}

	return nil
}
