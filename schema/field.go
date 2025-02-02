package schema

type FieldDefinition struct {
	Name []byte
	Arguments []*ArgumentDefinition
	Type *FieldType
	Directives []*Directive
	Default []byte
}

func (f *FieldDefinition) Location() *Location {
	return &Location{
		Name: []byte("FIELD_DEFINITION"),
	}
}

type FieldDefinitions []*FieldDefinition

func (f FieldDefinitions) Last(name string) *FieldDefinition {
	var res *FieldDefinition
	for _, field := range f {
		if string(field.Name) == name {
			res = field
		}
	}

	return res
}

func (f FieldDefinitions) Has(name string) bool {
	for _, field := range f {
		if string(field.Name) == name {
			return true
		}
	}
	
	return false
}
