package kubeadmclient

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"errors"

	"github.com/debarshibasak/go-kubeadmclient/sshclient"
	"github.com/google/uuid"
)

type MasterNode struct {
	*Node
}

func NewMasterNode(username string, ipOrHost string, privateKeyLocation string) *MasterNode {
	return &MasterNode{
		&Node{
			username:           username,
			ipOrHost:           ipOrHost,
			privateKeyLocation: privateKeyLocation,
			clientID:           uuid.New().String(),
		},
	}
}

func (n *MasterNode) getAllWorkerNodeNames() ([]string, error) {
	var hostnames []string

	out, err := n.sshClientWithTimeout(1 * time.Minute).Collect(`sudo KUBECONFIG=/etc/kubernetes/admin.conf kubectl get nodes --selector="node-type=worker"`)
	if err != nil {
		return hostnames, errors.New("error while fetching list")
	}

	for _, hostname := range strings.Split(strings.TrimSpace(out), "\n") {
		trueHostName := strings.TrimSpace(strings.Split(hostname, " ")[0])
		if trueHostName != "NAME" && trueHostName != "" {
			hostnames = append(hostnames, trueHostName)
		}
	}

	return hostnames, nil
}

func (n *MasterNode) getMaster() error {
	return n.run("sudo KUBECONFIG=/etc/kubernetes/admin.conf kubectl get nodes --selector='node-role.kubernetes.io/master='")
}

func (n *MasterNode) taintAsMaster() error {
	return n.run("sudo KUBECONFIG=/etc/kubernetes/admin.conf kubectl taint nodes --selector=kubernetes.io/hostname=`hostname` node-role.kubernetes.io/master-")
}

func (n *MasterNode) applyFile(file string) error {
	return n.run("sudo KUBECONFIG=/etc/kubernetes/admin.conf kubectl apply -f " + file)
}

func (n *MasterNode) getToken() (string, error) {

	sh := sshclient.SSHConnection{
		Username:    n.username,
		IP:          n.ipOrHost,
		KeyLocation: n.privateKeyLocation,
	}

	out, err := sh.Collect("sudo kubeadm token list -o json")
	if err != nil {
		return "", err
	}

	c := make(map[string]interface{})

	err = json.Unmarshal([]byte(out), &c)
	if err != nil {
		return "", err
	}

	return c["token"].(string), nil
}

func (n *MasterNode) run(shell string) error {
	return n.sshClient().Run([]string{shell})
}

func (n *MasterNode) ctlCommand(cmd string) error {
	return n.run("sudo KUBECONFIG=/etc/kubernetes/admin.conf " + cmd)
}

func (n *MasterNode) ctlCommandCollect(cmd string) (string, error) {
	return n.sshClient().Collect("sudo KUBECONFIG=/etc/kubernetes/admin.conf " + cmd)
}

func (n *MasterNode) getKubeConfig() (string, error) {
	return n.sshClient().Collect("sudo cat /etc/kubernetes/admin.conf")
}

type IPHost struct {
	IP   string
	Host string
}

func (n *MasterNode) getAllMasterNodeNames() ([]string, error) {
	var hostnames []string
	out, err := n.sshClientWithTimeout(1 * time.Minute).Collect(`sudo KUBECONFIG=/etc/kubernetes/admin.conf kubectl get nodes --selector="node-role.kubernetes.io/master="`)
	if err != nil {
		return hostnames, errors.New("error while fetching list")
	}

	for _, hostname := range strings.Split(strings.TrimSpace(out), "\n") {
		trueHostName := strings.TrimSpace(strings.Split(hostname, " ")[0])
		if trueHostName != "NAME" && trueHostName != "" {
			hostnames = append(hostnames, trueHostName)
		}
	}

	return hostnames, nil
}

func (n *MasterNode) getJoinCommand() (string, error) {
	return n.sshClient().Collect("sudo kubeadm token create --print-join-command")
}

func (n *MasterNode) installAndFetchCommand(kubeadm Kubeadm, vip string) (string, error) {

	osType, err := n.determineOS()
	if err != nil {
		return "", err
	}

	err = n.sshClient().Run(osType.Commands())
	if err != nil {
		return "", err
	}

	err = n.sshClient().ScpToWithData([]byte(generateKubeadmConfig(vip, kubeadm)), "/tmp/kubeadm-config.yaml")
	if err != nil {
		return "", err
	}

	out, err := n.sshClientWithTimeout(30 * time.Minute).Collect("sudo kubeadm init --config /tmp/kubeadm-config.yaml --upload-certs")
	if err != nil {
		log.Println(out)
		return "", err
	}

	return getControlPlaneJoinCommand(out), nil
}

func (n *MasterNode) install(kubeadm Kubeadm, availability *highAvailability) error {

	osType, err := n.determineOS()
	if err != nil {
		return err
	}

	err = n.sshClientWithTimeout(30 * time.Minute).Run(osType.Commands())
	if err != nil {
		return err
	}

	var s string

	if availability != nil {
		s = "sudo " + availability.JoinCommand
	} else {
		s = "sudo kubeadm init --pod-network-cidr=" + kubeadm.PodNetwork + " --service-cidr=" + kubeadm.ServiceNetwork + " --service-dns-domain=" + kubeadm.DNSDomain
	}

	return n.sshClientWithTimeout(30 * time.Minute).Run([]string{s})
}

func (n *MasterNode) deleteNode(hostname string) error {
	return n.ctlCommand("kubectl delete node " + hostname)
}

func (n *MasterNode) reset() error {
	return n.ctlCommand("kubeadm reset -f")
}

func getControlPlaneJoinCommand(data string) string {
	var cmd string

	for _, line := range strings.Split(data, "\n") {

		if strings.HasPrefix(strings.TrimSpace(line), "kubeadm") {
			cmd = cmd + strings.ReplaceAll(line, "\\", "")
		}

		if strings.HasPrefix(strings.TrimSpace(line), "--discovery") {
			cmd = cmd + strings.ReplaceAll(line, "\\", "")
		}

		if strings.HasPrefix(strings.TrimSpace(line), "--control-plane") {
			cmd = cmd + strings.ReplaceAll(line, "\\", "")
			return cmd
		}
	}

	return cmd
}
