# zapx

Extensions for [`uber-go/zap`](https://github.com/uber-go/zap), a structured
logging library for [Go](https://golang.org/).

[![GoDoc](https://godoc.org/go.bobheadxi.dev/zapx?status.svg)](https://godoc.org/go.bobheadxi.dev/zapx)
[![Build Status](https://dev.azure.com/bobheadxi/bobheadxi/_apis/build/status/bobheadxi.zapx?branchName=master)](https://dev.azure.com/bobheadxi/bobheadxi/_build/latest?definitionId=6&branchName=master)
[![Go Report Card](https://goreportcard.com/badge/go.bobheadxi.dev/zapx)](https://goreportcard.com/report/go.bobheadxi.dev/zapx)
[![Benchmarks](https://img.shields.io/website/https/zapx.bobheadxi.dev/benchmarks?down_color=lightgrey&down_message=offline&label=benchmarks&up_color=green&up_message=available)](https://zapx.bobheadxi.dev/benchmarks)

## Usage

Each subpackage is a separate module to minimize their dependency trees. To use
the package you want, just import them using their respective package names -
for example:

```sh
go get go.bobheadxi.dev/zapx/zapx
go get go.bobheadxi.dev/zapx/ztest
go get go.bobheadxi.dev/zapx/zhttp
```

Refer to the [godoc](https://godoc.org/go.bobheadxi.dev/zapx) for a complete
listing of available packages and their functionality.

## Development

A few Makefile targets are available to help with development:

```sh 
make mod  # updates module definitions for all submodule
make test # runs tests for each submodule
```

Refer to the [Makefile](./Makefile) for more details.
