## Dependencies

- [golang](https://golang.org/doc/install)

## Run Tests

```
$ go test ./...
```

## Running Locally

```
$ go run cmd/api/main.go
```

Now point your browser to: http://localhost:8080

## Build

```
$ scripts/build
```

## Static checks

```
$ scripts/lint
```

## REST API

Documentation TBD, see [service source](services/sets.go) for now.

## TODO

See [TODO](TODO.md).

## Style Guide

- Use https://jsonapi.org for JSON API conventions.
