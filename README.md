<div align="center">

# Auth0 Exporter

[![release pipeline](https://github.com/tfadeyi/auth0-simple-exporter/actions/workflows/release.yml/badge.svg?style=flat-square)](https://github.com/tfadeyi/auth0-simple-exporter/actions/workflows/release.yml)
[![Continuous Integration](https://github.com/tfadeyi/auth0-simple-exporter/actions/workflows/ci.yml/badge.svg?style=flat-square)](https://github.com/tfadeyi/auth0-simple-exporter/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache_2.0-yellowgreen.svg?style=flat-square)](https://github.com/tfadeyi/auth0-simple-exporter/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/language-Go-blue.svg?style=flat-square)](https://github.com/tfadeyi/auth0-simple-exporter)
[![GitHub release](https://img.shields.io/badge/release-0.0.2-green.svg?style=flat-square)](https://github.com/tfadeyi/auth0-simple-exporter/releases)

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
Run exporter's container with TLS disabled.

```shell
$ export TOKEN="< auth0 management API static static token >"
$ export DOMAIN="< auth0 tenant domain >"
$ docker run --network host -u $(id -u):$(id -g) -e TOKEN="$TOKEN" -e DOMAIN="$DOMAIN" ghcr.io/tfadeyi/auth0-simple-exporter:latest export --tls.disabled
```

## Installation

* ### Download Binaries

    Binaries can be downloaded from [Releases](https://github.com/tfadeyi/auth0-simple-exporter/releases) page.

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
    $ export TOKEN="< auth0 management API static static token >"
    $ export DOMAIN="< auth0 tenant domain >"
    ```
    ```shell
    # Installing by passing in secret directly
    helm upgrade --install --create-namespace -n auth0-exporter auth0-exporter \
      https://tfadeyi.github.io/charts \
      --set auth0.domain="$DOMAIN" --set auth0.token="$TOKEN" \
      --set exporter.tls.disabled=true
    ```

    More info on the helm deployment can be found [here](deploy/charts/auth0-exporter/README.md).

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

Monitor the percentage of successful logins:

```
(tenant_successful_login_operations_total / (on job,instance) (tenant_successful_login_operations_total + tenant_failed_login_operations_total)) * 100
```

Monitor the current logged-in users:

```
(on job,instance) (tenant_successful_login_operations_total - tenant_successful_logout_operations_total)
```

## Metrics

| Metric                                           | Meaning                                                  | Labels |
|--------------------------------------------------|----------------------------------------------------------|--------|
| `tenant_successful_sign_up_total`           | The number of successful signup operations. (codes: ss)  |        |
| `tenant_failed_sign_up_total`               | The number of failed signup operations. (codes: fs)      ||
| `tenant_successful_login_operations_total`  | The number of successful login operations. (codes: s)    |        |
| `tenant_failed_login_operations_total`      | The number of failed login operations. (codes: f,fp,fu)  | code   |
| `tenant_successful_logout_operations_total` | The number of successful logout operations. (codes: slo) |        |
| `tenant_failed_logout_operations_total`    | The number of failed logout operations. (codes: flo)     |        |

## Known Issues

### API Rate Limits
When the Prometheus scraping job interval is too frequent the exporter might encounter api-rate limit from Auth0.
To mitigate this try increasing the scraping interval for the job.

### Not all logs/events are available
Currently, not all logs/events from Auth0 are exposed, if a metric is not exposed, feel free to open a feature request.


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

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for information on how to get started.

Feedback is always appreciated, whether it's a bug or feature request, feel free to open an issue using one of the templates.

## License
Apache 2.0, see [LICENSE.md](./LICENSE).
