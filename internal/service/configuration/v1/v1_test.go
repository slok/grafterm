package v1_test

import (
	"io"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/configuration"
)

var (
	goodJSON = `
{
  "version": "v1",
  "datasources": {
  	"gitlab": {
      "prometheus": {
        "address": "https://dashboards.gitlab.com/api/datasources/proxy/6/"
      }
	},
	"ds": {
      "prometheus": {
        "address": "http://127.0.0.1:9090"
      }
    }
  },
  "dashboard": {
    "variables": {
      "env": {
        "constant": {
          "value": "gprd"
        }
      },
      "interval": {
        "interval": {
          "steps": 50
        }
      }
	},
	"widgets": [
		{
        "title": "widget1",
        "gridPos": { "w": 5, "x": 20, "y": 30 },
        "gauge": {
          "query": {
            "datasourceID": "gitlab",
            "expr": "avg_over_time(probe_success{env=\"{{.env}}\",monitor=\"default\",instance=\"https://gitlab.com\", job=\"blackbox-tls-redirect\"}[{{.interval}}])"
		  },
		  "percentValue": true,
		  "min": 10,
		  "max": 20,
          "thresholds": [
            { "color": "#d44a3a", "startValue": 10 },
            { "color": "#2dc937", "startValue": 30 }
          ]
        }
      },
	  {
        "title": "widget2",
        "gridPos": { "w": 10, "x": 10, "y": 10 },
        "singlestat": {
					"unit": "second",
					"decimals": 2,
          "query": {
            "datasourceID": "gitlab",
            "expr": "avg_over_time(probe_success{env=\"{{.env}}\",monitor=\"default\",instance=\"https://gitlab.com\", job=\"blackbox-tls-redirect\"}[{{.interval}}])"
          },
          "valueText": "{{ if (lt .value 1.0) }}DOWN{{else}}UP{{end}}",
          "thresholds": [
            { "color": "#d44a3a" },
            { "color": "#2dc937", "startValue": 1 }
          ]
        }
	  },
	  {
        "title": "widget3",
        "gridPos": { "w": 55, "x": 66, "y": 77 },
		"graph": {
		  "visualization": {
		    "legend": {
		      "disable": true,
		      "rightSide": true
			},
			"seriesOverride": [
              {
                "regex": "p99",
								"color": "#c15c17",
								"nullPointMode": "connected"
              },
              {
                "regex": "p95",
								"color": "#f2c96d",
								"nullPointMode": "null"
              },
              {
                "regex": "p50",
								"color": "#f9ba8f",
								"nullPointMode": "zero"
              }
			],
			"yAxis": {
				"unit": "second",
				"decimals": 1
			}
		  },
		  "queries": [
			{
              "datasourceID": "ds",
              "expr": "max(handler:http_request_duration_seconds_bucket:sum_rate2m_histogram_quantile_perc99)",
              "legend": "p99"
            },
            {
              "datasourceID": "ds",
              "expr": "max(handler:http_request_duration_seconds_bucket:sum_rate2m_histogram_quantile_perc95)",
              "legend": "p95"
            },
            {
              "datasourceID": "ds",
              "expr": "max(handler:http_request_duration_seconds_bucket:sum_rate2m_histogram_quantile_perc50)",
              "legend": "p50"
            }
		  ]
		}
	  }
	]
  }
}`

	goodDashboard = model.Dashboard{
		Grid: model.Grid{
			MaxWidth: 100,
		},
		Variables: []model.Variable{
			{
				Name: "env",
				VariableSource: model.VariableSource{Constant: &model.ConstantVariableSource{
					Value: "gprd",
				}},
			},
			{
				Name: "interval",
				VariableSource: model.VariableSource{Interval: &model.IntervalVariableSource{
					Steps: 50,
				}},
			},
		},
		Widgets: []model.Widget{
			{
				Title:   "widget1",
				GridPos: model.GridPos{W: 5, X: 20, Y: 30},
				WidgetSource: model.WidgetSource{Gauge: &model.GaugeWidgetSource{
					Query: model.Query{
						DatasourceID: "gitlab",
						Expr:         `avg_over_time(probe_success{env="{{.env}}",monitor="default",instance="https://gitlab.com", job="blackbox-tls-redirect"}[{{.interval}}])`,
					},
					PercentValue: true,
					Max:          20,
					Min:          10,
					Thresholds: []model.Threshold{
						{Color: "#d44a3a", StartValue: 10},
						{Color: "#2dc937", StartValue: 30},
					},
				}},
			},
			{
				Title:   "widget2",
				GridPos: model.GridPos{W: 10, X: 10, Y: 10},
				WidgetSource: model.WidgetSource{Singlestat: &model.SinglestatWidgetSource{
					Query: model.Query{
						DatasourceID: "gitlab",
						Expr:         `avg_over_time(probe_success{env="{{.env}}",monitor="default",instance="https://gitlab.com", job="blackbox-tls-redirect"}[{{.interval}}])`,
					},
					ValueText: "{{ if (lt .value 1.0) }}DOWN{{else}}UP{{end}}",
					ValueRepresentation: model.ValueRepresentation{
						Unit:     "second",
						Decimals: 2,
					},
					Thresholds: []model.Threshold{
						{Color: "#d44a3a"},
						{Color: "#2dc937", StartValue: 1},
					},
				}},
			},
			{
				Title:   "widget3",
				GridPos: model.GridPos{W: 55, X: 66, Y: 77},
				WidgetSource: model.WidgetSource{Graph: &model.GraphWidgetSource{
					Visualization: model.GraphVisualization{
						Legend: model.Legend{
							Disable:   true,
							RightSide: true,
						},
						SeriesOverride: []model.SeriesOverride{
							{Regex: "p99", Color: "#c15c17", CompiledRegex: regexp.MustCompile("p99"), NullPointMode: model.NullPointModeConnected},
							{Regex: "p95", Color: "#f2c96d", CompiledRegex: regexp.MustCompile("p95"), NullPointMode: model.NullPointModeAsNull},
							{Regex: "p50", Color: "#f9ba8f", CompiledRegex: regexp.MustCompile("p50"), NullPointMode: model.NullPointModeAsZero},
						},
						YAxis: model.YAxis{
							ValueRepresentation: model.ValueRepresentation{
								Unit:     "second",
								Decimals: 1,
							},
						},
					},
					Queries: []model.Query{
						{
							DatasourceID: "ds",
							Expr:         `max(handler:http_request_duration_seconds_bucket:sum_rate2m_histogram_quantile_perc99)`,
							Legend:       "p99",
						},
						{
							DatasourceID: "ds",
							Expr:         `max(handler:http_request_duration_seconds_bucket:sum_rate2m_histogram_quantile_perc95)`,
							Legend:       "p95",
						},
						{
							DatasourceID: "ds",
							Expr:         `max(handler:http_request_duration_seconds_bucket:sum_rate2m_histogram_quantile_perc50)`,
							Legend:       "p50",
						},
					},
				}},
			},
		},
	}
	goodDatasources = []model.Datasource{
		{
			ID: "ds",
			DatasourceSource: model.DatasourceSource{Prometheus: &model.PrometheusDatasource{
				Address: "http://127.0.0.1:9090",
			}},
		},
		{
			ID: "gitlab",
			DatasourceSource: model.DatasourceSource{Prometheus: &model.PrometheusDatasource{
				Address: "https://dashboards.gitlab.com/api/datasources/proxy/6/",
			}},
		},
	}
)

func TestLoadConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		config         func() io.Reader
		loader         func() configuration.Loader
		expDashboard   model.Dashboard
		expDatasources []model.Datasource
		expErr         bool
	}{
		{
			name: "Invalid JSON should return an error",
			loader: func() configuration.Loader {
				return &configuration.JSONLoader{}
			},
			config: func() io.Reader {
				return strings.NewReader(`{"version": "v1",}`)
			},
			expErr: true,
		},
		{
			name: "Valid JSON should return an correct dashboards",
			loader: func() configuration.Loader {
				return &configuration.JSONLoader{}
			},
			config: func() io.Reader {
				return strings.NewReader(goodJSON)
			},
			expDashboard:   goodDashboard,
			expDatasources: goodDatasources,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			loader := test.loader()
			gotcfg, err := loader.Load(test.config())

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				gotDashboard, err := gotcfg.Dashboard()
				require.NoError(err)
				sort.Slice(gotDashboard.Variables, func(i, j int) bool { return gotDashboard.Variables[i].Name < gotDashboard.Variables[j].Name })
				sort.Slice(gotDashboard.Widgets, func(i, j int) bool { return gotDashboard.Widgets[i].Title < gotDashboard.Widgets[j].Title })
				assert.Equal(test.expDashboard, gotDashboard)

				gotDatasources, err := gotcfg.Datasources()
				// Sort arrays before test.
				sort.Slice(gotDatasources, func(i, j int) bool { return gotDatasources[i].ID < gotDatasources[j].ID })
				require.NoError(err)
				assert.Equal(test.expDatasources, gotDatasources)
			}
		})
	}
}
