package dsl

import (
	. "github.com/ThoughtWorksStudios/datagen/test_helpers"
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
	expectedErrorMsg := `1:1 (0): no match found, expected: "def", "generate", [ \t\r\n] or EOF`
	ExpectsError(t, expectedErrorMsg, err)
}

func TestReservedRulesRestrictions(t *testing.T) {
	keyWords := map[string]string{
		"def generate":              `1:5 (4): no match found, expected: !"generate" or [ \t\r\n]`,
		"generate def(2)":           `1:10 (9): no match found, expected: !"def" or [ \t\r\n]`,
		"generate def":              `1:10 (9): no match found, expected: !"def" or [ \t\r\n]`,
		"def integer":               `1:5 (4): no match found, expected: !"integer" or [ \t\r\n]`,
		"generate string(2)":        `1:10 (9): no match found, expected: !"string" or [ \t\r\n]`,
		"def decimal":               `1:5 (4): no match found, expected: !"decimal" or [ \t\r\n]`,
		"def date":                  `1:5 (4): no match found, expected: !"date" or [ \t\r\n]`,
		"def dict":                  `1:5 (4): no match found, expected: !"dict" or [ \t\r\n]`,
		"generate null(2)":          `1:10 (9): no match found, expected: !"null" or [ \t\r\n]`,
		"def true":                  `1:5 (4): no match found, expected: !"true" or [ \t\r\n]`,
		"def false":                 `1:5 (4): no match found, expected: !"false" or [ \t\r\n]`,
		"def t { def string }":      `1:9 (8): no match found, expected: !"def", "}" or [ \t\r\n]`,
		"def t { generate string }": `1:9 (8): no match found, expected: !"generate", "}" or [ \t\r\n]`,
		"def t { e def }":           `1:11 (10): no match found, expected: "-", "0", "\"", "date", "decimal", "dict", "false", "integer", "null", "string", "true", [ \t\r\n], [0-9] or [1-9]`,
		"def t { e generate }":      `1:11 (10): no match found, expected: "-", "0", "\"", "date", "decimal", "dict", "false", "integer", "null", "string", "true", [ \t\r\n], [0-9] or [1-9]`,
	}

	for keyWord, expectedErrMessage := range keyWords {
		_, err := Parse("", []byte(keyWord))
		ExpectsError(t, expectedErrMessage, err)

	}
}

func TestShouldGiveErrorWhenNoCountIsGivenToGenerate(t *testing.T) {
	expectedErrMessage := `1:14 (13): no match found, expected: "(", [ \t\r\n] or [a-z0-9_]i`
	_, err := Parse("", []byte("generate Blah"))
	ExpectsError(t, expectedErrMessage, err)
}

func TestShouldGiveErrorForUnknowFieldTypes(t *testing.T) {
	specs := map[string]string{
		"generate(1) t { e eek }": `1:9 (8): no match found, expected: [ \t\r\n] or [a-z_]i`,
		"def t { e blah }":        `1:11 (10): no match found, expected: "-", "0", "\"", "date", "decimal", "dict", "false", "integer", "null", "string", "true", [ \t\r\n], [0-9] or [1-9]`,
	}

	for spec, expectedErrMessage := range specs {
		_, err := Parse("", []byte(spec))
		ExpectsError(t, expectedErrMessage, err)

	}
}

func TestEntityNameCannotStartWithInteger(t *testing.T) {
	expectedErrMessage := `1:5 (4): no match found, expected: [ \t\r\n] or [a-z_]i`
	_, err := Parse("", []byte("def 0 { }"))
	ExpectsError(t, expectedErrMessage, err)
}
