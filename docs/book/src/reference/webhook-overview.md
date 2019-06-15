# Webhook

Webhooks are HTTP callbacks, providing a way for notifications to be delivered
to an external web server. A web application implementing webhooks will send an
HTTP request (typically POST) to other application when certain event happens.
In the kubernetes world, there are 3 kinds of webhooks:
[admission webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#admission-webhooks),
[authorization webhook](https://kubernetes.io/docs/reference/access-authn-authz/webhook/) and
[CRD conversion webhook](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definition-versioning/#webhook-conversion).

In [controller-runtime](https://godoc.org/sigs.k8s.io/controller-runtime/pkg/webhook)
libraries, we support admission webhooks and CRD conversion webhooks.

Admission webhook feature is beta since kubernetes 1.9. CRD conversion webhook
feature is alpha since kubernetes 1.13 and beta in kubernetes 1.15.
