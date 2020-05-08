all: build
	git commit -a -m "regenerate web app"

build:
	gobenchdata web generate ./benchmarks --web.title='zapx continuous benchmarks'
