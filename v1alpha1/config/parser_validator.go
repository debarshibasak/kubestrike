package config

import (
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/ghodss/yaml"
)

var (
	errKind       = errors.New("kind not supported")
	errAPIKind    = errors.New("api version is not valid")
	errMultipass  = errors.New("provider is set to multipass but configuration is not set")
	errBaremetal  = errors.New("provider is set to multipass but configuration is not set")
	errNetworking = errors.New("networking configurations are not set")
)

type Parser struct {
	useStrictAPIVersionCheck bool
}

func New(useStrictAPIVersionCheck bool) *Parser {

	return &Parser{useStrictAPIVersionCheck: useStrictAPIVersionCheck}
}

func validateAPIVersion(apiVersion string) error {

	req, err := http.NewRequest(http.MethodGet, apiVersion, nil)
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

func (p *Parser) Parse(config []byte) (*ClusterOrchestrator, error) {

	var clusterOrchestrator ClusterOrchestrator
	err := yaml.Unmarshal(config, &clusterOrchestrator)
	if err != nil {
		return nil, err
	}

	if p.useStrictAPIVersionCheck {
		if err := validateAPIVersion(clusterOrchestrator.APIVersion); err != nil {
			return nil, err
		}
	}

	return &clusterOrchestrator, nil
}
