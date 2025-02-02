package schema


type EnumDefinition struct {
	Name []byte
	Values []*EnumElement
	Extentions []*EnumDefinition
	Directives []*Directive
}

func (e *EnumDefinition) Location() *Location {
	return &Location{
		Name: []byte("ENUM"),
	}
}

type EnumDefinitions []*EnumDefinition

func (e EnumDefinitions) Has(name string) bool {
	for _, enum := range e {
		if string(enum.Name) == name {
			return true
		}
	}

	return false
}

type EnumElement struct {
	Name []byte
	Value []byte
	Directives []*Directive
}

func (e *EnumElement) Location() *Location {
	return &Location{
		Name: []byte("ENUM_VALUE"),
	}
}