package configuration

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/slok/grafterm/internal/service/configuration/meta"
	v1 "github.com/slok/grafterm/internal/service/configuration/v1"
)

// JSONLoader will load configuration in JSON format.
// It autodetects the version configuration so the user
// doesn't know what version of configuration is loading.
type JSONLoader struct{}

// Load satisfies configuration.Loader interface.
func (j JSONLoader) Load(r io.Reader) (Configuration, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	cfg, err := newConfig(bs)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %s", err)
	}

	return cfg, nil
}

// newConfig will get the correct object configuration
// based on the version of the configuration file.
func newConfig(cfgData []byte) (Configuration, error) {
	cfgVersion := &struct {
		meta.Meta
	}{}
	err := json.Unmarshal(cfgData, cfgVersion)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %s", err)
	}

	var cfg Configuration
	switch cfgVersion.Version {
	case v1.Version:
		cfg = &v1.Configuration{}
	default:
		return nil, fmt.Errorf("%s is not a valid configuration version", cfgVersion.Version)
	}

	return cfg, nil
}
