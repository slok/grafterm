package view

import (
	"regexp"
	"testing"

	"github.com/slok/grafterm/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestWidgetColorManager(t *testing.T) {
	tests := []struct {
		name     string
		getColor func(wcm widgetColorManager) (string, error)
		expColor string
		expErr   bool
	}{
		{
			name: "Get default colors in order.",
			getColor: func(wcm widgetColorManager) (string, error) {
				wcm.GetDefaultColor()
				wcm.GetDefaultColor()
				wcm.GetDefaultColor()
				return wcm.GetDefaultColor(), nil
			},
			expColor: "#EF843C",
		},
		{
			name: "Get threshold colors.",
			getColor: func(wcm widgetColorManager) (string, error) {
				thresholds := []model.Threshold{
					model.Threshold{
						StartValue: 0,
						Color:      "#111111",
					},
					model.Threshold{
						StartValue: 30,
						Color:      "#222222",
					},
					model.Threshold{
						StartValue: 75,
						Color:      "#333333",
					},
				}
				return wcm.GetColorFromThresholds(thresholds, 50)
			},
			expColor: "#222222",
		},
		{
			name: "Geting threshold colors without thresholds should error.",
			getColor: func(wcm widgetColorManager) (string, error) {
				thresholds := []model.Threshold{}
				return wcm.GetColorFromThresholds(thresholds, 50)
			},
			expErr: true,
		},
		{
			name: "Get threshold colors.",
			getColor: func(wcm widgetColorManager) (string, error) {
				thresholds := []model.Threshold{
					model.Threshold{
						StartValue: 0,
						Color:      "#111111",
					},
					model.Threshold{
						StartValue: 30,
						Color:      "#222222",
					},
					model.Threshold{
						StartValue: 75,
						Color:      "#333333",
					},
				}
				return wcm.GetColorFromThresholds(thresholds, 50)
			},
			expColor: "#222222",
		},
		{
			name: "Get color from series legend.",
			getColor: func(wcm widgetColorManager) (string, error) {
				cfg := model.GraphWidgetSource{
					Visualization: model.GraphVisualization{
						SeriesOverride: []model.SeriesOverride{
							model.SeriesOverride{
								CompiledRegex: regexp.MustCompile("a-.*"),
								Color:         "#111111",
							},
							model.SeriesOverride{
								CompiledRegex: regexp.MustCompile("b-.*"),
								Color:         "#222222",
							},
							model.SeriesOverride{
								CompiledRegex: regexp.MustCompile("c-.*"),
								Color:         "#333333",
							},
							model.SeriesOverride{
								CompiledRegex: regexp.MustCompile("d-.*"),
								Color:         "#444444",
							},
							model.SeriesOverride{
								CompiledRegex: regexp.MustCompile("e-.*"),
								Color:         "#444444",
							},
						},
					},
				}
				return wcm.GetColorFromSeriesLegend(cfg, "d-12345"), nil
			},
			expColor: "#444444",
		},
		{
			name: "Get default color when there's no match with series legend regexes.",
			getColor: func(wcm widgetColorManager) (string, error) {
				cfg := model.GraphWidgetSource{
					Visualization: model.GraphVisualization{
						SeriesOverride: []model.SeriesOverride{
							model.SeriesOverride{
								CompiledRegex: regexp.MustCompile("a-.*"),
								Color:         "#111111",
							},
						},
					},
				}
				return wcm.GetColorFromSeriesLegend(cfg, "d-12345"), nil
			},
			expColor: "#7EB26D",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var wcm widgetColorManager
			color, err := test.getColor(wcm)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expColor, color)
			}
		})
	}
}
