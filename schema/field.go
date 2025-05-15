package schema

type FieldDefinition struct {
	Name       []byte
	Arguments  []*ArgumentDefinition
	Type       *FieldType
	Directives []*Directive
	Default    []byte
	Location   *Location
}

func (f *FieldDefinition) IsPrimitive() bool {
	typeName := string(f.Type.Name)

	if typeName == "ID" || typeName == "String" || typeName == "Int" || typeName == "Float" || typeName == "Boolean" {
		return true
	}

	return false
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
