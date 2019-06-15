# Admission Webhook

Admission webhooks are HTTP callbacks that receive admission requests, process
them and return admission responses.

Kubernetes provides the following types of admission webhooks:

- **Mutating Admission Webhook**: It can mutate the object before it is stored.
It can be used to default fields in a resource requests, e.g. fields in
Deployment that are not specified by the user. It can be used to inject sidecar
containers.

- **Validating Admission Webhook**: It can validate the object and reject it if
invalid. It allows more complex validation than pure schema-based validation.
e.g. cross-field validation and pod image whitelisting.

The apiserver by default doesn't authenticate itself to the webhooks. But if you
want to authenticate the clients, you need to configure the apiserver to use
basic auth, bearer token, or a cert to authenticate itself to the webhooks. You
can find detailed steps
[here](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#authenticate-apiservers).
