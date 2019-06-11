# What about builtin types?

Implementing webhooks for builtin-types requires a little more work to wire
things together, since you can implement neither the `Defaulter` interface nor
the `Validator` interface for a builtin type.

The steps you will need to take are:

1. Create the webhook server.
1. Ensure the server has been added in the manager.
1. Create a handler for your webhook.
1. Register the handler with a path in your server.

There is an [example](https://github.com/kubernetes-sigs/controller-runtime/tree/master/examples/builtins)
in the controller-runtime repo.
