default: build

build: 
	pigeon -o dsl/dsl.go dsl/dsl.peg 
	go build
	go test ./...
