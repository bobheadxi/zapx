GO111MODULE = on

PACKAGES = benchmarks util zapx zazure zgcp zgql zhttp zmetrics zpgx ztest
.PHONY: $(PACKAGES)

mod: $(PACKAGES)

$(PACKAGES):
	(cd $@ ; go mod tidy)
