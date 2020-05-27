.PHONY: all build test lint clean

all: lint test build

protobuf:
	protoc -I schema/protobuf/sync_pb/ schema/protobuf/sync_pb/*.proto --go_out=schema/protobuf/sync_pb/

build:
	go run main.go

test:
	go test -v ./...

lint:
	golangci-lint run -E gofmt -E golint --exclude-use-default=false

clean:
	rm -f sync-server

docker:
	docker-compose build

docker-up:
	docker-compose up

docker-test:
	docker-compose -f docker-compose.yml run --rm dev make test

instrumented:
	gowrap gen -p github.com/brave/go-sync/datastore -i Datastore -t ./.prom-gowrap.tmpl -o ./datastore/instrumented_datastore.go
