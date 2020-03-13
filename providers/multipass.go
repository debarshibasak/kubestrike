package providers

import (
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/debarshibasak/go-multipass/multipass"
)

type Multipass struct {
	Worker int
	Master int
}

func (m *Multipass) Provision() ([]*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, *kubeadmclient.HaProxyNode, error) {

	var (
		masters   []string
		workers   []string
		haproxyIP string

		masterNodes []*kubeadmclient.MasterNode
		workerNodes []*kubeadmclient.WorkerNode
		haproxy     *kubeadmclient.HaProxyNode

		publicKeyLocation  string
		privateKeyLocation string
		err                error
	)

	publicKeyLocation, privateKeyLocation, err = kubeadmclient.PublicKeyExists()
	if err != nil {
		return masterNodes, workerNodes, haproxy, err
	}

	publicKey, err := ioutil.ReadFile(publicKeyLocation)
	if err != nil {
		return masterNodes, workerNodes, haproxy, err
	}

	if m.Master > 1 {
		instance, err := multipass.Launch(&multipass.LaunchReq{
			CPU:  2,
			Name: "haproxy",
		})
		if err != nil {
			return masterNodes, workerNodes, haproxy, err
		}

		err = multipass.Exec(&multipass.ExecRequest{
			Name:    instance.Name,
			Command: "sh -c 'echo " + strings.TrimSpace(string(publicKey)) + " >> /home/ubuntu/.ssh/authorized_keys'",
		})
		if err != nil {
			return masterNodes, workerNodes, haproxy, err
		}

		haproxyIP = instance.IP
	}

	for i := 0; i < m.Master; i++ {
		instance, err := multipass.Launch(&multipass.LaunchReq{
			CPU: 2,
		})
		if err != nil {
			return masterNodes, workerNodes, haproxy, err
		}

		err = multipass.Exec(&multipass.ExecRequest{
			Name:    instance.Name,
			Command: "sh -c 'echo " + strings.TrimSpace(string(publicKey)) + " >> /home/ubuntu/.ssh/authorized_keys'",
		})
		if err != nil {
			return masterNodes, workerNodes, haproxy, err
		}

		masters = append(masters, instance.IP)
	}

	var workerWaitGroup sync.WaitGroup

	for i := 0; i < m.Worker; i++ {

		workerWaitGroup.Add(1)

		go func(workerWaitGroup *sync.WaitGroup) {
			defer workerWaitGroup.Done()

			instance, err := multipass.Launch(&multipass.LaunchReq{
				CPU: 2,
			})
			if err != nil {
				log.Println(err)
			}

			err = multipass.Exec(&multipass.ExecRequest{
				Name:    instance.Name,
				Command: "sh -c 'echo " + strings.TrimSpace(string(publicKey)) + " >> /home/ubuntu/.ssh/authorized_keys'",
			})
			if err != nil {
				log.Println(err)
			}

			workers = append(workers, instance.IP)

		}(&workerWaitGroup)
	}

	workerWaitGroup.Wait()

	for _, master := range masters {
		masterNodes = append(masterNodes, kubeadmclient.NewMasterNode("ubuntu", master, privateKeyLocation))
	}

	if haproxyIP != "" {
		haproxy = kubeadmclient.NewHaProxyNode("ubuntu", haproxyIP, privateKeyLocation)
	}

	for _, worker := range workers {
		workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode("ubuntu", worker, privateKeyLocation))
	}

	return masterNodes, workerNodes, haproxy, nil
}
