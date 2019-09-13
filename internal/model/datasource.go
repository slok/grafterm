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
	Graphite   *GraphiteDatasource   `json:"graphite,omitempty"`
	InfluxDB   *InfluxDBDatasource   `json:"influxdb,omitempty"`
}

// FakeDatasource is the fake datasource.
type FakeDatasource struct{}

// PrometheusDatasource is the Prometheus kind datasource.
type PrometheusDatasource struct {
	Address string `json:"address,omitempty"`
}

// GraphiteDatasource is the Graphite kind datasource.
type GraphiteDatasource struct {
	Address string `json:"address,omitempty"`
}

// InfluxDBDatasource is the Graphite kind datasource.
type InfluxDBDatasource struct {
	Address  string `json:"address,omitempty"`
	Insecure bool   `json:"insecure,omitempty"`
	Database string `json:"database,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
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
	case d.Graphite != nil:
		err = d.Graphite.validate()
	case d.InfluxDB != nil:
		err = d.InfluxDB.validate()
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

func (g GraphiteDatasource) validate() error {
	if g.Address == "" {
		return fmt.Errorf("Graphite API address can't be empty")
	}

	return nil
}

func (g InfluxDBDatasource) validate() error {
	if g.Address == "" {
		return fmt.Errorf("InfluxDB API address can't be empty")
	}

	return nil
}
