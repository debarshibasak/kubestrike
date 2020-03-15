package kubeadmclient

import (
	"os"
)

type highAvailability struct {
	JoinCommand string
}

func PublicKeyExists() (string, string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	publicKeyLocation := homeDir + "/.ssh/id_rsa.pub"
	privateKeyLocation := homeDir + "/.ssh/id_rsa"
	if _, err := os.Stat(publicKeyLocation); err == nil {
		if _, err := os.Stat(privateKeyLocation); err == nil {
			return publicKeyLocation, privateKeyLocation, nil
		}
		return "", "", err
	}

	return "", "", err
}

//https://www.jordyverbeek.nl/nieuws/kubernetes-ha-cluster-installation-guide
func generateKubeadmConfig(ip string, kubeadm Kubeadm) string {
	return `
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: "1.17.3"
apiServer:
   certSANs:
   - "` + ip + `"
controlPlaneEndpoint: "` + ip + `:6443"
networking:
  podSubnet: ` + kubeadm.PodNetwork + `
  serviceSubnet: ` + kubeadm.ServiceNetwork + `
  clusterName: "` + kubeadm.ClusterName + `"
`
}
