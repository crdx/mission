mod release

set quiet := true
set shell := ["bash", "-cu", "-o", "pipefail"]

BIN := 'mission'

[private]
help:
    just --list --unsorted

fmt:
    go fmt ./...

lint:
    unbuffer go vet ./... | gostack
    unbuffer golangci-lint run --color never | gostack

fix:
    unbuffer golangci-lint run --color never --fix | gostack

test:
    unbuffer go test -cover ./... | gostack --test
