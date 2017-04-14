all: fmt deps test

deps:
	go get .

clean:
	-rm -rf bin

fmt:
	gofmt -w .

test: fmt
	go test ./...
