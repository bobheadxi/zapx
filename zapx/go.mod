module go.bobheadxi.dev/zapx/zapx

go 1.12

replace (
	go.bobheadxi.dev/zapx/util => ../util
	go.bobheadxi.dev/zapx/ztest => ../ztest
)

require (
	github.com/stretchr/testify v1.4.0
	go.bobheadxi.dev/zapx/util v0.6.4
	go.bobheadxi.dev/zapx/ztest v0.6.4
	go.uber.org/zap v1.10.0
)
