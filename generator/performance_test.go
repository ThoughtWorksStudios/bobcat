package generator

import "testing"

func setup(b *testing.B) *Generator {
	g := NewGenerator("thing", nil)
	g.WithField("name", "string", 10)
	g.WithField("age", "decimal", [2]float64{2, 4})
	g.WithStaticField("species", "human")

	b.ResetTimer()
	return g
}

func BenchmarkGenerateOneThousand(b *testing.B) {
	g := setup(b)
	g.Generate(1000)
}

func BenchmarkGenerateTenThousand(b *testing.B) {
	g := setup(b)
	g.Generate(10000)
}

func BenchmarkGenerateOneHundredThousand(b *testing.B) {
	g := setup(b)
	g.Generate(100000)
}

func BenchmarkGenerateFiveHundredThousand(b *testing.B) {
	g := setup(b)
	g.Generate(500000)
}

func BenchmarkGenerateOneMillion(b *testing.B) {
	g := setup(b)
	g.Generate(1000000)
}
