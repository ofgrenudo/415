MODULE     := github.com/ofgrenudo/415
BIN_DIR    := bin
BINARIES   := $(notdir $(wildcard cmd/*))

# Overridable from the environment/CLI, e.g. by CI after release-please
# has computed the next version:
#   export VERSION=$(cat .version) && make build
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT     ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE       ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
AUTHOR     ?= $(shell git config user.name 2>/dev/null || echo unknown)

LDFLAGS := -X '$(MODULE)/pkg/version.Version=$(VERSION)' \
           -X '$(MODULE)/pkg/version.Commit=$(COMMIT)' \
           -X '$(MODULE)/pkg/version.Date=$(DATE)' \
           -X '$(MODULE)/pkg/version.Author=$(AUTHOR)'

.PHONY: all build clean $(BINARIES)

all: build

build: $(BINARIES)

$(BINARIES):
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$@ ./cmd/$@

clean:
	rm -rf $(BIN_DIR)
