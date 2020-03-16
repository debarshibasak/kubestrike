package config

type Kind string

const (
	CreateClusterKind Kind = "CreateCluster"
	AddNodeKind       Kind = "AddNode"
	RemoveNodeKind    Kind = "RemoveNode"
	DeleteClusterKind Kind = "DeleteCluster"
)

type Base struct {
	APIVersion  string   `yaml:"apiVersion" json:"apiVersion"`
	Kind        Kind     `yaml:"kind" json:"kind"`
	Provider    Provider `yaml:"provider" json:"provider"`
	ClusterName string   `yaml:"clusterName" json:"clusterName"`
}

type ClusterOperation interface {
	Run(verbose bool) error
	Validate() error
	Parse(config []byte) (ClusterOperation, error)
}
