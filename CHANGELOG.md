# Changelog

## [0.2.1] - 2021-05-07

### Fixed

- Updated Testify to v1.7.0 to resovle issue #60; a CVE 6.5 present in the Yaml package that Testify v1.4.0 uses.

## [0.2.0] - 2019-07-26

### Added

- ARM release builds.
- Add `nullPointMode` series override setting (with `zero` and `connected` strategies).
- Graphite datasource.
- Milliseconds unit conversion.
- Quit grafterm with `Esc` key.
- User defined datasources via flag and/or env var.
- Alias flag to override dashboard datasource ID using user datasource IDs.
- Fallback dashboard referenced datasources to user datasources.

### Fixed

- Gauges that had color thresholds not being show.

## [0.1.0] - 2019-05-13

### Added

- `start` and `end` flags to visualize fixed time range graphs.
- `var` repeatable flag to override dashboard variables.
- Unit formatters for time, RPS, percent and ratios.
- Unit and decimals on the graph widget Y-axis.
- Unit and decimals on the singlestat widget.
- MaxWidth option that sets the horizontal scale of the grid.
- Widget grid fixed mode.
- Widget grid adaptive mode.
- Grid implementation for widgets.
- Dynamic X-axis time labels based on time range and steps.
- Templated queries using variables.
- Const and autointerval variables.
- Color override on graph series based on legend regex.
- Templated legends on graph widget.
- Legend on graph widget.
- Graph widget.
- Single metric gather.
- Metric range gather.
- Allow multiple datasources in the same dashboard.
- Debug flag that will set a verbose logger (will break UI rendering but prints errors and infos).
- Termdash render engine implementation for widgets.
- Singlestat widget.
- Gauge widget.
- Main term application.
- Fake metrics gatherer.

[unreleased]: https://github.com/slok/grafterm/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/slok/grafterm/compare/v0.1.0...0.2.0
[0.1.0]: https://github.com/slok/grafterm/releases/tag/v0.1.0
