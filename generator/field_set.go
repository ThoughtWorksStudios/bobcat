package generator

type FieldSet struct {
	fields map[string]*FieldEntry
	names  []string
}

type FieldEntry struct {
	Name  string
	Field *Field
}

func NewFieldSet() *FieldSet {
	return &FieldSet{names: make([]string, 0), fields: make(map[string]*FieldEntry)}
}

func (f *FieldSet) AddField(name string, field *Field) {
	if !f.HasField(name) {
		f.names = append(f.names, name)
	}

	f.fields[name] = &FieldEntry{name, field}
}

func (f *FieldSet) HasField(name string) bool {
	_, ok := f.fields[name]
	return ok
}

func (f *FieldSet) GetField(name string) *Field {
	if !f.HasField(name) {
		return nil
	}
	return f.fields[name].Field
}

func (f *FieldSet) AllFields() []*FieldEntry {
	fields := make([]*FieldEntry, len(f.names))
	for i, name := range f.names {
		fields[i] = f.fields[name]
	}
	return fields
}

func (f *FieldSet) Names() []string {
	return f.names
}
