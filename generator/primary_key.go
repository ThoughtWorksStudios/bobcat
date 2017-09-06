package generator

type PrimaryKey struct {
	name string
	kind string
}

func (pk *PrimaryKey) Field() *Field {
	switch pk.kind {
	case "serial":
		return &Field{fieldType: &SerialType{}}
	default:
		return &Field{fieldType: &MongoIDType{}}
	}
}

func (pk *PrimaryKey) Inherit(target *Generator, source *Generator) {
	target.fields[pk.name] = &Field{fieldType: &ReferenceType{referred: source, fieldName: source.PrimaryKeyName()}}
	target.pkey = pk
}

func (pk *PrimaryKey) Attach(target *Generator) {
	target.fields[pk.name] = pk.Field()
	target.pkey = pk
}

func NewPrimaryKeyConfig(name, kind string) *PrimaryKey {
	return &PrimaryKey{name: name, kind: kind}
}
