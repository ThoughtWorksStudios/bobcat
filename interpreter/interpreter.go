package interpreter

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/dsl"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"github.com/ThoughtWorksStudios/bobcat/generator"
	"os"
	"strings"
	"time"
)

// Might be useful to pull these out into another file
var UNIX_EPOCH time.Time
var NOW time.Time

const (
	PK_FIELD_CONFIG = "$PK_FIELD"
)

func init() {
	UNIX_EPOCH, _ = time.Parse("2006-01-02", "1970-01-01")
	NOW = time.Now()
}

var AnonExtendNames NamespaceCounter = make(NamespaceCounter)

type Interpreter struct {
	disableMetadata bool
	basedir         string
	emitter         Emitter
	dryRun          bool
}

func New(emitter Emitter, disableMetadata bool) *Interpreter {
	return &Interpreter{
		emitter:         emitter,
		basedir:         ".",
		disableMetadata: disableMetadata,
	}
}

func (i *Interpreter) ConfigureDryRun() {
	i.dryRun = true
}

func (i *Interpreter) SetCustomDictonaryPath(path string) {
	generator.CustomDictPath = path
}

func (i *Interpreter) importFile(importNode *Node, scope *Scope, deferred bool) (interface{}, error) {
	if result, err := i.LoadFile(importNode.ValStr(), scope, deferred); err != nil {
		return nil, importNode.WrapErr(err)
	} else {
		return result, nil
	}
}

func (i *Interpreter) LoadFile(filename string, scope *Scope, deferred bool) (interface{}, error) {
	scope.PredefinedDefaults(SymbolTable{
		"NOW":        NOW,
		"UNIX_EPOCH": UNIX_EPOCH,
	})

	original := i.basedir
	realpath, re := Resolve(filename, original)

	if re != nil {
		return nil, re
	}

	if alreadyImported, e := scope.Imports.HaveSeen(realpath); e == nil {
		if alreadyImported {
			return nil, nil
		}
	} else {
		return nil, e
	}

	if base, e := Basedir(filename, original); e == nil {
		i.basedir = base
		defer func() { i.basedir = original }()
	} else {
		return nil, e
	}

	if parsed, pe := parseFile(realpath); pe == nil {
		ast := parsed.(*Node)
		scope.Imports.MarkSeen(realpath) // optimistically mark before walking ast in case the file imports itself

		return i.Visit(ast, scope, deferred)
	} else {
		return nil, pe
	}
}

func (i *Interpreter) CheckFile(filename string) error {
	_, errors := parseFile(filename)
	return errors
}

/**
 * yes, this is practically the exact implementation of dsl.ParseFile(), with the exception
 * of named return values; I believe it is this difference that accounts for parse errors
 * being swallowed by the generated dsl.ParseFile(). we should submit a PR for this.
 */
func parseFile(filename string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()

	return dsl.ParseReader(filename, f, dsl.GlobalStore("filename", filename))
}

func (i *Interpreter) Visit(node *Node, scope *Scope, deferred bool) (interface{}, error) {
	switch node.Kind {
	case "root", "sequential":
		if deferred {
			return i.Compile(node.Children, scope)
		}
		return i.Eval(node.Children, scope)
	case "lambda":
		closure := func(scope *Scope) (interface{}, error) {
			if fn, err := i.Compile(node.Children, scope); err == nil {
				boundArgs := make([]string, len(node.Args))
				for idx, arg := range node.Args {
					boundArgs[idx] = arg.ValStr()
				}

				lambda := &Lambda{Name: node.Name, Params: boundArgs, Executor: fn}

				if node.Name != "" {
					symbol := node.Name

					// TODO: refactor, DRY?
					if s := scope.DefinedInScope(symbol); s == scope {
						Warn("%v Symbol %q has already been declared in this scope", node.Ref, symbol)
					}

					scope.SetSymbol(symbol, lambda)
				}
				return lambda, nil
			} else {
				return nil, err
			}
		}
		if deferred {
			return closure, nil
		}

		return closure(scope)
	case "call":
		lambdaNode := node.ValNode()
		closure := func(scope *Scope) (interface{}, error) {
			if callable, err := i.Visit(lambdaNode, scope, false); err == nil {
				if lambda, ok := callable.(*Lambda); ok {
					if args, err := i.AllValuesFromNodeSet(node.Args, scope, false); err == nil {
						return lambda.Call(scope, args.([]interface{})...)
					} else {
						return nil, err
					}
				} else {
					return nil, lambdaNode.Err("Expected a lambda, but got %v (%T)", callable, callable)
				}
			} else {
				return nil, err
			}
		}

		if deferred {
			return closure, nil
		}

		return closure(scope)
	case "atomic":
		return i.Visit(node.ValNode(), scope, deferred)
	case "binary":
		return i.resolveBinaryNode(node, scope, deferred)
	case "range":
		// currently this only takes literals, so no need to defer.
		// ideally, it should accept expressions (or at least identifiers), and when
		// that happens, we will need to handle a deferral
		return i.RangeFromNode(node, scope)
	case "entity":
		closure := func(scope *Scope) (interface{}, error) {
			return i.EntityFromNode(node, scope, false)
		}

		if deferred {
			return closure, nil
		}

		return closure(scope)
	case "generation":
		closure := func(scope *Scope) (interface{}, error) {
			return i.GenerateFromNode(node, scope, false)
		}
		if deferred {
			return closure, nil
		}

		return closure(scope)
	case "identifier":
		closure := func(scope *Scope) (interface{}, error) {
			if entry, err := i.ResolveIdentifier(node, scope); err == nil {
				return entry, nil
			} else {
				return nil, err
			}
		}

		if deferred {
			return closure, nil
		}
		return closure(scope)
	case "assignment":
		symbol := node.Children[0].ValStr()
		valNode := node.Children[1]
		var value interface{}
		if v, err := i.Visit(valNode, scope, deferred); err == nil {
			value = v
		} else {
			return nil, err
		}

		closure := func(scope *Scope) (interface{}, error) {
			if s := scope.DefinedInScope(symbol); s != nil {
				if v, ok := value.(DeferredResolver); ok {
					if val, err := v(scope); err != nil {
						return nil, err
					} else {
						value = val
					}
				}
				/**
				 * must set in the scope where symbol is defined, which is not
				 * necessarily the current scope. the ability to assign a value
				 * to a symbol in a parent scope is intentional. if you instead
				 * want variable shadowing, use a variable declaration in the
				 * present scope, NOT an assignment expression.
				 */
				s.SetSymbol(symbol, value)
				return value, nil
			} else {
				return nil, node.Err("Cannot assign value; symbol %q has not yet been declared in scope hierarchy", symbol)
			}
		}

		if deferred {
			return closure, nil
		} else {
			return closure(scope)
		}
	case "variable":
		symbol := node.Name
		var value interface{}

		if nil != node.Value {
			valNode := node.ValNode()
			if v, err := i.Visit(valNode, scope, deferred); err == nil {
				value = v
			} else {
				return nil, err
			}
		}

		closure := func(scope *Scope) (interface{}, error) {
			if s := scope.DefinedInScope(symbol); s == scope {
				Warn("%v Symbol %q has already been declared in this scope", node.Ref, symbol)
			}

			if v, ok := value.(DeferredResolver); ok {
				if val, err := v(scope); err != nil {
					return nil, err
				} else {
					value = val
				}
			}

			scope.SetSymbol(symbol, value)
			return value, nil
		}

		if deferred {
			return closure, nil
		} else {
			return closure(scope)
		}
	case "literal-collection":
		return i.AllValuesFromNodeSet(node.Children, scope, deferred)
	case "literal-int":
		return node.ValInt(), nil
	case "literal-float":
		return node.ValFloat(), nil
	case "literal-string":
		return node.ValStr(), nil
	case "literal-bool":
		return node.ValBool(), nil
	case "literal-date":
		return node.ValTime(), nil
	case "literal-null":
		return nil, nil
	case "import":
		// Currently, Import stateents aren't deferrable as they are only allowed
		// in the top level context; imports cannot occur within any other expression
		// that is deferrable (e.g. Entity, Lambda, etc)
		return i.importFile(node, scope, deferred)
	case "primary-key":
		// currently, we don't support deferred eval on pk statements as they are only
		// allowed at top-level and within entity decalarations. If this changes, we need
		// to make some minor modifications here
		if nameVal, err := i.Visit(node.ValNode(), scope, deferred); err != nil {
			return nil, err
		} else {
			if name, ok := nameVal.(string); ok {
				kind := node.Related.ValStr()
				val := generator.NewPrimaryKeyConfig(name, kind)
				scope.SetSymbol(PK_FIELD_CONFIG, val)
				return val, nil
			} else {
				return nil, node.ValNode().Err("Expected a string, but got %v", nameVal)
			}
		}

	case "field":
		//TODO: Change this...
		return node, nil
	default:
		return nil, node.Err("Unexpected token type %s %v", node.Kind, node)
	}
}
func (i *Interpreter) resolveBinaryNode(node *Node, scope *Scope, deferred bool) (interface{}, error) {
	lhs, e1 := i.Visit(node.ValNode(), scope, deferred)
	if e1 != nil {
		return nil, e1
	}

	rhs, e2 := i.Visit(node.Related, scope, deferred)
	if e2 != nil {
		return nil, e2
	}

	return i.ApplyOperator(node.Name, lhs, rhs, scope, deferred)
}

func (i *Interpreter) ApplyOperator(op string, left, right interface{}, scope *Scope, deferred bool) (interface{}, error) {
	switch op {
	case "+", "-":
		switch left.(type) {
		case int64:
			return i.addToInt(op, left.(int64), right, scope, deferred)
		case float64:
			return i.addToFloat(op, left.(float64), right, scope, deferred)
		case string:
			return i.addToString(op, left.(string), right, scope, deferred)
		case bool:
			return i.addToBool(op, left.(bool), right, scope, deferred)
		case DeferredResolver:
			if !deferred {
				if lhs, err := left.(DeferredResolver)(scope); err == nil {
					return i.ApplyOperator(op, lhs, right, scope, false)
				} else {
					return nil, err
				}
			}

			return i.handleDeferredLHS(op, left.(DeferredResolver), right), nil
		default:
			return nil, incompatible(op)
		}
	case "*", "/":
		switch left.(type) {
		case int64:
			return i.multByInt(op, left.(int64), right, scope, deferred)
		case float64:
			return i.multByFloat(op, left.(float64), right, scope, deferred)
		case string:
			return i.multByString(op, left.(string), right, scope, deferred)
		case DeferredResolver:
			if !deferred {
				if lhs, err := left.(DeferredResolver)(scope); err == nil {
					return i.ApplyOperator(op, lhs, right, scope, false)
				} else {
					return nil, err
				}
			}
			return i.handleDeferredLHS(op, left.(DeferredResolver), right), nil
		default:
			return nil, incompatible(op)
		}
	default:
		return nil, fmt.Errorf("Unknown operator %q", op)
	}
}

func (i *Interpreter) AllValuesFromNodeSet(ns NodeSet, scope *Scope, deferred bool) (interface{}, error) {
	result := make([]interface{}, len(ns))
	containsDeferred := false

	for index, child := range ns {
		if item, e := i.Visit(child, scope, deferred); e == nil {
			if _, ok := item.(DeferredResolver); ok {
				containsDeferred = true
			}
			result[index] = item
		} else {
			return nil, e
		}
	}

	if containsDeferred {
		closure := func(scope *Scope) (interface{}, error) {
			resolved := make([]interface{}, len(result))
			for i, item := range result {
				if _, ok := item.(DeferredResolver); ok {
					if r, e := item.(DeferredResolver)(scope); e == nil {
						resolved[i] = r
					} else {
						return nil, e
					}
				} else {
					resolved[i] = item
				}
			}
			return resolved, nil
		}

		return closure, nil
	}

	return result, nil
}

func (i *Interpreter) Eval(expressions NodeSet, scope *Scope) (interface{}, error) {
	var val interface{}
	var err error
	for _, node := range expressions {
		if val, err = i.Visit(node, scope, false); err != nil {
			return nil, err
		}
	}
	return val, nil
}

func (i *Interpreter) Compile(expressions NodeSet, scope *Scope) (DeferredResolver, error) {
	queue := make([]interface{}, len(expressions))
	for idx, node := range expressions {
		if item, err := i.Visit(node, scope, true); err != nil {
			return nil, err
		} else {
			queue[idx] = item
		}
	}

	return (&ExecQueue{expr: queue}).Run, nil
}

func (i *Interpreter) RangeFromNode(node *Node, scope *Scope) (*CountRange, error) {
	bounds := make([]int64, 2)

	for idx, n := range node.Children {
		if !n.Is("literal-int") {
			return nil, n.Err("Range bounds must be integers")
		}

		bounds[idx] = n.ValInt()
	}

	return &CountRange{Min: bounds[0], Max: bounds[1]}, nil // TODO: support generic range instead of CountRange?
}

func (i *Interpreter) defaultArgumentFor(fieldType string) (interface{}, error) {
	switch fieldType {
	case "string":
		return int64(5), nil
	case "integer":
		return [2]int64{1, 10}, nil
	case "decimal":
		return [2]float64{1, 10}, nil
	case "date":
		return []interface{}{UNIX_EPOCH, NOW, ""}, nil
	case "entity", "identifier":
		return nil, nil
	case "bool", "serial", "uid":
		return nil, nil
	default:
		return nil, fmt.Errorf("Field of type `%s` requires arguments", fieldType)
	}
}

func (i *Interpreter) EntityFromNode(node *Node, scope *Scope, deferred bool) (*generator.Generator, error) {
	// create child scope for entities - much like JS function scoping
	parentScope := scope
	scope = ExtendScope(scope)

	body := node.ValNode()

	var pk *generator.PrimaryKey

	if nil != body.Related {
		var err error
		if pk, err = i.expectsPrimaryKeyStatement(body.Related, scope, deferred); err != nil {
			return nil, err
		}
	}

	var entity *generator.Generator
	formalName := node.Name

	if node.HasRelation() {
		symbol := node.Related.ValStr()
		if parent, e := i.ResolveEntity(node.Related, scope); nil == e {

			if formalName == "" {
				formalName = strings.Join([]string{"$" + AnonExtendNames.NextAsStr(symbol), symbol}, "::")
			}

			entity = generator.ExtendGenerator(formalName, parent, pk, i.disableMetadata)
		} else {
			return nil, node.Err("Cannot resolve parent entity %q for entity %q", symbol, formalName)
		}
	} else {
		if formalName == "" {
			formalName = "$" + AnonExtendNames.NextAsStr("$")
		}

		if nil == pk {
			pk = i.ResolvePrimaryKeyConfig(scope)
		}
		entity = generator.NewGenerator(formalName, pk, i.disableMetadata)
	}

	// Add entity to symbol table before iterating through field defs so fields can reference
	// the current entity. Currently, though, this will be problematic as we don't have a nullable
	// option for fields. The workaround is to inline override.
	parentScope.SetSymbol(formalName, entity)

	if nil != body.Value {
		fieldsetNode := body.ValNode()

		if !fieldsetNode.Is("field-set") {
			return nil, fieldsetNode.Err("Expected a fieldset, but got %q", fieldsetNode.Kind)
		}

		for _, field := range fieldsetNode.Children {
			if !field.Is("field") && !field.Is("distribution") {
				return nil, field.Err("Expected a `field` declaration, but instead got `%s`", field.Kind) // should never get here
			}

			fieldType := field.ValNode().Kind

			switch {
			case "distribution" == fieldType:
				if err := i.withDistributionField(entity, field, scope, deferred); err != nil {
					return nil, field.WrapErr(err)
				}
			case "binary" == fieldType:
				var fieldVal interface{}
				var err error

				if fieldVal, err = i.resolveBinaryNode(field.ValNode(), scope, true); err != nil {
					return nil, field.WrapErr(err)
				}

				if err = i.withExpressionField(entity, field.Name, fieldVal); err != nil {
					return nil, field.WrapErr(err)
				}
			case "identifier" == fieldType:
				if symbol, ok := field.ValNode().Value.(string); ok {
					if entity.HasField(symbol) {
						closure := func(scope *Scope) (interface{}, error) {
							if s := scope.DefinedInScope(symbol); nil != s {
								return s.ResolveSymbol(symbol), nil
							}
							return nil, fmt.Errorf("Cannot resolve symbol %q", symbol)
						}
						if err := i.withExpressionField(entity, field.Name, closure); err != nil {
							return nil, field.WrapErr(err)
						}
						continue
					}
				}
				fallthrough
			case "entity" == fieldType:
				fallthrough
			case "builtin" == fieldType:
				if err := i.withDynamicField(entity, field, scope, deferred); err != nil {
					return nil, field.WrapErr(err)
				}
			case strings.HasPrefix(fieldType, "literal-"):
				if err := i.withExpressionField(entity, field.Name, field.ValNode().Value); err != nil {
					return nil, field.WrapErr(err)
				}
			default:
				return nil, field.Err("Unexpected field type %s; field declarations must be either a built-in type or a literal value", fieldType)
			}
		}
	}

	return entity, nil
}

type Validator func(v interface{}, index int) error

func assertValStr(v interface{}, index int) error {
	if _, ok := v.(string); !ok {
		return fmt.Errorf("Expected %v to be a string, but was %T.", v, v)
	}
	return nil
}

func assertCollection(v interface{}, index int) error {
	if _, ok := v.([]interface{}); !ok {
		return fmt.Errorf("Expected %v to be a collection, but was %T.", v, v)
	}
	return nil
}

func assertValInt(v interface{}, index int) error {
	if _, ok := v.(int64); !ok {
		return fmt.Errorf("Expected %v to be an integer, but was %T.", v, v)
	}
	return nil
}

func assertValFloat(v interface{}, index int) error {
	if _, ok := v.(float64); !ok {
		return fmt.Errorf("Expected %v to be a decimal, but was %T.", v, v)
	}
	return nil
}

func assertValTime(v interface{}, index int) error {
	if _, ok := v.(time.Time); !ok {
		return fmt.Errorf("Expected %v to be a datetime, but was %T.", v, v)
	}
	return nil
}

func assertDateFieldArgs(v interface{}, index int) error {
	if index < 2 {
		return assertValTime(v, index)
	}
	return assertValStr(v, index)
}

func expectsArgs(atLeast, atMost int, fn Validator, fieldType string, args []interface{}) error {
	var er error
	var size int

	if nil == args {
		size = 0
	} else {
		size = len(args)
	}

	if atLeast == atMost {
		if atLeast != size {
			return fmt.Errorf("Field type `%s` expected %d args, but %d found.", fieldType, atLeast, size)
		}
	} else {
		if size < atLeast || size > atMost {
			return fmt.Errorf("Field type `%s` expected %d - %d args, but %d found.", fieldType, atLeast, atMost, size)
		}
	}

	if size > 0 && nil != fn {
		for i, val := range args {
			if er = fn(val, i); er != nil {
				return er
			}
		}
	}

	return er
}

func (i *Interpreter) withExpressionField(entity *generator.Generator, fieldName string, fieldValue interface{}) error {
	var err error

	switch val := fieldValue.(type) {
	case DeferredResolver:
		err = entity.WithDeferredField(fieldName, val)
	default:
		err = entity.WithLiteralField(fieldName, val)
	}

	return err
}

func (i *Interpreter) parseArgsForField(fieldType string, args []interface{}) interface{} {
	switch fieldType {
	case "integer":
		return [2]int64{args[0].(int64), args[1].(int64)}
	case "decimal":
		return [2]float64{args[0].(float64), args[1].(float64)}
	case "string":
		return args[0].(int64)
	case "dict":
		return args[0].(string)
	case "date":
		format := ""
		if 3 == len(args) {
			format = args[2].(string)
		}
		return []interface{}{args[0].(time.Time), args[1].(time.Time), format}
	case "enum":
		return args[0].([]interface{})
	default:
		return nil
	}
}

func (i *Interpreter) withDistributionField(entity *generator.Generator, field *Node, scope *Scope, deferred bool) error {
	fieldVal := field.ValNode()
	fieldType := fieldVal.ValStr()
	if 0 == len(field.Args) {
		return field.Err("Distributions require a domain")
	}
	a, err := i.AllValuesFromNodeSet(field.Args, scope, false)
	if err != nil {
		return field.WrapErr(err)
	}

	args, _ := a.([]interface{})
	weights := make([]float64, len(args))
	argTypes := make([]string, len(args))
	arguments := make([]interface{}, len(args))

	for p := 0; p < len(args); p++ {
		arg := args[p].(*Node)
		argVal := arg.ValNode()
		var argType string
		weights[p] = args[p].(*Node).Weight

		if argVal.Is("builtin") {
			argType = argVal.ValStr()
		} else {
			argType = argVal.Kind
		}

		argTypes[p] = argType

		switch {
		case strings.HasPrefix(argType, "literal-"):
			argTypes[p] = "static"
			arguments[p] = argVal.Value
		case argType == "identifier":
			if entry, err := i.ResolveIdentifier(argVal, scope); err != nil {
				return arg.WrapErr(err)
			} else {
				argTypes[p] = "entity"
				arguments[p] = entry
			}
		case argType == "entity":
			if nested, e := i.expectsEntity(argVal, scope, deferred); e != nil {
				return fieldVal.WrapErr(e)
			} else {
				arguments[p] = nested
			}
		default:
			if len(arg.Args) == 0 {
				arguments[p], _ = i.defaultArgumentFor(argTypes[p])
			} else {
				fieldArgs, err := i.AllValuesFromNodeSet(arg.Args, scope, false)
				if err != nil {
					return arg.WrapErr(err)
				}

				arguments[p] = i.parseArgsForField(argTypes[p], fieldArgs.([]interface{}))
			}
		}
	}

	return entity.WithDistribution(field.Name, fieldType, argTypes, arguments, weights)
}

func (i *Interpreter) withDynamicField(entity *generator.Generator, field *Node, scope *Scope, deferred bool) error {
	var err error

	fieldVal := field.ValNode()
	var fieldType string

	if fieldVal.Is("builtin") {
		fieldType = fieldVal.ValStr()
	} else {
		fieldType = fieldVal.Kind
	}

	var countRange *CountRange

	if nil != field.CountRange {
		if countRange, err = i.expectsRange(field.CountRange, scope); err != nil {
			return err
		}

		if err = countRange.Validate(); err != nil {
			return field.CountRange.WrapErr(err)
		}
	}

	if 0 == len(field.Args) {
		arg, e := i.defaultArgumentFor(fieldType)
		if e != nil {
			return fieldVal.WrapErr(e)
		} else {
			if fieldVal.Is("builtin") {
				return entity.WithField(field.Name, fieldType, arg, countRange, field.Unique)
			}

			if nested, e := i.expectsEntity(fieldVal, scope, deferred); e != nil {
				return fieldVal.WrapErr(e)
			} else {
				return entity.WithEntityField(field.Name, nested, arg, countRange)
			}
		}
	}

	a, e := i.AllValuesFromNodeSet(field.Args, scope, false)

	if e != nil {
		return e
	}

	args, _ := a.([]interface{})

	switch fieldType {
	case "integer":
		if err = expectsArgs(2, 2, assertValInt, fieldType, args); err == nil {
			// return entity.WithField(field.Name, fieldType, [2]int64{args[0].(int64), args[1].(int64)}, countRange, field.Unique)
			return entity.WithField(field.Name, fieldType, i.parseArgsForField(fieldType, args), countRange, field.Unique)
		}
	case "decimal":
		if err = expectsArgs(2, 2, assertValFloat, fieldType, args); err == nil {
			return entity.WithField(field.Name, fieldType, i.parseArgsForField(fieldType, args), countRange, field.Unique)
		}
	case "string":
		if err = expectsArgs(1, 1, assertValInt, fieldType, args); err == nil {
			return entity.WithField(field.Name, fieldType, i.parseArgsForField(fieldType, args), countRange, field.Unique)
		}
	case "dict":
		if err = expectsArgs(1, 1, assertValStr, fieldType, args); err == nil {
			return entity.WithField(field.Name, fieldType, i.parseArgsForField(fieldType, args), countRange, field.Unique)
		}
	case "date":
		if err = expectsArgs(2, 3, assertDateFieldArgs, fieldType, args); err == nil {
			return entity.WithField(field.Name, fieldType, i.parseArgsForField(fieldType, args), countRange, field.Unique)
		}
	case "bool":
		if err = expectsArgs(0, 0, nil, fieldType, args); err == nil {
			return entity.WithField(field.Name, fieldType, nil, countRange, field.Unique)
		}
	case "enum":
		if err = expectsArgs(1, 1, assertCollection, fieldType, args); err == nil {
			return entity.WithField(field.Name, fieldType, i.parseArgsForField(fieldType, args), countRange, field.Unique)
		} else {
			return field.Err("Expected a collection, but got %v", args[0])
		}
	case "serial": // in the future, consider 1 arg for starting point for sequence
		if err = expectsArgs(0, 0, nil, fieldType, args); err == nil {
			return entity.WithField(field.Name, fieldType, nil, countRange, false)
		}
	case "uid":
		if err = expectsArgs(0, 0, nil, fieldType, args); err == nil {
			return entity.WithField(field.Name, fieldType, nil, countRange, false)
		}
	case "identifier", "entity":
		if nested, e := i.expectsEntity(fieldVal, scope, deferred); e != nil {
			return fieldVal.WrapErr(e)
		} else {
			if err = expectsArgs(0, 0, nil, "entity", args); err == nil {
				return entity.WithEntityField(field.Name, nested, nil, countRange)
			}
		}
	}
	return fieldVal.WrapErr(err)
}

func (i *Interpreter) expectsRange(rangeNode *Node, scope *Scope) (*CountRange, error) {
	switch rangeNode.Kind {
	case "range":
		return i.RangeFromNode(rangeNode, scope)
	case "literal-int":
		return &CountRange{Min: rangeNode.ValInt(), Max: rangeNode.ValInt()}, nil
	case "identifier":
		if v, e := i.ResolveIdentifier(rangeNode, scope); e != nil {
			return nil, e
		} else {
			switch v.(type) {
			case int64:
				return &CountRange{Min: v.(int64), Max: v.(int64)}, nil
			case *CountRange:
				cr, _ := v.(*CountRange)
				return cr, nil
			}
		}
	}

	return nil, rangeNode.Err("Expected a range")
}

func (i *Interpreter) expectsEntity(entityRef *Node, scope *Scope, deferred bool) (*generator.Generator, error) {
	switch entityRef.Kind {
	case "identifier":
		return i.ResolveEntity(entityRef, scope)
	case "entity":
		return i.EntityFromNode(entityRef, scope, deferred)
	default:
		if x, e := i.Visit(entityRef, scope, deferred); e != nil {
			return nil, e
		} else {
			if g, ok := x.(*generator.Generator); ok {
				return g, nil
			} else {
				return nil, entityRef.Err("Expected an entity, but got %v", x)
			}
		}
	}
}

func (i *Interpreter) expectsPrimaryKeyStatement(pkNode *Node, scope *Scope, deferred bool) (*generator.PrimaryKey, error) {
	if !pkNode.Is("primary-key") {
		return nil, pkNode.Err("Expected a primary key statement, but got %q", pkNode.Kind)
	}

	if res, err := i.Visit(pkNode, scope, deferred); err != nil {
		return nil, err
	} else {
		if pk, ok := res.(*generator.PrimaryKey); ok {
			return pk, nil
		} else {
			return nil, pkNode.Err("Expected a primary key specification, but got %v", res)
		}
	}
}

func (i *Interpreter) expectsInteger(intNode *Node, scope *Scope, deferred bool) (int64, error) {
	if result, err := i.Visit(intNode, scope, deferred); err != nil {
		return 0, err
	} else {
		if val, ok := result.(int64); ok {
			return val, nil
		} else {
			return 0, intNode.Err("Expected an integer, but got %v", result)
		}
	}
}

func (i *Interpreter) ResolvePrimaryKeyConfig(scope *Scope) *generator.PrimaryKey {
	if r := scope.ResolveSymbol(PK_FIELD_CONFIG); r == nil {
		return generator.DEFAULT_PK_CONFIG
	} else {
		pk, _ := r.(*generator.PrimaryKey)
		return pk
	}
}

/*
 * A convenience wrapper for ResolveIdentifier, which casts to *generator.Generator. Currently, this
 * is the only type of value that is in the symbol table, but we may support other types in the future
 */
func (i *Interpreter) ResolveEntity(identifierNode *Node, scope *Scope) (*generator.Generator, error) {
	if resolved, err := i.ResolveIdentifier(identifierNode, scope); err != nil {
		return nil, err
	} else {
		if entity, ok := resolved.(*generator.Generator); ok {
			return entity, nil
		} else {
			return nil, identifierNode.Err("identifier %q should refer to an entity, but instead was %v", identifierNode.ValStr(), resolved)
		}
	}
}

func (i *Interpreter) ResolveIdentifier(identiferNode *Node, scope *Scope) (interface{}, error) {
	if scope == nil {
		return nil, identiferNode.Err("Scope is missing! This should be impossible.")
	}

	if !identiferNode.Is("identifier") {
		return nil, identiferNode.Err("Expected an identifier, but got %s", identiferNode.Kind)
	}

	symbol := identiferNode.ValStr()

	if s := scope.DefinedInScope(symbol); nil != s {
		return s.ResolveSymbol(symbol), nil
	}

	return nil, identiferNode.Err("Cannot resolve symbol %q", symbol)
}

func (i *Interpreter) GenerateFromNode(generationNode *Node, scope *Scope, deferred bool) (interface{}, error) {
	if i.dryRun {
		return []interface{}{}, nil
	}

	var entityGenerator *generator.Generator

	entity := generationNode.Args[1]
	if g, e := i.expectsEntity(entity, scope, deferred); e != nil {
		return nil, e

	} else {
		entityGenerator = g
	}

	count, err := i.expectsInteger(generationNode.Args[0], scope, deferred)
	if err != nil {
		return nil, err
	}

	if count < int64(1) {
		return nil, generationNode.Err("Must generate at least 1 %v entity", entityGenerator)
	}

	if err := entityGenerator.EnsureGeneratable(count); err != nil {
		return nil, generationNode.Err(err.Error())
	}

	return entityGenerator.Generate(count, i.emitter.NextEmitter(i.emitter.Receiver(), entityGenerator.Type(), true), scope), nil
}
