# Brave Sync Server v2

A sync server implemented in go to communicate with Brave sync clients using
[components/sync/protocol/sync.proto](https://cs.chromium.org/chromium/src/components/sync/protocol/sync.proto).
Chromium version for sync protocol buffer files: Chromium 80.0.3987.132.

## Developer Setup
1. [Install Go 1.14](https://golang.org/doc/install)
2. [Install GolangCI-Lint](https://github.com/golangci/golangci-lint#install)
3. Clone this repo
4. [Install protobuf protocol compiler](https://github.com/protocolbuffers/protocolbuffers/protobuf#protocol-compiler-installation) if you need to compile protobuf files
5. Build via `make`

## Local development using Docker and DynamoDB Local
1. Clone this repo
2. Run `make docker`
3. Run `make docker-up`
4. For running unit tests, run `make docker-test`
