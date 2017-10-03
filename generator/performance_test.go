package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/builtins"
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/emitter"
	"testing"
	"time"
)

func addBuiltin(g *Generator, name, builtinName string, args ...interface{}) {
	builtin, _ := NewBuiltin(builtinName)
	g.WithField(name, NewDeferredType(func(_ *Scope) (interface{}, error) {
		return builtin.Call(args...)
	}), nil)
}

func addBuiltinWithCount(g *Generator, name, builtinName string, count *CountRange, args ...interface{}) {
	builtin, _ := NewBuiltin(builtinName)
	g.WithField(name, NewDeferredType(func(_ *Scope) (interface{}, error) {
		return builtin.Call(args...)
	}), count)
}

func addLiteral(g *Generator, name string, value interface{}) {
	g.WithField(name, NewLiteralType(value), nil)
}

func setup(b *testing.B) *Generator {
	g := NewGenerator("thing", nil, false)
	addBuiltin(g, "name", STRING_TYPE, int64(10))
	addBuiltin(g, "age", FLOAT_TYPE, float64(2), float64(4))
	addLiteral(g, "species", "human")
	return g
}

func resetTimerAndGenerateX(b *testing.B, g *Generator, x int64) {
	b.ResetTimer()
	g.Generate(x, NewDummyEmitter(), NewRootScope())
}

func Benchmark_Generate_OneThousand(b *testing.B) {
	resetTimerAndGenerateX(b, setup(b), 1000)
}

func Benchmark_Generate_TenThousand(b *testing.B) {
	resetTimerAndGenerateX(b, setup(b), 10000)
}

func Benchmark_Generate_OneHundredThousand(b *testing.B) {
	resetTimerAndGenerateX(b, setup(b), 100000)
}

func Benchmark_Generate_FiveHundredThousand(b *testing.B) {
	resetTimerAndGenerateX(b, setup(b), 500000)
}

func Benchmark_Generate_OneMillion(b *testing.B) {
	resetTimerAndGenerateX(b, setup(b), 1000000)
}

func Benchmark_Generate_OneThousandWithEntityField(b *testing.B) {
	g := setup(b)
	generator := NewGenerator("Person", nil, false)
	addBuiltin(g, "name", STRING_TYPE, int64(10))
	generator.WithEntityField("pet", g, nil)

	resetTimerAndGenerateX(b, generator, 1000)
}

func Benchmark_Generate_OneHundredThousandWithEntityField(b *testing.B) {
	g := setup(b)
	generator := NewGenerator("Person", nil, false)
	addBuiltin(g, "name", STRING_TYPE, int64(10))
	generator.WithEntityField("pet", g, nil)

	resetTimerAndGenerateX(b, generator, 100000)
}

func Benchmark_Generate_FiveHundredThousandWithEntityField(b *testing.B) {
	g := setup(b)
	generator := NewGenerator("Person", nil, false)
	addBuiltin(g, "name", STRING_TYPE, int64(10))
	generator.WithEntityField("pet", g, nil)

	resetTimerAndGenerateX(b, generator, 500000)
}

func Benchmark_Generate_OneMillionWithEntityField(b *testing.B) {
	g := setup(b)
	generator := NewGenerator("Person", nil, false)
	addBuiltin(g, "name", STRING_TYPE, int64(10))
	generator.WithEntityField("pet", g, nil)

	resetTimerAndGenerateX(b, generator, 1000000)
}

func Benchmark_Generate_OneThousandWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", nil, false)
	addBuiltin(g, "name", STRING_TYPE, int64(10))
	generator.WithEntityField("pet", g, nil)
	generator.WithEntityField("vet", g2, nil)

	resetTimerAndGenerateX(b, generator, 1000)
}

func Benchmark_Generate_OneHundredThousandWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", nil, false)
	addBuiltin(g, "name", STRING_TYPE, int64(10))
	generator.WithEntityField("pet", g, nil)
	generator.WithEntityField("vet", g2, nil)

	resetTimerAndGenerateX(b, generator, 100000)
}

func Benchmark_Generate_FiveHundredThousandWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", nil, false)
	addBuiltin(g, "name", STRING_TYPE, int64(10))
	generator.WithEntityField("pet", g, nil)
	generator.WithEntityField("vet", g2, nil)

	resetTimerAndGenerateX(b, generator, 500000)
}

func Benchmark_Generate_OneMillionWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", nil, false)
	addBuiltin(g, "name", STRING_TYPE, int64(10))
	generator.WithEntityField("pet", g, nil)
	generator.WithEntityField("vet", g2, nil)

	resetTimerAndGenerateX(b, generator, 1000000)
}

func Benchmark_Field_GenerateValue_For_OneMillion_Integers(b *testing.B) {
	builtin, _ := NewBuiltin(INT_TYPE)
	f := &Field{fieldType: NewDeferredType(func(_ *Scope) (interface{}, error) {
		return builtin.Call(int64(1), int64(100))
	}), count: &CountRange{Min: 1000000, Max: 1000000}}
	f.GenerateValue(nil, NewDummyEmitter(), NewRootScope())
}

func Benchmark_Field_GenerateValue_For_OneMillion_Floats(b *testing.B) {
	builtin, _ := NewBuiltin(FLOAT_TYPE)
	f := &Field{fieldType: NewDeferredType(func(_ *Scope) (interface{}, error) {
		return builtin.Call(float64(1), float64(100))
	}), count: &CountRange{Min: 1000000, Max: 1000000}}
	f.GenerateValue(nil, NewDummyEmitter(), NewRootScope())
}

func Benchmark_Field_GenerateValue_For_OneMillion_Literals(b *testing.B) {
	f := &Field{fieldType: NewLiteralType("blah"), count: &CountRange{Min: 1000000, Max: 1000000}}
	f.GenerateValue(nil, NewDummyEmitter(), NewRootScope())
}

func Benchmark_Field_GenerateValue_For_OneMillion_Bools(b *testing.B) {
	builtin, _ := NewBuiltin(BOOL_TYPE)
	f := &Field{fieldType: NewDeferredType(func(_ *Scope) (interface{}, error) {
		return builtin.Call()
	}), count: &CountRange{Min: 1000000, Max: 1000000}}
	f.GenerateValue(nil, NewDummyEmitter(), NewRootScope())
}

func Benchmark_Field_GenerateValue_For_OneMillion_Dates(b *testing.B) {
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	builtin, _ := NewBuiltin(DATE_TYPE)
	f := &Field{fieldType: NewDeferredType(func(_ *Scope) (interface{}, error) {
		return builtin.Call(timeMin, timeMax, "")
	}), count: &CountRange{Min: 1000000, Max: 1000000}}

	b.ResetTimer()
	f.GenerateValue(nil, NewDummyEmitter(), NewRootScope())
}

func Benchmark_Field_GenerateValue_For_OneMillion_MongoIDs(b *testing.B) {
	builtin, _ := NewBuiltin(UID_TYPE)
	f := &Field{fieldType: NewDeferredType(func(_ *Scope) (interface{}, error) {
		return builtin.Call()
	}), count: &CountRange{Min: 1000000, Max: 1000000}}

	f.GenerateValue(nil, NewDummyEmitter(), NewRootScope())
}

func Benchmark_Field_GenerateValue_For_OneMillion_Strings(b *testing.B) {
	builtin, _ := NewBuiltin(STRING_TYPE)
	f := &Field{fieldType: NewDeferredType(func(_ *Scope) (interface{}, error) {
		return builtin.Call(int64(100))
	}), count: &CountRange{Min: 1000000, Max: 1000000}}

	f.GenerateValue(nil, NewDummyEmitter(), NewRootScope())
}
