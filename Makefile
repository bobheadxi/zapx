all: build
	git commit -a -m "regenerate web app"

build:
	gobenchdata-web -o benchmarks --title 'zapx continuous benchmarks' --desc 'benchmarks for <a href="https://github.com/bobheadxi/zapx">bobheadxi/zapx</a>, a library of extensions for <a href="https://github.com/uber-go/zap">uber-go/zap</a>' --source github.com/bobheadxi/zapx --canonical-import go.bobheadxi.dev/zapx

