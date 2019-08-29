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

`cloudevents` has one command at the moment: `invoke`. It will read yaml files, convert them to http 
and send them to the target given.  

```shell
cloudevents invoke http://localhost:8080 -f ./yaml/v0.3
```
