GO=go

# go install gotest.tools/gotestsum@latest
ifeq ($(shell which gotestsum),)
GOTEST=$(GO) test
else
GOTEST=gotestsum -f testname --
endif

TAGS?=
OUTDIR?=bin

# https://go.dev/doc/gdb
# https://github.com/golang/vscode-go/blob/master/docs/debugging.md
ifndef GO_BUILD_OPTS
ifdef DEBUG
GO_BUILD_OPTS=-gcflags=all="-N -l"
else
GO_BUILD_OPTS=-trimpath
endif
endif

.PHONY: dd
dd: bin/dd

# 1. Detect unformatted Go files
# 2. Run shellcheck (shell scripts linter)
# 3. Download latest web platform
# 4. Rebuild protocol buffer stubs
# 5. Build the entire Go codebase
# 6. Run golangci-lint (Go linters)
# 7. Build AK binary with version and/or debug info
# 8. Run all automated tests (unit + integration)
all: lint build bin/dd test

.PHONY: clean
clean:
	rm -rf $(OUTDIR)

.PHONY: bin
bin: bin/dd

.PHONY: bin/dd
bin/dd:
	$(GO) build --tags "${TAGS}" -o "$@" -ldflags="$(LDFLAGS)" $(GO_BUILD_OPTS) ./cmd/$(shell basename $@)

.PHONY: build
build:
	mkdir -p $(OUTDIR)
	$(GO) build $(GO_BUILD_OPTS) ./...


golangci_lint=$(shell which golangci-lint)

# Based on: https://golangci-lint.run/welcome/install/#other-ci
# See: https://github.com/golangci/golangci-lint/releases
$(OUTDIR)/tools/golangci-lint:
	mkdir -p $(OUTDIR)/tools
ifeq ($(golangci_lint),)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(OUTDIR)/tools" v1.64.5
else
	ln -fs $(golangci_lint) $(OUTDIR)/tools/golangci-lint
endif

.PHONY: lint
lint: $(OUTDIR)/tools/golangci-lint
	$(OUTDIR)/tools/golangci-lint run

scripts=$(shell find . -name \*.sh -not -path "*/.venv/*")

.PHONY: shellcheck
shellcheck:
ifneq ($(scripts),)
	docker run --rm -v $(shell pwd):/src -w /src koalaman/shellcheck:stable -a $(scripts) -x
endif

.PHONY: test
test:
	$(GOTEST) -race ./...
