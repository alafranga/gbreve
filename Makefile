NAME     := $(notdir $(PWD))
VERSION  := $(shell git describe --tags 2>/dev/null || echo 'v0.0.0')
BUILD    := $(shell date +%FT%T%z)

BINARIES := $(patsubst cmd/%/main.go,%,$(wildcard cmd/*/main.go))
LDFLAGS  := "-s -w -X main.version=$(VERSION) -X main.build=$(BUILD)"

# Build
build:
	@for bin in $(BINARIES); do go build -ldflags $(LDFLAGS) -o bin/$$bin cmd/$$bin/main.go; done

# Clean
clean:
	@for bin in $(BINARIES); do rm -f bin/$$bin; done
	@test -d bin && rmdir --ignore-fail-on-non-empty bin || true
	@rm -f $(NAME).tgz

# Compact
compact:
	@command -v upx >/dev/null && for bin in $(BINARIES); do upx -f --brute bin/$$bin; done || true

# Lint
lint:
	golangci-lint run

# Test
test:
	go test -parallel=4 ./...

# Zip
zip:
	sleep 0.1; tar czvf $(NAME).tgz -C bin .

.DEFAULT_GOAL := test
.PHONY: build clean compact lint test
