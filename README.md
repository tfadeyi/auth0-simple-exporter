# Auth0 Simple Log Exporter

Exports Prometheus metrics of Auth0 Log Events.

## Usage

```shell
Usage:
  exporter export [flags]

Flags:
      --auth0.client-id string       Auth0 management api client-id
      --auth0.client-secret string   Auth0 management api static token.
      --auth0.domain string          Auth0 tenant's domain.
      --auth0.token string           Auth0 management api static token
  -h, --help                         help for export
      --profiling                    Enabled pprof profiling on the exporter on port :6060. (help: https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/)
      --web.listen-address int       Port where the server will listen. (default 8081)
      --web.metrics-path string      URL Path under which to expose metrics. (default "metrics")
```

## Metrics


## Example queries

