.PHONY: list

# --------
# for those that like to go against the grain :-)
ifndef GOPATH
GOPATH:=$(shell go env GOPATH)
else
GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
endif

ifndef GOBIN
GOBIN:=$(GOPATH)/bin
endif

EXAMPLE_FILE:=examples/example.lang

# --------

default: run

run: clean build test smoke

# list available targets in Makefile
list:
	@echo ""
	@echo "Make targets:"
	@echo '  ' `$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^(default|setup|ci|exectest|execsmoke)$$' -e '^$@$$' | xargs`
	@echo ""

# one-time setup for local environments
setup:
	@echo "===== Setup ====="
	@which go > /dev/null 2>&1 && echo "<`go version`> is already installed" || brew install golang
	@echo 'Ensure the following environment variables are set if you have not already done so:'
	@echo '    GOPATH=$(GOPATH)'
	@echo '    PATH=$(GOBIN):$$PATH'
	@mkdir -p $(GOPATH)/src/github.com/ThoughtWorksStudios
	@test -e $(GOPATH)/src/github.com/ThoughtWorksStudios/bobcat || ln -s `pwd` $(GOPATH)/src/github.com/ThoughtWorksStudios/bobcat
	@echo ""

# one-time automation of dev setup for local environments
local: setup clean depend build test smoke

# automate run for ci
ci: depend run

# add and run werkcer cli locally
wercker:
	brew tap wercker/wercker
	brew install wercker-cli
	wercker build

# get dependencies requires by the application
depend:
	go get github.com/mna/pigeon
	go get -d

prepare:
	@$(GOBIN)/pigeon -o dsl/dsl.go dsl/dsl.peg

compile:
	@echo "===== Compiling ====="
	@for platform in `test -z "$$WERCKER_ROOT" && echo "darwin" || echo "darwin linux windows"`; do \
	  echo "Building binary for $$platform"; GOOS=$$platform GOARCH=amd64 go build -o bobcat-$$platform; \
	done
	@echo ""

exectest:
	@echo "===== Unit tests ====="
	go test ./common ./dictionary ./dsl ./emitter ./generator ./interpreter .
	@echo ""

execsmoke:
	@echo "===== Smoke test ====="
	@echo "Running binary on $(EXAMPLE_FILE)"
	@test "Darwin" = `uname -s` && ./bobcat-darwin $(EXAMPLE_FILE) || ./bobcat-linux $(EXAMPLE_FILE)
	@echo ""

build: prepare compile

# just run tests
test: prepare exectest

# smoke tests
smoke: build execsmoke

# Runs benchmarks
performance:
	go test -bench=. ./dictionary ./interpreter ./generator

# remove junk files
clean:
	@echo "===== Cleaning up ====="
	go clean
	rm -f dsl/dsl.go
	rm -f bobcat bobcat-*
	find . -type f -name \*.json -delete
	@echo ""

# create a release tarball
release: depend build
	@echo ""
	@echo "===== Packaging release ====="
	tar czf bobcat.tar.gz bobcat-* examples/example.lang
	@echo ""
