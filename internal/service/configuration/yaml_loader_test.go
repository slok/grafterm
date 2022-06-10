package configuration_test

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/grafterm/internal/service/configuration"
)

func TestLoadYAML(t *testing.T) {
	tests := []struct {
		name       string
		config     func() io.Reader
		loader     func() configuration.Loader
		expVersion string
		expErr     bool
	}{
		{
			name: "Invalid YAML should return an error.",
			loader: func() configuration.Loader {
				return &configuration.YAMLLoader{}
			},
			config: func() io.Reader {
				return strings.NewReader(`{"version": "v1"}`)
			},
			expVersion: "v1",
		},
		{
			name: "Unknown YAML version should error.",
			loader: func() configuration.Loader {
				return &configuration.YAMLLoader{}
			},
			config: func() io.Reader {
				return strings.NewReader(`version: v0.987654321`)
			},
			expErr: true,
		},
		{
			name: "Nested 'meta.version' should fail",
			loader: func() configuration.Loader {
				return &configuration.YAMLLoader{}
			},
			config: func() io.Reader {
				return strings.NewReader(`#
---
meta:
  version: v1
`)
			},
			expErr: true,
		},
		{
			name: "Valid YAML V1 load.",
			loader: func() configuration.Loader {
				return &configuration.YAMLLoader{}
			},
			config: func() io.Reader {
				return strings.NewReader(`#
---
version: v1
`)
			},
			expVersion: "v1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			loader := test.loader()
			gotcfg, err := loader.Load(test.config())

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expVersion, gotcfg.Version())
			}
		})
	}
}
