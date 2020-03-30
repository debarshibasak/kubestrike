package engine

type Orchestrator interface {
	CreateCluster() error
}
