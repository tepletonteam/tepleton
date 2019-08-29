all: test install

NOVENDOR = go list github.com/tepleton/basecoin/... | grep -v /vendor/

build:
	go build github.com/tepleton/basecoin/cmd/...

install:
	go install github.com/tepleton/basecoin/cmd/...

test:
	go test `${NOVENDOR}`
	#go run tests/tepleton/*.go

get_deps:
	go get -d github.com/tepleton/basecoin/...

update_deps:
	go get -d -u github.com/tepleton/basecoin/...

get_vendor_deps:
	go get github.com/Masterminds/glide
	glide install

build-docker:
	docker run -it --rm -v "$(PWD):/go/src/github.com/tepleton/basecoin" -w "/go/src/github.com/tepleton/basecoin" -e "CGO_ENABLED=0" golang:alpine go build ./cmd/basecoin
	docker build -t "tepleton/basecoin" .

clean:
	@rm -f ./basecoin

.PHONY: all build install test get_deps update_deps get_vendor_deps build-docker clean
