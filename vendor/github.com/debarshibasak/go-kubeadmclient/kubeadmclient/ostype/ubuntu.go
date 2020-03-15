package osType

type Ubuntu struct {
}

func (u *Ubuntu) Commands() []string {
	var cmds []string

	cmds = append(cmds, u.InstallDocker()...)
	cmds = append(cmds, u.InstallKubernetes()...)
	cmds = append(cmds, u.TurnOffSwaps()...)

	return cmds
}

func (u *Ubuntu) InstallDocker() []string {
	return []string{
		"sudo apt-get update",
		"sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
		`sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"`,
		"sudo apt update",
		"apt-cache policy docker-ce",
		"sudo apt install docker-ce -y",
	}
}

func (u *Ubuntu) InstallKubernetes() []string {
	return []string{
		"sudo apt-get update",
		"sudo apt-get install -y iptables arptables ebtables",
		//"sudo update-alternatives --set iptables /usr/sbin/iptables-legacy",
		//"sudo update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy",
		//"sudo update-alternatives --set arptables /usr/sbin/arptables-legacy",
		//"sudo update-alternatives --set ebtables /usr/sbin/ebtables-legacy",
		"sudo apt-get update && sudo apt-get install -y apt-transport-https curl",
		"curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -",
		`cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
`,
		"sudo apt-get update",
		"sudo apt-get install -y kubelet kubeadm kubectl",
		"sudo apt-mark hold kubelet kubeadm kubectl",
	}
}

func (u *Ubuntu) TurnOffSwaps() []string {
	return []string{
		"sudo swapoff -a",
		`sudo sed -i "/ swap / s/^\(.*\)$/#\1/g" /etc/fstab`,
		"sudo sysctl net.bridge.bridge-nf-call-iptables=1",
	}
}
