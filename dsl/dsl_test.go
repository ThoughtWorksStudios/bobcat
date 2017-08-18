package dsl

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	re "regexp"
	"testing"
)

func runParser(script string) (interface{}, error) {
	return Parse("testScript", []byte(script), Recover(false))
}

func testEntity(name, parent string, body NodeSet) *Node {
	var parentNode *Node

	if "" != parent {
		parentNode = IdNode(nil, parent)
	}

	if "" != name {
		return AssignNode(nil, IdNode(nil, name), EntityNode(nil, IdNode(nil, name), parentNode, body))
	} else {
		return EntityNode(nil, nil, parentNode, body)
	}
}

func tableSpecForReservedWords() map[string]string {
	result := make(map[string]string)
	keyWords := []string{"date", "decimal", "dict", "false", "generate", "integer", "string"}
	for _, kw := range keyWords {
		result[fmt.Sprintf(`%s =`, kw)] = fmt.Sprintf(`Illegal identifier: %q is a reserved word`, kw)
		result[fmt.Sprintf(`t = { %s: string }`, kw)] = fmt.Sprintf(`Illegal identifier: %q is a reserved word`, kw)
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
	testRoot := RootNode(nil, NodeSet{testEntity("Bird", "", NodeSet{})})
	actual, err := runParser("Bird = {  }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestCanParseMultipleEntities(t *testing.T) {
	bird1 := testEntity("Bird", "", NodeSet{})
	bird2 := testEntity("Bird2", "", NodeSet{})
	testRoot := RootNode(nil, NodeSet{bird1, bird2})
	actual, err := runParser("Bird = {  }\nBird2 = { }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParsesChildEntity(t *testing.T) {
	entity := testEntity("Robin", "", NodeSet{})
	entity.Children[1].Related = IdNode(nil, "Bird")
	testRoot := RootNode(nil, NodeSet{entity})
	actual, err := runParser("Robin = Bird {  }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParsesBasicGenerationStatement(t *testing.T) {
	args := NodeSet{IntLiteralNode(nil, 1)}
	genBird := GenNode(nil, IdNode(nil, "Bird"), args)
	testRoot := RootNode(nil, NodeSet{genBird})
	actual, err := runParser("generate(1, Bird)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestCanParseMultipleGenerationStatements(t *testing.T) {
	arg := IntLiteralNode(nil, 1)
	genBird := GenNode(nil, IdNode(nil, "Bird"), NodeSet{arg})
	bird2Gen := GenNode(nil, IdNode(nil, "Bird2"), NodeSet{arg})
	testRoot := RootNode(nil, NodeSet{genBird, bird2Gen})

	actual, err := runParser("generate(1, Bird)\ngenerate(1, Bird2)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestCanOverrideFieldInGenerateStatement(t *testing.T) {
	arg := IntLiteralNode(nil, 1)
	value := StrLiteralNode(nil, "birdie")
	field := StaticFieldNode(nil, IdNode(nil, "name"), value, nil)
	genBird := GenNode(nil, testEntity("", "Bird", NodeSet{field}), NodeSet{arg})
	testRoot := RootNode(nil, NodeSet{genBird})
	actual, err := runParser("generate(1, Bird { name: \"birdie\" })")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestCanOverrideMultipleFieldsInGenerateStatement(t *testing.T) {
	value1 := StrLiteralNode(nil, "birdie")
	field1 := StaticFieldNode(nil, IdNode(nil, "name"), value1, nil)
	value2 := BuiltinNode(nil, "integer")
	arg1 := IntLiteralNode(nil, 1)
	arg2 := IntLiteralNode(nil, 2)
	field2 := DynamicFieldNode(nil, IdNode(nil, "age"), value2, NodeSet{arg1, arg2}, nil)

	arg := IntLiteralNode(nil, 1)
	genBird := GenNode(nil, testEntity("", "Bird", NodeSet{field1, field2}), NodeSet{arg})
	testRoot := RootNode(nil, NodeSet{genBird})
	actual, err := runParser("generate(1, Bird { name: \"birdie\", age: integer(1,2) })")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParsedBothBasicEntityAndGenerationStatement(t *testing.T) {
	args := NodeSet{IntLiteralNode(nil, 1)}
	genBird := GenNode(nil, IdNode(nil, "Bird"), args)
	bird := testEntity("Bird", "", NodeSet{})
	testRoot := RootNode(nil, NodeSet{bird, genBird})
	actual, err := runParser("Bird = {}\ngenerate (1, Bird)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithDynamicFieldWithBound(t *testing.T) {
	value := BuiltinNode(nil, "string")
	count := RangeNode(nil, IntLiteralNode(nil, 1), IntLiteralNode(nil, 8))
	field := DynamicFieldNode(nil, IdNode(nil, "name"), value, NodeSet{}, count)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("Bird = { name: string<1..8> }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithDynamicFieldWithoutArgs(t *testing.T) {
	value := BuiltinNode(nil, "string")
	field := DynamicFieldNode(nil, IdNode(nil, "name"), value, NodeSet{}, nil)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("Bird = { name: string }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithDynamicFieldWithArgs(t *testing.T) {
	value := BuiltinNode(nil, "string")
	args := NodeSet{IntLiteralNode(nil, 1)}
	field := DynamicFieldNode(nil, IdNode(nil, "name"), value, args, nil)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("Bird = { name: string(1) }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntitywithDynamicFieldWithMultipleArgs(t *testing.T) {
	value := BuiltinNode(nil, "integer")
	arg1 := IntLiteralNode(nil, 1)
	arg2 := IntLiteralNode(nil, 5)
	args := NodeSet{arg1, arg2}
	field := DynamicFieldNode(nil, IdNode(nil, "name"), value, args, nil)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("Bird = { name: integer(1, 5) }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithMultipleFields(t *testing.T) {
	value := BuiltinNode(nil, "string")
	arg := IntLiteralNode(nil, 1)
	field1 := DynamicFieldNode(nil, IdNode(nil, "name"), value, NodeSet{arg}, nil)

	value = BuiltinNode(nil, "integer")
	arg1 := IntLiteralNode(nil, 1)
	arg2 := IntLiteralNode(nil, 5)
	args := NodeSet{arg1, arg2}
	field2 := DynamicFieldNode(nil, IdNode(nil, "age"), value, args, nil)

	bird := testEntity("Bird", "", NodeSet{field1, field2})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("Bird = { name: string(1), age: integer(1, 5) }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithStaticField(t *testing.T) {
	value := StrLiteralNode(nil, "birdie")
	field := StaticFieldNode(nil, IdNode(nil, "name"), value, nil)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("Bird = { name: \"birdie\" }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithEntityDeclarationField(t *testing.T) {
	args := NodeSet{}
	goatValue := StrLiteralNode(nil, "billy")
	goatField := StaticFieldNode(nil, IdNode(nil, "name"), goatValue, nil)
	goat := testEntity("Goat", "", NodeSet{goatField})
	field := DynamicFieldNode(nil, IdNode(nil, "pet"), goat, args, nil)
	person := testEntity("Person", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{person})
	actual, err := runParser("Person = { pet: Goat = { name: \"billy\" } }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithEntityReferenceField(t *testing.T) {
	args := NodeSet{}
	goatValue := StrLiteralNode(nil, "billy")
	goatField := StaticFieldNode(nil, IdNode(nil, "name"), goatValue, nil)
	goat := testEntity("Goat", "", NodeSet{goatField})
	value := IdNode(nil, "Goat")
	field := DynamicFieldNode(nil, IdNode(nil, "pet"), value, args, nil)
	person := testEntity("Person", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{goat, person})
	actual, err := runParser("Goat = { name: \"billy\" } Person = { pet: Goat }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestVariableAssignment(t *testing.T) {
	value := StrLiteralNode(nil, "hello")
	name := StaticFieldNode(nil, IdNode(nil, "name"), value, nil)
	foo := testEntity("Foo", "", NodeSet{name})
	testRoot := RootNode(nil, NodeSet{
		foo,
		AssignNode(nil,
			IdNode(nil, "foos"),
			GenNode(nil, IdNode(nil, "Foo"), NodeSet{IntLiteralNode(nil, 3)}),
		),
	})
	actual, err := runParser("Foo = { name: \"hello\" } foos = generate(3, Foo)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
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
	_, err := runParser("Blah = { name: }")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestEntityDefinitionRequiresCurlyBrackets(t *testing.T) {
	expectedErrMessage := `Missing right-hand of assignment expression "Bird ="`
	_, err := runParser("Bird =")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestFieldListWithoutCommas(t *testing.T) {
	expectedErrMessage := `Multiple field declarations must be delimited with a comma`
	_, err := runParser("Bird = { h: string b: string }")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestIllegalIdentifiers(t *testing.T) {
	specs := map[string]string{
		"4 = { }":            `Illegal identifier "4"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"$eek = { }":         `Illegal identifier "$eek"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"generate (1, $eek)": `Illegal identifier "$eek"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"e$ek = { }":         `Illegal identifier "e$ek"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"generate = { }":     `Illegal identifier: "generate" is a reserved word`,
	}

	for spec, expectedErrMessage := range specs {
		_, err := runParser(spec)
		ExpectsError(t, expectedErrMessage, removeLocationInfo(err))

	}
}
