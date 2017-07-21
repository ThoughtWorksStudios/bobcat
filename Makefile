.PHONY: list

default: run

run: build test

# list available targets in Makefile
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e 'default|setup' -e '^$@$$' | xargs

# automate dev setup for local environments
local: setup depend build test

# setup for local environments
setup:
	brew install golang
	@echo 'Add this to your shell startup file:'
	@echo '    export GOPATH=`go env GOPATH`'
	@echo '    export PATH=$$GOPATH/bin:$$PATH'
	mkdir -p ~/go/src/github.com/ThoughtWorksStudios
	ln -s `pwd` ~/go/src/github.com/ThoughtWorksStudios/datagen

# automate dev setup using docker
docker:
	docker run -d -h development -it --rm -v `pwd`:/go/src/github.com/ThoughtWorksStudios/datagen kyleolivo/datagen

# get dependencies requires by the application
depend:
	go get github.com/mna/pigeon
	go get -d

# build and install the application
build:
	`go env GOPATH`/bin/pigeon -o dsl/dsl.go dsl/dsl.peg
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