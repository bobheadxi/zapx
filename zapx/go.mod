module go.bobheadxi.dev/zapx/zapx

go 1.12

replace go.bobheadxi.dev/zapx/internal => ../internal

require (
	github.com/pkg/errors v0.8.1 // indirect
	github.com/stretchr/testify v1.4.0
	go.bobheadxi.dev/zapx/internal v0.0.0-20190823225728-d75f7d59f073
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
)
