[![CI](https://github.com/tfadeyi/auth0-simple-exporter/actions/workflows/ci.yml/badge.svg)](https://github.com/tfadeyi/auth0-simple-exporter/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache_2.0-yellow.svg)](https://github.com/tfadeyi/auth0-simple-exporter/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/language-Go-blue.svg)](https://github.com/tfadeyi/auth0-simple-exporter)
[![GitHub release](https://img.shields.io/badge/release-0.0.6-green.svg)](https://github.com/tfadeyi/auth0-simple-exporter/releases)
# Auth0 Simple Log Exporter

A simple Prometheus exporter for Auth0 log [Events](https://auth0.com/docs/api/management/v2#!/Logs/get_logs), for a simple
way to monitor Auth0 from a Prometheus monitoring stack.

> Development is in progress.

## Pre-Requisites

* Auth0 tenant [management API](https://auth0.com/docs/api#management-api) client credentials.
* (Optional) Auth0 tenant management API [static token](https://auth0.com/docs/secure/tokens/access-tokens/management-api-access-tokens).

More info on how to get the credentials can be found [here](./docs/auth0.md).

## TL;DR
Run exporter's container with TLS disabled.

```shell
$ export TOKEN="< auth0 management API static static token >"
$ export DOMAIN="< auth0 tenant domain >"
$ docker run --network host -u $(id -u):$(id -g) -e TOKEN="$TOKEN" -e DOMAIN="$DOMAIN" ghcr.io/tfadeyi/auth0-simple-exporter:latest export --tls.disabled
```

## Get this image
The recommended way to get the Docker Image is to pull the prebuilt image from the project's Github Container Registry.
```shell
$ docker pull ghcr.io/tfadeyi/auth0-simple-exporter:latest
```
To use a specific version, you can pull a versioned tag.
```shell
$ docker pull ghcr.io/tfadeyi/auth0-simple-exporter:[TAG]
```

## Download

Binary can be downloaded from [Releases](https://github.com/tfadeyi/auth0-simple-exporter/releases) page.

## Helm
This is loading the secret obtained from create-service-account step 
`export HELM_SECRET="$(cat credentials.json)"`

Run exporter with TLS disabled.

```console
# Installing by passing in secret directly
helm upgrade --install --create-namespace -n jetstack-secure jetstack-agent \
  oci://eu.gcr.io/jetstack-secure-enterprise/charts/jetstack-agent \
  --set config.organisation="strange-jones" --set config.cluster="<CLUSTER_NAME>" \
  --set authentication.createSecret=true --set authentication.secretValue="$HELM_SECRET"
```

More info on the helm deployment can be found [here](./charts/auth0-exporter/README.md).

## Usage

```
Usage:
  exporter export [flags]

Flags:
      --auth0.client-id string       Auth0 management api client-id
      --auth0.client-secret string   Auth0 management api static token.
      --auth0.domain string          Auth0 tenant's domain. (i.e: <tenant_name>.eu.auth0.com)
      --auth0.token string           Auth0 management api static token
  -h, --help                         help for export
      --namespace string             Exporter's namespace
      --pprof                        Enabled pprof profiling on the exporter on port :6060. (help: https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/)
      --tls.cert-file string         The certificate file for the exporter.
      --tls.disabled
      --tls.key-file string          The key file for the exporter.
      --tls.auto                  Allow the exporter manage its own certificates.
      --web.listen-address int       Port where the server will listen. (default 8081)
      --web.metrics-path string      URL Path under which to expose metrics. (default "metrics")
```

Environment variables: 
* TOKEN, Auth0 management API static token.
* CLIENT_SECRET, Auth0 management API client-secret.
* CLIENT_ID, Auth0 management API client-id.
* DOMAIN, Auth0 tenant domain.

## Metrics

| Metric                                           | Meaning                                                  | Labels |
|--------------------------------------------------|----------------------------------------------------------|--------|
| `auth0_tenant_successful_sign_up_total`             | The number of successful signup operations. (codes: ss)  |        |
| `auth0_tenant_failed_sign_up_total`              | The number of failed signup operations. (codes: fs)      ||
| `auth0_tenant_successful_login_operations_total` | The number of successful login operations. (codes: s)    |        |
| `auth0_tenant_failed_login_operations_total` | The number of failed login operations. (codes: f,fp,fu)  | code   |


## Example queries

Retrieve the percentage of successful logins:

## Known Issues

When the Prometheus scraping job interval is too low the exporter might encounter api-rate limit from Auth0.
To mitigate this try increasing the scraping interval for the job.  

## Development

#### Nix
To start the development environment:
```shell
source env-dev.sh && develop
```
This will boot up a Nix devshell with a Prometheus instance running in the background,
`http://localhost:9090`.

## Release
The repo uses [goreleaser](https://goreleaser.com/) and [ko](https://ko.build/) to release the different artifacts.
To make a new release just create a new git tag, this will trigger a new Github action release [workflow](https://github.com/tfadeyi/auth0-simple-exporter/blob/main/.github/workflows/release.yml).

```shell
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```

## License
Apache 2.0, see [LICENSE.md](./LICENSE).
