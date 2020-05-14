GIT_COMMIT = $(shell git rev-parse HEAD)
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" | tr -d '\n')
GO_VERSION = $(shell go version | awk {'print $$3'})

.PHONY: build assets docker-build linux-binary

build:
	go build

fmt:
	go fmt ./...

test:
	go test -cover -race ./...

assets:
	go-assets-builder static -p assets -o assets/assets.go

docker-build:
	docker build -t wg-registry .

linux-binary: docker-build
	docker run -it -d --name=tmp wg-registry bash
	docker cp tmp:/artifacts/wg-registry ./dist/
	docker rm -f tmp