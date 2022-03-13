#!/bin/bash

set -xe
go test ./...
go build -o backoff ./cmd/backoff
