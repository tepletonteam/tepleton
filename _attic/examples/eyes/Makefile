LINKER_FLAGS:="-X github.com/tepleton/tepleton-sdk/client/commands.CommitHash=`git rev-parse --short HEAD`"

install:
	@go install -ldflags $(LINKER_FLAGS) ./cmd/...

test: test_unit test_cli

test_unit:
	@go test `glide novendor`

test_cli:
	./tests/cli/eyes.sh

.PHONY: install test test_unit test_cli
