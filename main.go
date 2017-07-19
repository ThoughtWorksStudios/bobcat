// Package word provides utilities for word games.
package main

import (
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/interpreter"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func parseSpec(filename string) (interface{}, error) {
	f, _ := os.Open(filename)
	return dsl.ParseReader(filename, f, dsl.GlobalStore("filename", filename))
}

func fileDoesNotExist(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("You must pass in a file")
	}

	filename := os.Args[1]
	if fileDoesNotExist(filename) {
		log.Fatalf("File passed '%v' does not exist\n", filename)
	}

	if tree, err := parseSpec(filename); err != nil {
		log.Fatalf("Error parsing %s: %v", filename, err)
	} else {
		if errors := interpreter.New().Visit(tree.(dsl.Node)); errors != nil {
			log.Fatalln(errors)
		}
	}
}
