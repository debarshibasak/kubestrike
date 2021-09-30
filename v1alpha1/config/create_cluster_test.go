package config_test

import (
	"fmt"
	"testing"

	"github.com/debarshibasak/kubestrike/v1alpha1/config"
	"github.com/ghodss/yaml"
)

func TestParsing(t *testing.T) {

	kubeadm := `
apiVersion: kubestrike.debarshi.github.com/master/v1alpha1
kind: CreateClusterKind
provider: Multipass
multipass: 
  masterCount: 1
  workerCount: 1
kubeadm: 
  networking: 
    plugin: flannel
    podCidr: 10.233.0.0/16
`

	var a config.CreateCluster
	err := yaml.Unmarshal([]byte(kubeadm), &a)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(a.APIVersion)
	fmt.Println(a.KubeadmEngine)
	fmt.Println(a.ClusterName)

	k3sclient := `
apiVersion: kubestrike.debarshi.github.com/master/v1alpha1
kind: CreateClusterKind
provider: Multipass
multipass: 
  masterCount: 1
  workerCount: 1
k3s: 
  networking: 
    backend: vxlan
    podCidr: 10.233.0.0/16
`

	err = yaml.Unmarshal([]byte(k3sclient), &a)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(a.APIVersion)
	fmt.Println(a.KubeadmEngine)
	fmt.Println(a.K3sEngine)
	fmt.Println(a.ClusterName)

}
