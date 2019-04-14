package main

import (
	"flag"
	"os"
)

const (
	defConfig          = "dashboard.json"
	defRefreshInterval = "10s"
	defLogPath         = "meterm.log"
)

type flags struct {
	cfg               string
	debug             bool
	version           bool
	refreshInterval   string
	logPath           string
	relativeTimeRange string
}

func newFlags() (*flags, error) {
	flags := &flags{}
	fl := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Register flags.
	fl.StringVar(&flags.cfg, "cfg", defConfig, "the path to the configuration file")
	fl.StringVar(&flags.refreshInterval, "refresh-interval", defRefreshInterval, "the interval to refresh the dashboard")
	fl.StringVar(&flags.logPath, "log-path", defLogPath, "the path where the log output will be written")
	fl.StringVar(&flags.relativeTimeRange, "relative-time-range", "", "optional relative time range (from now) for the dashboard time range")
	fl.BoolVar(&flags.debug, "debug", false, "enable debug mode, on debug mode it will print logs to the desired output")
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
