package main

import (
	"flag"
	"os"
)

const (
	defConfig          = "dashboard.json"
	defRefreshInterval = "10s"
)

type flags struct {
	cfg             string
	debug           bool
	version         bool
	refreshInterval string
}

func newFlags() (*flags, error) {
	flags := &flags{}
	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Register flags.
	fl.StringVar(&flags.cfg, "cfg", defConfig, "the path to the configuration file")
	fl.StringVar(&flags.refreshInterval, "refresh-interval", defRefreshInterval, "the interval to refresh the dashboard")
	fl.BoolVar(&flags.debug, "debug", false, "enable debug mode")
	fl.BoolVar(&flags.version, "version", false, "print version")

	fl.Parse(os.Args[1:])

	if err := flags.validate(); err != nil {
		return nil, err
	}

	return flags, nil
}

func (f *flags) validate() error {
	return nil
}
