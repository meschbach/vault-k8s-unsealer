#!/bin/bash

set -xe
go test ./...
go build -o backoff ./cmd/backoff
go build -o k8s-watcher ./cmd/k8s-watcher
