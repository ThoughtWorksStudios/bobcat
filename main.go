// Package word provides utilities for word games.
package main

import "fmt"
import "os"

func parseSpec(filename string) (interface{}, error) {
	f, _ := os.Open(filename)
	return ParseReader(filename, f)
}

func main() {
	tree, err := parseSpec("person.lang")
	fmt.Println("ERR", err)
	if err != nil {
		fmt.Println("got an error", err)
	} else {
		fmt.Println("parse tree", tree)
	}
}
