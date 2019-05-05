all: build
	git commit -a -m "regenerate web app"

build:
	gobenchdata-web --title "zapx continuous benchmarks" --desc ""
