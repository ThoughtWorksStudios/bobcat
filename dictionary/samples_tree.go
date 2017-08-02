package dictionary

//go:generate go get github.com/mjibson/esc
//go:generate esc -o data.go -pkg fake data
type samplesTree map[string]map[string][]string

func (st samplesTree) hasKeyPath(lang, cat string) bool {
	if _, ok := st[lang]; ok {
		if _, ok = st[lang][cat]; ok {
			return true
		}
	}
	return false
}
