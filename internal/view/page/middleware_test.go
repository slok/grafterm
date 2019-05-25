package page

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/grafterm/internal/view/sync"
	"github.com/slok/grafterm/internal/view/template"
)

type mockWidget struct {
	calledReq *sync.Request
}

func (m *mockWidget) Sync(_ context.Context, r *sync.Request) error {
	m.calledReq = r
	return nil
}

func TestWidgetDataMiddleware(t *testing.T) {
	tests := map[string]struct {
		data         template.Data
		overrideData template.Data
		syncReq      *sync.Request
		expData      template.Data
	}{
		"Storing static data should add that data on every call to the sync.": {
			data: map[string]interface{}{
				"name":     "Batman",
				"realName": "Bruce",
				"lastName": "Wayne",
				"location": "Gotham",
			},
			syncReq: &sync.Request{
				TemplateData: map[string]interface{}{
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
			syncReq: &sync.Request{
				TemplateData: map[string]interface{}{
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
			syncReq: &sync.Request{
				TemplateData: map[string]interface{}{
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
			w.Sync(context.TODO(), test.syncReq)

			assert.Equal(t, test.expData, mw.calledReq.TemplateData)
		})
	}
}
