BIN := 'mission'

@_help:
    just --list --unsorted

# build binary
@make:
    mkif dist/{{ BIN }} $(find -type f) -x 'just remake'

# rebuild binary
remake:
    #!/bin/bash
    set -e
    go fmt ./...
    go vet ./...
    export CGO_ENABLED=0
    go build -trimpath -ldflags '-s -w' -o dist/{{ BIN }}

# remove build
@clean:
    rm -fv dist/{{ BIN }}

# run tests
@test *args:
    go test -cover ./... {{ args }}

# run linter
@lint:
    golangci-lint run
