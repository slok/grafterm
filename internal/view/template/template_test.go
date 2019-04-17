package template_test

import (
	"testing"

	"github.com/slok/meterm/internal/view/template"
	"github.com/stretchr/testify/assert"
)

func newData() template.Data {
	return template.Data{
		Dashboard: template.Dashboard{
			Range: "10m",
		},
		Query: template.Query{
			DatasourceID: "ds1",
			Labels: map[string]string{
				"code":    "401",
				"method":  "GET",
				"handler": "/test/:id/status",
			},
		},
	}
}

func TestDataRender(t *testing.T) {
	tests := []struct {
		name string
		data template.Data
		tpl  string
		exp  string
	}{
		{
			name: "Dashboard range",
			data: newData(),
			tpl:  "range: {{ .Dashboard.Range }}",
			exp:  "range: 10m",
		},
		{
			name: "Query datasource",
			data: newData(),
			tpl:  "datasource: {{ .Query.DatasourceID }}",
			exp:  "datasource: ds1",
		},
		{
			name: "Query labels",
			data: newData(),
			tpl:  "data: [{{ .Query.Labels.method }}] {{ .Query.Labels.handler }} {{ .Query.Labels.code }}",
			exp:  "data: [GET] /test/:id/status 401",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.data.Render(test.tpl)
			assert.Equal(t, test.exp, got)
		})
	}
}

func TestDataCopy(t *testing.T) {
	tests := []struct {
		name           string
		data           template.Data
		transform      func(data template.Data) template.Data
		expTransformed template.Data
		expOriginal    template.Data
	}{
		{
			name: "Dashboard",
			data: newData(),
			transform: func(data template.Data) template.Data {
				return data.WithDashboard(template.Dashboard{
					Range: "15m",
				})
			},
			expOriginal: newData(),
			expTransformed: template.Data{
				Dashboard: template.Dashboard{
					Range: "15m",
				},
				Query: template.Query{
					DatasourceID: "ds1",
					Labels: map[string]string{
						"code":    "401",
						"method":  "GET",
						"handler": "/test/:id/status",
					},
				},
			},
		},
		{
			name: "Query",
			data: newData(),
			transform: func(data template.Data) template.Data {
				return data.WithQuery(template.Query{
					DatasourceID: "ds2",
					Labels: map[string]string{
						"code":     "402",
						"otherKey": "otherValue",
					},
				})
			},
			expOriginal: newData(),
			expTransformed: template.Data{
				Dashboard: template.Dashboard{
					Range: "10m",
				},
				Query: template.Query{
					DatasourceID: "ds2",
					Labels: map[string]string{
						"code":     "402",
						"otherKey": "otherValue",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			got := test.transform(test.data)
			assert.Equal(test.expTransformed, got)
			assert.Equal(test.data, test.expOriginal)
		})
	}
}
