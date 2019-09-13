# Configuration file

First of all there are dashboard examples [here][dashboard-examples]

The configuration file format is JSON and is splitted in two main blocks, `datasources` and `dashboard`.

```json
{
  "version": "v1",
  "datasources": [],
  "dashboard": {}
}
```

## Datasources

This main block contains a list of the datasources being used by the dashboard, depending on the datasource type it will have different options. The dashboard widgets will reference the datasource by the `id`.

```json
  "datasources": [
    {
      "id": "ds1",
      "prometheus": {
        "address": "http://127.0.0.1:9090"
      }
    },
    {
      "id": "ds2",
      "prometheus": {
        "address": "http://127.0.0.1:9091"
      }
    },
    {
      "id": "ds3",
      "graphite": {
        "address": "http://127.0.0.1:7123"
      }
    }
  ],
```

### Types

#### [Prometheus]

This will gather metrics from Prometheus.

Options:

- `address`: Address to Prometheus API

#### [Graphite]

This will gather metrics from Graphite API backends.

Options:

- `address`: Address to Graphite API

#### [InfluxDB]

This will gather metrics from InfluxDB backends.

Options:

- `address`: Address to InfluxDB API
- `database`: Database to use
- `username`: Username for basic auth
- `password`: Password for basic auth
- `insecure`: True to allow insecure https

## Dashboard

The dashboard contains the dashboard configuration and is composed of multiple smaller configuration blocks.

```json
  "dashboard": {
    "grid": {},
    "variables": [],
    "widgets": []
  }
```

### Grid

The grid has configuration of how the grid of the dashboard will behave.

```json

"grid": {
    "fixedWidgets": true,
    "maxWidth": 100
}
```

#### `maxWidth`

Set the scale of the widgets width, by default is 100, so the widgets should have a width based on 100 by default.

For example, 100 means all the row and 50 means half of the row, if you change the `maxWith` you will need to change the width of the widgets.

#### `fixedWidgets`

The grid has 2 wais of behave, adaptive and fixed, by default this setting is `false` so it means that the grid is adaptive.

Adaptive grids ignore widget's `gridPos.x` and `gridPos.y` and only check the width of the widget (`gridPos.w`), this means that it will fill the row until the next widget doesn't fit on that row and will create a new row.

Fixed grids need that the widget have the `x`, `y` and `w`, are more flexible because you can leave spaces between widgets but need all the data so the widget can be placed on the grid.

### Variables

Some strings on the dashboard can be templated (for now only queries and graph labels), here comes the utility of the variables.

Variables are variables that have dynamic values depending on how and when the dashboard has been loaded, there are of different types.

#### Default variables

These are variables that grafterm loads automatically on the dashboard and can be used by the templating system:

- `__range`: The range duration the dashboard is loading.
- `__refresInterval`: The dashboard refresh interval selected.

#### Constant

Constant variables are constants, they don't change, are used to set the value on one place and use it on many queries, for example `environment`.

```json
"variables": [
    {
        "name": "env",
        "constant": {
            "value": "production"
        }
    }
]
```

#### Interval

Interval sets on a variable a dynamic interval based on the range loaded using optional `steps` value. This is handy to have smoother graphs when the range is big because based on the steps the interval would be bigger also and would remove the spikes

```json
"variables": [
    {
        "name": "interval",
        "interval": {
            "steps": 50
        }
    }
]
```

### Widgets

All widgets have some common settings and then custom settings that differ one from the others depending on the kind of widget.

#### Common

```json
"widgets": [
    {
        "title": "Some widget",
        "gridPos": {
            "x": 0,
            "y": 0,
            "w": 5,
        }
    }
]
```

##### `title`

The tittle of the widget.

##### `gridPos`

This argument describes the where and size of the widget. if using adaptive grid `x` and `y` will be ignored. check `Grid` section to know how this works.

#### Gauge

This widget is for realtime metrics, doens't show a range of metrics it shows the last point in time (now) of the metric, this means that only accepts one query.

It has fixed minimum and maximum values to be rendered on the terminal and can work as a percentage based widget (from 0% to 100%).

It can have configured the thresholds so it changes the color of the widget based on them.

```json
{
  "gauge": {
    "query": {},
    "percentValue": true,
    "max": 60,
    "min": 0,
    "thresholds": [
      {
        "color": "#37872D"
      },
      {
        "color": "#FA6400",
        "startValue": 50
      },
      {
        "color": "#C4162A",
        "startValue": 75
      }
    ]
  }
}
```

##### `percentValue`

If `true` it will show the percent value based on the `max` setting. For example if `max: 60` then a value of `60` would show `100%`.

If `false` it will show the value itself. For example if `max: 60` would show `60` as value but the widget would render the gauge bar at the highest point.

##### `thresholds`

Is a list of thresholds, if no `startValue` it will be taken as the base color, if no more thresholds this will be the color of the widget. If more thresholds are on the list then it will set the color based on the range of the thresholds from `startValue` until the next `startValue`.

#### Singlestat

The singlestat acts similar to the Gauge, it's realtime and accepts thresholds but id renders the value itself and not a visual representation of fixed boundaries.

```json
"singlestat": {
    "query": {},
    "unit": "bytes",
    "decimals": 3,
    "valueText": "{{.value}}",
    "thresholds": [
        {
            "color": "#299c46"
        },
        {
            "color": "#FF780A",
            "startValue": 350
        },
        {
            "color": "#d44a3a",
            "startValue": 600
        }
    ]
}
```

##### `valueText`

This is what will be rendered on the singlestat panel, is a go template and the value of the metric is accessible in `.value`. By default it will print the metric value as it is (`{{.value}}`)

Examples:

- Print with 2 decimals: `{{ printf "%.2f" .value }}`
- Print DOWN if value `<1` and UP on `>=1`: `{{ if (lt .value 1.0) }}DOWN{{else}}UP{{end}}`

##### `unit`

Will convert the value to the unit text representation. Check `unit` section in this same doc.

##### `decimals`

The number of decimals used for the representation when the unit format is used.

#### Graph

This widget graphs different metric series in a range. It accepts multiple queries that will be aggregated on the same graph. A single query can be rendered with multiple series (depending on the returned results).

```json
"graph": {
    "visualization": {
        "legend": {
            "disable": false,
            "rightSide": true
        },
        "yAxis": {
          "unit": "seconds",
          "decimals": 0
        },
        "seriesOverride": [
            {
                "regex": "p99",
                "color": "#c15c17",
                "nullPointMode": "connected"
            },
            {
                "regex": "p95",
                "color": "#f2c96d",
                "nullPointMode": "zero"
            },
            {
                "regex": "p50",
                "color": "#f9ba8f"
            }
        ]
    },
    "queries": []
}
```

##### `visualization.legend`

The legend of the graph visualization can be enabled or disabled. It's enabled by default and can be set on the right of the graph, by default it's on the bottom of it.

##### `visualization.seriesOverride`

Each of the graph series can be override based on the legend displayed using a regex, this means that multiple series can be override using the same options.

The setting that can be override at this moment are:

- `color`: The color of the displayed series.
- `nullPointMode`: This will fill the datapoints on the graph that are missing with different strategies, this setting is useful for graphs that don't have sufficent metrics or are very spaced. The strategies are:
  - `connected`: Will use an already near known value and use this.
  - `zero`: Will fill the data point value with 0s.

##### `visualization.yAxis`

In this block the settings that represent the format of the y axis are customized.

- `unit`: Will convert the value to the unit text representation. Check `unit` section in this same doc.
- `decimals`: The number of decimals used for the representation when the unit format is used.

### Templating

Templating of strings use golang built in template. You can use variables of different kinds on different parts of the dashboard.

Examples of using templating system:

- `"{{ .backend }}"`
- `sum(rate(http_request_duration_seconds_count[{{.interval}}]))`

### Query

Query acts differently depending on the widget. If the widget is a real-time widget or a range based widget.

Is composed of a datasource, a legend representation and an expression. The expression will be used on the datasource referenced and the legend will represent the different metrics obtained.

```json
{
  "datasourceID": "ds",
  "expr": "sum(rate(http_request_duration_seconds_count[{{.interval}}])) by (code)",
  "legend": "{{ .code }}"
}
```

The legend has the ability to use templating and has inside loaded the metric labels obtained by the datasource kind.

### Units

Some widgets have unit formatting support, these are the ones that can be used:

- ``(Default): It will fallback to `short`.
- `short`: Will make the values in short format, e.g:
  - `1000`:`1 K`
  - `1000000`: `1 Mil`
- `none`: The value as it is.
- `percent`: Will add a `%` prefix, e.g:
  - `100`: `100%`
  - `20`: `20%`
- `ratio`: Conversion values from 0-1 to percent format, e.g:
  - `0.5`: `50%`
  - `2.1`: `210%`
- `seconds`: Will convert to a single unit pretty duration form, is based on second unit, e.g:
  - `300`: `5m`
  - `10`: `10s`
  - `0.3`: `300ms`
  - `61200`: `17h`
  - `432000`: `5d`
- `reqps`: Will add the `reqps` suffix.
- `bytes`: based value in bytes will convert to a pretty format, e.g:
  - `1024`: `1 KiB`
  - `1.405e+8`: `134 MiB`
  - `4.508e+13`: `41 TiB`

Units come in combination with the `decimals` settings.

[dashboard-examples]: /dashboard-examples
[prometheus]: http://prometheus.io
[graphite]: http://graphiteapp.org
