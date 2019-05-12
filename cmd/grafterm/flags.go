package main

import (
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	// Version is the application version.
	Version = "dev"
)

const (
	defConfig          = "dashboard.json"
	defRefreshInterval = "10s"
	defLogPath         = "grafterm.log"

	// flag descriptions.
	descCfg               = "the path to the configuration file"
	descRefreshInterval   = "the interval to refresh the dashboard"
	descLogPath           = "the path where the log output will be written"
	descRelativeTimeRange = "relative time range (from now) for the dashboard time range"
	descDebug             = "enable debug mode, on debug mode it will print logs to the desired output"
	descVar               = "repeatable flag that will override the variable defined on the dashboard (in 'key=value' form)"
)

type flags struct {
	variables         map[string]string
	cfg               string
	debug             bool
	version           bool
	refreshInterval   string
	logPath           string
	relativeTimeRange string
}

func newFlags() (*flags, error) {
	flags := &flags{
		variables: map[string]string{},
	}

	// Create app.
	app := kingpin.New("grafterm", "graph metrics on the terminal")
	app.Version(Version)

	// Register flags.
	app.Flag("cfg", descCfg).Default(defConfig).Short('c').StringVar(&flags.cfg)
	app.Flag("refresh-interval", descRefreshInterval).Default(defRefreshInterval).Short('r').StringVar(&flags.refreshInterval)
	app.Flag("log-path", descLogPath).Default(defLogPath).StringVar(&flags.logPath)
	app.Flag("relative-time-range", descLogPath).StringVar(&flags.relativeTimeRange)
	app.Flag("var", descVar).Short('v').StringMapVar(&flags.variables)
	app.Flag("debug", descDebug).BoolVar(&flags.debug)
	app.Parse(os.Args[1:])

	if err := flags.validate(); err != nil {
		return nil, err
	}

	return flags, nil
}

func (f *flags) validate() error {
	return nil
}
