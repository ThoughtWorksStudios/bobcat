package dictionary

// NOTE: this package is a fork of sorts of https://github.com/icrowley/fake

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:generate go get github.com/mjibson/esc
//go:generate esc -o data.go -pkg fake data

// cat/subcat/lang/samples
type samplesTree map[string]map[string][]string

var samplesLock sync.Mutex
var samplesCache = make(samplesTree)
var r = rand.New(&rndSrc{src: rand.NewSource(time.Now().UnixNano())})
var lang = "en"
var useExternalData = false
var enFallback = true
var availLangs = GetLangs()
var customDataLocation = ""

var (
	// ErrNoLanguageFn is the error that indicates that given language is not available
	ErrNoLanguageFn = func(lang string) error { return fmt.Errorf("The language passed (%s) is not available", lang) }
	// ErrNoSamplesFn is the error that indicates that there are no samples for the given language
	ErrNoSamplesFn = func(lang string) error { return fmt.Errorf("No samples found for language: %s", lang) }
)

// Seed uses the provided seed value to initialize the internal PRNG to a
// deterministic state.
func Seed(seed int64) {
	r.Seed(seed)
}

type rndSrc struct {
	mtx sync.Mutex
	src rand.Source
}

func (s *rndSrc) Int63() int64 {
	s.mtx.Lock()
	n := s.src.Int63()
	s.mtx.Unlock()
	return n
}

func (s *rndSrc) Seed(n int64) {
	s.mtx.Lock()
	s.src.Seed(n)
	s.mtx.Unlock()
}

// GetLangs returns a slice of available languages
func GetLangs() []string {
	var langs []string
	for k, v := range data {
		if v.isDir && k != "/" && k != "/data" {
			langs = append(langs, strings.Replace(k, "/data/", "", 1))
		}
	}
	return langs
}

// SetLang sets the language in which the data should be generated
// returns error if passed language is not available
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
func join(parts ...string) string {
	var filtered []string
	for _, part := range parts {
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	return strings.Join(filtered, " ")
}

func compositeFormat(format string) string {
	var compositeResult string
	for _, ref := range strings.Split(format, "|") {
		r := tryLookup(ref)
		if r == "" {
			compositeResult += string(ref)
		} else {
			if strings.HasSuffix(ref, "_format") {
				compositeResult += valueFromFormat(r)
			} else {
				compositeResult += r
			}
		}
	}
	return compositeResult
}

func valueFromFormat(format string) string {
	var result string
	for _, ru := range compositeFormat(format) {
		if ru != '#' {
			result += string(ru)
		} else {
			result += strconv.Itoa(r.Intn(10))
		}
	}

	return result
}

func formatLookup(lang, cat string, fallback bool) string {
	format := tryLookup(cat + "_format")
	return valueFromFormat(format)
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

func SetCustomDataLocation(location string) {
	customDataLocation = location
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

func readFile(lang, cat string) ([]byte, error) {
	fullpath := fullPath(lang, cat)
	file, err := FS(useExternalData).Open(fullpath)
	if err != nil {
		return nil, ErrNoSamplesFn(lang)
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

func UseExternalData(flag bool) {
	useExternalData = flag
}

func EnFallback(flag bool) {
	enFallback = flag
}

func (st samplesTree) hasKeyPath(lang, cat string) bool {
	if _, ok := st[lang]; ok {
		if _, ok = st[lang][cat]; ok {
			return true
		}
	}
	return false
}
