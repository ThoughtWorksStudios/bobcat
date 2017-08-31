package dictionary

import (
	"bufio"
	"math"
	"strings"
)

type calculatorFunc func(scanner *bufio.Scanner) (int64, error)

func CalculatePossibilities(cat string) int64 {
	result, err := calculatePossibilities(lang, cat, dictionaryCalculator)
	if err != nil {
		result, _ = calculatePossibilities(lang, cat+FORMAT_SUFFIX, formatCalculator)
	}
	return result
}

func calculatePossibilities(lang, cat string, calc calculatorFunc) (int64, error) {
	useExternalData = true
	fullpath := fullPath(lang, cat)
	file, err := FS(useExternalData).Open(fullpath)
	useExternalData = false
	if err != nil {
		fullpath = fullPath(lang, cat)
		file, err = FS(useExternalData).Open(fullpath)
		if err != nil {
			return 0, err
		}
	}
	defer file.Close()

	return calc(bufio.NewScanner(file))
}

func dictionaryCalculator(scanner *bufio.Scanner) (int64, error) {
	var result int64 = 0
	for scanner.Scan() {
		result++
	}
	return result, nil
}

func formatCalculator(scanner *bufio.Scanner) (int64, error) {
	var result int64 = 0
	for scanner.Scan() {
		lineResult := calculateFormatPossibilities(scanner.Text())
		if lineResult == -1 {
			return -1, nil
		}
		result += lineResult

	}
	return result, nil
}

func calculateFormatPossibilities(format string) int64 {
	var result int64 = 1
	for _, ref := range strings.Split(format, FORMAT_SEP) {
		var subPossibilities int64 = 0
		if strings.Contains(ref, NUMERIC_SIG) {
			subPossibilities = calculateNumericPossibilities(ref)
		} else if ref != " " {
			subPossibilities = CalculatePossibilities(ref)
		}
		if subPossibilities != 0 {
			result *= subPossibilities
		}
		if result <= 0 {
			//The number of possibilities exceeds the maximum value that int64
			// can contain. So, we treat it as practically infinite
			return -1
		}
	}
	return result
}

func calculateNumericPossibilities(format string) int64 {
	slots := strings.Count(format, NUMERIC_SIG)
	return int64(math.Pow(9.0, float64(slots)))
}
