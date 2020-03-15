package networking

type Networking struct {
	Manifests string
	Name      string
}

func LookupNetworking(cni string) *Networking {
	switch cni {
	case "flannel":
		return Flannel
	case "cilium":
		return Cilium
	case "calico":
		return Calico
	case "weavenet":
		return WeaveNet
	default:
		return nil
	}
}

var (
	Flannel = &Networking{
		Manifests: "https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
		Name:      "flannel",
	}

	//Canal = &Networking{
	//	Manifests: "https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml",
	//	Name:      "canal",
	//}

	Cilium = &Networking{
		Manifests: "https://raw.githubusercontent.com/cilium/cilium/v1.6/install/kubernetes/quick-install.yaml",
		Name:      "cilium",
	}

	Calico = &Networking{
		Manifests: "https://docs.projectcalico.org/v3.11/manifests/calico.yaml",
		Name:      "calico",
	}

	WeaveNet = &Networking{
		Manifests: `kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"`,
		Name:      "weavenet",
	}
)
