all: fmt deps test

deps:
	@go list golang.org/x/tools/cmd/goimports || go get golang.org/x/tools/cmd/goimports
	go generate -x
	go get .

fmt:
	goimports -w .
	go vet .


clean:
	-rm -rf bin

test: fmt
	go test ./...
