# Metadata about this makefile and position
MKFILE_PATH := $(lastword $(MAKEFILE_LIST))
CURRENT_DIR := $(dir $(realpath $(MKFILE_PATH)))
CURRENT_DIR := $(CURRENT_DIR:/=)

# Get the project metadata
GOVERSION := 1.8
PROJECT := github.com/catsby/go-twitch
OWNER := $(dir $(PROJECT))
OWNER := $(notdir $(OWNER:/=))
NAME := $(notdir $(PROJECT))
EXTERNAL_TOOLS = github.com/ajg/form \
github.com/dnaeon/go-vcr/cassette \
github.com/dnaeon/go-vcr/recorder \
github.com/hashicorp/go-cleanhttp \
github.com/mitchellh/mapstructure \
gopkg.in/yaml.v2 \
gopkg.in/check.v1

# List of tests to run
TEST ?= ./...

# List all our actual files, excluding vendor
GOFILES = $(shell go list $(TEST) | grep -v /vendor/)

# Tags specific for building
GOTAGS ?=

# Number of procs to use
GOMAXPROCS ?= 4

default: test vet

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@echo 'go vet ./...'
	@go vet ./... ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

# bootstrap installs the necessary go tools for development or build
bootstrap:
	@echo "==> Bootstrapping ${PROJECT}..."
	@for t in ${EXTERNAL_TOOLS}; do \
		echo "--> Installing $$t" ; \
		go get -u "$$t"; \
	done

# test runs the test suite
test:
	@echo "==> Testing ${PROJECT}..."
	@go test -timeout=60s -parallel=20 -tags="${GOTAGS}" ${TEST} ${TESTARGS}

# test runs the test suite, verbosely
testv:
	@echo "==> Testing ${PROJECT}..."
	@go test -timeout=60s -parallel=20 -tags="${GOTAGS}" ${TEST} -v ${TESTARGS}

pre-commit:
	@echo "==> Installing pre-commit hook..."
	cp scripts/pre-commit .git/hooks/

.PHONY: bootstrap test 
