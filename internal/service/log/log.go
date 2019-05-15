package log

import (
	"io"

	"github.com/rs/zerolog"
)

// Logger knows how to log.
type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// Dummy is a dummy logger that doesn't log anything.
var Dummy = &dummy{}

type dummy struct{}

func (d dummy) Infof(format string, args ...interface{})  {}
func (d dummy) Warnf(format string, args ...interface{})  {}
func (d dummy) Errorf(format string, args ...interface{}) {}

// Config is the Logger configuration
type Config struct {
	Output io.Writer
}

// New returns a new logger.
func New(cfg Config) Logger {
	return newZero(cfg)
}

func newZero(cfg Config) Logger {
	return &zero{
		logger: zerolog.New(cfg.Output).With().
			Timestamp().
			Logger(),
	}
}

type zero struct {
	logger zerolog.Logger
}

func (z zero) Infof(format string, args ...interface{}) {
	z.logger.Info().Msgf(format, args...)
}
func (z zero) Warnf(format string, args ...interface{}) {
	z.logger.Warn().Msgf(format, args...)
}
func (z zero) Errorf(format string, args ...interface{}) {
	z.logger.Error().Msgf(format, args...)
}
