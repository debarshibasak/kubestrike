package config_test

import (
	"fmt"
	"testing"

	"github.com/debarshibasak/kubestrike/v1alpha1/config"
	"github.com/ghodss/yaml"
)

func TestParsing(t *testing.T) {

	data := `
apiVersion: kubestrike.debarshi.github.com/master/v1alpha1
kind: ClusterOrchestration
provider: Multipass
multipass:
  masterCount: 1
  workerCount: 1
networking:
  podCidr: 10.233.0.0/16
  plugin: flannel`

	var a config.ClusterOrchestrator
	err := yaml.Unmarshal([]byte(data), &a)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(a.ApiVersion)

}
