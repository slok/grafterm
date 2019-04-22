// Code generated by mockery v1.0.0. DO NOT EDIT.

package render

import mock "github.com/stretchr/testify/mock"
import model "github.com/slok/grafterm/internal/model"

// GaugeWidget is an autogenerated mock type for the GaugeWidget type
type GaugeWidget struct {
	mock.Mock
}

// GetWidgetCfg provides a mock function with given fields:
func (_m *GaugeWidget) GetWidgetCfg() model.Widget {
	ret := _m.Called()

	var r0 model.Widget
	if rf, ok := ret.Get(0).(func() model.Widget); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.Widget)
	}

	return r0
}

// SetColor provides a mock function with given fields: hexColor
func (_m *GaugeWidget) SetColor(hexColor string) error {
	ret := _m.Called(hexColor)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(hexColor)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Sync provides a mock function with given fields: isPercent, value
func (_m *GaugeWidget) Sync(isPercent bool, value float64) error {
	ret := _m.Called(isPercent, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(bool, float64) error); ok {
		r0 = rf(isPercent, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
