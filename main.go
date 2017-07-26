package main

import (
	"flag"
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/interpreter"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func debug(format string, tokens ...interface{}) {
	format = format + "\n"
	fmt.Fprintf(os.Stderr, format, tokens...)
}

func parseSpec(filename string) (interface{}, error) {
	f, _ := os.Open(filename)
	return dsl.ParseReader(filename, f, dsl.GlobalStore("filename", filename), dsl.Recover(false))
}

func fileDoesNotExist(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

func defHelpMessage() {
	flag.CommandLine.Usage = func() {
		log.Print("Usage: ./datagen [ options ] spec_file.lang")
		log.Print("\nOptions:")
		flag.CommandLine.PrintDefaults()
	}
}

func main() {
	defHelpMessage()
	outputFile := flag.CommandLine.String("dest", "entities.json", "destination file for generated content (NOTE that -dest and -split-output are mutually exclusize; the -dest flag will be ignored)")
	filePerEntity := flag.CommandLine.Bool("split-output", false, "Create a seperate output file per definition with the filename being the definition's name. (NOTE that -split-output and -dest are mutually exclusize; the -dest flag will be ignored)")

	//everything except the executable itself
	flag.CommandLine.Parse(os.Args[1:])

	//flag.CommandLine.Args() returns anything passed that doesn't start with a "-"
	if len(flag.CommandLine.Args()) == 0 {
		log.Print("You must pass in a file")
		flag.CommandLine.Usage()
	}

	filename := flag.CommandLine.Args()[0]
	if fileDoesNotExist(filename) {
		log.Printf("File passed '%v' does not exist\n", filename)
		flag.CommandLine.Usage()
	}

	if tree, err := parseSpec(filename); err != nil {
		log.Fatalf("Error parsing %s: %v", filename, err)
	} else {
		inter := interpreter.New()
		if errors := inter.Visit(tree.(dsl.Node)); errors != nil {
			log.Fatalln(errors)
		}
		if errors := inter.WriteGeneratedContent(*outputFile, *filePerEntity); errors != nil {
			log.Fatalln(errors)
		}
	}
}
