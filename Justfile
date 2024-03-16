set quiet := true

BIN := 'mission'

[private]
help:
    just --list --unsorted

fmt:
    go fmt ./...
    just --fmt

lint:
    #!/bin/bash
    set -eo pipefail
    unbuffer go vet ./... | gostack
    unbuffer golangci-lint --color never run | gostack

test:
    #!/bin/bash
    set -eo pipefail
    unbuffer go test -cover ./... | gostack --test
