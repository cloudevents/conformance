# CloudEvents Conformance Testing

`cloudevents` is a tool for testing CloudEvents receivers.

[![GoDoc](https://godoc.org/github.com/cloudevents/conformance?status.svg)](https://godoc.org/github.com/cloudevents/conformance)
[![Go Report Card](https://goreportcard.com/badge/cloudevents/conformance)](https://goreportcard.com/report/cloudevents/conformance)

_Work in progress._

## Installation

`cloudevents` can be installed via:

```shell
go get github.com/cloudevents/conformance/cmd/cloudevents
```

To update your installation:

```shell
go get -u github.com/cloudevents/conformance/cmd/cloudevents
```

## Usage

`cloudevents` has three commands at the moment: `send`, `invoke` and `listen`.

### Send

`send` will do a one-off creation of a cloudevent and send to a given target.

```shell script
cloudevents send http://localhost:8080 --id abc-123 --source cloudevents.conformance.tool --type foo.bar
```

### Invoke

`invoke` will read yaml files, convert them to http and send them to the given
target.

```shell script
cloudevents invoke http://localhost:8080 -f ./yaml/v0.3
```

### Listen

`listen` will accept http request and write the converted yaml to stdout.

```shell script
cloudevents listen -v > got.yaml
```

Optionally, you can forward the incoming request to a target.

```shell script
cloudevents listen -v -t http://localhost:8181 > got.yaml
```

## Advanced Usage

If you would like to produce a pre-produced event yaml file, you can use
`listen` to collect requests. This works with both running event producers that
can be directed at the `listen` port or directly with `send`.

## TODO

Feature requests:

- [x] Add a `yaml` command that allows the same format of send but direct to
      yaml.
- [x] Add a `--yaml` flag to `send` that outputs what the definition of the sent
      event would look like in yaml.
- [x] Add `-f -` support for reading from STDIN as the file for yaml, this
      allows for `send` | `invoke` chaining.
- [ ] Add format support to `send` to select binary or structured.
