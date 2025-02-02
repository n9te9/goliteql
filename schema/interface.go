package schema

import "bytes"

type InterfaceDefinition struct {
	Name []byte
	Fields FieldDefinitions
	Extentions []*InterfaceDefinition
	Directives []*Directive
}

func (i *InterfaceDefinition) Location () *Location {
	return &Location{
		Name: []byte("INTERFACE"),
	}
}

func (i *InterfaceDefinition) GetFieldByName(name []byte) *FieldDefinition {
	for _, field := range i.Fields {
		if bytes.Equal(field.Name, name) {
			return field
		}
	}

	return nil
}

func (i *InterfaceDefinition) TypeName() []byte {
	return i.Name
}