# Available Metrics

Exposes the Auth0 metrics collected by the exporter in a prometheus format.

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

### Voice call operations
| Metric                                           | Meaning                                                             | Labels |
|--------------------------------------------------|---------------------------------------------------------------------|--------|
| `tenant_send_voice_call_operations_total`        | The total number of voice_call operations.                          | client |
| `tenant_failed_send_voice_call_operations_total` | The number of voice_call operations. (codes: gd_send_voice_failure) | client |

### SMS operations
| Metric                                    | Meaning                                                                | Labels |
|-------------------------------------------|------------------------------------------------------------------------|--------|
| `tenant_send_sms_operations_total`        | The total number of successful send_sms operations.                    | client |
| `tenant_failed_send_sms_operations_total` | The number of failed send_sms operations. (codes: gd_send_sms_failure) | client |

### Email operations
| Metric                               | Meaning                                                                          | Labels |
|--------------------------------------|----------------------------------------------------------------------------------|--------|
| `tenant_send_email_operations_total` | The number of successful send email operations. (codes: gd_send_email)           | client |
| `tenant_change_email_total`          | The total number of change_user_email operations on the tenant. (codes: sce,fce) | client |
| `tenant_failed_change_email_total`   | The number of failed change user email operations on the tenant. (codes: fce)    | client |

### Push notification operations
| Metric                                  | Meaning                                                                                  | Labels |
|-----------------------------------------|------------------------------------------------------------------------------------------|--------|
| `tenant_push_notification_total`        | The total number of push_notification operations. (codes: gd_send_pn,gd_send_pn_failure) | client |
| `tenant_failed_push_notification_total` | The number of failed push_notification operations. (codes: gd_send_pn_failure)           | client |

### Password operations
| Metric                                          | Meaning                                                                                         | Labels |
|-------------------------------------------------|-------------------------------------------------------------------------------------------------|--------|
| `tenant_post_change_password_hook_total`        | The total number of post change user password hook operations on the tenant. (codes: scph,fcph) | client |
| `tenant_failed_post_change_password_hook_total` | The number of failed post change user password hook operations on the tenant. (codes: fcph)     | client |
| `tenant_change_password_request_total`          | The total number of change_password_request operations. (codes: scpr,fcpr)                      | client |
| `tenant_failed_change_password_request_total`   | The number of failed change password request operations. (codes: fcpr)                          | client |
| `tenant_change_password_total`                  | The total number of change_user_password operations on the tenant. (codes: scp,fcp)             | client |        |
| `tenant_failed_change_password_total`           | The number of failed change_user_password operations on the tenant. (codes: fcp)                | client |

### Passwordless Send code link
| Metric                        | Meaning                                                  | Labels      |
|-------------------------------|----------------------------------------------------------|-------------|
| `tenant_send_code_link_total` | The number of send_code_link operations. (codes: cls,cs) | type,client |

### Delete user operations
| Metric                            | Meaning                                                                   | Labels |
|-----------------------------------|---------------------------------------------------------------------------|--------|
| `tenant_delete_user_total`        | The total number of delete user operations on the tenant. (codes: du,fdu) | client |
| `tenant_failed_delete_user_total` | The number of failed delete user operations on the tenant. (codes: fdu)   | client |

### Change phone number operations
| Metric                                    | Meaning                                                                              | Labels |
|-------------------------------------------|--------------------------------------------------------------------------------------|--------|
| `tenant_change_phone_number_total`        | The total number of change_phone_number operations on the tenant. (codes: scpn,fcpn) | client |
| `tenant_failed_change_phone_number_total` | The number of failed change phone number operations on the tenant. (codes: fcpn)     | client |

### API operations
| Metric                               | Meaning                                                              | Labels |
|--------------------------------------|----------------------------------------------------------------------|--------|
| `tenant_api_operations_total`        | The total number of API operations on the tenant. (codes: sapi,fapi) | client |
| `tenant_failed_api_operations_total` | The number of failed API operations on the tenant. (codes: fapi)     | client |
