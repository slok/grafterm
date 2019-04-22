package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/slok/grafterm/internal/service/configuration"
)

// JSONLoader will load JSON V1 configuration.
type JSONLoader struct{}

// Load satisfies configuration.Loader interface.
func (j JSONLoader) Load(r io.Reader) (configuration.Configuration, error) {
	cfg := &Configuration{}

	// load cfg.
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bs, cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %s", err)
	}

	// Validate cfg.
	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating the loaded configuration: %s", err)
	}

	return cfg, nil
}
