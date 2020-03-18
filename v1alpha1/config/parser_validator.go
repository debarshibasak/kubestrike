package config

import (
	"log"
	"net/http"
	"time"

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

	log.Println("https://" + apiVersion)
	req, err := http.NewRequest(http.MethodGet, "https://"+apiVersion, nil)
	if err != nil {
		return err
	}

	var client http.Client

	client.Timeout = 10 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errAPIKind
	}

	return nil
}

func (p *Parser) Parse(config []byte) (ClusterOperation, error) {
	var base Base

	err := yaml.Unmarshal(config, &base)
	if err != nil {
		return nil, errors.New("error while parsing configuration")
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
		return &DeleteNode{}
	default:
		log.Fatal("kind not supported")
		return nil
	}
}
