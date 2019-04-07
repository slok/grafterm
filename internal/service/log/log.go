package log

import (
	"fmt"
	"log"
)

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

// STD is a logger that uses the standard go logger, mainly for debugging.
var STD = &std{}

type std struct{}

func (s std) Infof(format string, args ...interface{}) {
	f := fmt.Sprintf("[INFO] %s", format)
	log.Printf(f, args...)
}
func (s std) Errorf(format string, args ...interface{}) {
	f := fmt.Sprintf("[ERROR] %s", format)
	log.Printf(f, args...)
}
