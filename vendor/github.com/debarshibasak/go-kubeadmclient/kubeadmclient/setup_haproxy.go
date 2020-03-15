package kubeadmclient

func (k *Kubeadm) setupHAPRoxy() error {
	var masterIP []string
	for _, master := range k.MasterNodes {
		masterIP = append(masterIP, master.ipOrHost)
	}
	return k.HaProxyNode.install(masterIP)
}
