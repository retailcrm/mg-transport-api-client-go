# Usage examples

This directory contains examples on how to use the library properly. Currently two examples are available:
- [Simple Telegram transport (only text messages are supported)](telegram)
- [MG Webhook processing example](webhooks)

## `telegram`

How to run the example:
1. Copy `config.json.dist` to `config.json`.
2. Replace placeholder values in the `config.json` with your data. You'll need a service like ngrok.io for that.
3. Navigate to `telegram` directory via terminal and run `go run ./...`

The sample will automatically register itself in the target system and will begin transferring the messages between Telegram and MessageGateway.

## `webhook`

You can run this example by executing `go run /...`. Transport API webhooks from MessageGateway should be sent to the `/api/v1/webhook` route.
