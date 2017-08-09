package dictionary

import (
	"testing"
)

func resetCache(b *testing.B) {
	samplesCache = make(samplesTree)
	b.ResetTimer()
}

func Benchmark_Simple_ValueFromDictionary(b *testing.B) {
	resetCache(b)
	ValueFromDictionary("first_names")
}

func Benchmark_Simple_ValueFromDictionary_OneThousand_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000; i++ {
		ValueFromDictionary("first_names")
	}
}

func Benchmark_Simple_ValueFromDictionary_OneHundredThousand_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 100000; i++ {
		ValueFromDictionary("first_names")
	}
}

func Benchmark_Simple_ValueFromDictionary_OneMillion_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000000; i++ {
		ValueFromDictionary("first_names")
	}
}

func Benchmark_NumericFormat_ValueFromDictionary(b *testing.B) {
	resetCache(b)
	ValueFromDictionary("phone_numbers")
}

func Benchmark_NumericFormat_ValueFromDictionary_OneThousand_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000; i++ {
		ValueFromDictionary("phone_numbers")
	}
}

func Benchmark_NumericFormat_ValueFromDictionary_OneHundredThousand_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 100000; i++ {
		ValueFromDictionary("phone_numbers")
	}
}

func Benchmark_NumericFormat_ValueFromDictionary_OneMillion_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000000; i++ {
		ValueFromDictionary("phone_numbers")
	}
}

func Benchmark_CompositeFormat_ValueFromDictionary(b *testing.B) {
	resetCache(b)
	ValueFromDictionary("full_names")
}

func Benchmark_CompositeFormat_ValueFromDictionary_OneThousand_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000; i++ {
		ValueFromDictionary("full_names")
	}
}

func Benchmark_CompositeFormat_ValueFromDictionary_OneHundredThousand_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 100000; i++ {
		ValueFromDictionary("full_names")
	}
}

func Benchmark_CompositeFormat_ValueFromDictionary_OneMillion_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000000; i++ {
		ValueFromDictionary("full_names")
	}
}

func Benchmark_CustomDict_ValueFromDictionary(b *testing.B) {
	resetCache(b)
	ValueFromDictionary("testdata/custom")
}

func Benchmark_CustomDict_ValueFromDictionary_OneThousand_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000; i++ {
		ValueFromDictionary("testdata/custom")
	}
}

func Benchmark_CustomDict_ValueFromDictionary_OneHundredThousand_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 100000; i++ {
		ValueFromDictionary("testdata/custom")
	}
}

func Benchmark_CustomDict_ValueFromDictionary_OneMillion_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000000; i++ {
		ValueFromDictionary("testdata/custom")
	}
}

func Benchmark_CustomCompositeDict_ValueFromDictionary(b *testing.B) {
	resetCache(b)
	ValueFromDictionary("testdata/custom_composite")
}

func Benchmark_CustomCompositeDict_ValueFromDictionary_OneThousand_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000; i++ {
		ValueFromDictionary("testdata/custom_composite")
	}
}

func Benchmark_CustomCompositeDict_ValueFromDictionary_OneHundredThousand_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 100000; i++ {
		ValueFromDictionary("testdata/custom_composite")
	}
}

func Benchmark_CustomCompositeDict_ValueFromDictionary_OneMillion_Times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000000; i++ {
		ValueFromDictionary("testdata/custom_composite")
	}
}

func Benchmark_valueFromFormat_NumericFormat(b *testing.B) {
	valueFromFormat("####")
}

func Benchmark_valueFromFormat_NumericFormat_OneThousand_times(b *testing.B) {
	for i := 1; i <= 1000; i++ {
		valueFromFormat("####")
	}
}

func Benchmark_valueFromFormat_NumericFormat_OneHundredThousand_times(b *testing.B) {
	for i := 1; i <= 100000; i++ {
		valueFromFormat("####")
	}
}

func Benchmark_valueFromFormat_NumericFormat_OneMillion_times(b *testing.B) {
	for i := 1; i <= 1000000; i++ {
		valueFromFormat("####")
	}
}

func Benchmark_valueFromFormat_CompositeFormat(b *testing.B) {
	resetCache(b)
	valueFromFormat("first_names| |last_names")
}

func Benchmark_valueFromFormat_CompositeFormat_OneThousand_times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000; i++ {
		valueFromFormat("first_names| |last_names")
	}
}

func Benchmark_valueFromFormat_CompositeFormat_OneHundredThousand_times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 100000; i++ {
		valueFromFormat("first_names| |last_names")
	}
}

func Benchmark_valueFromFormat_CompositeFormat_OneMillion_times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000000; i++ {
		valueFromFormat("first_names| |last_names")
	}
}

func Benchmark_valueFromFormat_CompositeNumericFormat(b *testing.B) {
	resetCache(b)
	valueFromFormat("first_names| |last_names| |####")
}

func Benchmark_valueFromFormat_CompositeNumericFormat_OneThousand_times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000; i++ {
		valueFromFormat("first_names| |last_names| |####")
	}
}

func Benchmark_valueFromFormat_CompositeNumericFormat_OneHundredThousand_times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 100000; i++ {
		valueFromFormat("first_names| |last_names| |####")
	}
}

func Benchmark_valueFromFormat_CompositeNumericFormat_OneMillion_times(b *testing.B) {
	resetCache(b)
	for i := 1; i <= 1000000; i++ {
		valueFromFormat("first_names| |last_names| |####")
	}
}
