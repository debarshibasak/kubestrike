package engine

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/networking"
)

type KubeadmEngine struct {
	Networking  Networking                  `yaml:"networking" json:"networking"`
	Masters     []*kubeadmclient.MasterNode `yaml:"-" json:"-"`
	Workers     []*kubeadmclient.WorkerNode `yaml:"-" json:"-"`
	HAProxy     *kubeadmclient.HaProxyNode  `yaml:"-" json:"-"`
	ClusterName string                      `yaml:"-" json:"-"`
}

type Networking struct {
	Plugin      string `yaml:"plugin" json:"plugin"`
	PodCidr     string `yaml:"podCidr" json:"podCidr"`
	ServiceCidr string `yaml:"serviceCidr" json:"serviceCidr"`
}

func (k *KubeadmEngine) CreateCluster() error {

	log.Println("[kubestrike] engine to be used - kubeadm")

	var networkingPlugin networking.Networking

	cni := strings.TrimSpace(k.Networking.Plugin)
	if cni == "" {
		networkingPlugin = *networking.Flannel
	} else {
		v := networking.LookupNetworking(cni)
		networkingPlugin = *v
		if networkingPlugin.Name == "" {
			return errors.New("network plugin in empty")
		}
	}

	kubeadmClient := kubeadmclient.Kubeadm{
		ClusterName:    k.ClusterName,
		HaProxyNode:    k.HAProxy,
		MasterNodes:    k.Masters,
		WorkerNodes:    k.Workers,
		VerboseMode:    false,
		Networking:     &networkingPlugin,
		PodNetwork:     k.Networking.PodCidr,
		ServiceNetwork: k.Networking.ServiceCidr,
	}

	err := kubeadmClient.CreateCluster()
	if err != nil {
		return err
	}

	kubeConfig, err := kubeadmClient.GetKubeConfig()
	if err != nil {
		return err
	}

	u, _ := user.Current()

	kubeconfigLocation := u.HomeDir + "/.kubeconfig_" + k.ClusterName
	if err := ioutil.WriteFile(kubeconfigLocation, []byte(kubeConfig), os.FileMode(0777)); err != nil {
		return err
	}

	log.Println("[kubestrike] You can access the cluster now")
	fmt.Println("")
	fmt.Println("KUBECONFIG=" + kubeconfigLocation + " kubectl get nodes")
	return nil
}
