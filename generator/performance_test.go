package generator

import (
	. "github.com/ThoughtWorksStudios/bobcat/common"
	. "github.com/ThoughtWorksStudios/bobcat/test_helpers"
	"testing"
	"time"
)

func setup(b *testing.B) *Generator {
	g := NewGenerator("thing", false)
	g.WithField("name", "string", 10, nil)
	g.WithField("age", "decimal", [2]float64{2, 4}, nil)
	g.WithStaticField("species", "human")
	return g
}

func resetTimerAndGenerateX(b *testing.B, g *Generator, x int64) {
	b.ResetTimer()
	g.Generate(x, NewTestEmitter())
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
	generator := NewGenerator("Person", false)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)

	resetTimerAndGenerateX(b, generator, 1000)
}

func Benchmark_Generate_OneHundredThousandWithEntityField(b *testing.B) {
	g := setup(b)
	generator := NewGenerator("Person", false)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)

	resetTimerAndGenerateX(b, generator, 100000)
}

func Benchmark_Generate_FiveHundredThousandWithEntityField(b *testing.B) {
	g := setup(b)
	generator := NewGenerator("Person", false)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)

	resetTimerAndGenerateX(b, generator, 500000)
}

func Benchmark_Generate_OneMillionWithEntityField(b *testing.B) {
	g := setup(b)
	generator := NewGenerator("Person", false)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)

	resetTimerAndGenerateX(b, generator, 1000000)
}

func Benchmark_Generate_OneThousandWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", false)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)
	generator.WithEntityField("vet", g2, 1, nil)

	resetTimerAndGenerateX(b, generator, 1000)
}

func Benchmark_Generate_OneHundredThousandWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", false)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)
	generator.WithEntityField("vet", g2, 1, nil)

	resetTimerAndGenerateX(b, generator, 100000)
}

func Benchmark_Generate_FiveHundredThousandWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", false)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)
	generator.WithEntityField("vet", g2, 1, nil)

	resetTimerAndGenerateX(b, generator, 500000)
}

func Benchmark_Generate_OneMillionWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", false)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)
	generator.WithEntityField("vet", g2, 1, nil)

	resetTimerAndGenerateX(b, generator, 1000000)
}

func Benchmark_Field_GenerateValue_For_OneMillion_Integers(b *testing.B) {
	f := &Field{fieldType: &IntegerType{min: 1, max: 100}, count: &CountRange{Min: 1000000, Max: 1000000}}
	f.GenerateValue("", NewTestEmitter())
}

func Benchmark_Field_GenerateValue_For_OneMillion_Floats(b *testing.B) {
	f := &Field{fieldType: &FloatType{min: float64(1), max: float64(100)}, count: &CountRange{Min: 1000000, Max: 1000000}}
	f.GenerateValue("", NewTestEmitter())
}

func Benchmark_Field_GenerateValue_For_OneMillion_Literals(b *testing.B) {
	f := &Field{fieldType: &LiteralType{value: "blah"}, count: &CountRange{Min: 1000000, Max: 1000000}}
	f.GenerateValue("", NewTestEmitter())
}

func Benchmark_Field_GenerateValue_For_OneMillion_Bools(b *testing.B) {
	f := &Field{fieldType: &BoolType{}, count: &CountRange{Min: 1000000, Max: 1000000}}
	f.GenerateValue("", NewTestEmitter())
}

func Benchmark_Field_GenerateValue_For_OneMillion_Dates(b *testing.B) {
	timeMin, _ := time.Parse("2006-01-02", "1945-01-01")
	timeMax, _ := time.Parse("2006-01-02", "1945-01-02")
	f := &Field{fieldType: &DateType{min: timeMin, max: timeMax}, count: &CountRange{Min: 1000000, Max: 1000000}}
	b.ResetTimer()
	f.GenerateValue("", NewTestEmitter())
}

func Benchmark_Field_GenerateValue_For_OneMillion_MongoIDs(b *testing.B) {
	f := &Field{fieldType: &MongoIDType{}, count: &CountRange{Min: 1000000, Max: 1000000}}
	f.GenerateValue("", NewTestEmitter())
}

func Benchmark_Field_GenerateValue_For_OneMillion_Strings(b *testing.B) {
	f := &Field{fieldType: &StringType{length: 100}, count: &CountRange{Min: 1000000, Max: 1000000}}
	f.GenerateValue("", NewTestEmitter())
}
