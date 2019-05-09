package configuration

import (
	"io"

	"github.com/slok/grafterm/internal/model"
)

// Configuration is the interface that the different configurations need to implement.
type Configuration interface {
	// Version gets the version of the configuration.
	Version() string
	// Dashboard gets the domain model dashboard from the configuration.
	Dashboard() (model.Dashboard, error)
	// Dashboard gets the domain model datasources from the configuration.
	Datasources() ([]model.Datasource, error)
}

// Loader knows how to load different configuration versions.
type Loader interface {
	// Load loads a configuration.
	Load(r io.Reader) (Configuration, error)
}
