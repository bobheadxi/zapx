module go.bobheadxi.dev/zapx/zhttp

go 1.12

replace (
	go.bobheadxi.dev/zapx/util => ../util
	go.bobheadxi.dev/zapx/zapx => ../zapx
	go.bobheadxi.dev/zapx/ztest => ../ztest
)

require (
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/stretchr/testify v1.4.0
	go.bobheadxi.dev/res v0.2.0
	go.bobheadxi.dev/zapx/util v0.0.0-20190823225728-d75f7d59f073
	go.bobheadxi.dev/zapx/zapx v0.0.0-20190823225728-d75f7d59f073
	go.bobheadxi.dev/zapx/ztest v0.0.0-20190823225728-d75f7d59f073
	go.uber.org/zap v1.10.0
	golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7 // indirect
)
