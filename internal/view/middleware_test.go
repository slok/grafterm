package view

import (
	"context"
	"testing"

	"github.com/slok/grafterm/internal/view/template"
	"github.com/stretchr/testify/assert"
)

type mockWidget struct {
	calledCfg syncConfig
}

func (m *mockWidget) sync(_ context.Context, cfg syncConfig) error {
	m.calledCfg = cfg
	return nil
}

func TestWidgetDataMiddleware(t *testing.T) {
	tests := map[string]struct {
		data         template.Data
		overrideData template.Data
		cfg          syncConfig
		expData      template.Data
	}{
		"Storing static data should add that data on every call to the sync.": {
			data: map[string]interface{}{
				"name":     "Batman",
				"realName": "Bruce",
				"lastName": "Wayne",
				"location": "Gotham",
			},
			cfg: syncConfig{
				templateData: map[string]interface{}{
					"location":  "Arkham asylum",
					"transport": "batmobile",
				},
			},
			expData: map[string]interface{}{
				"name":      "Batman",
				"realName":  "Bruce",
				"lastName":  "Wayne",
				"location":  "Arkham asylum",
				"transport": "batmobile",
			},
		},
		"Storing override data should add that data on every call to the sync.": {
			overrideData: map[string]interface{}{
				"name":     "Batman",
				"realName": "Bruce",
				"lastName": "Wayne",
				"location": "Gotham",
			},
			cfg: syncConfig{
				templateData: map[string]interface{}{
					"location":  "Arkham asylum",
					"transport": "batmobile",
				},
			},
			expData: map[string]interface{}{
				"name":      "Batman",
				"realName":  "Bruce",
				"lastName":  "Wayne",
				"location":  "Gotham",
				"transport": "batmobile",
			},
		},
		"Override data should be merged and have priority.": {
			data: map[string]interface{}{
				"name":       "Batman",
				"realName":   "Bruce",
				"lastName":   "Wayne",
				"worstEnemy": "Joker",
			},
			overrideData: map[string]interface{}{
				"name":     "Batman2",
				"realName": "Bruce",
				"lastName": "Wayne2",
				"location": "Gotham",
			},
			cfg: syncConfig{
				templateData: map[string]interface{}{
					"location":  "Arkham asylum",
					"transport": "batmobile",
				},
			},
			expData: map[string]interface{}{
				"name":       "Batman2",
				"realName":   "Bruce",
				"lastName":   "Wayne2",
				"location":   "Gotham",
				"transport":  "batmobile",
				"worstEnemy": "Joker",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mw := &mockWidget{}
			w := withWidgetDataMiddleware(test.data, test.overrideData, mw)
			w.sync(context.TODO(), test.cfg)

			assert.Equal(t, test.expData, mw.calledCfg.templateData)
		})
	}
}
