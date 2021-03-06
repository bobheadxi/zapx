module go.bobheadxi.dev/zapx/zgcp

go 1.12

replace go.bobheadxi.dev/zapx/util => ../util

require (
	cloud.google.com/go v0.44.3
	github.com/pkg/errors v0.8.1 // indirect
	github.com/stretchr/testify v1.4.0
	go.bobheadxi.dev/zapx/util v0.6.4
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	google.golang.org/api v0.9.0
)
