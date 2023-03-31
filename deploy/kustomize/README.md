# auth0-exporter

To deploy with kustomize, make sure to update the `secrets/.env.secrets` file with the Auth0 token
or client credentials information then run:

```shell
kubectl -k apply -f .
```