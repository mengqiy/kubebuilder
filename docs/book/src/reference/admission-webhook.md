# Admission Webhook

Admission webhooks are HTTP callbacks that receive admission requests, process them and return admission responses.
There are two types of admission webhooks: mutating admission webhook and validating admission webhook.
With mutating admission webhooks, you may change the request object before it is stored (e.g. for implementing defaulting of fields)
With validating admission webhooks, you may not change the request, but you can reject it (e.g. for implementing validation of the request).

### Why Admission Webhooks are Important

Admission webhooks are the mechanism to enable kubernetes extensibility through CRD.
- Mutating admission webhook is the only way to do defaulting for CRDs.
- Validating admission webhook allows for more complex validation than pure schema-based validation.
e.g. cross-field validation or cross-object validation.

It can also be used to add custom logic in the core kubernetes API.

### Mutating Admission Webhook

A mutating admission webhook receives an admission request which contains an object.
The webhook can either decline the request directly or returning JSON patches for modifying the original object.
- If admitting the request, the webhook is responsible for generating JSON patches and send them back in the
admission response.
- If declining the request, a reason message should be returned in the admission response.

### Validating Admission Webhook

A validating admission webhook receives an admission request which contains an object.
The webhook can either admit or decline the request.
A reason message should be returned in the admission response if declining the request.

### Authentication

The apiserver by default doesn't authenticate itself to the webhooks.
That means the webhooks don't authenticate the identities of the clients.

But if you want to authenticate the clients, you need to configure the apiserver to use basic auth, bearer token,
or a cert to authenticate itself to the webhooks. You can find detailed steps
[here](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#authenticate-apiservers).
