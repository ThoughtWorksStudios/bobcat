default: build

build: 
	go get
	pigeon -o dsl/dsl.go dsl/dsl.peg
	go build
	go test ./...
