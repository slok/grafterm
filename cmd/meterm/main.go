package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"

	"github.com/slok/meterm/internal/controller"
	"github.com/slok/meterm/internal/model"
	"github.com/slok/meterm/internal/service/log"
	"github.com/slok/meterm/internal/service/metric"
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
		m.logger = log.STD
	}

	// Create gatherer.
	gatherer := &metric.FakeGatherer{}

	// Load dashboard file.
	f, err := os.Open(m.flags.cfg)
	if err != nil {
		return err
	}
	defer f.Close()
	dashboard := &model.Dashboard{}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bs, dashboard); err != nil {
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
		app := view.NewApp(ctrl, renderer, m.logger)

		g.Add(
			func() error {
				err := app.Run(ctx, *dashboard)
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
