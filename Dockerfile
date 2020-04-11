FROM golang:1.14

ADD . /go/src/github.com/sosedoff/wg-registry
WORKDIR /go/src/github.com/sosedoff/wg-registry

RUN \
  GOOS=linux \
  GOARCH=amd64 \
  CGO_ENABLED=1 \
  go build -o /artifacts/wg-registry