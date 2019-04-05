/*
Package mocks will have all the mocks of the library.
*/
package mocks // import "github.com/slok/meterm/internal/mocks"

//go:generate mockery -output ./service/metric -outpkg metric -dir ../service/metric -name Gatherer
