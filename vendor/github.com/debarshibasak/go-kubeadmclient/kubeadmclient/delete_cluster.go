package kubeadmclient

import (
	"errors"
)

func (k *Kubeadm) DeleteCluster() error {

	if len(k.MasterNodes) == 0 {
		return errors.New("no master specified")
	}

	nodelist, err := k.MasterNodes[0].getAllWorkerNodeNames()
	if err != nil {
		return err
	}

	masterNodeList, err := k.MasterNodes[0].getAllMasterNodeNames()
	if err != nil {
		return err
	}

	err = k.validateDeleteCluster(nodelist, masterNodeList)
	if err != nil {
		return err
	}

	if k.ResetOnDeleteCluster {
		err := k.RemoveNode()
		if !k.SkipWorkerFailure {
			if err != nil {
				return err
			}
		}
	} else {
		if err := k.deleteNodes(nodelist); err != nil {
			return err
		}
	}

	if len(masterNodeList) > 0 {

		if err := k.deleteNodes(masterNodeList); err != nil {
			return err
		}

		if k.ResetOnDeleteCluster {
			//TODO parallelize
			for _, master := range k.MasterNodes {
				if err := master.reset(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (k *Kubeadm) validateDeleteCluster(nodelist []string, masterNodeList []string) error {
	if k.ResetOnDeleteCluster && len(k.WorkerNodes) < len(nodelist) {
		return errors.New("will not be able to reset nodes as the nodelist is greater than worker nodes. This hints that some node details are missing")
	}
	if k.ResetOnDeleteCluster && len(k.MasterNodes) < len(masterNodeList) {
		return errors.New("will not be able to reset nodes as the nodelist is greater than master nodes. This hints that some node details are missing")
	}
	if len(masterNodeList) == 0 && len(nodelist) == 0 {
		return errors.New("looks like an empty cluster")
	}
	return nil
}
