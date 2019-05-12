package main

import (
	"os"
	"time"

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
	defRelativeDur     = "1h"

	// flag descriptions.
	descCfg             = "the path to the configuration file"
	descRefreshInterval = "the interval to refresh the dashboard"
	descLogPath         = "the path where the log output will be written"
	descRelativeDur     = "the relative duration from now to load the graph."
	descStart           = "the time the dashboard will start in time. Accepts 2 formats, relative time from now based on duration(e.g.: 24h, 15m), or fixed duration in ISO 8601 (e.g.: 2019-05-12T09:35:11+00:00). If set it disables relative duration flag."
	descEnd             = "the time the dashboard will end in time. Accepts 2 formats, relative time from now based on duration(e.g.: 24h, 15m), or fixed duration in ISO 8601 (e.g.: 2019-05-12T09:35:11+00:00)."
	descDebug           = "enable debug mode, on debug mode it will print logs to the desired output"
	descVar             = "repeatable flag that will override the variable defined on the dashboard (in 'key=value' form)"
)

type flags struct {
	variables       map[string]string
	cfg             string
	debug           bool
	version         bool
	refreshInterval time.Duration
	logPath         string
	start           string
	relativeDur     time.Duration
	end             string
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
	app.Flag("refresh-interval", descRefreshInterval).Default(defRefreshInterval).Short('r').DurationVar(&flags.refreshInterval)
	app.Flag("log-path", descLogPath).Default(defLogPath).StringVar(&flags.logPath)
	app.Flag("relative-duration", descRelativeDur).Short('d').DurationVar(&flags.relativeDur)
	app.Flag("start", descStart).Short('s').StringVar(&flags.start)
	app.Flag("end", descEnd).Short('e').StringVar(&flags.end)
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
