GO111MODULE = on
PACKAGES = benchmarks util zapx zazure zgcp zgql zhttp zmetrics zpgx ztest

MOD_VENDOR = off
MOD_TARGETS = $(addprefix mod_, $(PACKAGES))
mod: $(MOD_TARGETS)
$(MOD_TARGETS): mod_%: %
	@echo "[INFO] updating package '$<'..."
	@(cd $< && go mod tidy)
	@if [ "$(MOD_VENDOR)" = "on" ]; then \
		(cd $< && go mod vendor) \
	fi

TEST_FLAGS = -v -race
TEST_TARGETS = $(addprefix test_, $(PACKAGES))
test: $(TEST_TARGETS)
$(TEST_TARGETS): test_%: %
	@echo "[INFO] testing package '$<'..."
	@(cd $< && go test $(TEST_FLAGS) ./...)

CODECOV_TARGETS = $(addprefix codecov_, $(PACKAGES))
codecov: $(CODECOV_TARGETS)
	@echo "[INFO] uploading coverage report for package '$<'..."
	@(cd $< && bash <(curl -s https://codecov.io/bash) -t $(CODECOV_TOKEN))

clean:
	@rm -rf \
		go.mod go.sum \
		vendor */vendor \
		coverage.txt */coverage.txt
