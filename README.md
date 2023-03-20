# Brave Sync Server v2

A sync server implemented in go to communicate with Brave sync clients using
[components/sync/protocol/sync.proto](https://cs.chromium.org/chromium/src/components/sync/protocol/sync.proto).
Current Chromium version for sync protocol buffer files used in this repo is Chromium 109.0.5414.46.

This server supports endpoints as bellow.
- The `POST /v2/command/` endpoint handles Commit and GetUpdates requests from sync clients and return corresponding responses both in protobuf format. Detailed of requests and their corresponding responses are defined in `schema/protobuf/sync_pb/sync.proto`. Sync clients are responsible for generating valid access tokens and present them to the server in the Authorization header of requests.

Currently we use dynamoDB as the datastore, the schema could be found in `schema/dynamodb/table.json`.

## Developer Setup
1. [Install Go 1.14](https://golang.org/doc/install)
2. [Install GolangCI-Lint](https://github.com/golangci/golangci-lint#install)
3. Clone this repo
4. [Install protobuf protocol compiler](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation) if you need to compile protobuf files, which could be built using `make protobuf`.
5. Build via `make`

## Local development using Docker and DynamoDB Local
1. Clone this repo
2. Run `make docker`
3. Run `make docker-up`
4. For running unit tests, run `make docker-test`
