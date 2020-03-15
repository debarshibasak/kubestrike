package kubeadmclient

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
)

// CreateCluster will take the incoming Kubeadm struct.
// It will create an HA or NonHA cluster based on incoming Kubeadm struct.
// In non-HA mode it will provision a master, provision workers and add them to the cluster
// In HA mode, it will provision HAProxy, master and workers and then add them to the cluster.
// You can also specify the network plugin, pod cidr range, service cidr range and dns domain that should used in the cluster.
// There parameters are options and default set it picked up on initialization.
func (k *Kubeadm) CreateCluster() error {

	var (
		joinCommand string
		err         error
	)

	if k.ClusterName == "" {
		return errors.New("cluster name is not set")
	}

	err = k.validateAndUpdateDefault()
	if err != nil {
		return err
	}

	startTime := time.Now()

	log.Println("total master - " + fmt.Sprintf("%v", len(k.MasterNodes)))
	log.Println("total workers - " + fmt.Sprintf("%v", len(k.WorkerNodes)))

	if k.HaProxyNode != nil {
		log.Println("total haproxy - " + fmt.Sprintf("%v", 1))
	}

	masterCreationStartTime := time.Now()
	joinCommand, err = k.setupMaster(k.determineSetup())
	if err != nil {
		return err
	}

	log.Printf("time taken to create masters = %v", time.Since(masterCreationStartTime))

	workerCreationTime := time.Now()

	if err := k.setupWorkers(joinCommand); err != nil {
		return err
	}

	log.Printf("time taken to create workers = %v", time.Since(workerCreationTime))

	for _, file := range k.ApplyFiles {
		err := k.MasterNodes[0].applyFile(file)
		if err != nil {
			return err
		}
	}

	if k.Networking != nil {
		log.Printf("installing networking plugin = %v", k.Networking.Name)
		err := k.MasterNodes[0].applyFile(k.Networking.Manifests)
		if err != nil {
			return err
		}
	}

	log.Printf("Time taken to create cluster %v\n", time.Since(startTime).String())

	return nil
}
