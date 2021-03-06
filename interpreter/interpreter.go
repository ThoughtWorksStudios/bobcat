package interpreter

import (
	"fmt"
	. "github.com/ThoughtWorksStudios/bobcat/builtins"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	"github.com/ThoughtWorksStudios/bobcat/dsl"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"github.com/ThoughtWorksStudios/bobcat/generator"
	"os"
	"strings"
	"time"
)

const (
	PK_FIELD_CONFIG = "$PK_FIELD"
)

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
	CustomDictPath = path
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

func unwrapAtomic(node *Node) *Node {
	for "atomic" == node.Kind {
		node = node.ValNode()
	}
	return node
}

func unwrapSequential(node *Node) *Node {
	for 1 == len(node.Children) && node.Children[0].Kind == "sequential" {
		node.Children = node.Children[0].Children
	}
	return node
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
	case "builtin":
		if builtin, err := NewBuiltin(node.Name); err != nil {
			return nil, node.WrapErr(err)
		} else {
			if deferred {
				return func(scope *Scope) (interface{}, error) {
					return builtin, nil
				}, nil
			}

			return builtin, nil
		}

	case "lambda":
		node = unwrapSequential(node)

		var fn DeferredResolver
		var err error

		if fn, err = i.Compile(node.Children, scope); err != nil {
			return nil, node.WrapErr(err)
		}

		boundArgs := make([]string, len(node.Args))

		for idx, arg := range node.Args {
			boundArgs[idx] = arg.ValStr()
		}

		symbol := node.Name

		closure := func(scope *Scope) (interface{}, error) {
			lambda := NewLambda(node.Name, boundArgs, fn, scope)
			if symbol != "" {
				// TODO: refactor, DRY?
				if s := scope.DefinedInScope(symbol); s == scope {
					Warn("%v Symbol %q has already been declared in this scope", node.Ref, symbol)
				}

				scope.SetSymbol(symbol, lambda)
			}
			return lambda, nil
		}

		if deferred {
			return closure, nil
		}

		return closure(scope)
	case "call":
		callableNode := unwrapAtomic(node.ValNode())

		closure := func(scope *Scope) (interface{}, error) {
			if result, err := i.Visit(callableNode, scope, false); err == nil {
				if callable, ok := result.(Callable); ok {
					return i.BindAndInvokeCallable(callable, node, scope)
				} else {
					return nil, callableNode.Err("Expected a lambda, but got %v (%T)", result, result)
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
		return i.Visit(unwrapAtomic(node.ValNode()), scope, deferred)
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
			if entry, err := i.ResolveIdentifierFromNode(node, scope); err == nil {
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
		symbol := node.Name
		valNode := node.ValNode()

		var precompiledValue interface{}
		if v, err := i.Visit(valNode, scope, deferred); err == nil {
			precompiledValue = v // TREAT THIS AS IMMUTABLE HEREAFTER!
		} else {
			return nil, err
		}

		closure := func(scope *Scope) (interface{}, error) {
			var result interface{}
			if s := scope.DefinedInScope(symbol); s != nil {
				if resolver, ok := precompiledValue.(DeferredResolver); ok {
					if val, err := resolver(scope); err != nil {
						return nil, err
					} else {
						result = val
					}
				} else {
					result = precompiledValue
				}
				/**
				 * must set in the scope where symbol is defined, which is not
				 * necessarily the current scope. the ability to assign a value
				 * to a symbol in a parent scope is intentional. if you instead
				 * want variable shadowing, use a variable declaration in the
				 * present scope, NOT an assignment expression.
				 */
				s.SetSymbol(symbol, result)

				return result, nil
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

		var precompiledValue interface{}

		if nil != node.Value {
			valNode := node.ValNode()
			if v, err := i.Visit(valNode, scope, deferred); err == nil {
				precompiledValue = v // TREAT THIS AS IMMUTABLE HEREAFTER!
			} else {
				return nil, err
			}
		}

		closure := func(scope *Scope) (interface{}, error) {
			if s := scope.DefinedInScope(symbol); s == scope {
				Warn("%v Symbol %q has already been declared in this scope", node.Ref, symbol)
			}

			var result interface{} = nil

			if nil != precompiledValue {
				if resolver, ok := precompiledValue.(DeferredResolver); ok {
					if val, err := resolver(scope); err != nil {
						return nil, err
					} else {
						result = val
					}
				} else {
					result = precompiledValue
				}
			}

			scope.SetSymbol(symbol, result)
			return result, nil
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
				kind := node.Related.Name
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
	Msg("[%s, %T(%v), %T(%v)]", op, left, left, right, right)

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
		case *TimeWithFormat:
			return i.addToTime(op, left.(*TimeWithFormat), right, scope, deferred)
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

	return NewExecQueue(queue).Run, nil
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
		if parent, e := i.ResolveEntityFromNode(node.Related, scope); nil == e {

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
			if !field.Is("field") && !field.Is(DIST_TYPE) {
				return nil, field.Err("Expected a `field` declaration, but instead got `%s`", field.Kind) // should never get here
			}

			fieldVal := field.ValNode()

			fieldVal = unwrapAtomic(fieldVal)
			field.Value = fieldVal

			countRange, err := i.CountRangeFromNode(field.CountRange, scope, false)

			if err != nil {
				return nil, err
			}

			args, err := i.FieldArgumentsFromNodeSet(field.Args, scope)

			if err != nil {
				return nil, err
			}

			switch fieldVal.Kind {
			case DIST_TYPE:
				if err := i.withDistributionField(entity, field, scope, deferred); err != nil {
					return nil, field.WrapErr(err)
				}
			case "identifier":
				symbol := fieldVal.ValStr()

				if entity.HasField(symbol) {
					closure := func(scope *Scope) (interface{}, error) {
						if value, err := i.ResolveIdentifier(symbol, scope); err == nil {
							return value, nil
						} else {
							return nil, fieldVal.WrapErr(err)
						}
					}

					if err := i.withExpressionField(entity, field.Name, closure); err != nil {
						return nil, fieldVal.WrapErr(err)
					}
					continue
				}

				value, err := i.ResolveIdentifier(symbol, scope)

				if err != nil {
					return nil, fieldVal.WrapErr(err)
				}

				switch value.(type) {
				case Callable:
					return nil, fieldVal.Err("Field values cannot be lambda types; you probably meant to call the lambda to generate the field value")
				case *generator.Generator:
					nested := value.(*generator.Generator)

					if len(args) != 0 {
						return nil, fieldVal.Err("entity field types do not take arguments")
					}

					if err = entity.WithEntityField(field.Name, nested, countRange); err != nil {
						return nil, fieldVal.WrapErr(err)
					}
				default:
					if len(args) != 0 {
						return nil, fieldVal.Err("value field types do not take arguments")
					}

					if err = i.withExpressionField(entity, field.Name, value); err != nil {
						return nil, fieldVal.WrapErr(err)
					}
				}
			case "entity":
				if nested, err := i.expectsEntity(fieldVal, scope, false); err != nil {
					return nil, fieldVal.WrapErr(err)
				} else {
					if len(args) != 0 {
						return nil, fieldVal.Err("entity field types do not take arguments")
					}

					if err = entity.WithEntityField(field.Name, nested, countRange); err != nil {
						return nil, fieldVal.WrapErr(err)
					}
				}
			case "binary":
				var value interface{}
				var err error

				if value, err = i.resolveBinaryNode(fieldVal, scope, true); err != nil {
					return nil, fieldVal.WrapErr(err)
				}

				if err = i.withExpressionField(entity, field.Name, value); err != nil {
					return nil, fieldVal.WrapErr(err)
				}
			case "lambda", "builtin":
				return nil, fieldVal.Err("Field values cannot be lambda types; you probably meant to call the lambda to generate the field value")
			case "call":
				callableNode := unwrapAtomic(fieldVal.ValNode())

				callableResolver, err := i.Visit(callableNode, scope, true)

				if err != nil {
					return nil, callableNode.WrapErr(err)
				}

				resolver, ok := callableResolver.(DeferredResolver)

				if !ok {
					return nil, callableNode.Err("Expected a lambda, but got %v", callableResolver)
				}

				closure := func(scope *Scope) (interface{}, error) {

					if result, err := resolver(scope); err == nil {
						if callable, ok := result.(Callable); ok {
							return i.BindAndInvokeCallable(callable, fieldVal, scope)
						} else {
							return nil, callableNode.Err("Expected a lambda, but got %v", result)
						}
					} else {
						return nil, callableNode.WrapErr(err)
					}
				}

				if err := i.withExpressionField(entity, field.Name, closure); err != nil {
					return nil, fieldVal.WrapErr(err)
				}
			default:
				if len(args) != 0 {
					return nil, fieldVal.Err("value field types do not take arguments")
				}

				if value, err := i.Visit(fieldVal, scope, false); err == nil {
					if err = i.withExpressionField(entity, field.Name, value); err != nil {
						return nil, fieldVal.WrapErr(err)
					}
				} else {
					return nil, fieldVal.WrapErr(err)
				}
			}
		}
	}

	return entity, nil
}

func (i *Interpreter) BindAndInvokeCallable(callable Callable, callNode *Node, scope *Scope) (val interface{}, err error) {
	if val, err = i.AllValuesFromNodeSet(callNode.Args, scope, false); err == nil {
		if val, err = callable.Call(val.([]interface{})...); err == nil {
			return
		}
	}
	return nil, callNode.WrapErr(err)
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
	case INT_TYPE:
		return [2]int64{args[0].(int64), args[1].(int64)}
	case FLOAT_TYPE:
		return [2]float64{args[0].(float64), args[1].(float64)}
	case STRING_TYPE:
		return args[0].(int64)
	case DICT_TYPE:
		return args[0].(string)
	case DATE_TYPE:
		format := ""
		if 3 == len(args) {
			format = args[2].(string)
		}
		return []interface{}{args[0].(time.Time), args[1].(time.Time), format}
	case ENUM_TYPE:
		return args[0].([]interface{})
	default:
		return nil
	}
}

func (i *Interpreter) withDistributionField(entity *generator.Generator, field *Node, scope *Scope, deferred bool) error {
	distNode := field.ValNode()
	distFn := distNode.Name // normal, percent, etc

	numArgs := len(distNode.Args)

	if 0 == numArgs {
		return distNode.Err("distributions require a domain")
	}

	if distFn == NORMAL_DIST && 2 != len(distNode.Args) {
		return distNode.Err("%v requires exactly 2 arguments to define bounds", distFn)
	}

	weights := make([]float64, numArgs)
	fields := make([]generator.FieldType, numArgs)

	for idx, argNode := range distNode.Args {
		var valNode *Node

		if argNode.Is("associative-arg") {
			if distFn == NORMAL_DIST {
				return argNode.Err("%v distributions do not take weights", distFn)
			}

			if wt, err := i.Visit(argNode.ValNode(), scope, false); err != nil { // immediately resolve; cannot use variables declared in transient scope
				return argNode.Related.WrapErr(err)
			} else {
				switch wt.(type) {
				case int64:
					weights[idx] = float64(wt.(int64))
				case float64:
					weights[idx] = wt.(float64)
				default:
					return argNode.Related.Err("weights must be numeric values; instead, got %v", wt)
				}
			}

			valNode = unwrapAtomic(argNode.Related)
		} else {
			if distFn != NORMAL_DIST {
				return argNode.Err("%v distributions must have weights")
			}
			valNode = unwrapAtomic(argNode)
		}

		switch valNode.Kind {
		case "lambda", "builtin":
			return valNode.Err("cannot distribute over lambdas")
		default:
			if value, err := i.Visit(valNode, scope, distFn != NORMAL_DIST); err != nil {
				return valNode.WrapErr(err)
			} else {
				switch value.(type) {
				case DeferredResolver: // impossible to reach when NORMAL_DIST
					fields[idx] = generator.NewDeferredType(value.(DeferredResolver))
				case int64:
					if distFn == NORMAL_DIST {
						weights[idx] = float64(value.(int64))
					} else {
						fields[idx] = generator.NewLiteralType(value)
					}
				case float64:
					if distFn == NORMAL_DIST {
						weights[idx] = value.(float64)
					} else {
						fields[idx] = generator.NewLiteralType(value)
					}
				default:
					if distFn == NORMAL_DIST {
						return valNode.Err("%v boundaries must be numeric", distFn)
					}
					fields[idx] = generator.NewLiteralType(value)
				}
			}
		}
	}

	if distribution, err := generator.NewDistribution(distFn, weights, fields); err == nil {
		entity.WithField(field.Name, distribution, nil)
	} else {

		return field.WrapErr(err)
	}

	return nil
}

func (i *Interpreter) CountRangeFromNode(node *Node, scope *Scope, deferred bool) (*CountRange, error) {
	if nil != node {
		if countRange, err := i.expectsRange(node, scope); err == nil {
			if err = countRange.Validate(); err != nil {
				return nil, node.WrapErr(err)
			}
			return countRange, nil
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) FieldArgumentsFromNodeSet(argNodes NodeSet, scope *Scope) ([]interface{}, error) {
	if 0 == len(argNodes) {
		return make([]interface{}, 0), nil
	}

	if a, e := i.AllValuesFromNodeSet(argNodes, scope, false); e == nil {
		args, _ := a.([]interface{})
		return args, nil
	} else {
		return nil, e
	}
}

func (i *Interpreter) expectsRange(rangeNode *Node, scope *Scope) (*CountRange, error) {
	switch rangeNode.Kind {
	case "range":
		return i.RangeFromNode(rangeNode, scope)
	case "literal-int":
		return &CountRange{Min: rangeNode.ValInt(), Max: rangeNode.ValInt()}, nil
	case "identifier":
		if v, e := i.ResolveIdentifierFromNode(rangeNode, scope); e != nil {
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
		return i.ResolveEntityFromNode(entityRef, scope)
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
func (i *Interpreter) ResolveEntityFromNode(identifierNode *Node, scope *Scope) (*generator.Generator, error) {
	if resolved, err := i.ResolveIdentifierFromNode(identifierNode, scope); err != nil {
		return nil, err
	} else {
		if entity, ok := resolved.(*generator.Generator); ok {
			return entity, nil
		} else {
			return nil, identifierNode.Err("identifier %q should refer to an entity, but instead was %v", identifierNode.ValStr(), resolved)
		}
	}
}

func (i *Interpreter) ResolveIdentifierFromNode(identiferNode *Node, scope *Scope) (interface{}, error) {
	if scope == nil {
		return nil, identiferNode.Err("Scope is missing! This should be impossible.")
	}

	if !identiferNode.Is("identifier") {
		return nil, identiferNode.Err("Expected an identifier, but got %s", identiferNode.Kind)
	}

	symbol := identiferNode.ValStr()

	if result, err := i.ResolveIdentifier(symbol, scope); err == nil {
		return result, nil
	} else {
		return nil, identiferNode.WrapErr(err)
	}
}

func (i *Interpreter) ResolveIdentifier(symbol string, scope *Scope) (interface{}, error) {
	if s := scope.DefinedInScope(symbol); nil != s {
		return s.ResolveSymbol(symbol), nil
	}

	return nil, fmt.Errorf("Cannot resolve symbol %q", symbol)
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

	return entityGenerator.Generate(count, i.emitter.NextEmitter(i.emitter.Receiver(), entityGenerator.Type(), true), scope)
}
