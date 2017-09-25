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

func testEntity(name, parent string, fields NodeSet) *Node {
	var parentNode *Node

	if "" != parent {
		parentNode = IdNode(nil, parent)
	}

	var body *Node

	if len(fields) > 0 {
		body = EntityBodyNode(nil, nil, FieldSetNode(nil, fields))
	} else {
		body = EntityBodyNode(nil, nil, nil)
	}

	if "" != name {
		return EntityNode(nil, IdNode(nil, name), parentNode, body)
	} else {
		return EntityNode(nil, nil, parentNode, body)
	}
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
	actual, err := runParser("entity Bird {  }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseBinaryExpression(t *testing.T) {
	expected := RootNode(nil, NodeSet{
		&Node{
			Name: "+", Kind: "binary",
			Value: AtomicNode(nil, IntLiteralNode(nil, 1)), Related: &Node{
				Name: "+", Kind: "binary",
				Value: IntLiteralNode(nil, 2), Related: IntLiteralNode(nil, 3),
			},
		},
	})

	actual, err := runParser(`1 + 2 + 3`)

	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, expected.String(), actual.(*Node).String())
}

func TestBinaryExpressionOperatorPrecedence(t *testing.T) {
	expected := RootNode(nil, NodeSet{
		&Node{
			Name: "+", Kind: "binary",
			Value: &Node{
				Name: "*", Kind: "binary",
				Value: AtomicNode(nil, IntLiteralNode(nil, 1)), Related: IntLiteralNode(nil, 2),
			}, Related: IntLiteralNode(nil, 3),
		},
	})

	actual, err := runParser(`1 * 2 + 3`)
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, expected.String(), actual.(*Node).String())

	expected = RootNode(nil, NodeSet{
		&Node{
			Name: "+", Kind: "binary",
			Value: AtomicNode(nil, IntLiteralNode(nil, 1)), Related: &Node{
				Name: "*", Kind: "binary",
				Value: IntLiteralNode(nil, 2), Related: IntLiteralNode(nil, 3),
			},
		},
	})

	actual, err = runParser(`1 + 2 * 3`)

	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, expected.String(), actual.(*Node).String())
}

func TestCanParseMultipleEntities(t *testing.T) {
	bird1 := testEntity("Bird", "", NodeSet{})
	bird2 := testEntity("Bird2", "", NodeSet{})
	testRoot := RootNode(nil, NodeSet{bird1, bird2})
	actual, err := runParser("entity Bird {  }\nentity Bird2 { }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParsesChildEntity(t *testing.T) {
	entity := testEntity("Robin", "", NodeSet{})
	entity.Related = IdNode(nil, "Bird")
	testRoot := RootNode(nil, NodeSet{entity})
	actual, err := runParser("entity Robin << Bird {  }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParsesBasicGenerationStatement(t *testing.T) {
	args := NodeSet{IntLiteralNode(nil, 1), IdNode(nil, "Bird")}
	genBird := GenNode(nil, args)
	testRoot := RootNode(nil, NodeSet{genBird})
	actual, err := runParser("generate(1, Bird)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestCanParseMultipleGenerationStatements(t *testing.T) {
	arg := IntLiteralNode(nil, 1)
	genBird := GenNode(nil, NodeSet{arg, IdNode(nil, "Bird")})
	bird2Gen := GenNode(nil, NodeSet{arg, IdNode(nil, "Bird2")})
	testRoot := RootNode(nil, NodeSet{genBird, bird2Gen})

	actual, err := runParser("generate(1, Bird)\ngenerate(1, Bird2)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestCanOverrideFieldInGenerateStatement(t *testing.T) {
	arg := IntLiteralNode(nil, 1)
	value := StrLiteralNode(nil, "birdie")
	field := ExpressionFieldNode(nil, IdNode(nil, "name"), value, nil)
	genBird := GenNode(nil, NodeSet{arg, testEntity("", "Bird", NodeSet{field})})
	testRoot := RootNode(nil, NodeSet{genBird})
	actual, err := runParser("generate(1, Bird << { name: \"birdie\" })")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestCanOverrideMultipleFieldsInGenerateStatement(t *testing.T) {
	value1 := StrLiteralNode(nil, "birdie")
	field1 := ExpressionFieldNode(nil, IdNode(nil, "name"), value1, nil)
	value2 := BuiltinNode(nil, INT_TYPE)
	arg1 := IntLiteralNode(nil, 1)
	arg2 := IntLiteralNode(nil, 2)
	field2 := DynamicFieldNode(nil, IdNode(nil, "age"), value2, NodeSet{arg1, arg2}, nil, false)

	arg := IntLiteralNode(nil, 1)
	genBird := GenNode(nil, NodeSet{arg, testEntity("", "Bird", NodeSet{field1, field2})})
	testRoot := RootNode(nil, NodeSet{genBird})
	actual, err := runParser("generate(1, Bird << { name: \"birdie\", age: $int(1,2) })")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParsedBothBasicEntityAndGenerationStatement(t *testing.T) {
	args := NodeSet{IntLiteralNode(nil, 1), IdNode(nil, "Bird")}
	genBird := GenNode(nil, args)
	bird := testEntity("Bird", "", NodeSet{})
	testRoot := RootNode(nil, NodeSet{bird, genBird})
	actual, err := runParser("entity Bird {}\ngenerate (1, Bird)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithDynamicFieldWithBound(t *testing.T) {
	value := BuiltinNode(nil, STRING_TYPE)
	count := RangeNode(nil, IntLiteralNode(nil, 1), IntLiteralNode(nil, 8))
	field := DynamicFieldNode(nil, IdNode(nil, "name"), value, NodeSet{}, count, false)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { name: $str()<1..8> }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithDynamicFieldWithoutArgs(t *testing.T) {
	value := BuiltinNode(nil, STRING_TYPE)
	field := DynamicFieldNode(nil, IdNode(nil, "name"), value, NodeSet{}, nil, false)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { name: $str() }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithDistributedFieldWithoutArgs(t *testing.T) {
	value := BuiltinNode(nil, INT_TYPE)
	distField := NodeSet{DynamicFieldNode(nil, IdNode(nil, ""), value, NodeSet{}, nil, false)}

	field := DistributionFieldNode(nil, IdNode(nil, "age"), DistributionNode(nil, "normal"), distField)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { age: distribution(normal, $int()) }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithDistributedStaticField(t *testing.T) {
	value := StrLiteralNode(nil, "blah")
	distField := DynamicFieldNode(nil, IdNode(nil, ""), value, NodeSet{}, nil, false)
	distField.Weight = 10.0

	field := DistributionFieldNode(nil, IdNode(nil, "age"), DistributionNode(nil, "percent"), NodeSet{distField})
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { age: distribution(percent, 10 => \"blah\") }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithDistributedFieldWithArgs(t *testing.T) {
	value := BuiltinNode(nil, INT_TYPE)
	arg1 := IntLiteralNode(nil, 1)
	arg2 := IntLiteralNode(nil, 50)
	args := NodeSet{arg1, arg2}
	distField := NodeSet{DynamicFieldNode(nil, IdNode(nil, ""), value, args, nil, false)}
	field := DistributionFieldNode(nil, IdNode(nil, "age"), DistributionNode(nil, "normal"), distField)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { age: distribution(normal, $int(1,50))}")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithSubDistributionShouldError(t *testing.T) {
	_, err := runParser("entity Bird { age: distribution(percent, 10 => distribution(normal, $float(1.0,100.0)) }")
	ExpectsError(t, "Cannot use distributions as an argument to another distribution", err)
}

func TestParseEntityWithDistributedFieldWithPercentArgs(t *testing.T) {
	value := BuiltinNode(nil, INT_TYPE)
	arg1 := IntLiteralNode(nil, 1)
	arg2 := IntLiteralNode(nil, 50)
	args := NodeSet{arg1, arg2}
	f := DynamicFieldNode(nil, IdNode(nil, ""), value, args, nil, false)
	f.Weight = 18.0
	distField := NodeSet{f}
	field := DistributionFieldNode(nil, IdNode(nil, "age"), DistributionNode(nil, "percent"), distField)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { age: distribution(percent, 18 => $int(1,50))}")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithDistributedFieldWithWeightedArgs(t *testing.T) {
	value := BuiltinNode(nil, INT_TYPE)
	arg1 := IntLiteralNode(nil, 1)
	arg2 := IntLiteralNode(nil, 50)
	args := NodeSet{arg1, arg2}
	f := DynamicFieldNode(nil, IdNode(nil, ""), value, args, nil, false)
	f.Weight = 18.0
	distField := NodeSet{f}
	field := DistributionFieldNode(nil, IdNode(nil, "age"), DistributionNode(nil, "weighted"), distField)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { age: distribution(weighted, 18% => $int(1,50))}")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithDynamicFieldWithUniqueFlag(t *testing.T) {
	value := BuiltinNode(nil, STRING_TYPE)
	args := NodeSet{IntLiteralNode(nil, 1)}
	field := DynamicFieldNode(nil, IdNode(nil, "name"), value, args, nil, true)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { name: $str(1) unique }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithExpressionFieldWithUniqueFlag(t *testing.T) {
	_, err := runParser("entity Bird { name: \"blah\" unique }")
	expectedErrorMsg := "Expression fields do not support uniqueness"
	ExpectsError(t, expectedErrorMsg, removeLocationInfo(err))
}

func TestParseEntityWithDynamicFieldWithArgs(t *testing.T) {
	value := BuiltinNode(nil, STRING_TYPE)
	args := NodeSet{IntLiteralNode(nil, 1)}
	field := DynamicFieldNode(nil, IdNode(nil, "name"), value, args, nil, false)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { name: $str(1) }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithBuiltinFieldWithMultipleArgs(t *testing.T) {
	value := BuiltinNode(nil, INT_TYPE)
	arg1 := IntLiteralNode(nil, 1)
	arg2 := IntLiteralNode(nil, 5)
	args := NodeSet{arg1, arg2}
	field := DynamicFieldNode(nil, IdNode(nil, "name"), value, args, nil, false)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { name: $int(1, 5) }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithMultipleFields(t *testing.T) {
	value := BuiltinNode(nil, STRING_TYPE)
	arg := IntLiteralNode(nil, 1)
	field1 := DynamicFieldNode(nil, IdNode(nil, "name"), value, NodeSet{arg}, nil, false)

	value = BuiltinNode(nil, INT_TYPE)
	arg1 := IntLiteralNode(nil, 1)
	arg2 := IntLiteralNode(nil, 5)
	args := NodeSet{arg1, arg2}
	field2 := DynamicFieldNode(nil, IdNode(nil, "age"), value, args, nil, false)

	bird := testEntity("Bird", "", NodeSet{field1, field2})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { name: $str(1), age: $int(1, 5) }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithStaticField(t *testing.T) {
	value := StrLiteralNode(nil, "birdie")
	field := ExpressionFieldNode(nil, IdNode(nil, "name"), value, nil)
	bird := testEntity("Bird", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{bird})
	actual, err := runParser("entity Bird { name: \"birdie\" }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithEntityDeclarationField(t *testing.T) {
	args := NodeSet{}
	goatValue := StrLiteralNode(nil, "billy")
	goatField := ExpressionFieldNode(nil, IdNode(nil, "name"), goatValue, nil)
	goat := testEntity("Goat", "", NodeSet{goatField})
	field := DynamicFieldNode(nil, IdNode(nil, "pet"), goat, args, nil, false)
	person := testEntity("Person", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{person})
	actual, err := runParser("entity Person { pet: entity Goat { name: \"billy\" } }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestParseEntityWithEntityReferenceField(t *testing.T) {
	args := NodeSet{}
	goatValue := StrLiteralNode(nil, "billy")
	goatField := ExpressionFieldNode(nil, IdNode(nil, "name"), goatValue, nil)
	goat := testEntity("Goat", "", NodeSet{goatField})
	value := IdNode(nil, "Goat")
	field := DynamicFieldNode(nil, IdNode(nil, "pet"), value, args, nil, false)
	person := testEntity("Person", "", NodeSet{field})
	testRoot := RootNode(nil, NodeSet{goat, person})
	actual, err := runParser("entity Goat { name: \"billy\" } entity Person { pet: Goat }")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestVariableAssignment(t *testing.T) {
	value := StrLiteralNode(nil, "hello")
	name := ExpressionFieldNode(nil, IdNode(nil, "name"), value, nil)
	foo := testEntity("Foo", "", NodeSet{name})
	testRoot := RootNode(nil, NodeSet{
		foo,
		AssignNode(nil,
			IdNode(nil, "foos"),
			GenNode(nil, NodeSet{IntLiteralNode(nil, 3), IdNode(nil, "Foo")}),
		),
	})
	actual, err := runParser("entity Foo { name: \"hello\" } foos = generate(3, Foo)")
	AssertNil(t, err, "Didn't expect to get an error: %v", err)
	AssertEqual(t, testRoot.String(), actual.(*Node).String())
}

func TestRequiresValidStatements(t *testing.T) {
	_, err := runParser("!")
	expectedErrorMsg := `Don't know how to evaluate "!"`
	ExpectsError(t, expectedErrorMsg, removeLocationInfo(err))
}

func TestGenerateWithNoArgumentsProducesError(t *testing.T) {
	expectedErrMessage := "`generate` statement \"generate Blah\" requires arguments `(count, entity)`"
	_, err := runParser("generate Blah")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestEntityFieldRequiresType(t *testing.T) {
	expectedErrMessage := `Missing field type for field declaration "name"`
	_, err := runParser("entity Blah { name: }")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestAssignmentMissingRightHandSide(t *testing.T) {
	expectedErrMessage := `Missing right-hand of assignment expression "Bird ="`
	_, err := runParser("Bird =")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestEntityDefinitionRequiresCurlyBrackets(t *testing.T) {
	expectedErrMessage := `Unterminated entity expression (missing closing curly brace)`
	_, err := runParser("entity Bird {")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestFieldListWithoutCommas(t *testing.T) {
	expectedErrMessage := `Multiple field declarations must be delimited with a comma`
	_, err := runParser("entity Bird { h: string b: string }")
	ExpectsError(t, expectedErrMessage, removeLocationInfo(err))
}

func TestIllegalIdentifiers(t *testing.T) {
	specs := map[string]string{
		"let 2hot = true":           `Illegal identifier "2hot"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"1luv = [1, 2]":             `Illegal identifier "1luv"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"entity a {$ok: 0}":         `Illegal identifier "$ok"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"entity 4fun { }":           `Illegal identifier "4fun"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"entity $eek { }":           `Illegal identifier "$eek"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"generate (1, $a)":          `Illegal identifier "$a"; identifiers start with a letter or underscore, followed by zero or more letters, underscores, and numbers`,
		"entity generate { }":       `Illegal identifier "generate"; reserved words cannot be used as identifiers`,
		"entity pk { }":             `Illegal identifier "pk"; reserved words cannot be used as identifiers`,
		"entity t {false: string }": `Illegal identifier "false"; reserved words cannot be used as identifiers`,
		"entity = [1, 2]":           `Illegal identifier "entity"; reserved words cannot be used as identifiers`,
		"generate(1, generate)":     `Illegal identifier "generate"; reserved words cannot be used as identifiers`,
		"[let]":                     `Illegal identifier "let"; reserved words cannot be used as identifiers`,
		"entity t << generate {}":   `Illegal identifier "generate"; reserved words cannot be used as identifiers`,
		"generate << {}":            `Illegal identifier "generate"; reserved words cannot be used as identifiers`,
	}

	for spec, expectedErrMessage := range specs {
		_, err := runParser(spec)
		ExpectsError(t, expectedErrMessage, removeLocationInfo(err))

	}
}
