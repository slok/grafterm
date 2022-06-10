package configuration

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/slok/grafterm/internal/service/configuration/meta"
	v1 "github.com/slok/grafterm/internal/service/configuration/v1"
)

// YAMLLoader will load configuration in YAML format.
// It autodetects the version configuration so the user
// doesn't know what version of configuration is loading.
type YAMLLoader struct{}

// Load satisfies configuration.Loader interface.
func (j YAMLLoader) Load(r io.Reader) (Configuration, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	cfg, err := newYAMLConfig(bs)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bs, cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling yaml: %s", err)
	}

	return cfg, nil
}

// newYAMLConfig will get the correct object configuration
// based on the version of the configuration file.
func newYAMLConfig(cfgData []byte) (Configuration, error) {
	cfgVersion := &struct {
		meta.Meta `yaml:",inline"`
	}{}
	err := yaml.Unmarshal(cfgData, cfgVersion)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling yaml: %s", err)
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
