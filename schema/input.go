package schema

type InputDefinition struct {
	Name       []byte
	Fields     FieldDefinitions
	tokens     Tokens
	Extentions []*InputDefinition
}

func (i *InputDefinition) Location() *Location {
	return &Location{
		Name: []byte("INPUT_OBJECT"),
	}
}
