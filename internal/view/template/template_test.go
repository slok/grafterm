package template_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/grafterm/internal/view/template"
)

func newData() template.Data {
	return map[string]interface{}{
		"__interval": "2m",
		"__range":    "10m",
		"__start":    "2019-04-19T08:38:59+02:00",
		"__end":      "2019-04-19T10:38:59+02:00",
		"custom":     "test",
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
			name: "template data",
			data: newData(),
			tpl:  "test: [{{ .__interval }}] {{ .custom }}",
			exp:  "test: [2m] test",
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
			name: "Variables",
			data: newData(),
			transform: func(data template.Data) template.Data {
				return data.WithData(map[string]interface{}{
					"custom": "customized-on-test",
					"newkey": "newVar",
				})
			},
			expOriginal: newData(),
			expTransformed: map[string]interface{}{
				"__interval": "2m",
				"__range":    "10m",
				"__start":    "2019-04-19T08:38:59+02:00",
				"__end":      "2019-04-19T10:38:59+02:00",
				"custom":     "customized-on-test",
				"newkey":     "newVar",
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
