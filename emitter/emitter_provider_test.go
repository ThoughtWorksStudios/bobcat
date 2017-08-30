package emitter

import (
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestValidateFileNameTemplate(t *testing.T) {
	p := &PerTypeEmitterProvider{}
	ExpectsError(t, "You must provide a filename template", p.validateFilenameTemplate(""))
	ExpectsError(t, "Filename template must have a `.json` extension", p.validateFilenameTemplate("foojson"))
	ExpectsError(t, "Filename template must have a basename before the `.json` extension", p.validateFilenameTemplate(".json"))
	ExpectsError(t, "Filename template must have a basename before the `.json` extension", p.validateFilenameTemplate("foo/.json"))

	// Tests for directory existence
	ExpectsError(t, "Directory \"does/not/exist\" is not a directory", p.validateFilenameTemplate("./does/not/exist/file.json"))

	// Tests that basedir is a directory node
	ExpectsError(t, "Directory \"../dsl/dsl.peg\" is not a directory", p.validateFilenameTemplate("../dsl/dsl.peg/file.json"))

	AssertNil(t, p.validateFilenameTemplate("ok.json"), "Should allow just a filename with `.json` extension")
	AssertNil(t, p.validateFilenameTemplate("ok.JsOn"), "Should ignore case when checking for presence of `.json` extension")
	AssertNil(t, p.validateFilenameTemplate("../dsl/ok.json"), "Should allow path before filename")
	AssertNil(t, p.validateFilenameTemplate("../ok.json"), "Should allow ../ paths")
}

func TestPathFromType(t *testing.T) {
	p := &PerTypeEmitterProvider{basedir: "foo", basename: "bar"}
	AssertEqual(t, "foo/bar-myType.json", p.PathFromType("myType"), "Should derive file path from type, basedir, and basename")
}

func TestNewPerTypeEmitterProvider(t *testing.T) {
	ep, e := NewPerTypeEmitterProvider("stuff.json")
	p := ep.(*PerTypeEmitterProvider)
	AssertNil(t, e, "Should not receive error on creation")
	AssertEqual(t, p.basedir, ".", "Should set basedir on construction, and default to $(pwd)")
	AssertEqual(t, p.basename, "stuff", "Should set basename on construction")
}

func TestGet(t *testing.T) {
	basename := strconv.Itoa(rand.Int())

	p, _ := NewPerTypeEmitterProvider(basename + ".json")
	emitter, e := p.Get("fooType")

	defer func() {
		emitter.Finalize()
		os.Remove(basename + "-fooType.json")
	}()

	AssertNil(t, e, "Should not have received error when instantiating emitter")

	_, ok := emitter.(*FlatEmitter)
	Assert(t, ok, "Emitter instance should be a FlatEmitter")
}
