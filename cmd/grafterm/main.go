package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"

	"github.com/slok/grafterm/internal/controller"
	"github.com/slok/grafterm/internal/model"
	"github.com/slok/grafterm/internal/service/configuration"
	"github.com/slok/grafterm/internal/service/log"
	"github.com/slok/grafterm/internal/service/metric"
	metricdatasource "github.com/slok/grafterm/internal/service/metric/datasource"
	metricmiddleware "github.com/slok/grafterm/internal/service/metric/middleware"
	"github.com/slok/grafterm/internal/view"
	"github.com/slok/grafterm/internal/view/render/termdash"
)

// Main is the main application.
type Main struct {
	flags  *flags
	logger log.Logger
}

// Run runs the main app.
func (m *Main) Run() error {
	if m.flags.version {
		fmt.Fprint(os.Stdout, Version)
		return nil
	}

	// If debug mode then use a verbose logger.
	m.logger = log.Dummy
	if m.flags.debug {
		f, err := os.OpenFile(m.flags.logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			return err
		}
		defer f.Close()

		m.logger = log.New(log.Config{
			Output: f,
		})
	}

	// Load Dashboard.
	cfg, err := loadConfiguration(m.flags.cfg)
	if err != nil {
		return err
	}

	ddss, err := cfg.Datasources()
	if err != nil {
		return err
	}

	udss, err := m.loadUserDatasources()
	if err != nil {
		return err
	}

	gatherer, err := m.createGatherer(ddss, udss)
	if err != nil {
		return err
	}

	// Create controller.
	ctrl := controller.NewController(gatherer)

	// Create renderer.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	renderer, err := termdash.NewTermDashboard(cancel, m.logger)
	if err != nil {
		return err
	}
	defer renderer.Close()

	// Prepare app for running.
	var g run.Group

	// Capture signals.
	{
		sigC := make(chan os.Signal, 1)
		exitC := make(chan struct{})
		signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)

		g.Add(
			func() error {
				select {
				case <-sigC:
				case <-exitC:
				}

				return nil
			},
			func(e error) {
				close(exitC)
			})
	}

	// Run application.
	{
		appcfg := view.AppConfig{
			RefreshInterval:   m.flags.refreshInterval,
			RelativeTimeRange: m.flags.relativeDur,
			OverrideVariables: m.flags.variables,
		}

		// Only set fixed time if start set.
		if m.flags.start != "" {
			start, err := timeFromFlag(m.flags.start)
			if err != nil {
				return fmt.Errorf("error parsing start flag: %s", err)
			}
			end, err := timeFromFlag(m.flags.end)
			if err != nil {
				return fmt.Errorf("error parsing end flag: %s", err)
			}

			appcfg.TimeRangeStart = start
			appcfg.TimeRangeEnd = end

			// Check times are correct.
			if !appcfg.TimeRangeEnd.IsZero() && appcfg.TimeRangeEnd.Before(appcfg.TimeRangeStart) {
				return fmt.Errorf("end timestamp can't be before start timestamp")
			}
		}

		app := view.NewApp(appcfg, ctrl, renderer, m.logger)
		ds, err := cfg.Dashboard()
		if err != nil {
			return err
		}

		g.Add(
			func() error {
				err := app.Run(ctx, ds)
				if err != nil {
					return err
				}
				defer cancel()
				return nil
			},
			func(e error) {
				cancel()
			})
	}

	return g.Run()
}

func loadConfiguration(cfgPath string) (configuration.Configuration, error) {
	// Load dashboard file.
	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg, err := configuration.JSONLoader{}.Load(f)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (m *Main) loadUserDatasources() ([]model.Datasource, error) {
	// If we could not load user datasources do not fail.
	f, err := os.Open(m.flags.userDSPath)
	if err != nil {
		m.logger.Warnf("could not load '%s' user datasources file: %s", m.flags.userDSPath, err)
		return []model.Datasource{}, nil
	}
	defer f.Close()

	cfg, err := configuration.JSONLoader{}.Load(f)
	if err != nil {
		return nil, err
	}

	return cfg.Datasources()
}

func (m *Main) createGatherer(dashboardDss, userDss []model.Datasource) (metric.Gatherer, error) {
	gatherer, err := metricdatasource.NewGatherer(metricdatasource.ConfigGatherer{
		DashboardDatasources: dashboardDss,
		UserDatasources:      userDss,
		Aliases:              m.flags.aliases,
	})
	if err != nil {
		return nil, err
	}
	gatherer = metricmiddleware.Logger(m.logger, gatherer)

	return gatherer, nil
}

// timeFromFlag gets the time from a flag based on a duration or on a
// fixed time stamp.
func timeFromFlag(v string) (time.Time, error) {
	var t time.Time

	// Try parsing using duration.
	d, err := time.ParseDuration(v)
	if err == nil {
		t = time.Now().UTC().Add(-1 * d)
	} else {
		// Try parsing as ISO 8601.
		parsedTime, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return t, fmt.Errorf("'%s' is not a valid timestamp or duration string", v)
		}
		t = parsedTime
	}

	return t, nil
}

func main() {
	flags, err := newFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %s\n", err)
		os.Exit(1)
	}

	m := Main{
		flags: flags,
	}

	if err := m.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error executing program: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
