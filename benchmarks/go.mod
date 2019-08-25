module go.bobheadxi.dev/zapx/benchmarks

go 1.12

replace (
	go.bobheadxi.dev/zapx/internal => ../internal
	go.bobheadxi.dev/zapx/zapx => ../zapx
	go.bobheadxi.dev/zapx/zgcp => ../zgcp
)

require (
	github.com/stretchr/testify v1.4.0
	go.bobheadxi.dev/res v0.2.0
	go.bobheadxi.dev/zapx/internal v0.6.1
	go.bobheadxi.dev/zapx/zapx v0.6.1
	go.bobheadxi.dev/zapx/zgcp v0.6.1
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	google.golang.org/api v0.9.0
)
