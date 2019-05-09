package model

import "fmt"

// Datasource is where the data will be retrieved.
type Datasource struct {
	ID               string
	DatasourceSource `json:",inline"`
}

// DatasourceSource represents the datasource.
type DatasourceSource struct {
	Fake       *FakeDatasource       `json:"fake,omitempty"`
	Prometheus *PrometheusDatasource `json:"prometheus,omitempty"`
}

// FakeDatasource is the fake datasource.
type FakeDatasource struct{}

// PrometheusDatasource is the Prometheus kind datasource..
type PrometheusDatasource struct {
	Address string `json:"address,omitempty"`
}

// Validate validates the object model is correct.
func (d Datasource) Validate() error {
	if d.ID == "" {
		return fmt.Errorf("datasource ID is required")
	}

	// Check sources.
	var err error
	switch {
	case d.Prometheus != nil:
		err = d.Prometheus.validate()
	case d.Fake != nil:
	default:
		err = fmt.Errorf("declared datasource %s can't be empty", d.ID)
	}
	if err != nil {
		return err
	}

	return nil
}

func (p PrometheusDatasource) validate() error {
	if p.Address == "" {
		return fmt.Errorf("prometheus address can't be empty")
	}

	return nil
}
