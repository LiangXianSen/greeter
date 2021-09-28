GO 		:= go
ID      := greeter
REPO    := github.com/LiangXianSen/$(ID)
VERSION := v0.0.1
SRC 	:= $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PKG 	:= $(shell go list ./...|grep -v /vendor/)
TARGET 	:= client server

# === This section provides application running utility === #
.PHONY: all build clean check lint test benchmark

all: build

build: $(TARGET)

$(TARGET): $(SRC)
	@$(GO) build $(REPO)/cmd/$@

test: check
	@$(GO) test -race $(PKG) -v -p 1 -coverprofile=.coverage.out
	@$(GO) tool cover -func=.coverage.out
	@rm -f .coverage.out

benchmark: check
	@$(GO) test -benchmem -bench=. -count=3 $(PKG)

check:
	@$(GO) vet -composites=false $(PKG)

lint:
	@golint -set_exit_status $(PKG)

lint-runner:
	@golangci-lint run ./...

doc:
	@godoc -http=localhost:6060 -play

clean:
	@rm $(TARGET)
