BIN := 'mission'

@_help:
    just --list --unsorted

# build binary
@make:
    mkif dist/{{ BIN }} $(find src -type f) -x 'just remake'

# rebuild binary
remake:
    #!/bin/bash
    set -e
    cd src
    go fmt ./...
    go vet ./...
    export CGO_ENABLED=0
    go build -trimpath -ldflags '-s -w' -o ../dist/{{ BIN }}

# remove build
@clean:
    rm -fv dist/*

# run tests
@test *args:
    cd src && go test ./... {{ args }}

# run linter
@lint:
    cd src && golangci-lint run
