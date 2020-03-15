package kubeadmclient

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type HaProxyNode struct {
	*Node
}

func NewHaProxyNode(username string, ipOrHost string, privateKey string) *HaProxyNode {
	return &HaProxyNode{&Node{
		username:           username,
		ipOrHost:           ipOrHost,
		privateKeyLocation: privateKey,
		verboseMode:        false,
		clientID:           uuid.New().String(),
	}}
}

func (n *HaProxyNode) generateConfig(masterIPs []string) string {

	var serverCheckBlock string
	for i, ip := range masterIPs {
		serverCheckBlock = serverCheckBlock + fmt.Sprintf("  server k8s-api-%v %v:6443 check\n", i, ip)
	}

	return `
frontend k8s-api
  log-format %hr\ %ST\ %B\ %Ts
  bind ` + n.ipOrHost + `:6443
  bind 127.0.0.1:6443
  mode tcp
  timeout client 10m
  option tcplog
  default_backend k8s-api

backend k8s-api
  mode tcp
  option tcplog
  option tcp-check
  balance roundrobin
  timeout connect  30s
  timeout server  10m
  default-server inter 10s downinter 5s rise 2 fall 2 slowstart 60s maxconn 250 maxqueue 256 weight 100
` + serverCheckBlock
}

func (n *HaProxyNode) install(masterIPs []string) error {

	osType, err := n.determineOS()
	if err != nil {
		return err
	}

	if err := n.sshClientWithTimeout(20 * time.Minute).Run(osType.InstallDocker()); err != nil {
		return err
	}

	if err := n.sshClient().ScpToWithData([]byte(n.generateConfig(masterIPs)), "/tmp/haproxy.cfg"); err != nil {
		return err
	}

	if err := n.sshClient().Run([]string{
		"sudo mkdir -p /usr/local/etc/haproxy/",
		"sudo chown $USER:$USER /usr/local/etc/haproxy/",
		"sudo cp /tmp/haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg",
	}); err != nil {
		return err
	}

	_, err = n.sshClientWithTimeout(20 * time.Minute).Collect("sudo docker run -d --network=host -it --restart=always -v /usr/local/etc/haproxy/:/usr/local/etc/haproxy/:ro --name kubernetes-apiserver-haproxy haproxy:1.7")
	if err != nil {
		return err
	}

	return nil

}
