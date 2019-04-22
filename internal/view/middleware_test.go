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
	tests := []struct {
		name    string
		data    template.Data
		cfg     syncConfig
		expData template.Data
	}{
		{
			name: "storing static data should add that data on every call to the sync",
			data: map[string]string{
				"name":     "Batman",
				"realName": "Bruce",
				"lastName": "Wayne",
				"location": "Gotham",
			},
			cfg: syncConfig{
				templateData: map[string]string{
					"location":  "Arkham asylum",
					"transport": "batmobile",
				},
			},
			expData: map[string]string{
				"name":      "Batman",
				"realName":  "Bruce",
				"lastName":  "Wayne",
				"location":  "Arkham asylum",
				"transport": "batmobile",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mw := &mockWidget{}
			w := withWidgetDataMiddleware(test.data, mw)
			w.sync(context.TODO(), test.cfg)

			assert.Equal(t, test.expData, mw.calledCfg.templateData)
		})
	}
}
