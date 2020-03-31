package provider

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"strings"
	"sync"
	"time"

	"github.com/debarshibasak/machina"

	"errors"

	"github.com/debarshibasak/go-kubeadmclient/kubeadmclient"
	"github.com/debarshibasak/go-multipass/multipass"
)

type MultipassCreateCluster struct {
	MasterCount int `yaml:"masterCount" json:"masterCount"`
	WorkerCount int `yaml:"workerCount" json:"workerCount"`
}

type MultiPassDeleteCluster struct {
	OnlyKube bool     `yaml:"onlyKube" json:"onlyKube"`
	MasterIP []string `yaml:"master" json:"master"`
	WorkerIP []string `yaml:"workers" json:"workers"`
}

type MultiPassAddNode struct {
	WorkerCount int      `yaml:"workerCount" json:"workerCount"`
	Master      []string `yaml:"master" json:"master"`
}

type MultiPassDeleteNode struct {
	WorkerCount []string `yaml:"workers" json:"workers"`
	Master      []string `yaml:"master" json:"master"`
}

func (node *MultiPassDeleteNode) GetNodesForDeletion() (*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, error) {

	var workers []*kubeadmclient.WorkerNode
	var master *kubeadmclient.MasterNode

	_, pvkey, err := kubeadmclient.PublicKeyExists()
	if err != nil {
		return master, workers, err
	}

	for _, n := range node.WorkerCount {
		workers = append(workers, kubeadmclient.NewWorkerNode("ubuntu", n, pvkey))
	}

	return kubeadmclient.NewMasterNode("ubuntu", node.Master[0], pvkey), workers, nil
}

func (node *MultiPassAddNode) GetNodes() (*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, error) {

	var workers []*kubeadmclient.WorkerNode
	publicKeyLocation, privateKeyLocation, err := kubeadmclient.PublicKeyExists()
	if err != nil {
		return nil, workers, err
	}

	publicKey, err := ioutil.ReadFile(publicKeyLocation)
	if err != nil {
		return nil, workers, err
	}

	done := make(chan struct{})
	defer close(done)
	log.Print("[kubestrike] creating instances...")

	go func() {
		log.Print("[kubestrike] waiting...")
		for {
			select {
			default:
				time.Sleep(1 * time.Second)
				fmt.Print(".")
			case <-done:
				return
			}
		}
	}()

	for i := 0; i < node.WorkerCount; i++ {
		instance, err := multipass.Launch(&multipass.LaunchReq{
			CPU: 2,
		})
		if err != nil {
			return nil, workers, err
		}

		err = multipass.Exec(&multipass.ExecRequest{
			Name:    instance.Name,
			Command: "sh -c 'echo " + strings.TrimSpace(string(publicKey)) + " >> /home/ubuntu/.ssh/authorized_keys'",
		})

		workers = append(workers, kubeadmclient.NewWorkerNode("ubuntu", instance.IP, privateKeyLocation))
	}

	log.Println("[kubestrike] acquired instances")

	return kubeadmclient.NewMasterNode("ubuntu", node.Master[0], privateKeyLocation), workers, nil
}

func (m *MultiPassDeleteCluster) DeleteInstances() ([]*kubeadmclient.MasterNode, []*kubeadmclient.WorkerNode, error) {

	var masterNodes []*kubeadmclient.MasterNode
	var workerNodes []*kubeadmclient.WorkerNode

	if !m.OnlyKube {
		instances, err := multipass.List()
		if err != nil {
			return masterNodes, workerNodes, err
		}

		for _, instance := range instances {
			if err := multipass.Delete(&multipass.DeleteRequest{Name: instance.Name}); err != nil {
				return masterNodes, workerNodes, err
			}
		}
	} else {

		usr, _ := user.Current()

		for _, ip := range m.MasterIP {
			masterNodes = append(masterNodes, kubeadmclient.NewMasterNode("ubuntu", ip, usr.HomeDir+"/.ssh/id_rsa"))
		}

		for _, ip := range m.WorkerIP {
			workerNodes = append(workerNodes, kubeadmclient.NewWorkerNode("ubuntu", ip, usr.HomeDir+"/.ssh/id_rsa"))
		}
	}

	return masterNodes, workerNodes, nil
}

func (m *MultipassCreateCluster) Provision() ([]*machina.Node, []*machina.Node, *machina.Node, error) {

	var (
		masters []string
		workers []string

		masterNodes []*machina.Node
		workerNodes []*machina.Node
		haproxy     *machina.Node

		publicKeyLocation string
		err               error
	)

	publicKeyLocation, _, err = kubeadmclient.PublicKeyExists()
	if err != nil {
		return masterNodes, workerNodes, haproxy,
			errors.New("id_rsa and id_rsa.pub does not exist. Please generate them before you proceed - " + err.Error())
	}

	done := make(chan struct{})
	defer close(done)
	log.Print("[kubestrike] creating vm...")

	go func() {
		log.Print("[kubestrike] waiting...")
		for {
			select {
			default:
				time.Sleep(1 * time.Second)
				fmt.Print(".")
			case <-done:
				return
			}
		}
	}()

	publicKeyLocation, privateKeyLocation, err := kubeadmclient.PublicKeyExists()
	if err != nil {
		return masterNodes, workerNodes, haproxy, err
	}

	publicKey, err := ioutil.ReadFile(publicKeyLocation)
	if err != nil {
		return masterNodes, workerNodes, haproxy, err
	}

	if m.MasterCount > 1 {
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

		haproxy = machina.NewNode("ubuntu", instance.IP, privateKeyLocation)
	}

	for i := 0; i < m.MasterCount; i++ {
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

	for i := 0; i < m.WorkerCount; i++ {

		workerWaitGroup.Add(1)

		go func(workerWaitGroup *sync.WaitGroup) {
			defer workerWaitGroup.Done()

			instance, err := multipass.Launch(&multipass.LaunchReq{
				CPU: 2,
			})
			if err != nil {
				log.Println(err)
			}

			if instance.State == "Stopped" {
				log.Println("instance is stopped")
				return
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
		masterNodes = append(masterNodes, machina.NewNode("ubuntu", master, privateKeyLocation))
	}

	for _, worker := range workers {
		workerNodes = append(workerNodes, machina.NewNode("ubuntu", worker, privateKeyLocation))
	}

	return masterNodes, workerNodes, haproxy, nil
}
