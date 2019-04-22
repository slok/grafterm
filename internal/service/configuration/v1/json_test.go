package v1_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/configuration/meta"
	v1 "github.com/slok/grafterm/internal/service/configuration/v1"
)

func TestJSONLoaderLoad(t *testing.T) {
	tests := []struct {
		name       string
		jsonConfig string
		expConfig  *v1.Configuration
		expErr     bool
	}{
		{
			name:       "Invalid JSON should return an error",
			jsonConfig: `{"version": "v1",}`,
			expErr:     true,
		},
		{
			name:       "Invalid JSON version should return an error",
			jsonConfig: `{"version": "v2"}`,
			expErr:     true,
		},
		{
			name: "A correct configuration should load correctly datasources.",
			jsonConfig: `{
	"version": "v1",
	"datasources": [
		{
			"id": "ds1",
			"fake": {}
		},
		{
			"id": "ds2",
			"prometheus": {
				"address": "http://127.0.0.1:9090"
			}
		}
	]
}`,
			expConfig: &v1.Configuration{
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
			},
		},
		{
			name: "A correct configuration should load correctly dashboard.",
			jsonConfig: `{
	"version": "v1",
	"dashboard": {
		"rows": [
			{
				"title": "row1",
				"border": true,
				"widgets": [
					{
						"title": "widget1",
						"gauge": {
							"percentValue": true,
							"max": 60,
							"query": {
								"expr": "testquery"
							},
							"thresholds": [
								{
									"color": "#37872D"
								},
								{
									"color": "#FA6400",
									"startValue": 50
								},
								{
									"color": "#C4162A",
									"startValue": 75
								}
							]
						}
					},
					{
						"title": "widget2",
						"singlestat": {
							"textFormat": "%.02f",
							"query": {
								"expr": "testquery"
							},
							"thresholds": [
								{
									"color": "#FFF000"
								}
							]
						}
					}
				]
			}
		]
	}
}`,
			expConfig: &v1.Configuration{
				Meta: meta.Meta{Version: "v1"},
				Dashboard: v1.Dashboard{
					Rows: []model.Row{
						model.Row{
							Title:  "row1",
							Border: true,
							Widgets: []model.Widget{
								model.Widget{
									Title: "widget1",
									WidgetSource: model.WidgetSource{
										Gauge: &model.GaugeWidgetSource{
											PercentValue: true,
											Max:          60,
											Query: model.Query{
												Expr: "testquery",
											},
											Thresholds: []model.Threshold{
												model.Threshold{
													Color: "#37872D",
												},
												model.Threshold{
													Color:      "#FA6400",
													StartValue: 50,
												},
												model.Threshold{
													Color:      "#C4162A",
													StartValue: 75,
												},
											},
										},
									},
								},
								model.Widget{
									Title: "widget2",
									WidgetSource: model.WidgetSource{
										Singlestat: &model.SinglestatWidgetSource{
											TextFormat: "%.02f",
											Query: model.Query{
												Expr: "testquery",
											},
											Thresholds: []model.Threshold{
												model.Threshold{
													Color: "#FFF000",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var l v1.JSONLoader

			r := strings.NewReader(test.jsonConfig)
			gotcfg, err := l.Load(r)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expConfig, gotcfg)
			}
		})
	}
}
