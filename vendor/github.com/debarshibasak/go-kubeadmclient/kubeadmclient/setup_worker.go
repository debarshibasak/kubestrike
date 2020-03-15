package kubeadmclient

import (
	"github.com/pkg/errors"
)

var errWhileAddWorker = errors.New("error while adding worker")

type workerError struct {
	worker *WorkerNode
	err    error
}

func (k *Kubeadm) setupWorkers(joinCommand string) error {
	errc := make(chan workerError)

	if len(k.WorkerNodes) > 0 {
		for i, workerNode := range k.WorkerNodes {

			go func(node *WorkerNode, i int) {
				err := node.install(joinCommand)
				errc <- workerError{worker: node, err: err}

				if i == len(k.WorkerNodes)-1 {
					close(errc)
				}

			}(workerNode, i)
		}
	}

	return k.workerErrorManager(errc)
}
