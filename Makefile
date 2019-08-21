.PHONY: all test get_deps

all: test install

install: get_deps
	go install github.com/tepleton/basecoin/cmd/...

test:
	go test github.com/tepleton/basecoin/...

get_deps:
	go get -d github.com/tepleton/basecoin/...

update_deps:
	go get -d -u github.com/tepleton/basecoin/...
