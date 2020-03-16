package config

import (
	"fmt"
	"testing"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient/networking"
)

func TestCreateCluster_Validate(t *testing.T) {
	p := NewParser(false)

	op, err := p.Parse([]byte(`
apiVersion: github.com/debarshibasak/kubestrike/tree/master/v1alpha1
kind: CreateCluster
provider: Multipass
clusterName: testcluster
multipass:
  masterCount: 1
  workerCount: 1
networking:
  podCidr: 10.233.0.0/16
  plugin: flannel
`))

	if err != nil {
		t.Fatal(err)
	}

	if op == nil {
		t.Fatal("nil ops")
	}

	c := op.(*CreateCluster)

	fmt.Printf("%+v", networking.LookupNetworking(c.Networking.Plugin))
}
