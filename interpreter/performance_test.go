package interpreter

import "testing"

func BenchmarkLoadFileForOneThousandEntities(b *testing.B) {
	i := New()
	i.LoadFile("testdata/performance/1_thousand.lang", NewRootScope())
}

func BenchmarkLoadFileForTenThousandEntities(b *testing.B) {
	i := New()
	i.LoadFile("testdata/performance/10_thousand.lang", NewRootScope())
}

func BenchmarkLoadFileForFiftyThousandEntities(b *testing.B) {
	i := New()
	i.LoadFile("testdata/performance/50_thousand.lang", NewRootScope())
}

func BenchmarkLoadFileForSeventyThousandEntities(b *testing.B) {
	i := New()
	i.LoadFile("testdata/performance/70_thousand.lang", NewRootScope())
}

func BenchmarkLoadFileForEightyThousandEntities(b *testing.B) {
	i := New()
	i.LoadFile("testdata/performance/80_thousand.lang", NewRootScope())
}

func BenchmarkLoadFileForNinetyThousandEntities(b *testing.B) {
	i := New()
	i.LoadFile("testdata/performance/90_thousand.lang", NewRootScope())
}

func BenchmarkLoadFileForOnehundredThousandEntities(b *testing.B) {
	i := New()
	i.LoadFile("testdata/performance/100_thousand.lang", NewRootScope())
}

func BenchmarkLoadFileForOneMillionEntities(b *testing.B) {
	i := New()
	i.LoadFile("testdata/performance/1_million.lang", NewRootScope())
}
