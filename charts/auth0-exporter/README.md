# auth0-exporter

Auth0 Prometheus exporter

![Version: 0.0.1](https://img.shields.io/badge/Version-0.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0.7](https://img.shields.io/badge/AppVersion-0.0.7-informational?style=flat-square)

## Additional Information

The Auth0 exporter helm chart installs the exporter on the target cluster, the collected metrics will exposed
on the `/metrics` endpoint.

## Installing the Chart

### Deploying the chart

There are two ways to install the chart:

#### Method 1: create secret manually

```
# pre-create secret
kubectl create secret -n auth0-exporter "<SOME_SECRET_NAME>" --from-file=credentials.json
# Install refering to secret
helm upgrade --install --create-namespace -n auth0-exporter auth0-exporter \
  oci://eu.gcr.io/jetstack-secure-enterprise/charts/jetstack-agent \
  --set config.organisation="strange-jones"  --set config.cluster="<CLUSTER_NAME>" \
  --set authentication.secretName="<SOME_SECRET_NAME>"
```

#### Method 2: Pass secret to chart as a value, it creates the secret

*This is loading the secret obtained from create-service-account step [above](#obtaining-credentials)
`export HELM_SECRET="$(cat credentials.json)"`*

```console
# Installing by passing in secret directly
helm upgrade --install --create-namespace -n jetstack-secure jetstack-agent \
  oci://eu.gcr.io/jetstack-secure-enterprise/charts/jetstack-agent \
  --set config.organisation="strange-jones" --set config.cluster="<CLUSTER_NAME>" \
  --set authentication.createSecret=true --set authentication.secretValue="$HELM_SECRET"
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| auth0 | object | `{"clientId":"","clientSecret":"","createSecret":true,"domain":"<change_me>.eu.auth0.com","secretName":"auth0-credentials","token":""}` | Exporter's Auth0 client configuration |
| auth0.clientId | string | `""` | Auth0 management api client-id. (do not set if static token is already set) |
| auth0.clientSecret | string | `""` | Auth0 management api client-secret. (do not set if static token is already set) |
| auth0.domain | string | `"<change_me>.eu.auth0.com"` | Auth0 tenant's domain. (i.e: <tenant_name>.eu.auth0.com) |
| auth0.token | string | `""` | Auth0 management api static token. (the token can be used instead of client credentials) |
| exporter | object | `{"logLevel":1,"metricsEndpoint":"metrics","namespace":"","port":9301,"pprof":false,"tls":{"auto":false,"certFile":"","createSecret":false,"disabled":false,"hosts":[],"keyFile":"","secretKey":"","secretName":""}}` | Exporter's configuration |
| exporter.metricsEndpoint | string | `"metrics"` | URL Path under which to expose the collected auth0 metrics. |
| exporter.port | int | `9301` | Port where the server will listen. |
| exporter.pprof | bool | `false` | Enabled pprof profiling on the exporter on port :6060. (help: https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/) |
| exporter.tls | object | `{"auto":false,"certFile":"","createSecret":false,"disabled":false,"hosts":[],"keyFile":"","secretKey":"","secretName":""}` | Exporter's TLS configuration |
| exporter.tls.auto | bool | `false` | Allow the exporter to use autocert to renew its certificates with letsencrypt. (can only be used if the exporter is publicly accessible by the internet) |
| exporter.tls.certFile | string | `""` | The certificate file for the exporter TLS connection. |
| exporter.tls.disabled | bool | `false` | Run exporter without TLS. |
| exporter.tls.keyFile | string | `""` | The key file for the exporter TLS connection. |
| fullnameOverride | string | `""` | Helm default setting, use this to shorten install name |
| image | object | `{"pullPolicy":"IfNotPresent","repository":"ghcr.io/tfadeyi/auth0-simple-exporter","tag":"v0.0.7"}` | image settings |
| imagePullSecrets | list | `[]` | specify credentials if pulling from a customer registry |
| labels | object | `{}` |  |
| nameOverride | string | `""` | Helm default setting to override release name, leave blank |
| nodeSelector | object | `{}` |  |
| podAnnotations | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| replicaCount | int | `1` |  |
| resources.limits.cpu | string | `"100m"` |  |
| resources.limits.memory | string | `"128Mi"` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"128Mi"` |  |
| securityContext.capabilities.drop[0] | string | `"ALL"` |  |
| securityContext.readOnlyRootFilesystem | bool | `true` |  |
| securityContext.runAsNonRoot | bool | `true` |  |
| securityContext.runAsUser | int | `1000` |  |
| service.port | int | `9301` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` |  |
| tolerations | list | `[]` |  |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs vv1.11.0](https://github.com/norwoodj/helm-docs/releases/vv1.11.0)
