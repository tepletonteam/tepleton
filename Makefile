.PHONY: all test get_deps

all: test install

install: get_deps
	go install github.com/tepleton/blackstar/cmd/...

test:
	go test github.com/tepleton/blackstar/...

get_deps:
	go get -d github.com/tepleton/blackstar/...
