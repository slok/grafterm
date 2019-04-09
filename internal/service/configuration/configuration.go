package configuration

import (
	"io"

	"github.com/slok/meterm/internal/model"
)

// Configuration is the interface that the different configurations need to implement.
type Configuration interface {
	// Validate validates the configuration in case there are configuration errors.
	Validate() error
	// GetDashboard gets the domain model dashboard from the configuration.
	GetDashboard() model.Dashboard
	// GetDashboard gets the domain model datasources from the configuration.
	GetDatasources() []model.Datasource
}

// Loader knows how to load different configuration versions.
type Loader interface {
	// Load loads a configuration.
	Load(r io.Reader) (Configuration, error)
}
