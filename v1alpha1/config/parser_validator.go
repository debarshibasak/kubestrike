package config

import (
	"log"

	"errors"

	"github.com/ghodss/yaml"
)

var (
	errKind               = errors.New("kind not supported")
	errAPIKind            = errors.New("api version is not valid")
	errMultipass          = errors.New("provider is set to multipass but configuration is not set")
	errBaremetal          = errors.New("provider is set to multipass but configuration is not set")
	errNetworking         = errors.New("networking configurations are not set")
	errClusterNameIsEmpty = errors.New("cluster name is empty")
)

type Parser struct {
	useStrictAPIVersionCheck bool
}

func NewParser(useStrictAPIVersionCheck bool) *Parser {
	return &Parser{useStrictAPIVersionCheck: useStrictAPIVersionCheck}
}

func validateAPIVersion(apiVersion string) error {

	if apiVersion != "v1" {
		return errors.New("unsupported api version " + apiVersion)
	}

	return nil
}

func (p *Parser) Parse(config []byte) (ClusterOperation, error) {
	var base Base

	err := yaml.Unmarshal(config, &base)
	if err != nil {
		return nil, errors.New("error while parsing first configuration - " + err.Error())
	}

	if p.useStrictAPIVersionCheck {
		if err := validateAPIVersion(base.APIVersion); err != nil {
			return nil, err
		}
	}

	return getOperation(base.Kind).Parse(config)
}

func getOperation(kind Kind) ClusterOperation {
	switch kind {
	case CreateClusterKind:
		return &CreateCluster{}
	case DeleteClusterKind:
		return &DeleteCluster{}
	case AddNodeKind:
		return &AddNode{}
	case RemoveNodeKind:
		return &RemoveNode{}
	default:
		log.Fatal("kind not supported")
		return nil
	}
}
