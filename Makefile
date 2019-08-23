GO111MODULE = on
PACKAGES = benchmarks util zapx zazure zgcp zgql zhttp zmetrics zpgx ztest

# mod manages dependencies for submodules
MOD_VENDOR = off
MOD_TARGETS = $(addprefix mod_, $(PACKAGES))
mod: $(MOD_TARGETS)
$(MOD_TARGETS): mod_%: %
	@echo "[INFO] updating package '$<'..."
	@(cd $< && go mod tidy)
	@if [ "$(MOD_VENDOR)" = "on" ]; then \
		(cd $< && go mod vendor) \
	fi

# test runs tests for submodules
TEST_FLAGS = -v -race
TEST_TARGETS = $(addprefix test_, $(PACKAGES))
test: $(TEST_TARGETS)
$(TEST_TARGETS): test_%: %
	@echo "[INFO] testing package '$<'..."
	@(cd $< && go test $(TEST_FLAGS) ./...)

# lint performs static analyses on submodules
LINT_TARGETS = $(addprefix lint_, $(PACKAGES))
lint: $(LINT_TARGETS)
$(LINT_TARGETS): lint_%: %
	@echo "[INFO] linting package '$<'..."
	cd $< && \
	go vet ./... && \
	go fmt ./... && \
	golint $(go list ./... | grep -v /vendor/) && \
	go test -run xxxx

# codecov uploads coverage reports to codecov.io for all submodules
CODECOV_TARGETS = $(addprefix codecov_, $(PACKAGES))
codecov: $(CODECOV_TARGETS)
$(CODECOV_TARGETS): codecov_%: %
	@echo "[INFO] uploading coverage report for package '$<'..."
	@(cd $< && bash <(curl -s https://codecov.io/bash) -t $(CODECOV_TOKEN))

# clean removes various bits and pieces
clean:
	@rm -rf \
		go.mod go.sum \
		vendor */vendor \
		coverage.txt */coverage.txt
