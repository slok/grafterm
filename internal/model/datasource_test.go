package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/grafterm/internal/model"
)

func getBaseDatasource() model.Datasource {
	return model.Datasource{
		ID: "test",
		DatasourceSource: model.DatasourceSource{
			Fake: &model.FakeDatasource{},
		},
	}
}

func TestValidateDatasource(t *testing.T) {
	tests := []struct {
		name   string
		ds     func() model.Datasource
		expErr bool
	}{
		{
			name: "All ok.",
			ds: func() model.Datasource {
				return getBaseDatasource()
			},
			expErr: false,
		},
		{
			name: "A datasources without ID should error.",
			ds: func() model.Datasource {
				d := getBaseDatasource()
				d.ID = ""
				return d
			},
			expErr: true,
		},
		{
			name: "A declared datasource without datasource should error.",
			ds: func() model.Datasource {
				d := getBaseDatasource()
				d.Fake = nil
				return d
			},
			expErr: true,
		},
		{
			name: "A Prometheus datasource without address should error.",
			ds: func() model.Datasource {
				d := getBaseDatasource()
				d.Prometheus = &model.PrometheusDatasource{
					Address: "",
				}
				return d
			},
			expErr: true,
		},
		{
			name: "A Graphite datasource without address should error.",
			ds: func() model.Datasource {
				d := getBaseDatasource()
				d.Graphite = &model.GraphiteDatasource{
					Address: "",
				}
				return d
			},
			expErr: true,
		},
		{
			name: "A InfluxDB datasource without address should error.",
			ds: func() model.Datasource {
				d := getBaseDatasource()
				d.InfluxDB = &model.InfluxDBDatasource{
					Address: "",
				}
				return d
			},
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			err := test.ds().Validate()
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
