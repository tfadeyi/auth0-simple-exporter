<div align="center">

# Auth0 Exporter

[![Continuous Integration](https://img.shields.io/github/actions/workflow/status/tfadeyi/auth0-simple-exporter/ci.yml?branch=main&style=flat-square)](https://github.com/tfadeyi/auth0-simple-exporter/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache_2.0-yellowgreen.svg?style=flat-square)](https://github.com/tfadeyi/auth0-simple-exporter/blob/main/LICENSE)
[![Language](https://img.shields.io/github/go-mod/go-version/tfadeyi/auth0-simple-exporter?style=flat-square)](https://github.com/tfadeyi/auth0-simple-exporter)
[![GitHub release](https://img.shields.io/github/v/release/tfadeyi/auth0-simple-exporter?color=green&style=flat-square)](https://github.com/tfadeyi/auth0-simple-exporter/releases)
[![Code size](https://img.shields.io/github/languages/code-size/tfadeyi/auth0-simple-exporter?color=orange&style=flat-square)](https://github.com/tfadeyi/auth0-simple-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/tfadeyi/auth0-simple-exporter?style=flat-square)](https://goreportcard.com/report/github.com/tfadeyi/auth0-simple-exporter)
</div>

---

A simple Prometheus exporter for [Auth0](https://auth0.com/) log [events](https://auth0.com/docs/api/management/v2#!/Logs/get_logs),
which allows you to collect metrics from Auth0 and expose them in a format that can be consumed by Prometheus.

> Development is in progress.

## Motivation

It can be difficult to monitor **Auth0** tenant events on a Prometheus stack,
especially compared to other monitoring systems such as Datadog.
This Prometheus exporter aims to simplify this, making it easier to expose tenant events.

## Prerequisites

* [Auth0](https://auth0.com/) account.
* Auth0 tenant [management API](https://auth0.com/docs/api#management-api) client credentials.
* *(Optional)* Auth0 tenant management API [static token](https://auth0.com/docs/secure/tokens/access-tokens/management-api-access-tokens).

## TL;DR

The quickest way to install the exporter is through Helm, make sure you have your Auth0 credentials at hand.

```shell
export TOKEN="< auth0 management API static static token >"
export DOMAIN="< auth0 tenant domain >"
```
```shell
  # Installing by passing in secret directly
  helm repo add auth0-exporter https://tfadeyi.github.io/auth0-simple-exporter
  helm upgrade --install --create-namespace -n auth0-exporter auth0-exporter/auth0-exporter \
  --set auth0.domain="$DOMAIN" --set auth0.token="$TOKEN" \
  --set exporter.tls.disabled=true
```
This will install the exporter running with TLS disabled.

## Installation

* ### Download Pre-built Binaries

    Binaries can be downloaded from [Releases](https://github.com/tfadeyi/auth0-simple-exporter/releases) page.
  * Download and run exporter's binary with TLS disabled.

    ```shell
    export TOKEN="< auth0 management API static static token >"
    export DOMAIN="< auth0 tenant domain >"

    curl -LJO https://github.com/tfadeyi/auth0-simple-exporter/releases/download/v0.0.2/auth0-simple-exporter-linux-amd64.tar.gz && \
    tar -xzvf auth0-simple-exporter-linux-amd64.tar.gz && \
    cd auth0-simple-exporter-linux-amd64

    ./auth0-simple-exporter export --tls.disabled
    ```
* ### Docker
    The recommended way to get the Docker Image is to pull the prebuilt image from the project's Github Container Registry.
    ```shell
    $ docker pull ghcr.io/tfadeyi/auth0-simple-exporter:latest
    ```
    To use a specific version, you can pull a versioned tag.
    ```shell
    $ docker pull ghcr.io/tfadeyi/auth0-simple-exporter:[TAG]
    ```

* ### Helm
    This shows a simple installation of the exporter helm chart, running with TLS disabled.
    ```shell
    export TOKEN="< auth0 management API static static token >"
    export DOMAIN="< auth0 tenant domain >"
    ```
    ```shell
    # Installing by passing in secret directly
    helm repo add auth0-exporter https://tfadeyi.github.io/auth0-simple-exporter
    helm upgrade --install --create-namespace -n auth0-exporter auth0-exporter/auth0-exporter \
      --set auth0.domain="$DOMAIN" --set auth0.token="$TOKEN" \
      --set exporter.tls.disabled=true
    ```

    More info on the helm deployment can be found [here](deploy/charts/auth0-exporter/README.md).

* ### Build from source
    From the repository root directory run:
    ```shell
    make build
    # or for multiple systems
    make build-all-platforms
    ```

* ### Nix
    The exporter can be used via Nix.
    ```shell
    nix run github:tfadeyi/auth0-simple-exporter export --tls.disabled
    ```

## Usage

```
Usage:
  exporter export [flags]

Flags:
      --auth0.client-id string       Auth0 management api client-id.
      --auth0.client-secret string   Auth0 management api client-secret.
      --auth0.domain string          Auth0 tenant's domain. (i.e: <tenant_name>.eu.auth0.com).
      --auth0.from string            Point in time from were to start fetching auth0 logs. (format: YYYY-MM-DD) (default "2023-04-02")
      --auth0.token string           Auth0 management api static token. (the token can be used instead of client credentials).
  -h, --help                         help for export
      --log.level string             Exporter log level (debug, info, warn, error). (default "warn")
      --namespace string             Exporter's namespace.
      --pprof.enabled                Enabled pprof profiling on the exporter on port :6060. (help: https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/).
      --pprof.listen-address int     Port where the pprof webserver will listen on. (default 6060)
      --probe.listen-address int     Port where the probe webserver will listen on. (default 8081)
      --probe.path string            URL Path under which to expose the probe metrics. (default "probe")
      --subsystem string             Exporter's subsystem.
      --tls.auto                     Allow the exporter to use autocert to renew its certificates with letsencrypt.
                                     (Can only be used if the exporter is publicly accessible by the internet)
      --tls.cert-file string         Path to the PEM encoded certificate for the auth0-exporter metrics to serve.
      --tls.disabled                 Run exporter without TLS. TLS is enabled by default.
      --tls.hosts strings            The different allowed hosts for the exporter. Only works when --tls.auto has been enabled.
      --tls.key-file string          Path to the PEM encoded key for the auth0-exporter metrics server.
      --web.listen-address int       Port where the exporter webserver will listen on. (default 8080)
      --web.path string              URL Path under which to expose the collected auth0 metrics. (default "metrics")
```

### Environment variables:
* ***TOKEN***, Auth0 management API static token.
* ***DOMAIN***, Auth0 tenant domain.
* ***CLIENT_SECRET***, Auth0 management API client-secret, (not required if setting the token).
* ***CLIENT_ID***, Auth0 management API client-id, (not required if setting the token).

## Example queries

Monitor the percentage of failed logins in the Auth0 tenant:

```
(tenant_failed_login_operations_total / tenant_login_operations_total) * 100
```

Monitor the number current logged-in users for a client application in Auth0 tenant:

```
(tenant_login_operations_total{client="ChatGPT"} - tenant_failed_login_operations_total{client="ChatGPT"}) - (tenant_logout_operations_total{client="ChatGPT"} - tenant_failed_logout_operations_total{client="ChatGPT"})
```

## Metrics

### Signup
| Metric                                   | Meaning                                             | Labels |
|------------------------------------------|-----------------------------------------------------|--------|
| `tenant_sign_up_operations_total`        | The total number of signup operations.              | client |
| `tenant_failed_sign_up_operations_total` | The number of failed signup operations. (codes: fs) | client |

### Login
| Metric                                 | Meaning                                                 | Labels |
|----------------------------------------|---------------------------------------------------------|--------|
| `tenant_login_operations_total`        | The total number of login operations.                   | client |
| `tenant_failed_login_operations_total` | The number of failed login operations. (codes: f,fp,fu) | client |

### Logout
| Metric                                  | Meaning                                              | Labels |
|-----------------------------------------|------------------------------------------------------|--------|
| `tenant_logout_operations_total`        | The total number of logout operations.               | client |
| `tenant_failed_logout_operations_total` | The number of failed logout operations. (codes: flo) | client |

**The other exposed metrics can be found [here](./docs/metrics.md).**

## Known Issues

* ### API Rate Limits
  When the Prometheus scraping job interval is too frequent the exporter might encounter api-rate limit from Auth0.
  To mitigate this try increasing the scraping interval for the job.

* ### Not all logs/events are available
  Currently, not all logs/events from Auth0 are exposed, if a metric is not exposed, feel free to open a feature request.

## Prometheus

Example Prometheus configuration for the exporter. Replace `AUTH0-EXPORTER-HOSTNAME` with your instance's hostname.
```yaml
scrape_configs:
  - job_name: auth0_exporter
    metrics_path: /metrics
    static_configs:
      - targets: ['<<AUTH0-EXPORTER-HOSTNAME>>:9301']
    relabel_configs:
      - source_labels: [ __address__ ]
        target_label: __param_target
      - source_labels: [ __param_target ]
        target_label: instance
      - target_label: __address__
        replacement: <<AUTH0-EXPORTER-HOSTNAME>>:9301
```

## Development

#### Makefile

Similar to other Golang projects, this projects makes use of make for building and testing the source code.

#### Nix
Before start development, add your tenant's Auth0 credentials to the `env-dev.sh`, this help when developing using Nix.
Once the credentials are added, you can start the development environment by:
```shell
$ source env-dev.sh
$ develop
```
This will boot up a Nix devshell with the need tools and information.

## Contributing

**Everyone** is welcome to contribute to the project.

Please see [CONTRIBUTING.md](.github/CONTRIBUTING.md) for information on how to get started.

Feedback is always appreciated, whether it's a bug or feature request, feel free to open an issue using one of the templates.

## License
Apache 2.0, see [LICENSE.md](./LICENSE).
