package generator

type GeneratedEntities []EntityResult

type EntityResult map[string]GeneratedValue

func NewGeneratedEntities(count int64) GeneratedEntities {
	return make([]EntityResult, count)
}

func (ge GeneratedEntities) Concat(newEntities GeneratedEntities) GeneratedEntities {
	for _, entity := range newEntities {
		ge = append(ge, entity)
	}
	return ge
}

type GeneratedValue interface{}

type GeneratedStringValue string
type GeneratedIntegerValue int
type GeneratedFloatValue float64
type GeneratedListValue []GeneratedValue
type GeneratedBoolValue bool
type GeneratedEntityValue EntityResult
