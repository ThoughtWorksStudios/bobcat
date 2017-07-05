// Package word provides utilities for word games.
package main

import "fmt"
import "os"

func parseSpec(filename string) (interface{}, error) {
	f, _ := os.Open(filename)
	return ParseReader(filename, f)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "no arguments passed")
		os.Exit(1)
	}
	filename := os.Args[1]
	if fileExists(filename) {
		fmt.Fprintf(os.Stderr, "File passed '%v' does not exist\n", filename)
		os.Exit(1)
	}
	tree, err := parseSpec(filename)
	fmt.Println("ERR", err)
	if err != nil {
		fmt.Println("got an error", err)
	} else {
		fmt.Println("parse tree", tree)
	}
}
