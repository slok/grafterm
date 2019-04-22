package v1_test

import (
	"regexp"
	"testing"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/configuration/meta"
	v1 "github.com/slok/grafterm/internal/service/configuration/v1"
	"github.com/stretchr/testify/assert"
)

func getBase() v1.Configuration {
	return v1.Configuration{
		Meta: meta.Meta{Version: "v1"},
		Datasources: []v1.Datasource{
			model.Datasource{
				ID:               "ds1",
				DatasourceSource: model.DatasourceSource{Fake: &model.FakeDatasource{}},
			},
			model.Datasource{
				ID: "ds2",
				DatasourceSource: model.DatasourceSource{Prometheus: &model.PrometheusDatasource{
					Address: "http://127.0.0.1:9090",
				}},
			},
		},
		Dashboard: v1.Dashboard{
			Rows: []model.Row{
				model.Row{
					Title:  "row1",
					Border: true,
					Widgets: []model.Widget{
						model.Widget{
							Title: "widget1",
							WidgetSource: model.WidgetSource{
								Graph: &model.GraphWidgetSource{
									Visualization: model.GraphVisualization{},
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name   string
		cfg    func() v1.Configuration
		exp    func() v1.Configuration
		expErr bool
	}{
		{
			name: "Everything correct.",
			cfg: func() v1.Configuration {
				return getBase()
			},
			exp: func() v1.Configuration {
				return getBase()
			},
		},
		{
			name: "graph visualization regex should autocomplete with the compiled the regex.",
			cfg: func() v1.Configuration {
				base := getBase()
				base.Dashboard.Rows[0].Widgets[0].Graph.Visualization.SeriesOverride = []model.SeriesOverride{
					model.SeriesOverride{
						Regex: ".*",
					},
				}
				return base
			},
			exp: func() v1.Configuration {
				base := getBase()
				base.Dashboard.Rows[0].Widgets[0].Graph.Visualization.SeriesOverride = []model.SeriesOverride{
					model.SeriesOverride{
						Regex:         ".*",
						CompiledRegex: regexp.MustCompile(".*"),
					},
				}
				return base
			},
		},
		{
			name: "Multiple datasource with the same ID should error.",
			cfg: func() v1.Configuration {
				base := getBase()
				base.Datasources = append(base.Datasources, model.Datasource{
					ID:               "ds1",
					DatasourceSource: model.DatasourceSource{Fake: &model.FakeDatasource{}},
				})

				return base
			},
			exp: func() v1.Configuration {
				return getBase()
			},
			expErr: true,
		},
		{
			name: "graph series visualization wrong regex should error.",
			cfg: func() v1.Configuration {
				base := getBase()
				base.Dashboard.Rows[0].Widgets[0].Graph.Visualization.SeriesOverride = []model.SeriesOverride{
					model.SeriesOverride{
						Regex: "8-(",
					},
				}

				return base
			},
			exp: func() v1.Configuration {
				return getBase()
			},
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			exp := test.exp()
			got := test.cfg()
			err := got.Validate()
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(exp, got)
			}
		})
	}
}
