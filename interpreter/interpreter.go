package interpreter

import (
	"fmt"
	"github.com/ThoughtWorksStudios/datagen/dsl"
	"github.com/ThoughtWorksStudios/datagen/generator"
	"os"
	"strconv"
	"strings"
	"time"
)

func debug(format string, tokens ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", tokens...)
}

type NamespaceCounter map[string]int

var AnonExtendNames NamespaceCounter = make(NamespaceCounter)

func (c NamespaceCounter) Next(key string) int {
	if ctr, hasKey := c[key]; hasKey {
		ctr += 1
		c[key] = ctr
		return ctr
	} else {
		c[key] = 1
		return 1
	}
}

func (c NamespaceCounter) NextAsStr(key string) string {
	return strconv.Itoa(c.Next(key))
}

type Interpreter struct {
	entities map[string]*generator.Generator // TODO: should probably be a more generic symbol table or possibly the parent scope
	output   GenerationOutput
}

func New() *Interpreter {
	return &Interpreter{
		entities: make(map[string]*generator.Generator),
		output:   GenerationOutput{},
	}
}

func (i *Interpreter) SetCustomDictonaryPath(path string) {
	generator.CustomDictPath = path
}

func (i *Interpreter) WriteGeneratedContent(dest string, filePerEntity bool) error {
	if filePerEntity {
		return i.output.writeFilePerKey()
	} else {
		return i.output.writeToFile(dest)
	}
}

func (i *Interpreter) LoadFile(filepath string, scope *Scope) error {
	if alreadyImported, e := scope.imports.HaveSeen(filepath); e == nil {
		if alreadyImported {
			return nil
		}

		if parsed, pe := i.parseFile(filepath); pe == nil {
			ast := parsed.(dsl.Node)
			if err := i.Visit(ast, scope); err == nil {
				scope.imports.MarkSeen(filepath)
				return nil
			} else {
				return err
			}
		} else {
			return pe
		}
	} else {
		return e
	}
}

func (i *Interpreter) CheckFile(filename string) error {
	_, errors := i.parseFile(filename)
	return errors
}

/**
 * yes, this is practically the exact implementation of dsl.ParseFile(), with the exception
 * of named return values; I believe it is this difference that accounts for parse errors
 * being swallowed by the generated dsl.ParseFile(). we should submit a PR for this.
 */
func (i *Interpreter) parseFile(filename string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()

	return dsl.ParseReader(filename, f, dsl.GlobalStore("filename", filename), dsl.Recover(true))
}

func (i *Interpreter) Visit(node dsl.Node, scope *Scope) error {
	switch node.Kind {
	case "root":
		var err error
		node.Children.Each(func(env *dsl.IterEnv, node dsl.Node) {
			if err = i.Visit(node, scope); err != nil {
				env.Halt()
			}
		})
		return err
	case "entity":
		_, err := i.EntityFromNode(node, scope)
		return err
	case "generation":
		return i.GenerateFromNode(node, scope)
	case "import":
		return i.LoadFile(node.ValStr(), scope)
	default:
		return node.Err("Unexpected token type %s", node.Kind)
	}
}

func (i *Interpreter) defaultArgumentFor(fieldType string) (interface{}, error) {
	switch fieldType {
	case "string":
		return 5, nil
	case "integer":
		return [2]int{1, 10}, nil
	case "decimal":
		return [2]float64{1, 10}, nil
	case "date":
		t1, _ := time.Parse("2006-01-02", "1945-01-01")
		t2, _ := time.Parse("2006-01-02", "2017-01-01")
		return [2]time.Time{t1, t2}, nil
	case "entity", "identifier":
		return 1, nil
	default:
		return nil, fmt.Errorf("Field of type `%s` requires arguments", fieldType)
	}
}

func (i *Interpreter) EntityFromNode(node dsl.Node, scope *Scope) (*generator.Generator, error) {
	// create child scope for entities - much like JS function scoping
	parentScope := scope
	scope = ExtendScope(scope)

	var entity *generator.Generator
	formalName := node.Name

	if node.HasRelation() {
		symbol := node.Related.ValStr()
		if parent, e := i.ResolveEntity(*node.Related, scope); nil == e {

			if formalName == "" {
				formalName = strings.Join([]string{"$" + AnonExtendNames.NextAsStr(symbol), symbol}, "::")
			}

			entity = generator.ExtendGenerator(formalName, parent)
			entity.Base = symbol
		} else {
			return nil, node.Err("Cannot resolve parent entity %q for entity %q", symbol, formalName)
		}
	} else {
		if formalName == "" {
			formalName = "$" + AnonExtendNames.NextAsStr("$")
		}
		entity = generator.NewGenerator(formalName, nil)
	}

	// Add entity to symbol table before iterating through field defs so fields can reference
	// the current entity. Currently, though, this will be problematic as we don't have a nullable
	// option for fields. The workaround is to inline override.
	parentScope.SetSymbol(formalName, "entity", entity)

	for _, field := range node.Children {
		if field.Kind != "field" {
			return nil, field.Err("Expected a `field` declaration, but instead got `%s`", field.Kind) // should never get here
		}

		fieldType := field.ValNode().Kind

		switch {
		case "identifier" == fieldType:
			fallthrough
		case "entity" == fieldType:
			fallthrough
		case "builtin" == fieldType:
			if err := i.withDynamicField(entity, field, scope); err != nil {
				return nil, field.WrapErr(err)
			}
		case strings.HasPrefix(fieldType, "literal-"):
			if err := i.withStaticField(entity, field); err != nil {
				return nil, field.WrapErr(err)
			}
		default:
			return nil, field.Err("Unexpected field type %s; field declarations must be either a built-in type or a literal value", fieldType)
		}
	}

	return entity, nil
}

func valStr(n dsl.Node) string {
	return n.Value.(string)
}

func valInt(n dsl.Node) int {
	return int(n.Value.(int64))
}

func valFloat(n dsl.Node) float64 {
	return n.Value.(float64)
}

func valTime(n dsl.Node) time.Time {
	return n.Value.(time.Time)
}

type Validator func(n dsl.Node) error

func assertValStr(n dsl.Node) error {
	if _, ok := n.Value.(string); !ok {
		return n.Err("Expected %v to be a string, but was %T.", n.Value, n.Value)
	}
	return nil
}

func assertValInt(n dsl.Node) error {
	if _, ok := n.Value.(int64); !ok {
		return n.Err("Expected %v to be an integer, but was %T.", n.Value, n.Value)
	}
	return nil
}

func assertValFloat(n dsl.Node) error {
	if _, ok := n.Value.(float64); !ok {
		return n.Err("Expected %v to be a decimal, but was %T.", n.Value, n.Value)
	}
	return nil
}

func assertValTime(n dsl.Node) error {
	if _, ok := n.Value.(time.Time); !ok {
		return n.Err("Expected %v to be a datetime, but was %T.", n.Value, n.Value)
	}
	return nil
}

func expectsArgs(num int, fn Validator, fieldType string, args dsl.NodeSet) error {
	if l := len(args); num != l {
		return args[0].Err("Field type `%s` expected %d args, but %d found.", fieldType, num, l)
	}

	var er error

	args.Each(func(env *dsl.IterEnv, node dsl.Node) {
		if er = fn(node); er != nil {
			env.Halt()
		}
	})

	return er
}

func (i *Interpreter) withStaticField(entity *generator.Generator, field dsl.Node) error {
	fieldValue := field.Value.(dsl.Node).Value
	return entity.WithStaticField(field.Name, fieldValue)
}

func (i *Interpreter) withDynamicField(entity *generator.Generator, field dsl.Node, scope *Scope) error {
	var err error

	fieldVal := field.ValNode()
	var fieldType string

	if fieldVal.Kind == "builtin" {
		fieldType = fieldVal.ValStr()
	} else {
		fieldType = fieldVal.Kind
	}

	if 0 == len(field.Args) {
		arg, e := i.defaultArgumentFor(fieldType)
		if e != nil {
			return field.WrapErr(e)
		} else {
			if fieldVal.Kind == "builtin" {
				return entity.WithField(field.Name, fieldType, arg)
			}

			if nested, e := i.expectEntity(fieldVal, scope); e != nil {
				return e
			} else {
				return entity.WithEntityField(field.Name, nested, arg)
			}
		}
	}

	switch fieldType {
	case "integer":
		if err = expectsArgs(2, assertValInt, fieldType, field.Args); err == nil {
			return entity.WithField(field.Name, fieldType, [2]int{valInt(field.Args[0]), valInt(field.Args[1])})
		}
	case "decimal":
		if err = expectsArgs(2, assertValFloat, fieldType, field.Args); err == nil {
			return entity.WithField(field.Name, fieldType, [2]float64{valFloat(field.Args[0]), valFloat(field.Args[1])})
		}
	case "string":
		if err = expectsArgs(1, assertValInt, fieldType, field.Args); err == nil {
			return entity.WithField(field.Name, fieldType, valInt(field.Args[0]))
		}
	case "dict":
		if err = expectsArgs(1, assertValStr, fieldType, field.Args); err == nil {
			return entity.WithField(field.Name, fieldType, valStr(field.Args[0]))
		}
	case "date":
		if err = expectsArgs(2, assertValTime, fieldType, field.Args); err == nil {
			return entity.WithField(field.Name, fieldType, [2]time.Time{valTime(field.Args[0]), valTime(field.Args[1])})
		}
	case "identifier", "entity":
		if nested, e := i.expectEntity(fieldVal, scope); e != nil {
			return e
		} else {
			/*
			 * TODO: rethink args for entity types because it's not consistent with other types; here,
			 * it serves as a way to generate multiple values, whereas in all other types it does not.
			 * we should have a consistent syntax for creating multi-value fields
			 */
			if err = expectsArgs(1, assertValInt, "entity", field.Args); err == nil {
				return entity.WithEntityField(field.Name, nested, valInt(field.Args[0]))
			}
		}
	}
	return err
}

func (i *Interpreter) expectEntity(entityRef dsl.Node, scope *Scope) (*generator.Generator, error) {
	switch entityRef.Kind {
	case "identifier":
		return i.ResolveEntity(entityRef, scope)
	case "entity":
		return i.EntityFromNode(entityRef, scope)
	default:
		return nil, entityRef.Err("Expected an entity expression or reference, but got %q", entityRef.Kind)
	}
}

/*
 * A convenience wrapper for ResolveIdentifier, which casts to *generator.Generator. Currently, this
 * is the only type of value that is in the symbol table, but we may support other types in the future
 */
func (i *Interpreter) ResolveEntity(identifierNode dsl.Node, scope *Scope) (*generator.Generator, error) {
	if resolved, err := i.ResolveIdentifier(identifierNode, scope); err != nil {
		return nil, err
	} else {
		if entity, ok := resolved.Value.(*generator.Generator); ok {
			return entity, nil
		} else {
			return nil, identifierNode.Err("identifier %q should refer to an entity, but instead was <type: %s, resolved: %v>", identifierNode.ValStr(), resolved.Type, resolved.Value)
		}
	}
}

func (i *Interpreter) ResolveIdentifier(identiferNode dsl.Node, scope *Scope) (*ScopeEntry, error) {
	if scope == nil {
		return nil, identiferNode.Err("Scope is missing! This should be impossible.")
	}

	if identiferNode.Kind != "identifier" {
		return nil, identiferNode.Err("Expected an identifier, but got %s", identiferNode.Kind)
	}

	if v := scope.ResolveSymbol(identiferNode.ValStr()); v != nil {
		return v, nil
	}

	return nil, identiferNode.Err("Cannot resolve symbol %q", identiferNode.ValStr())
}

func (i *Interpreter) GenerateFromNode(generationNode dsl.Node, scope *Scope) error {
	var entityGenerator *generator.Generator

	entity := generationNode.ValNode()
	switch entity.Kind {
	case "identifier":
		if g, e := i.ResolveEntity(entity, scope); nil != e {
			return e
		} else {
			entityGenerator = g
		}
	case "entity":
		if g, e := i.EntityFromNode(entity, scope); e != nil {
			return e
		} else {
			entityGenerator = g
		}
	default:
		return generationNode.Err("Unexpected node type %q; node is %v", entity.Kind, entity)
	}

	if 0 == len(generationNode.Args) {
		return generationNode.Err("generate requires an argument")
	}
	count, ok := generationNode.Args[0].Value.(int64)

	if !ok {
		return generationNode.Err("generate %q takes an integer count", entityGenerator.Name)
	}

	if count < int64(1) {
		return generationNode.Err("Must generate at least 1 `%s` entity", entityGenerator.Name)
	}

	i.output.addAndAppend(entityGenerator.Name, entityGenerator.Generate(count))
	return nil
}
