// Package word provides utilities for word games.
package main

import (
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/interpreter"
	"os"
)

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
		fmt.Fprintln(os.Stderr, "no arguments passed")
		os.Exit(1)
	}
	filename := os.Args[1]
	if fileDoesNotExist(filename) {
		fmt.Fprintf(os.Stderr, "File passed '%v' does not exist\n", filename)
		os.Exit(1)
	}
	tree, err := parseSpec(filename)
	if err != nil {
		fmt.Println("got an error", err)
	} else {
		errors := interpreter.New(nil).Visit(tree.(dsl.Node))

		if errors != nil {
			fmt.Println(errors)
			os.Exit(1)
		}
	}
}
