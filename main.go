package main

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/interpreter"
	"github.com/docopt/docopt-go"
	"log"
	"os"
	"path/filepath"
)

const (
	VERSION = "0.4.3"
	USAGE   = `
Usage: %s [-o DESTFILE] [-d DICTPATH] [-f | -s] [-cm] [--] INPUTFILE
  %s -v
  %s -h

Arguments:
  INPUTFILE  The file describing entities and generation statements
  DESTFILE   The output file (defaults to "entities.json")
  DICTPATH   The path to your user-defined dictionaries

Options:
  -h --help
  -v --version
  -c --check                           Check syntax of INPUTFILE
  -m --no-metadata                     Omit metadata in generated entities (e.g. $type, $extends, etc.)
  -o DESTFILE --output=DESTFILE        Specify output file [default: entities.json]
  -d DICTPATH --dictionaries=DICTPATH  Specify DICTPATH
  -f --flatten                         Flattens entity hierarchies into a flat array; entities are
                                         outputted in reverse order of dependency, and linked by "$id";
                                         cannot be combined with --split-output
  -s --split-output                    Outputs entities into files, separated by declared type; cannot
                                         be combined with --flatten

`
)

func init() {
	log.SetFlags(0)
}

func main() {

	progname := filepath.Base(os.Args[0])
	usage := fmt.Sprintf(USAGE, progname, progname, progname)
	autoExit := true // set to `true` to let docopt automatically exit; `false` to handle exit ourselves

	args, _ := docopt.Parse(usage, nil, true, VERSION, false, autoExit)

	outputFile, _ := args["--output"].(string)
	filePerEntity, _ := args["--split-output"].(bool)
	flattenOutput, _ := args["--flatten"].(bool)
	disableMetadata, _ := args["--no-metadata"].(bool)
	syntaxCheck, _ := args["--check"].(bool)

	filename, _ := args["INPUTFILE"].(string)

	i := interpreter.New(flattenOutput, disableMetadata)

	if customDicts, ok := args["--dictionaries"].(string); !ok {
		a, _ := filepath.Abs(filename)
		i.SetCustomDictonaryPath(filepath.Dir(a))
	} else {
		i.SetCustomDictonaryPath(customDicts)
	}

	if syntaxCheck {
		if errors := i.CheckFile(filename); errors != nil {
			log.Fatalf("Syntax check failed: %v\n", errors)
		}

		log.Println("Syntax OK")
		os.Exit(0)
	}

	if _, errors := i.LoadFile(filename, interpreter.NewRootScope()); errors != nil {
		log.Fatalln(errors)
	}

	if errors := i.WriteGeneratedContent(outputFile, filePerEntity, flattenOutput); errors != nil {
		log.Fatalln(errors)
	}
}
