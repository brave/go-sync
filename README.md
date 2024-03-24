# Brave Sync Server v2

A sync server implemented in go to communicate with Brave sync clients using
[components/sync/protocol/sync.proto](https://cs.chromium.org/chromium/src/components/sync/protocol/sync.proto).
Current Chromium version for sync protocol buffer files used in this repo is Chromium 116.0.5845.183.

This server supports endpoints as bellow.
- The `POST /v2/command/` endpoint handles Commit and GetUpdates requests from sync clients and return corresponding responses both in protobuf format. Detailed of requests and their corresponding responses are defined in `schema/protobuf/sync_pb/sync.proto`. Sync clients are responsible for generating valid access tokens and present them to the server in the Authorization header of requests.

Currently we use dynamoDB as the datastore, the schema could be found in `schema/dynamodb/table.json`.

## Developer Setup
1. [Install Go 1.18](https://golang.org/doc/install)
2. [Install GolangCI-Lint](https://github.com/golangci/golangci-lint#install)
3. [Install gowrap](https://github.com/hexdigest/gowrap#installation)
4. Clone this repo
5. [Install protobuf protocol compiler](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation) if you need to compile protobuf files, which could be built using `make protobuf`.
6. Build via `make`

## Local development using Docker and DynamoDB Local
1. Clone this repo
2. Run `make docker`
3. Run `make docker-up`
4. For running unit tests, run `make docker-test`

## Updating protocol definitions
1. Copy the `.proto` files from `components/sync/protocol` in `chromium` to `schema/protobuf/sync_pb` in `go-sync`.
2. Copy the `.proto` files from `components/sync/protocol` in `brave-core` to `schema/protobuf/sync_pb` in `go-sync`.
3. Run `make repath-proto` to set correct import paths in `.proto` files.
4. Run `make proto-go-module` to add the `go_module` option to `.proto` files.
5. Run `make protobuf` to generate the Go code from `.proto` definitions.

## Prometheus Instrumentation
The instrumented datastore and redis interfaces are generated, providing integration with Prometheus.  The following will re-generate the instrumented code, required when updating protocl definitions:

```
make instrumented
```

Changes to `datastore/datastore.go` or `cache/cache.go` should be followed with the above command.

## Self-hosting

### Setting up the servers with persistant storage
1. Run the folling commands to get two containers, `brave-sync:latest` and `brave-dynamo:latest`:
    ```
    GIT_VERSION=$(git describe --abbrev=8 --dirty --always --tags)
    GIT_COMMIT=$(git rev-parse --short HEAD)
    BUILD_TIME=$(date +%s)
    docker-compose build
    docker tag go-sync_web:latest brave-sync:latest
    docker tag go-sync_dynamo-local:latest brave-dynamo:latest
    ```
2. Copy the `docker-compose-self-host.yml` to wherever you wish to host your project as a `docker-compose.yml` file.
3. On your server, get a copy of the initial Brave Sync Dynamo DB out of the container:
    ```
    docker run --rm -t --name get-db -d brave-dynamo:latest
    mkdir -p /data/containers/brave/dynamo
    docker cp get-db:/db/shared-local-instance.db /data/containers/brave/dynamo/
    docker stop get-db
    chown 1000:1000 /data/containers/brave/dynamo/ -R
    ```
4. Either uncomment the `ports` section in the docker compose file and point your SSL proxy to `http://localhost:8295`, or run the SSL proxy inside the `brave` docker network and point it to `http://brave-sync:8295`.
5. Change the `SET_TO_SOMETHING_SECURE` value to something secure.
6. Run `docker compose up` from that new server.

### Using clients

#### Desktop
This command line option must be supplied every time you start Brave.
1. Start with `--sync-url="https://my.brave.sync.url/v2"`
2. Confirm the setting by visiting `brave://sync-internals/`

#### Android
This setting may not persist after a reboot on all devices, use at your own risk. Work is ongoing for a more formal way to [persist the sync URL](https://github.com/brave/brave-browser/issues/12314).
1. Enable `brave://flags/#enable-command-line-on-non-rooted-devices`
2. Create `/data/local/tmp/chrome-command-line` and add `--sync-url=https://my.brave.sync.url/v2` to it starting with an underscore over adb:
    ```
    adb shell
    echo "_ --sync-url=https://my.brave.sync.url/v2" > /data/local/tmp/chrome-command-line
    ```
3. Set up sync as normal on the device