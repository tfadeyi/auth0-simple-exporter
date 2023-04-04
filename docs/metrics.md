# Available Metrics

Exposes the Auth0 metrics collected by the exporter in a prometheus format.

| Metric                                          | Meaning                                                  | Labels |
|-------------------------------------------------|----------------------------------------------------------|--------|
| `auth0_tenant_successful_sign_up_total`         | The number of successful signup operations. (codes: ss)  |        |
| `auth0_tenant_failed_sign_up_total`             | The number of failed signup operations. (codes: fs)      ||
| `auth0_tenant_successful_login_operations_total` | The number of successful login operations. (codes: s)    |        |
| `auth0_tenant_failed_login_operations_total`    | The number of failed login operations. (codes: f,fp,fu)  | code   |
| `auth0_tenant_successful_logout_operations_total` | The number of successful logout operations. (codes: slo) |        |
| `auth0_tenant_failed_logout_operations_total`   | The number of failed logout operations. (codes: flo)     |        |
