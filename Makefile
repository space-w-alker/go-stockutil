
LOCALS := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PKGS   := $(wildcard *util)
COUNT  ?= 1

.PHONY: test deps docs $(PKGS)

.EXPORT_ALL_VARIABLES:
GO111MODULE = on

all: fmt deps test

deps:
	@go list github.com/mjibson/esc || go get github.com/mjibson/esc/...
	go generate -x ./...
	go get ./...
	-go mod tidy

fmt:
	$(info Formatting)
	@gofmt -w $(LOCALS)
	$(info Vetting code)
	@go vet ./...

docs:
	owndoc render --property rootpath=/go-stockutil/

$(PKGS):
	$(info Testing $(@))
	@go test -count=$(COUNT) ./$(@)/...

test: $(PKGS)

