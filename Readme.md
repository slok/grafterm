# Grafterm

Grafterm is an application to load dashboards on the terminal, you could think of grafterm as a simplified version of [Grafana] but for terminal.

## Features

- Multiple widgets (graph, singlestat, gauge).
- Multiple datasources usage (including aggregation).
- Extensible metrics implementation (Prometheus included).
- Templating of variables.
- Auto time interval adjustment for queries.
- Fixed and adaptive grid.
- Color customization on widgets.
- Configurable autorefresh.
- Single binary and easy usage/deployment.

## Installation

Download the binaries from [releases]

## Run options

- `--cfg`: Path for the dashboard configuration file to load.
- `--refresh-interval` refresh interval for the dashboard metrics.
- `--relative-time-range` relative time from now that will be used for the
- `--debug`: run in debug mode, it will log to `--log-path` output or `grafterm.log` by default.

## Dashboard

Check [this][cfg-md] section that explains how a dashboard is configured. Also check [dashboard examples][dashboard-examples]

## Datasources

Datasources are the way grafterm knows how to retrieve the metrics for the dashboard. these are the datasources supported at this moment:

- Prometheus

Open an issue or a PR to support more datasource types.

## Kudos

This project would not be possible without the effort of many people and projects but specially [Grafana] for the inspiration, ideas and the project itself, and [Termdash] for the rendering of all those fancy graphs on the terminal.

[grafana]: https://grafana.com/
[termdash]: https://github.com/mum4k/termdash
[releases]: https://github.com/slok/grafterm/releases
[cfg-md]: docs/cfg.md
[dashboard-examples]: dashboard-examples
