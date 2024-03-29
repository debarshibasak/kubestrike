package engine

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/debarshibasak/go-k3s/k3sclient/networking"

	"github.com/debarshibasak/go-k3s/k3sclient"
)

type K3SEngine struct {
	Networking     FlannelNetworking   `yaml:"networking" json:"networking"`
	Docker         bool                `yaml:"docker" json:"docker"`
	LoadBalancerIP string              `yaml:"loadbalancerIP" json:"loadbalancerIP"`
	Masters        []*k3sclient.Master `yaml:"-" json:"-"`
	Workers        []*k3sclient.Worker `yaml:"-" json:"-"`
	HAProxy        *k3sclient.HAProxy  `yaml:"-" json:"-"`
	ClusterName    string              `yaml:"-" json:"-"`
}

type FlannelNetworking struct {
	Backend       string `yaml:"backend" json:"backend"`
	PodCidr       string `yaml:"podCidr" json:"podCidr"`
	ServiceCidr   string `yaml:"serviceCidr" json:"serviceCidr"`
	ClusterDomain string `yaml:"clusterDomain" json:"clusterDomain"`
}

func (f *FlannelNetworking) generate() *networking.FlannelOptions {
	return &networking.FlannelOptions{
		Backend:       networking.GetBackend(f.Backend),
		PodCIDR:       f.PodCidr,
		ServiceCIDR:   f.ServiceCidr,
		ClusterDomain: f.ClusterDomain,
	}
}

func (k *K3SEngine) AddNode() error {

	if k.LoadBalancerIP == "" {
		k.LoadBalancerIP = k.Masters[0].GetIP()
	}

	k3sClient := k3sclient.K3sClient{
		ClusterName:    k.ClusterName,
		Master:         k.Masters,
		Worker:         k.Workers,
		UseDocker:      k.Docker,
		LoadBalancerIP: k.LoadBalancerIP,
	}

	return k3sClient.AddNode()
}

func (k *K3SEngine) CreateCluster() error {

	log.Println("[kubestrike] engine to be used - k3s")

	k3Client := k3sclient.K3sClient{
		ClusterName:    k.ClusterName,
		HAProxy:        k.HAProxy,
		Master:         k.Masters,
		Worker:         k.Workers,
		UseDocker:      k.Docker,
		FlannelOptions: k.Networking.generate(),
	}

	err := k3Client.CreateCluster()
	if err != nil {
		return err
	}

	kubeConfig, err := k3Client.GetKubeConfig()
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

func (k *K3SEngine) RemoveNode() error {

	k3sClient := k3sclient.K3sClient{
		Master: k.Masters,
		Worker: k.Workers,
	}

	return k3sClient.DeleteNode()
}

func (k *K3SEngine) DeleteCluster() error {

	k3sClient := k3sclient.K3sClient{
		Master: k.Masters,
		Worker: k.Workers,
	}

	return k3sClient.DeleteCluster()
}
