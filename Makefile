
LOCALS := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PKGS   := $(wildcard *util)

.PHONY: test deps $(PKGS)

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

$(PKGS):
	$(info Testing $(@))
	@go test -count=1 ./$(@)/...

test: $(PKGS)

