all: vendor fmt test

update:
	glide up

vendor:
	go list github.com/Masterminds/glide
	glide install

clean:
	-rm -rf vendor bin

fmt:
	gofmt -w .

test: fmt
	go test ./maputil
	go test ./pathutil
	go test ./sliceutil
	go test ./stringutil
