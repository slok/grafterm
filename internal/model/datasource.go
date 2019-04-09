package model

// Datasource is where the data will be retrieved.
type Datasource struct {
	ID               string `json:"id,omitempty"`
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
