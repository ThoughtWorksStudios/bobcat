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
# --------

default: run

run: build test

# list available targets in Makefile
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e 'default|setup' -e '^$@$$' | xargs

# one-time setup for local environments
setup:
	brew install golang
	@echo 'Add this to your shell startup file:'
	@echo '    export GOPATH=$(GOPATH)'
	@echo '    export PATH=$(GOBIN):$$PATH'
	mkdir -p $(GOPATH)/src/github.com/ThoughtWorksStudios
	test -e /Users/kyleolivo/go/src/github.com/ThoughtWorksStudios/datagen || ln -s `pwd` $(GOPATH)/src/github.com/ThoughtWorksStudios/datagen

# one-time automation of dev setup for local environments
local: setup depend build test

# start development environment using docker
docker:
	docker run -h development -it --rm -v `pwd`:/go/src/github.com/ThoughtWorksStudios/datagen kyleolivo/datagen

# automate run for ci
ci: depend run

# get dependencies requires by the application
depend:
	go get github.com/mna/pigeon
	go get -d

# build and install the application
build:
	$(GOBIN)/pigeon -o dsl/dsl.go dsl/dsl.peg
	go build

# test the application
test:
	go test ./...
	./datagen example.lang

# remove junk files
clean:
	rm -f dsl/dsl.go
	rm -f datagen
	find . -type f -name \*.json -delete
