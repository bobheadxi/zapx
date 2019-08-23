module go.bobheadxi.dev/zapx/benchmarks

go 1.12

replace (
	go.bobheadxi.dev/zapx/zapx => ../zapx
	go.bobheadxi.dev/zapx/zgcp => ../zgcp
)

require (
	github.com/pkg/errors v0.8.1 // indirect
	github.com/stretchr/testify v1.4.0
	go.bobheadxi.dev/res v0.2.0
	go.bobheadxi.dev/zapx/zapx v0.0.0-00010101000000-000000000000
	go.bobheadxi.dev/zapx/zgcp v0.0.0-00010101000000-000000000000
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0
	go.uber.org/zap v1.10.0
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/api v0.9.0
)
