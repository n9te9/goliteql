package schema

type InputDefinition struct {
	Name       []byte
	Fields     FieldDefinitions
	Extentions []*InputDefinition
}

func (i *InputDefinition) Location() *Location {
	return &Location{
		Name: []byte("INPUT_OBJECT"),
	}
}

func (i *InputDefinition) IsDefinition() bool {
	return true
}
