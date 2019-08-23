# zapx

Extensions for [`uber-go/zap`](https://github.com/uber-go/zap), a structured
logging library for [Go](https://golang.org/).

[![GoDoc](https://godoc.org/go.bobheadxi.dev/zapx?status.svg)](https://godoc.org/go.bobheadxi.dev/zapx)
[![Build Status](https://dev.azure.com/bobheadxi/bobheadxi/_apis/build/status/bobheadxi.zapx?branchName=master)](https://dev.azure.com/bobheadxi/bobheadxi/_build/latest?definitionId=6&branchName=master)
[![Go Report Card](https://goreportcard.com/badge/go.bobheadxi.dev/zapx)](https://goreportcard.com/report/go.bobheadxi.dev/zapx)
![Sourcegraph for Repo Reference Count](https://img.shields.io/sourcegraph/rrc/github.com/bobheadxi/zapx.svg)

## Usage

Each subpackage is a separate module to minimize their dependency trees. To use
the package you want, just import them using their respective package names -
for example:

```sh
go get go.bobheadxi.dev/zapx/zapx
go get go.bobheadxi.dev/zapx/ztest
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
