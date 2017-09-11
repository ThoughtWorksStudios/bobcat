package main

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"github.com/ThoughtWorksStudios/bobcat/interpreter"
	"github.com/docopt/docopt-go"
	"log"
	"os"
	"path/filepath"
)

const (
	VERSION = "0.4.4"
	USAGE   = `
Usage: %s [-o DESTFILE] [-d DICTPATH] [-cfms] [--] INPUTFILE
  %s -v
  %s -h

Arguments:
  INPUTFILE  The file describing entities and generation statements
  DESTFILE   The output file (defaults to "entities.json"); accepts "-" to output to STDOUT
  DICTPATH   The path to your user-defined dictionaries

Options:
  -h --help
  -v --version
  -c --check                           Check syntax of INPUTFILE
  -m --no-metadata                     Omit metadata in generated entities (e.g. $type, $extends, etc.)
  -o DESTFILE --output=DESTFILE        Specify output file [default: entities.json] (use "-" for DESTFILE
                                         to specify STDOUT)
  -d DICTPATH --dictionaries=DICTPATH  Specify DICTPATH
  -f --flatten                         Flattens entity hierarchies into a flat array; entities are
                                         outputted in reverse order of dependency, and linked by "$id"
  -s --split-output                    Aggregates entities by type into separate files; DESTFILE
                                         serves as the filename template, meaning each file has the
                                         entity type appended to its basename (i.e. before the ".json"
                                         extension, as in "entities-myType.json"). Implies --flatten.

`
)

func init() {
	log.SetFlags(0)
}

func createEmitter(filename string, config map[string]interface{}) Emitter {
	var err error
	var emitter Emitter

	splitOutput, _ := config["--split-output"].(bool)
	flatten, _ := config["--flatten"].(bool)

	if "-" == filename && splitOutput {
		log.Fatalln("Cannot use --split-output with STDOUT")
	}

	switch true {
	case splitOutput:
		emitter, err = SplitEmitterForFile(filename)
	case flatten:
		emitter, err = FlatEmitterForFile(filename)
	default:
		emitter, err = NestedEmitterForFile(filename)
	}

	if err != nil {
		log.Fatalln(err)
	}

	return emitter
}

func main() {

	progname := filepath.Base(os.Args[0])
	usage := fmt.Sprintf(USAGE, progname, progname, progname)
	autoExit := true // set to `true` to let docopt automatically exit; `false` to handle exit ourselves

	args, _ := docopt.Parse(usage, nil, true, VERSION, false, autoExit)

	outputFile, _ := args["--output"].(string)
	disableMetadata, _ := args["--no-metadata"].(bool)
	syntaxCheck, _ := args["--check"].(bool)

	filename, _ := args["INPUTFILE"].(string)

	emitter := createEmitter(outputFile, args)

	i := interpreter.New(emitter, disableMetadata)

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

	if errors := emitter.Init(); errors != nil {
		log.Fatalln(errors)
	}

	if _, errors := i.LoadFile(filename, interpreter.NewRootScope()); errors != nil {
		log.Fatalln(errors)
	}

	if errors := emitter.Finalize(); errors != nil {
		log.Fatalln(errors)
	}
}
