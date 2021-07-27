GO=go

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

TARGET = client server

project = github.com/LiangXianSen/greeter
packages := $(shell go list ./...|grep -v /vendor/)

.PHONY: check test lint

all: check build

build: $(TARGET)

$(TARGET): $(SRC)
	@$(GO) build $(project)/cmd/$@

test: check
	@$(GO) test -race $(packages) -v -coverprofile=.coverage.out
	@$(GO) tool cover -func=.coverage.out
	@rm -f .coverage.out

check:
	@$(GO) vet -composites=false $(packages)

lint:
	@golint -set_exit_status $(packages)

lint-runner:
	@golangci-lint run ./...

doc:
	@godoc -http=localhost:6060 -play

clean:
	@rm $(TARGET)
