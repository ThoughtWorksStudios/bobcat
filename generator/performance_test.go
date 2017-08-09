package generator

import (
	"testing"
)

func setup(b *testing.B) *Generator {
	g := NewGenerator("thing", nil)
	g.WithField("name", "string", 10, nil)
	g.WithField("age", "decimal", [2]float64{2, 4}, nil)
	g.WithStaticField("species", "human")
	return g
}

func resetTimerAndGenerateX(b *testing.B, g *Generator, x int64) {
	b.ResetTimer()
	g.Generate(x)
}

func BenchmarkGenerateOneThousand(b *testing.B) {
	resetTimerAndGenerateX(b, setup(b), 1000)
}

func BenchmarkGenerateTenThousand(b *testing.B) {
	resetTimerAndGenerateX(b, setup(b), 10000)
}

func BenchmarkGenerateOneHundredThousand(b *testing.B) {
	resetTimerAndGenerateX(b, setup(b), 100000)
}

func BenchmarkGenerateFiveHundredThousand(b *testing.B) {
	resetTimerAndGenerateX(b, setup(b), 500000)
}

func BenchmarkGenerateOneMillion(b *testing.B) {
	resetTimerAndGenerateX(b, setup(b), 1000000)
}

func BenchmarkGenerateOneThousandWithEntityField(b *testing.B) {
	g := setup(b)
	generator := NewGenerator("Person", nil)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)

	resetTimerAndGenerateX(b, generator, 1000)
}

func BenchmarkGenerateOneHundredThousandWithEntityField(b *testing.B) {
	g := setup(b)
	generator := NewGenerator("Person", nil)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)

	resetTimerAndGenerateX(b, generator, 100000)
}

func BenchmarkGenerateFiveHundredThousandWithEntityField(b *testing.B) {
	g := setup(b)
	generator := NewGenerator("Person", nil)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)

	resetTimerAndGenerateX(b, generator, 500000)
}

func BenchmarkGenerateOneMillionWithEntityField(b *testing.B) {
	g := setup(b)
	generator := NewGenerator("Person", nil)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)

	resetTimerAndGenerateX(b, generator, 1000000)
}

func BenchmarkGenerateOneThousandWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", nil)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)
	generator.WithEntityField("vet", g2, 1, nil)

	resetTimerAndGenerateX(b, generator, 1000)
}

func BenchmarkGenerateOneHundredThousandWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", nil)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)
	generator.WithEntityField("vet", g2, 1, nil)

	resetTimerAndGenerateX(b, generator, 100000)
}

func BenchmarkGenerateFiveHundredThousandWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", nil)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)
	generator.WithEntityField("vet", g2, 1, nil)

	resetTimerAndGenerateX(b, generator, 500000)
}

func BenchmarkGenerateOneMillionWithTwoEntityFields(b *testing.B) {
	g := setup(b)
	g2 := setup(b)
	generator := NewGenerator("Person", nil)
	generator.WithField("name", "string", 10, nil)
	generator.WithEntityField("pet", g, 1, nil)
	generator.WithEntityField("vet", g2, 1, nil)

	resetTimerAndGenerateX(b, generator, 1000000)
}
