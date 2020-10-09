# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGENERATE=$(GOCMD) generate
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
GOLIST=$(GOCMD) list
GOMOD=$(GOCMD) mod
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOTOOL=$(GOCMD) tool

BINARY_NAME=service
BASE_PACKAGE_NAME=github.com/muffix/relayr-challenge
CMD_PACKAGE_NAME=cmd/service

TEST_REPORT_OUTPUT=test-report.out
COVERAGE_OUTPUT=coverage.out
CPU_PROFILE_OUTPUT=profile.out
MEMORY_PROFILE_OUTPUT=memprofile.out

SRCS := $(shell find . -name '*.go')
LINTERS := \
	golang.org/x/lint/golint \
	golang.org/x/tools/cmd/goimports \
	github.com/kisielk/errcheck \
	honnef.co/go/tools/cmd/staticcheck

.EXPORT_ALL_VARIABLES:

GO111MODULE=on
CGO_ENABLED=0

.PHONY: compile build run deps updatedeps testdeps golint vet goimports goimports-check tidy tidy-check test test-coverprofile bench coverage clean

compile:
	$(GOBUILD) ./...

build:
	@mkdir -p build
	GOOS=linux $(GOBUILD) \
		-o build/$(BINARY_NAME) \
		-v \
		-ldflags="-w \
				  -s \
                  -X github.com/muffix/relayr-challenge/internal/httpapi.revision=${GITHUB_SHA} \
                  -X github.com/muffix/relayr-challenge/internal/httpapi.pipelineID=${GITHUB_RUN_ID} \
                  -X github.com/muffix/relayr-challenge/internal/httpapi.buildDate=$(shell date -u +%Y-%m-%dT%TZ)" \
        $(BASE_PACKAGE_NAME)/$(CMD_PACKAGE_NAME)

run:
	@mkdir -p build
	$(GOBUILD) \
		-o build/$(BINARY_NAME) \
		-v \
		-ldflags="-w \
				  -s \
				  -X github.com/muffix/relayr-challenge/internal/httpapi.revision=dev \
				  -X github.com/muffix/relayr-challenge/internal/httpapi.pipelineID=dev \
				  -X github.com/muffix/relayr-challenge/internal/httpapi.buildDate=$(shell date -u +%Y-%m-%dT%TZ)" \
		$(BASE_PACKAGE_NAME)/$(CMD_PACKAGE_NAME)
	build/$(BINARY_NAME)

deps:
	$(GOGET) -d -v ./...

updatedeps:
	$(GOGET) -d -v -u -f ./...

testdeps:
	$(GOGET) -d -v -t ./...
	$(GOGET) -v $(LINTERS)

golint:
	@for file in $(SRCS); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

govet:
	$(GOVET) ./...

goimports:
	goimports -format-only -w ${SRCS}

goimports-check:
	@if [ ! -z "$$(goimports -format-only -l ${SRCS})" ]; then \
      		echo "Found unformatted source files. Please run"; \
      		echo "  make goimports"; \
      		echo "To automatically format your files"; \
      		exit 1; \
      	fi

tidy:
	$(GOMOD) tidy

tidy-check: tidy
	@if [ -n "$$(git diff-index --exit-code --ignore-submodules --name-only HEAD | grep -E '^go.(mod|sum)$$')" ]; then \
		echo "go.mod or go.sum has changed after running go mod tidy for you."; \
		echo "Please make sure you review and commit the changes."; \
		exit 1; \
	fi

checks: golint goimports-check govet tidy-check

test:
	$(GOTEST) -cover ./...

test-coverprofile:
	$(GOTEST) -covermode=count -coverprofile=$(COVERAGE_OUTPUT) ./... -json > $(TEST_REPORT_OUTPUT)

bench:
	$(GOTEST) $(BASE_PACKAGE_NAME)/$(CMD_PACKAGE_NAME) -benchmem -cpuprofile $(CPU_PROFILE_OUTPUT) -memprofile $(MEMORY_PROFILE_OUTPUT)

coverage:
	$(GOTEST) -covermode=count -coverprofile=$(COVERAGE_OUTPUT) ./...
	$(GOTOOL) cover -html=$(COVERAGE_OUTPUT)

clean:
	rm -rf build/ $(TEST_REPORT_OUTPUT) $(COVERAGE_OUTPUT) $(CPU_PROFILE_OUTPUT) $(MEMORY_PROFILE_OUTPUT)
