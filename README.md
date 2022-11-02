# CloudEvents Conformance Testing

`cloudevents` is a tool for testing CloudEvents receivers.

[![GoDoc](https://godoc.org/github.com/cloudevents/conformance?status.svg)](https://godoc.org/github.com/cloudevents/conformance)
[![Go Report Card](https://goreportcard.com/badge/cloudevents/conformance)](https://goreportcard.com/report/cloudevents/conformance)

_Work in progress._

## Installation

The latest `cloudevents` can be installed via:

```shell
go install github.com/cloudevents/conformance/cmd/cloudevents@latest
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

### Diff

`diff` compares two yaml event files.

```shell script
cloudevents diff ./want.yaml ./got.yaml
```

`want.yaml` could have fewer fields specified to allow for fuzzy matching. 

Example, if you only wanted to compare on `type` and ignore additional fields:

```shell script
$ cat ./want.yaml
ContextAttributes:
  type: com.example.someevent
$ cat ./got.yaml
Mode: structured
ContextAttributes:
  specversion: 1.0
  type: com.example.someevent
  time: 2018-04-05T03:56:24Z
  id: 4321-4321-4321-a
  source: /mycontext/subcontext
  Extensions:
    comexampleextension1 : "value"
    comexampleextension2 : |
      {"othervalue": 5}
TransportExtensions:
  user-agent: "foo"
Data: |
  {"world":"hello"}

$ cloudevents diff ./want.yaml ./got.yaml --match type --ignore-additions
```

This validates that at least one event of type `com.example.someevent` is present in the `got.yaml` file.

## Advanced Usage

If you would like to produce a pre-produced event yaml file, you can use
`listen` to collect requests. This works with both running event producers that
can be directed at the `listen` port or directly with `send`.
