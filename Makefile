.PHONY: test deps

LOCALS :=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.EXPORT_ALL_VARIABLES:
GO111MODULE = on

all: fmt deps test

deps:
	@go list github.com/mjibson/esc || go get github.com/mjibson/esc/...
	go generate -x ./...
	go get ./...
	-go mod tidy

fmt:
	gofmt -w $(LOCALS)
	go vet ./...

test:
	go test -count=1 ./...

