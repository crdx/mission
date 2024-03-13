set quiet := true

BIN := 'mission'

[private]
help:
    just --list --unsorted

build:
    go build -trimpath -o dist/{{ BIN }}

fmt:
    go fmt ./...

lint:
    go vet ./...
    golangci-lint run

test:
    go test -cover ./...

clean:
    rm -vf dist/{{ BIN }}
