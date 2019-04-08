/*
Package mocks will have all the mocks of the library.
*/
package mocks // import "github.com/slok/meterm/internal/mocks"

// Controller mocks.
//go:generate mockery -output ./controller -outpkg controller -dir ../controller -name Controller

// Render mocks.
//go:generate mockery -output ./view/render -outpkg render -dir ../view/render -name Renderer
//go:generate mockery -output ./view/render -outpkg render -dir ../view/render -name GaugeWidget
//go:generate mockery -output ./view/render -outpkg render -dir ../view/render -name SinglestatWidget

// Services mocks.
//go:generate mockery -output ./service/metric -outpkg metric -dir ../service/metric -name Gatherer

// 3rd party
//go:generate mockery -output ./github.com/prometheus/client_golang/api/prometheus/v1 -outpkg v1 -dir ./thirdparty/github.com/prometheus/client_golang/api/prometheus/v1 -name API
