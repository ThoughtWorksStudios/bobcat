package dsl

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
	re "regexp"
	"testing"
)

func testEntityField(name string, location *Location, value interface{}, args NodeSet) Node {
	return Node{
		Kind:  "field",
		Name:  name,
		Ref:   location,
		Value: value,
		Args:  args,
	}
}

func testGenEntity(name string, location *Location, kids NodeSet, args NodeSet) Node {
	return Node{
		Kind:     "generation",
		Name:     name,
		Ref:      location,
		Args:     args,
		Children: kids,
	}
}

func testEntity(name string, location *Location, kids NodeSet) Node {
	return Node{
		Kind:     "definition",
		Ref:      location,
		Name:     name,
		Children: kids,
	}
}

func testRootNode(kids NodeSet) Node {
	return Node{
		Kind:     "root",
		Ref:      NewLocation("", 1, 1, 0),
		Children: kids,
	}
}

func tableSpecForReservedWords() map[string]string {
	result := make(map[string]string)
	keyWords := []string{"date", "decimal", "dict", "false", "def", "generate", "integer", "string"}
	for _, kw := range keyWords {
		result[fmt.Sprintf(`def %s`, kw)] = fmt.Sprintf(`Illegal identifier: %q is a reserved word`, kw)
		result[fmt.Sprintf(`def t { %s string }`, kw)] = fmt.Sprintf(`Illegal identifier: %q is a reserved word`, kw)
		result[fmt.Sprintf(`generate %s`, kw)] = fmt.Sprintf(`Illegal identifier: %q is a reserved word`, kw)
		result[fmt.Sprintf(`generate %s(3)`, kw)] = fmt.Sprintf(`Illegal identifier: %q is a reserved word`, kw)
	}
	return result
}

func removeLocationInfo(err error) error {
	prefix := re.MustCompile(`^\d+:\d+ \(\d+\):\s+(rule "[\w ]+":\s+)?`)
	return fmt.Errorf(prefix.ReplaceAllString(err.Error(), ""))
}

func TestParsesBasicEntity(t *testing.T) {
	testRoot := testRootNode(NodeSet{testEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{})})
	actual, err := Parse("", []byte("def Bird {  }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestCanParseMultipleEntities(t *testing.T) {
	bird1 := testEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{})
	bird2 := testEntity("Bird2", NewLocation("", 2, 1, 14), NodeSet{})
	testRoot := testRootNode(NodeSet{bird1, bird2})
	actual, err := Parse("", []byte("def Bird {  }\ndef Bird2 { }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParsesChildEntity(t *testing.T) {
	entity := testEntity("Robin", NewLocation("", 1, 1, 0), NodeSet{})
	entity.Parent = "Bird"
	testRoot := testRootNode(NodeSet{entity})
	actual, err := Parse("", []byte("def Robin:Bird {  }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParsesBasicGenerationStatement(t *testing.T) {
	args := NodeSet{Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 15, 14)}}
	genBird := testGenEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{}, args)
	testRoot := testRootNode(NodeSet{genBird})
	actual, err := Parse("", []byte("generate Bird(1)"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestCanParseMultipleGenerationStatements(t *testing.T) {
	arg := Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 15, 14)}
	genBird := testGenEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{}, NodeSet{arg})
	arg.Ref = NewLocation("", 2, 16, 32)
	bird2Gen := testGenEntity("Bird2", NewLocation("", 2, 1, 17), NodeSet{}, NodeSet{arg})
	bird2Gen.Name = "Bird2"
	testRoot := testRootNode(NodeSet{genBird, bird2Gen})
	actual, err := Parse("", []byte("generate Bird(1)\ngenerate Bird2(1)"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestCanOverrideFieldInGenerateStatement(t *testing.T) {
	arg := Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 15, 14)}
	value := Node{Kind: "literal-string", Value: "birdie", Ref: NewLocation("", 1, 25, 24)}
	field := testEntityField("name", NewLocation("", 1, 20, 19), value, nil)
	genBird := testGenEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{field}, NodeSet{arg})
	testRoot := testRootNode(NodeSet{genBird})
	actual, err := Parse("", []byte("generate Bird(1) { name \"birdie\" }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestCanOverrideMultipleFieldsInGenerateStatement(t *testing.T) {
	value1 := Node{Kind: "literal-string", Value: "birdie", Ref: NewLocation("", 1, 25, 24)}
	field1 := testEntityField("name", NewLocation("", 1, 20, 19), value1, nil)
	value2 := Node{Kind: "builtin", Value: "integer", Ref: NewLocation("", 1, 39, 38)}
	arg1 := Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 47, 46)}
	arg2 := Node{Kind: "literal-int", Value: 2, Ref: NewLocation("", 1, 49, 48)}
	field2 := testEntityField("age", NewLocation("", 1, 35, 34), value2, NodeSet{arg1, arg2})

	arg := Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 15, 14)}
	genBird := testGenEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{field1, field2}, NodeSet{arg})
	testRoot := testRootNode(NodeSet{genBird})
	actual, err := Parse("", []byte("generate Bird(1) { name \"birdie\", age integer(1,2) }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParsedBothBasicEntityAndGenerationStatement(t *testing.T) {
	args := NodeSet{Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 2, 15, 26)}}
	genBird := testGenEntity("Bird", NewLocation("", 2, 1, 12), NodeSet{}, args)
	bird := testEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{})
	testRoot := testRootNode(NodeSet{bird, genBird})
	actual, err := Parse("", []byte("def Bird {}\ngenerate Bird(1)"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntityWithDynamicFieldWithoutArgs(t *testing.T) {
	value := Node{Kind: "builtin", Ref: NewLocation("", 1, 17, 16), Value: "string"}
	field := testEntityField("name", NewLocation("", 1, 12, 11), value, NodeSet{})
	bird := testEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{})
	bird.Children = NodeSet{field}
	testRoot := testRootNode(NodeSet{bird})
	actual, err := Parse("", []byte("def Bird { name string }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntityWithDynamicFieldWithArgs(t *testing.T) {
	value := Node{Kind: "builtin", Ref: NewLocation("", 1, 17, 16), Value: "string"}
	args := NodeSet{Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 24, 23)}}
	field := testEntityField("name", NewLocation("", 1, 12, 11), value, args)
	bird := testEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{})
	bird.Children = NodeSet{field}
	testRoot := testRootNode(NodeSet{bird})
	actual, err := Parse("", []byte("def Bird { name string(1) }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntitywithDynamicFieldWithMultipleArgs(t *testing.T) {
	value := Node{Kind: "builtin", Ref: NewLocation("", 1, 17, 16), Value: "integer"}
	arg1 := Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 25, 24)}
	arg2 := Node{Kind: "literal-int", Value: 5, Ref: NewLocation("", 1, 28, 27)}
	args := NodeSet{arg1, arg2}
	field := testEntityField("name", NewLocation("", 1, 12, 11), value, args)
	bird := testEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{})
	bird.Children = NodeSet{field}
	testRoot := testRootNode(NodeSet{bird})
	actual, err := Parse("", []byte("def Bird { name integer(1, 5) }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntityWithMultipleFields(t *testing.T) {
	value := Node{Kind: "builtin", Ref: NewLocation("", 1, 17, 16), Value: "string"}
	arg := Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 24, 23)}
	field1 := testEntityField("name", NewLocation("", 1, 12, 11), value, NodeSet{arg})

	value = Node{Kind: "builtin", Ref: NewLocation("", 1, 32, 31), Value: "integer"}
	arg1 := Node{Kind: "literal-int", Value: 1, Ref: NewLocation("", 1, 40, 39)}
	arg2 := Node{Kind: "literal-int", Value: 5, Ref: NewLocation("", 1, 43, 42)}
	args := NodeSet{arg1, arg2}
	field2 := testEntityField("age", NewLocation("", 1, 28, 27), value, args)

	bird := testEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{})
	bird.Children = NodeSet{field1, field2}
	testRoot := testRootNode(NodeSet{bird})
	actual, err := Parse("", []byte("def Bird { name string(1), age integer(1, 5) }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntityWithStaticField(t *testing.T) {
	value := Node{Kind: "literal-string", Ref: NewLocation("", 1, 17, 16), Value: "birdie"}
	field := testEntityField("name", NewLocation("", 1, 12, 11), value, nil)
	bird := testEntity("Bird", NewLocation("", 1, 1, 0), NodeSet{field})
	testRoot := testRootNode(NodeSet{bird})
	actual, err := Parse("", []byte("def Bird { name \"birdie\" }"))
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestRequiresDefOrGenerateStatements(t *testing.T) {
	_, err := Parse("", []byte("eek"))
	expectedErrorMsg := `no match found, expected: "def", "generate", [ \t\r\n] or EOF`
	ExpectsError(t, expectedErrorMsg, removeLocationInfo(err))
}

func TestReservedRulesRestrictions(t *testing.T) {
	for keyWord, expectedErrMessage := range tableSpecForReservedWords() {
		_, err := Parse("", []byte(keyWord))

		ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
	}
}

// this may become obsolete once we start supporting nested entities and other value declarations
func TestShouldGiveErrorForUnknownFieldTypes(t *testing.T) {
	specs := map[string]string{
		"generate(1) t { e eek }": `no match found, expected: "date", "decimal", "def", "dict", "false", "generate", "integer", "null", "string", "true", [ \t\r\n] or [a-z_]i`,
		"def t { e blah }":        `no match found, expected: "-", "0", "\"", "\\0", "date", "decimal", "dict", "false", "integer", "null", "string", "t"i, "true", [ \t\r\n], [0-9] or [1-9]`,
	}

	for spec, expectedErrMessage := range specs {
		_, err := Parse("", []byte(spec))
		ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
	}
}

func TestShouldGiveErrorWhenNoCountIsGivenToGenerate(t *testing.T) {
	expectedErrMessage := `no match found, expected: "(", [ \t\r\n] or [a-z0-9_]i`
	_, err := Parse("", []byte("generate Blah"))
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestEntityFieldRequiresType(t *testing.T) {
	expectedErrMessage := `no match found, expected: "-", "0", "\"", "\\0", "date", "decimal", "dict", "false", "integer", "null", "string", "t"i, "true", [ \t\r\n], [0-9] or [1-9]`
	_, err := Parse("", []byte("def Blah { name }"))
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestEntityDefinitionRequiresCurlyBrackets(t *testing.T) {
	expectedErrMessage := `no match found, expected: ":", "{", [ \t\r\n] or [a-z0-9_]i`
	_, err := Parse("", []byte("def Bird"))
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestFieldListWithoutCommas(t *testing.T) {
	expectedErrMessage := `Multiple field declarations must be delimited with a comma`
	_, err := Parse("", []byte("def Bird { h string b string }"))
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestEntityNameMustBeAlphaNumericAndStartWithALetter(t *testing.T) {
	specs := map[string]string{
		"generate(1) 4": `no match found, expected: "date", "decimal", "def", "dict", "false", "generate", "integer", "null", "string", "true", [ \t\r\n] or [a-z_]i`,
		"def 4 { }":     `no match found, expected: "date", "decimal", "def", "dict", "false", "generate", "integer", "null", "string", "true", [ \t\r\n] or [a-z_]i`,
		"def $eek { }":  `no match found, expected: "date", "decimal", "def", "dict", "false", "generate", "integer", "null", "string", "true", [ \t\r\n] or [a-z_]i`,
		"generate $eek": `no match found, expected: "date", "decimal", "def", "dict", "false", "generate", "integer", "null", "string", "true", [ \t\r\n] or [a-z_]i`,
		"generate eek$": `no match found, expected: "(", [ \t\r\n] or [a-z0-9_]i`,
		"def e$ek { }":  `no match found, expected: ":", "{", [ \t\r\n] or [a-z0-9_]i`,
	}

	for spec, expectedErrMessage := range specs {
		_, err := Parse("", []byte(spec))
		ExpectsError(t, expectedErrMessage, removeLocationInfo(err))

	}
}
