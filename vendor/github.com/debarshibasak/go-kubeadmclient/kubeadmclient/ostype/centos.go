package osType

type Centos struct {
}

func (c *Centos) InstallDocker() []string {
	return []string{
		"sudo yum install -y yum-utils device-mapper-persistent-data lvm2",
		"sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo",
		"sudo yum install docker",
		"sudo systemctl start docker",
		"sudo systemctl enable docker",
	}
}

func (c *Centos) Commands() []string {
	cmds := []string{
		`cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF`,
		"setenforce 0",
		"sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config",
		"yum install -y kubelet kubeadm kubectl --disableexcludes=kubernetes",
		"systemctl enable --now kubelet",
		`cat <<EOF > /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF`,
		"sysctl --system",
		"lsmod | grep br_netfilter",
	}

	cmds = append(cmds, c.InstallDocker()...)

	return cmds
}
