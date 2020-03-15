// Package kubeadmclient provides kubernetes cluster operations function.
// You can create clusters, add nodes and remove nodes.
package kubeadmclient

import (
	"errors"
	"log"

	"time"
)

var (
	errMasterNotSpecified = errors.New("master node not specified")
	errWorkerNotSpecified = errors.New("worker not specified")
)

// AddNode will take the incoming Kubeadm struct.
// Fetch the joinCommand from master, provision nodes and add them to the cluster.
func (k *Kubeadm) AddNode() error {

	startTime := time.Now()

	if len(k.MasterNodes) == 0 {
		return errMasterNotSpecified
	}

	if len(k.WorkerNodes) == 0 {
		return errWorkerNotSpecified
	}

	joinCommand, err := k.MasterNodes[0].getJoinCommand()
	if err != nil {
		return err
	}

	if err := k.setupWorkers(joinCommand); err != nil {
		log.Println(err)
		if !k.SkipWorkerFailure {
			return err
		}

		return nil
	}

	log.Println("time taken = " + time.Since(startTime).String())

	return err
}
