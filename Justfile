set quiet := true

BIN := 'mission'

[private]
help:
    just --list --unsorted

fmt:
    go fmt ./...
    just --fmt

vet:
    #!/bin/bash
    set -eo pipefail
    unbuffer go vet ./... | gostack

lint:
    #!/bin/bash
    set -eo pipefail
    unbuffer golangci-lint --color never run | gostack

test:
    go test -cover ./...
