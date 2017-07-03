// Package word provides utilities for word games.
package main

import "fmt"

func main() {
	tree, err := ParseFile("person.lang")
	if err != nil {
		fmt.Println("got an error", err)
	} else {
		fmt.Println("parse tree", tree)
	}
}
