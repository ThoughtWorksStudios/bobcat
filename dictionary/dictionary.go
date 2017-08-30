package dictionary

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"sync"
)

const (
	FORMAT_SEP    = "|"
	FORMAT_SUFFIX = "_format"
	NUMERIC_SIG   = "#"
)

// NOTE: this package is a fork of sorts of https://github.com/icrowley/fake
var samplesLock sync.Mutex
var samplesCache = make(samplesTree)
var lang = "en"
var useExternalData = false
var enFallback = true
var availLangs = GetLangs()
var customDataLocation = ""

func ValueFromDictionary(cat string) string {
	s := tryLookup(cat)
	if s == "" {
		s = formatLookup(lang, cat, true)
	}
	return s
}

func tryLookup(cat string) string {
	useExternalData = true
	s := lookup(lang, cat, true)
	useExternalData = false
	if s == "" {
		s = lookup(lang, cat, true)
	}
	return s
}

func formatLookup(lang, cat string, fallback bool) string {
	format := tryLookup(cat + FORMAT_SUFFIX)
	return valueFromFormat(format)
}

//TODO: optimize this formats processing because it's slow
func valueFromFormat(format string) string {
	var result string
	for _, ref := range strings.Split(format, FORMAT_SEP) {
		if strings.Contains(ref, NUMERIC_SIG) {
			result += numericFormat(ref)
		} else if ref == " " {
			result += " "
		} else {
			result += compositeFormat(ref)
		}
	}
	return result
}

func compositeFormat(ref string) string {
	var result string
	r := tryLookup(ref)
	if r == "" {
		result += string(ref)
	} else if strings.HasSuffix(ref, FORMAT_SUFFIX) {
		result += valueFromFormat(r)
	} else {
		result += string(r)
	}
	return result
}

func numericFormat(format string) string {
	var result string
	for _, ru := range format {
		if ru == '#' {
			result += strconv.Itoa(r.Intn(10))
		} else {
			result += string(ru)
		}
	}
	return result
}

func lookup(lang, cat string, fallback bool) string {
	samplesLock.Lock()
	s := _lookup(lang, cat, fallback)
	samplesLock.Unlock()
	return s
}

func _lookup(lang, cat string, fallback bool) string {
	var samples []string

	if samplesCache.hasKeyPath(lang, cat) {
		samples = samplesCache[lang][cat]
	} else {
		var err error
		samples, err = populateSamples(lang, cat)
		if err != nil {
			if lang != "en" && fallback && enFallback && err.Error() == ErrNoSamplesFn(lang).Error() {
				return _lookup("en", cat, false)
			}
			return ""
		}
	}
	return samples[r.Intn(len(samples))]
}

func NumberOfPossibleValuesForDictionary(cat string) int64 {
	result, err := numberOfNewLinesForDictionary(lang, cat)
	if err != nil {
		result, _ = numberPossibilitiesPerNewLinesForFormat(lang, cat)
	}
	return result
}

func numberOfNewLinesForDictionary(lang, cat string) (int64, error) {
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

	var result int64 = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result++
	}
	return result, nil
}

func numberPossibilitiesPerNewLinesForFormat(lang, cat string) (int64, error) {
	useExternalData = true
	fullpath := fullPath(lang, cat+FORMAT_SUFFIX)
	file, err := FS(useExternalData).Open(fullpath)
	useExternalData = false
	if err != nil {
		fullpath = fullPath(lang, cat+FORMAT_SUFFIX)
		file, err = FS(useExternalData).Open(fullpath)
		if err != nil {
			return 0, err
		}
	}
	defer file.Close()

	var result int64 = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineResult := numberOfPossibilitiesForFormat(scanner.Text())
		if lineResult == -1 {
			return -1, nil
		}
		result += lineResult

	}
	return result, nil
}

func numberOfPossibilitiesForFormat(format string) int64 {
	var result int64 = 1
	for _, ref := range strings.Split(format, FORMAT_SEP) {
		var subPossibilities int64 = 0
		if strings.Contains(ref, NUMERIC_SIG) {
			subPossibilities = numberOfNumericPossibilities(ref)
		} else if ref != " " {
			subPossibilities = NumberOfPossibleValuesForDictionary(ref)
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

func numberOfNumericPossibilities(format string) int64 {
	slots := strings.Count(format, NUMERIC_SIG)
	return int64(math.Pow(9.0, float64(slots)))
}

func populateSamples(lang, cat string) ([]string, error) {
	data, err := readFile(lang, cat)
	if err != nil {
		return nil, err
	}

	if _, ok := samplesCache[lang]; !ok {
		samplesCache[lang] = make(map[string][]string)
	}

	samples := strings.Split(strings.TrimSpace(string(data)), "\n")

	samplesCache[lang][cat] = samples
	return samples, nil
}

func readFile(lang, cat string) ([]byte, error) {
	fullpath := fullPath(lang, cat)
	file, err := FS(useExternalData).Open(fullpath)
	if err != nil {
		return nil, ErrNoSamplesFn(lang)
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

func fullPath(lang, cat string) string {
	fullpath := fmt.Sprintf("/data/%s/%s", lang, cat)
	if useExternalData {
		if customDataLocation == "" {
			fullpath = cat
		} else {
			fullpath = fmt.Sprintf("%s/%s", customDataLocation, cat)
		}
	}
	return fullpath
}

func SetCustomDataLocation(location string) {
	customDataLocation = location
}

func UseExternalData(flag bool) {
	useExternalData = flag
}

func EnFallback(flag bool) {
	enFallback = flag
}

func GetLangs() []string {
	var langs []string
	for k, v := range data {
		if v.isDir && k != "/" && k != "/data" {
			langs = append(langs, strings.Replace(k, "/data/", "", 1))
		}
	}
	return langs
}

func SetLang(newLang string) error {
	found := false
	for _, l := range availLangs {
		if newLang == l {
			found = true
			break
		}
	}
	if !found {
		return ErrNoLanguageFn(newLang)
	}
	lang = newLang
	return nil
}
