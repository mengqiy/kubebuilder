# Implementing mutating and validating webhooks

If you want to implement mutating and validating webhooks for your CRD,
the only thing you need to do is to implement the `Defaulter` and (or)
the `Validator` interface.

Kubebuilder take care of the rest for you, such as

1. Creating the webhook server.
1. Ensuring the server has been added in the manager.
1. Creating handlers for your webhooks.
1. Registering each handler with a path in your server.

{{#literatego ./testdata/project/api/v1/cronjob_webhook.go}}

That was a doozy, but now we've got a working controller.  Let's test
against the cluster, then, if we don't have any issues, deploy it!
