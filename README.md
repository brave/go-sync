# Brave Sync Server v2

A sync server implemented in go to communicate with Brave sync clients using
[components/sync/protocol/sync.proto](https://cs.chromium.org/chromium/src/components/sync/protocol/sync.proto).
Current Chromium version for sync protocol buffer files used in this repo is Chromium 88.0.4324.96.

This server supports endpoints as bellow.
- The `POST /v2/command/` endpoint handles Commit and GetUpdates requests from sync clients and return corresponding responses both in protobuf format. Detailed of requests and their corresponding responses are defined in `schema/protobuf/sync_pb/sync.proto`. Sync clients are responsible for generating valid access tokens and present them to the server in the Authorization header of requests.

Currently we use dynamoDB as the datastore, the schema could be found in `schema/dynamodb/table.json`.

## Developer Setup

1. [Install Go 1.14](https://golang.org/doc/install)
2. [Install GolangCI-Lint](https://github.com/golangci/golangci-lint#install)
3. [Install gowrap](https://github.com/hexdigest/gowrap#installation)
4. Clone this repo
5. [Install protobuf protocol compiler](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation) if you need to compile protobuf files, which could be built using `make protobuf`.
6. [Install protoc-gen-go](https://developers.google.com/protocol-buffers/docs/reference/go-generated) for use with protobuf compiler to generate go code using `make protobuf`
7. Build via `make`

## Local development using Docker and DynamoDB Local

1. Clone this repo
2. Run `make docker`
3. Run `make docker-up`
4. For running unit tests, run `make docker-test`

## Prometheus Instrumentation

The instrumented datastore and redis interfaces are generated, providing integration with Prometheus.  The following will re-generate the instrumented code:

```
make instrumented
```

Changes to `datastore/datastore.go` or `cache/cache.go` should be followed with the above command.

## Protobuf files

Keeping protobuf files up to date requires synchronizing upstream changes from both `chromium` and `brave-core`.  `chromium` provides the base protobuf schemas, with `brave-core` providing Brave-specific schemas, as well as custom patches to the schemas provided by Chromium.

### Updating Protobuf Schemas

1. Clone `brave-browser`
2. [Install prerequisites](https://github.com/brave/brave-browser#install-prerequisites)
3. Run `npm install` from `brave-browser` root.
4. Run `npm run init` from `brave-browser` root to download `brave-core` and `chromium` source.
5. Apply patches from `brave-core` to protobuf schemas in Chromium source
6. Copy files from src/brave/src/components/sync/protocol/*.proto and src/brave/src/brave/components/sync/protocol/*.proto to go-sync/schemas/protobuf/sync_pb
7. Run find . | grep "\.proto$" | xargs sed -i "" 's/import \"components\/sync\/protocol\//import \"/g' to repath module imports
8. Run find . | grep "\.proto$" | xargs sed -i "" 's/import \"brave\/components\/sync\/protocol\//import \"/g' to repath module imports
