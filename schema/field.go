package schema

type FieldDefinition struct {
	Name       []byte
	Arguments  []*ArgumentDefinition
	Type       *FieldType
	Directives Directives
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

func (f *FieldDefinition) IsDeprecated() bool {
	for _, directive := range f.Directives {
		if string(directive.Name) == "deprecated" {
			return true
		}
	}
	return false
}

func (f *FieldDefinition) DeprecatedReason() string {
	for _, directive := range f.Directives {
		if string(directive.Name) == "deprecated" {
			if len(directive.Arguments) > 0 {
				for _, arg := range directive.Arguments {
					if string(arg.Name) == "reason" {
						return string(arg.Value)
					}
				}
			}
			return "No reason provided"
		}
	}

	return ""
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
