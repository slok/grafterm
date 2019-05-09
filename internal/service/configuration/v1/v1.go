package v1

import (
	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/configuration/meta"
)

const (
	// Version is the that represents the configuration version.
	Version = "v1"
)

// Datasource represents a configuration v1 datasource.
type Datasource = model.Datasource

// Dashboard represents a configuration v1 dashboard.
type Dashboard struct {
	Grid      model.Grid                 `json:"grid,omitempty"`
	Variables map[string]*model.Variable `json:"variables,omitempty"`
	Widgets   []model.Widget             `json:"widgets,omitempty"`
}

// Configuration is the v1 configuration.Satisfies configuration.Configuration interface.
type Configuration struct {
	meta.Meta     `json:",inline"`
	V1Datasources map[string]*Datasource `json:"datasources,omitempty"`
	V1Dashboard   Dashboard              `json:"dashboard,omitempty"`
}

// Version satisfies Configuration interface.
func (c *Configuration) Version() string {
	return Version
}

// Dashboard satisfies Configuration interface.
func (c *Configuration) Dashboard() (model.Dashboard, error) {
	// Transform to model.
	vars := []model.Variable{}
	for name, v := range c.V1Dashboard.Variables {
		v.Name = name
		vars = append(vars, *v)
	}
	dashboard := model.Dashboard{
		Grid:      c.V1Dashboard.Grid,
		Variables: vars,
		Widgets:   c.V1Dashboard.Widgets,
	}

	err := dashboard.Validate()
	if err != nil {
		return dashboard, err
	}

	return dashboard, nil
}

// Datasources satisfies Configuration interface.
func (c *Configuration) Datasources() ([]model.Datasource, error) {
	// Transform to model.
	dss := []model.Datasource{}
	for id, ds := range c.V1Datasources {
		ds.ID = id
		err := ds.Validate()
		if err != nil {
			return dss, err
		}

		dss = append(dss, *ds)
	}

	return dss, nil
}
