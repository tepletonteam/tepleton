.PHONY: all test get_deps

all: test install

NOVENDOR = go list github.com/tepleton/basecoin/... | grep -v /vendor/

build:
	go build github.com/tepleton/basecoin/cmd/...

install:
	go install github.com/tepleton/basecoin/cmd/...

test:
	go test --race `${NOVENDOR}`
	#go run tests/tepleton/*.go

get_deps:
	go get -d github.com/tepleton/basecoin/...

update_deps:
	go get -d -u github.com/tepleton/basecoin/...

get_vendor_deps:
	go get github.com/Masterminds/glide
	glide install

