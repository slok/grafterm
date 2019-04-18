package v1

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/configuration/meta"
)

const (
	v1Version = "v1"
)

// Datasource represents a configuration v1 datasource.
type Datasource = model.Datasource

// Dashboard represents a configuration v1 dashboard.
type Dashboard = model.Dashboard

// Configuration is the v1 configuration.Satisfies configuration.Configuration interface.
type Configuration struct {
	meta.Meta   `json:",inline"`
	Datasources []Datasource `json:"datasources,omitempty"`
	Dashboard   Dashboard    `json:"dashboard,omitempty"`
}

// Validate will validate and autocomplete the configuration.
func (c *Configuration) Validate() error {
	if strings.ToLower(c.Meta.Version) != v1Version {
		return fmt.Errorf("not a valid version for V1 configuration: %s", c.Meta.Version)
	}

	err := c.validateDatasources()
	if err != nil {
		return err
	}

	err = c.validateDashboard()
	if err != nil {
		return err
	}

	return nil
}

func (c *Configuration) validateDatasources() error {
	// Check there are multiple datasources with the same ID.
	dss := map[string]struct{}{}
	for _, ds := range c.Datasources {
		_, ok := dss[ds.ID]
		if ok {
			return fmt.Errorf("datasource %s ID registered multiple times", ds.ID)
		}
		dss[ds.ID] = struct{}{}
	}

	return nil
}

func (c *Configuration) validateDashboard() error {
	// Check dashboard.
	for _, row := range c.Dashboard.Rows {
		for _, widget := range row.Widgets {
			switch {
			// Check graphs.
			case widget.Graph != nil:
				for i, seriesOverride := range widget.Graph.Visualization.SeriesOverride {
					// Compile color regexes on override series.
					re, err := regexp.Compile(seriesOverride.Regex)
					if err != nil {
						return err
					}
					seriesOverride.CompiledRegex = re
					widget.Graph.Visualization.SeriesOverride[i] = seriesOverride
				}
			}
		}
	}

	return nil
}

// GetDashboard returns the model dashboard.
func (c *Configuration) GetDashboard() model.Dashboard {
	return c.Dashboard
}

// GetDatasources returns the model datasources.
func (c *Configuration) GetDatasources() []model.Datasource {
	return c.Datasources
}
