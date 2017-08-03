package main

import (
	"flag"
	"fmt"
	"github.com/ThoughtWorksStudios/bobcat/interpreter"
	"log"
	"os"
	"path/filepath"
)

func init() {
	log.SetFlags(0)
}

func debug(format string, tokens ...interface{}) {
	format = format + "\n"
	fmt.Fprintf(os.Stderr, format, tokens...)
}

func printHelpAndExit() {
	flag.CommandLine.Usage()
	os.Exit(1)
}

func main() {
	flag.CommandLine.Usage = func() {
		log.Print("Usage: ./bobcat [ options ] spec_file.lang")
		log.Print("\nOptions:")
		flag.CommandLine.PrintDefaults()
	}
	outputFile := flag.CommandLine.String("dest", "entities.json", "Destination file for generated content (NOTE that -dest and -split-output are mutually exclusize; the -dest flag will be ignored)")
	filePerEntity := flag.CommandLine.Bool("split-output", false, "Create a seperate output file per definition with the filename being the definition's name. (NOTE that -split-output and -dest are mutually exclusize; the -dest flag will be ignored)")
	syntaxCheck := flag.CommandLine.Bool("c", false, "Checks the syntax of the provided spec")
	customDicts := flag.CommandLine.String("d", "", "location of custom dictionary files ( e.g. ./bobcat -d=~/data/ examples/example.lang )")

	//everything except the executable itself
	flag.CommandLine.Parse(os.Args[1:])

	//flag.CommandLine.Args() returns anything passed that doesn't start with a "-"
	if len(flag.CommandLine.Args()) == 0 {
		log.Print("You must pass in a file")
		printHelpAndExit()
	}

	filename := flag.CommandLine.Args()[0]

	i := interpreter.New()

	if *customDicts == "" {
		a, _ := filepath.Abs(filename)
		i.SetCustomDictonaryPath(filepath.Dir(a))
	} else {
		i.SetCustomDictonaryPath(*customDicts)
	}

	if *syntaxCheck {
		if errors := i.CheckFile(filename); errors != nil {
			log.Fatalf("Syntax check failed: %v\n", errors)
		}

		log.Println("Syntax OK")
		os.Exit(0)
	}

	if errors := i.LoadFile(filename, interpreter.NewRootScope()); errors != nil {
		log.Fatalln(errors)
	}

	if errors := i.WriteGeneratedContent(*outputFile, *filePerEntity); errors != nil {
		log.Fatalln(errors)
	}
}
