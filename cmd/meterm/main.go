package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"

	"github.com/slok/meterm/internal/controller"
	"github.com/slok/meterm/internal/service/configuration"
	configurationv1 "github.com/slok/meterm/internal/service/configuration/v1"
	"github.com/slok/meterm/internal/service/log"
	metric "github.com/slok/meterm/internal/service/metric/datasource"
	"github.com/slok/meterm/internal/view"
	"github.com/slok/meterm/internal/view/render/termdash"
)

var (
	// Version is the application version.
	Version = "dev"
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

	// Load configuration.
	cfg, err := m.loadConfiguration()
	if err != nil {
		return err
	}

	gatherer, err := metric.NewGatherer(metric.ConfigGatherer{
		Datasources: cfg.GetDatasources(),
	})
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
		rd, err := time.ParseDuration(m.flags.refreshInterval)
		if err != nil {
			return err
		}

		var relTR time.Duration
		if m.flags.relativeTimeRange != "" {
			r, err := time.ParseDuration(m.flags.relativeTimeRange)
			if err != nil {
				return err
			}
			relTR = r
		}

		appcfg := view.AppConfig{
			RefreshInterval:   rd,
			RelativeTimeRange: relTR,
		}
		app := view.NewApp(appcfg, ctrl, renderer, m.logger)

		g.Add(
			func() error {
				err := app.Run(ctx, cfg.GetDashboard())
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

func (m *Main) loadConfiguration() (configuration.Configuration, error) {
	// Load dashboard file.
	f, err := os.Open(m.flags.cfg)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// For now only v1 supported.
	cfg, err := configurationv1.JSONLoader{}.Load(f)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func main() {
	flags, err := newFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %s", err)
		os.Exit(1)
	}

	m := Main{
		flags: flags,
	}

	if err := m.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error executing program: %s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
