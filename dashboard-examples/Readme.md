# Dashboard examples

## RED metrics

This dashboards shows [RED] metrics.

Useful override variables:

- `prefix`: Will be used as a prefix on all the metrics queries to Prometheus.
- `job`: Will be used as the `job` label to filter on all the Prometheus queries.

![](https://i.imgur.com/DOPeiWI.png)

## Go stats

This is a useful dashboard for go application.

Useful override variables:

- `job`: Will be used as the `job` label to filter on all the Prometheus queries.

![](https://i.imgur.com/dyiR7J6.png)

![](https://i.imgur.com/qeXRmOl.png)

## Gitlab

This is a gitlab based dashboard example.

![](https://i.imgur.com/RGlygHF.png)

## Kubernetes status

This is a port from [this](https://grafana.com/dashboards/5315) Grafana dashboard (with very small changes). Mainly shows the usage of the gauges.

![](https://i.imgur.com/N5jtCFT.png)

## Wikimedia

This is a port from [this](https://grafana.wikimedia.org/d/000000002/api-backend-summary) dashboard (with very small changes). It uses Graphite backend.

![](https://i.imgur.com/bJjGtyF.png)

[red]: https://www.weave.works/blog/the-red-method-key-metrics-for-microservices-architecture/
