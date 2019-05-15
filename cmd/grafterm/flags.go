package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/alecthomas/kingpin"
)

// Version is the application version.
var Version = "dev"

// Defaults settings required by flags..
const (
	defConfig          = "dashboard.json"
	defRefreshInterval = "10s"
	defLogPath         = "grafterm.log"
	defRelativeDur     = "1h"
	defGraftermDir     = "grafterm"
)

var defUserDatasourcePath = []string{defGraftermDir, "datasources.json"}

// Env vars.
const (
	envPrefix          = "GRAFTERM"
	envUserDatasources = envPrefix + "_USER_DATASOURCES"
)

// flag descriptions.
const (
	descCfg             = "the path to the configuration file"
	descRefreshInterval = "the interval to refresh the dashboard"
	descLogPath         = "the path where the log output will be written"
	descRelativeDur     = "the relative duration from now to load the graph."
	descStart           = "the time the dashboard will start in time. Accepts 2 formats, relative time from now based on duration(e.g.: 24h, 15m), or fixed duration in ISO 8601 (e.g.: 2019-05-12T09:35:11+00:00). If set it disables relative duration flag."
	descEnd             = "the time the dashboard will end in time. Accepts 2 formats, relative time from now based on duration(e.g.: 24h, 15m), or fixed duration in ISO 8601 (e.g.: 2019-05-12T09:35:11+00:00)."
	descDebug           = "enable debug mode, on debug mode it will print logs to the desired output"
	descVar             = "repeatable flag that will override the variable defined on the dashboard (in 'key=value' form)"
	descDSAlias         = "repeatable flag that maps dashboard ID datasources to user defined datasources in the form of 'dashboard=user' (in 'key=value' form)"
)

var descUserDS = fmt.Sprintf("path to a configuration file with user defined datasources, these datasources can override the dashboard datasources with the same ID and also can be used to alias them using datasource alias flags. It fallbacks to %s env var", envUserDatasources)

type flags struct {
	variables       map[string]string
	aliases         map[string]string
	cfg             string
	userDSPath      string
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
		aliases:   map[string]string{},
	}

	// Get default datasource path.
	userHome, _ := os.UserHomeDir()
	userDsPath := ""
	if userHome != "" {
		dsPath := []string{userHome}
		dsPath = append(dsPath, defUserDatasourcePath...)
		userDsPath = path.Join(dsPath...)
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
	app.Flag("ds-alias", descDSAlias).Short('a').StringMapVar(&flags.aliases)
	app.Flag("user-datasources", descUserDS).Default(userDsPath).Short('u').Envar(envUserDatasources).StringVar(&flags.userDSPath)
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
