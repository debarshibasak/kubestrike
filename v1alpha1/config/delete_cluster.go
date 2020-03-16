package config

type DeleteCluster struct {
	Base Base
}

func (d *DeleteCluster) Run(verbose bool) error {
	return nil
}

func (d *DeleteCluster) Validate() error {
	return nil
}

func (d *DeleteCluster) Parse(config []byte) (ClusterOperation, error) {
	return nil, nil
}
