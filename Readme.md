# Grafterm [![CircleCI][circleci-image]][circleci-url]

Visualize metrics dashboards on the terminal, like a simplified and minimalist version of [Grafana] for terminal.

![grafterm red dashboard](/img/grafterm-red-compressed.gif)

## Features

- Multiple widgets (graph, singlestat, gauge).
- Multiple datasources usage (including aggregation).
- Custom dashboards based on JSON configuration files.
- Extensible metrics datasource implementation (Prometheus included).
- Templating of variables.
- Auto time interval adjustment for queries.
- Auto unit formatting on widgets.
- Fixed and adaptive grid.
- Color customization on widgets.
- Configurable autorefresh.
- Single binary and easy usage/deployment.

## Installation

Download the binaries from [releases]

## Run examples

Simple run:

```bash
grafterm -c ./mydashboard.json
```

Setting a relative time duration:

```bash
grafterm -c ./mydashboard.json -d 48h
```

Setting a fixed time range to visualize the metrics using duration notation. In this example is start at `now-22h` and end at `now-20h`

```bash
grafterm -c ./mydashboard.json -s 22h -e 20h
```

Setting a fixed time range to visualize the metrics using timestamp [ISO 8601] notation.

```bash
grafterm -c ./mydashboard.json -s 2019-05-12T12:32:11+02:00 -e 2019-05-12T12:35:11+02:00
```

Replacing dashboard variables:

```bash
grafterm -c ./mydashboard.json -v env=prod -v job=envoy
```

## Dashboard

Check [this][cfg-md] section that explains how a dashboard is configured. Also check [dashboard examples][dashboard-examples]

## Datasources

Datasources are the way grafterm knows how to retrieve the metrics for the dashboard. these are the datasources supported at this moment:

- Prometheus

Open an issue or a PR to support more datasource types.

## Kudos

This project would not be possible without the effort of many people and projects but specially [Grafana] for the inspiration, ideas and the project itself, and [Termdash] for the rendering of all those fancy graphs on the terminal.

[circleci-image]: https://circleci.com/gh/slok/grafterm.svg?style=svg
[circleci-url]: https://circleci.com/gh/slok/grafterm
[grafana]: https://grafana.com/
[termdash]: https://github.com/mum4k/termdash
[releases]: https://github.com/slok/grafterm/releases
[cfg-md]: docs/cfg.md
[dashboard-examples]: dashboard-examples
[iso 8601]: https://en.wikipedia.org/wiki/ISO_8601
