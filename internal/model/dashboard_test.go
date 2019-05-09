package model_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/grafterm/internal/model"
)

func getBaseDashboard() model.Dashboard {
	return model.Dashboard{
		Grid: model.Grid{
			MaxWidth: 50,
		},
		Variables: []model.Variable{
			model.Variable{
				Name: "test-constant",
				VariableSource: model.VariableSource{Constant: &model.ConstantVariableSource{
					Value: "test-value",
				}},
			},
			model.Variable{
				Name: "test-interval",
				VariableSource: model.VariableSource{Interval: &model.IntervalVariableSource{
					Steps: 50,
				}},
			},
		},
		Widgets: []model.Widget{
			model.Widget{
				Title:   "test-gauge",
				GridPos: model.GridPos{W: 10, Y: 10, X: 10},
				WidgetSource: model.WidgetSource{Gauge: &model.GaugeWidgetSource{
					Query: model.Query{
						Expr:         "query",
						Legend:       "test",
						DatasourceID: "test",
					},
					Thresholds: []model.Threshold{
						model.Threshold{Color: "#FFFFFF"},
						model.Threshold{Color: "#FFF000", StartValue: 10},
						model.Threshold{Color: "#000FFF", StartValue: 50},
					},
				}},
			},
			model.Widget{
				Title:   "test-singlestat",
				GridPos: model.GridPos{W: 10, Y: 10, X: 10},
				WidgetSource: model.WidgetSource{Singlestat: &model.SinglestatWidgetSource{
					ValueText: "test",
					Query: model.Query{
						Expr:         "query",
						Legend:       "test",
						DatasourceID: "test",
					},
					Thresholds: []model.Threshold{
						model.Threshold{Color: "#FFFFFF"},
						model.Threshold{Color: "#FFF000", StartValue: 10},
						model.Threshold{Color: "#000FFF", StartValue: 50},
					},
				}},
			},
			model.Widget{
				Title:   "test-graph",
				GridPos: model.GridPos{W: 10, Y: 10, X: 10},
				WidgetSource: model.WidgetSource{Graph: &model.GraphWidgetSource{
					Queries: []model.Query{
						model.Query{Expr: "query", Legend: "test", DatasourceID: "test"},
						model.Query{Expr: "query2", Legend: "test2", DatasourceID: "test"},
						model.Query{Expr: "query3", Legend: "test3", DatasourceID: "test"},
					},
					Visualization: model.GraphVisualization{
						SeriesOverride: []model.SeriesOverride{
							model.SeriesOverride{Regex: "2..", Color: "#FFF000", CompiledRegex: regexp.MustCompile("2..")},
							model.SeriesOverride{Regex: "3..", Color: "#FFF001", CompiledRegex: regexp.MustCompile("3..")},
							model.SeriesOverride{Regex: "4..", Color: "#FFF002", CompiledRegex: regexp.MustCompile("4..")},
							model.SeriesOverride{Regex: "5..", Color: "#FFF003", CompiledRegex: regexp.MustCompile("5..")},
						},
					},
				}},
			},
		},
	}
}

func TestValidateDashboard(t *testing.T) {
	tests := []struct {
		name         string
		dashboard    func() model.Dashboard
		expDashboard func() model.Dashboard
		expErr       bool
	}{
		// Grid.
		{
			name: "Grid MaxWidth should set a default if doesn't have a value.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				d.Grid.MaxWidth = 0
				return d
			},
			expDashboard: func() model.Dashboard {
				d := getBaseDashboard()
				d.Grid.MaxWidth = 100
				return d
			},
		},

		// Variables.
		{
			name: "Variables should have a name.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				v := d.Variables[0]
				v.Name = ""
				d.Variables[0] = v
				return d
			},
			expErr: true,
		},
		{
			name: "Constant variables should have a value.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				d.Variables[0] = model.Variable{
					Name: "test",
					VariableSource: model.VariableSource{Constant: &model.ConstantVariableSource{
						Value: "",
					}},
				}
				return d
			},
			expErr: true,
		},
		{
			name: "variable without type should fail.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				d.Variables[0] = model.Variable{Name: "test"}
				return d
			},
			expErr: true,
		},
		{
			name: "Interval variables should have a valid step.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				d.Variables[0] = model.Variable{
					Name: "test",
					VariableSource: model.VariableSource{Interval: &model.IntervalVariableSource{
						Steps: 0,
					}},
				}
				return d
			},
			expErr: true,
		},

		// Widgets.
		{
			name: "A widget grid position width is required.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[0]
				w.GridPos.W = 0
				d.Widgets[0] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A widget grid position with fixed grid, Y is required.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				d.Grid.FixedWidgets = true
				w := d.Widgets[0]
				w.GridPos.X = 0
				d.Widgets[0] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A widget grid position with fixed grid, Y is required.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				d.Grid.FixedWidgets = true
				w := d.Widgets[0]
				w.GridPos.Y = 0
				d.Widgets[0] = w
				return d
			},
			expErr: true,
		},

		// Gauge widget.
		{
			name: "A gauge widget with a query should have an expression.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[0]
				w.Gauge.Query.Expr = ""
				d.Widgets[0] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A gauge widget with a query should have an datasource ID.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[0]
				w.Gauge.Query.DatasourceID = ""
				d.Widgets[0] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A gauge percent widget max should be greater than min value.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[0]
				w.Gauge.PercentValue = true
				w.Gauge.Min = 5
				w.Gauge.Max = -21
				d.Widgets[0] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A gauge widget thresholds can't have same start values.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[0]
				w.Gauge.Thresholds = []model.Threshold{
					model.Threshold{Color: "#FFFFFF", StartValue: 5},
					model.Threshold{Color: "#FFF000", StartValue: 5},
				}
				d.Widgets[0] = w
				return d
			},
			expErr: true,
		},

		// Singlestat widget.
		{
			name: "A singlestat widget should have a value text.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[1]
				w.Singlestat.ValueText = ""
				d.Widgets[1] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A singlestat widget with a query should have an expression.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[1]
				w.Singlestat.Query.Expr = ""
				d.Widgets[1] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A singlestat widget with a query should have an datasource ID.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[1]
				w.Singlestat.Query.DatasourceID = ""
				d.Widgets[1] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A singlestat widget thresholds can't have same start values.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[1]
				w.Singlestat.Thresholds = []model.Threshold{
					model.Threshold{Color: "#FFFFFF", StartValue: 5},
					model.Threshold{Color: "#FFF000", StartValue: 5},
				}
				d.Widgets[1] = w
				return d
			},
			expErr: true,
		},

		// Graph widget.
		{
			name: "A graph widget should have at least one query.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[2]
				w.Graph.Queries = []model.Query{}
				d.Widgets[2] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A graph widget with a query should have an expression.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[2]
				w.Graph.Queries = append(w.Graph.Queries, model.Query{Expr: "", Legend: "test", DatasourceID: "test"})
				d.Widgets[2] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A graph widget with a query should have an datasource ID.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[2]
				w.Graph.Queries = append(w.Graph.Queries, model.Query{Expr: "query", Legend: "test", DatasourceID: ""})
				d.Widgets[2] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A graph widget series override should have a regex.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[2]
				w.Graph.Visualization.SeriesOverride = append(w.Graph.Visualization.SeriesOverride, model.SeriesOverride{})
				d.Widgets[2] = w
				return d
			},
			expErr: true,
		},
		{
			name: "A graph widget series override should have the regexes compiled.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[2]
				w.Graph.Visualization.SeriesOverride = []model.SeriesOverride{
					model.SeriesOverride{Regex: "2..", Color: "#FFF000"},
					model.SeriesOverride{Regex: "3..", Color: "#FFF000"},
				}
				d.Widgets[2] = w
				return d
			},
			expDashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[2]
				w.Graph.Visualization.SeriesOverride = []model.SeriesOverride{
					model.SeriesOverride{Regex: "2..", Color: "#FFF000", CompiledRegex: regexp.MustCompile("2..")},
					model.SeriesOverride{Regex: "3..", Color: "#FFF000", CompiledRegex: regexp.MustCompile("3..")},
				}
				d.Widgets[2] = w
				return d
			},
		},
		{
			name: "A graph widget series override can't have more than one series with the same regex.",
			dashboard: func() model.Dashboard {
				d := getBaseDashboard()
				w := d.Widgets[2]
				w.Graph.Visualization.SeriesOverride = []model.SeriesOverride{
					model.SeriesOverride{Regex: "2..", Color: "#FFF000"},
					model.SeriesOverride{Regex: "2..", Color: "#FFF001"},
				}
				d.Widgets[2] = w
				return d
			},
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			got := test.dashboard()
			err := got.Validate()
			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				exp := test.expDashboard()
				assert.Equal(exp, got)
			}
		})
	}
}
