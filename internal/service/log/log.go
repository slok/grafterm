package log

// Logger knows how to log.
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// Dummy is a dummy logger that doesn't log anything.
var Dummy = &dummy{}

type dummy struct{}

func (d dummy) Infof(format string, args ...interface{})  {}
func (d dummy) Errorf(format string, args ...interface{}) {}
