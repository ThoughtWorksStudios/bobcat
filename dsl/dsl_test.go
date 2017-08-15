package dsl

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	re "regexp"
	"testing"
)

func runParser(script string) (interface{}, error) {
	return Parse("testScript", []byte(script), Recover(false))
}

func testField(name string, value interface{}, args NodeSet, countRange NodeSet) Node {
	return Node{
		Kind:       "field",
		Name:       name,
		Value:      value,
		Args:       args,
		CountRange: countRange,
	}
}

func testIdNode(name string) Node {
	return Node{Kind: "identifier", Value: name}
}

func testGenEntity(entity Node, args NodeSet) Node {
	return Node{
		Kind:  "generation",
		Value: entity,
		Args:  args,
	}
}

func testEntityOverride(name, parent string, body NodeSet) Node {
	id := testIdNode(parent)
	return Node{
		Kind:     "entity",
		Name:     name,
		Related:  &id,
		Children: body,
	}
}

func testEntity(name string, body NodeSet) Node {
	return testAssignNode(
		testIdNode(name),
		Node{
			Kind:     "entity",
			Name:     name,
			Children: body,
		},
	)
}

func testRootNode(kids NodeSet) Node {
	return Node{
		Kind:     "root",
		Children: kids,
	}
}

func testAssignNode(left, right Node) Node {
	return Node{
		Kind: "assignment",
		Children: NodeSet{left, right},
	}
}

func tableSpecForReservedWords() map[string]string {
	result := make(map[string]string)
	keyWords := []string{"date", "decimal", "dict", "false", "generate", "integer", "string"}
	for _, kw := range keyWords {
		result[fmt.Sprintf(`%s =`, kw)] = fmt.Sprintf(`Illegal identifier: %q is a reserved word`, kw)
		result[fmt.Sprintf(`t = { %s string }`, kw)] = fmt.Sprintf(`Illegal identifier: %q is a reserved word`, kw)
		result[fmt.Sprintf(`generate (1, %s)`, kw)] = fmt.Sprintf(`Illegal identifier: %q is a reserved word`, kw)
		result[fmt.Sprintf(`generate (3, %s)`, kw)] = fmt.Sprintf(`Illegal identifier: %q is a reserved word`, kw)
	}
	return result
}

func removeLocationInfo(err error) error {
	if nil == err {
		return nil
	}

	prefix := re.MustCompile(`^testScript:\d+:\d+ \(\d+\):\s+(?:rule (?:"[\w -]+"|\w+):\s+)?`)
	return fmt.Errorf(prefix.ReplaceAllString(err.Error(), ""))
}

func TestParsesBasicEntity(t *testing.T) {
	testRoot := testRootNode(NodeSet{testEntity("Bird", NodeSet{})})
	actual, err := runParser("Bird = {  }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestCanParseMultipleEntities(t *testing.T) {
	bird1 := testEntity("Bird", NodeSet{})
	bird2 := testEntity("Bird2", NodeSet{})
	testRoot := testRootNode(NodeSet{bird1, bird2})
	actual, err := runParser("Bird = {  }\nBird2 = { }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParsesChildEntity(t *testing.T) {
	entity := testEntity("Robin", NodeSet{})
	entity.Children[1].Related = &Node{Kind: "identifier", Value: "Bird"}
	testRoot := testRootNode(NodeSet{entity})
	actual, err := runParser("Robin = Bird {  }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParsesBasicGenerationStatement(t *testing.T) {
	args := NodeSet{Node{Kind: "literal-int", Value: 1}}
	genBird := testGenEntity(testIdNode("Bird"), args)
	testRoot := testRootNode(NodeSet{genBird})
	actual, err := runParser("generate(1, Bird)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestCanParseMultipleGenerationStatements(t *testing.T) {
	arg := Node{Kind: "literal-int", Value: 1}
	genBird := testGenEntity(testIdNode("Bird"), NodeSet{arg})
	bird2Gen := testGenEntity(testIdNode("Bird2"), NodeSet{arg})
	testRoot := testRootNode(NodeSet{genBird, bird2Gen})

	actual, err := runParser("generate(1, Bird)\ngenerate(1, Bird2)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestCanOverrideFieldInGenerateStatement(t *testing.T) {
	arg := Node{Kind: "literal-int", Value: 1}
	value := Node{Kind: "literal-string", Value: "birdie"}
	field := testField("name", value, nil, nil)
	genBird := testGenEntity(testEntityOverride("", "Bird", NodeSet{field}), NodeSet{arg})
	testRoot := testRootNode(NodeSet{genBird})
	actual, err := runParser("generate(1, Bird { name \"birdie\" })")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestCanOverrideMultipleFieldsInGenerateStatement(t *testing.T) {
	value1 := Node{Kind: "literal-string", Value: "birdie"}
	field1 := testField("name", value1, nil, nil)
	value2 := Node{Kind: "builtin", Value: "integer"}
	arg1 := Node{Kind: "literal-int", Value: 1}
	arg2 := Node{Kind: "literal-int", Value: 2}
	field2 := testField("age", value2, NodeSet{arg1, arg2}, nil)

	arg := Node{Kind: "literal-int", Value: 1}
	genBird := testGenEntity(testEntityOverride("", "Bird", NodeSet{field1, field2}), NodeSet{arg})
	testRoot := testRootNode(NodeSet{genBird})
	actual, err := runParser("generate(1, Bird { name \"birdie\", age integer(1,2) })")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParsedBothBasicEntityAndGenerationStatement(t *testing.T) {
	args := NodeSet{Node{Kind: "literal-int", Value: 1}}
	genBird := testGenEntity(testIdNode("Bird"), args)
	bird := testEntity("Bird", NodeSet{})
	testRoot := testRootNode(NodeSet{bird, genBird})
	actual, err := runParser("Bird = {}\ngenerate (1, Bird)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntityWithDynamicFieldWithBound(t *testing.T) {
	value := Node{Kind: "builtin", Value: "string"}
	count := NodeSet{Node{Kind: "literal-int", Value: 1}, Node{Kind: "literal-int", Value: 8}}
	field := testField("name", value, NodeSet{}, count)
	bird := testEntity("Bird", NodeSet{field})
	testRoot := testRootNode(NodeSet{bird})
	actual, err := runParser("Bird = { name string[1,8] }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntityWithDynamicFieldWithoutArgs(t *testing.T) {
	value := Node{Kind: "builtin", Value: "string"}
	field := testField("name", value, NodeSet{}, nil)
	bird := testEntity("Bird", NodeSet{field})
	testRoot := testRootNode(NodeSet{bird})
	actual, err := runParser("Bird = { name string }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntityWithDynamicFieldWithArgs(t *testing.T) {
	value := Node{Kind: "builtin", Value: "string"}
	args := NodeSet{Node{Kind: "literal-int", Value: 1}}
	field := testField("name", value, args, nil)
	bird := testEntity("Bird", NodeSet{field})
	testRoot := testRootNode(NodeSet{bird})
	actual, err := runParser("Bird = { name string(1) }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntitywithDynamicFieldWithMultipleArgs(t *testing.T) {
	value := Node{Kind: "builtin", Value: "integer"}
	arg1 := Node{Kind: "literal-int", Value: 1}
	arg2 := Node{Kind: "literal-int", Value: 5}
	args := NodeSet{arg1, arg2}
	field := testField("name", value, args, nil)
	bird := testEntity("Bird", NodeSet{field})
	testRoot := testRootNode(NodeSet{bird})
	actual, err := runParser("Bird = { name integer(1, 5) }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntityWithMultipleFields(t *testing.T) {
	value := Node{Kind: "builtin", Value: "string"}
	arg := Node{Kind: "literal-int", Value: 1}
	field1 := testField("name", value, NodeSet{arg}, nil)

	value = Node{Kind: "builtin", Value: "integer"}
	arg1 := Node{Kind: "literal-int", Value: 1}
	arg2 := Node{Kind: "literal-int", Value: 5}
	args := NodeSet{arg1, arg2}
	field2 := testField("age", value, args, nil)

	bird := testEntity("Bird", NodeSet{field1, field2})
	testRoot := testRootNode(NodeSet{bird})
	actual, err := runParser("Bird = { name string(1), age integer(1, 5) }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntityWithStaticField(t *testing.T) {
	value := Node{Kind: "literal-string", Value: "birdie"}
	field := testField("name", value, nil, nil)
	bird := testEntity("Bird", NodeSet{field})
	testRoot := testRootNode(NodeSet{bird})
	actual, err := runParser("Bird = { name \"birdie\" }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntityWithEntityDeclarationField(t *testing.T) {
	args := NodeSet{}
	goatValue := Node{Kind: "literal-string", Value: "billy"}
	goatField := testField("name", goatValue, nil, nil)
	goat := testEntity("Goat", NodeSet{goatField})
	field := testField("pet", goat, args, nil)
	person := testEntity("Person", NodeSet{field})
	testRoot := testRootNode(NodeSet{person})
	actual, err := runParser("Person = { pet Goat = { name \"billy\" } }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestParseEntityWithEntityReferenceField(t *testing.T) {
	args := NodeSet{}
	goatValue := Node{Kind: "literal-string", Value: "billy"}
	goatField := testField("name", goatValue, nil, nil)
	goat := testEntity("Goat", NodeSet{goatField})
	value := Node{Kind: "identifier", Value: "Goat"}
	field := testField("pet", value, args, nil)
	person := testEntity("Person", NodeSet{field})
	testRoot := testRootNode(NodeSet{goat, person})
	actual, err := runParser("Goat = { name \"billy\" } Person = { pet Goat }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestVariableAssignment(t *testing.T) {
	value := Node{Kind: "literal-string", Value: "hello"}
	name := testField("name", value, nil, nil)
	foo := testEntity("Foo", NodeSet{name})
	testRoot := testRootNode(NodeSet{
		foo,
		testAssignNode(
			testIdNode("foos"),
		    testGenEntity(testIdNode("Foo"), NodeSet{Node{Kind: "literal-int", Value: 3}}),
    	),
	})
	actual, err := runParser("Foo = { name \"hello\" } foos = generate(3, Foo)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(Node).String())
}

func TestRequiresValidStatements(t *testing.T) {
	_, err := runParser("eek")
	expectedErrorMsg := `Don't know how to evaluate "eek"`
	ExpectsError(t, expectedErrorMsg, removeLocationInfo(err))
}

func TestReservedRulesRestrictions(t *testing.T) {
	for keyWord, expectedErrMessage := range tableSpecForReservedWords() {
		_, err := runParser(keyWord)

		ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
	}
}

func TestShouldGiveErrorWhenNoCountIsGivenToGenerate(t *testing.T) {
	expectedErrMessage := "`generate` statement \"generate Blah\" requires arguments `(count, entity)`"
	_, err := runParser("generate Blah")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestEntityFieldRequiresType(t *testing.T) {
	expectedErrMessage := `Missing field type for field declaration "name"`
	_, err := runParser("Blah = { name }")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestEntityDefinitionRequiresCurlyBrackets(t *testing.T) {
	expectedErrMessage := `Missing right-hand of assignment expression "Bird ="`
	_, err := runParser("Bird =")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestFieldListWithoutCommas(t *testing.T) {
	expectedErrMessage := `Multiple field declarations must be delimited with a comma`
	_, err := runParser("Bird = { h string b string }")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestIllegalIdentifiers(t *testing.T) {
	specs := map[string]string{
		"4 = { }":             `Illegal identifier "4"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"$eek = { }":          `Illegal identifier "$eek"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"generate (1, $eek)": `Illegal identifier "$eek"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"e$ek = { }":          `Illegal identifier "e$ek"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"generate = { }":      `Illegal identifier: "generate" is a reserved word`,
	}

	for spec, expectedErrMessage := range specs {
		_, err := runParser(spec)
		ExpectsError(t, expectedErrMessage, removeLocationInfo(err))

	}
}
