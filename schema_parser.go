package ggparser

type Schema struct {
	Tokens []*SchemaToken
}

func ParseSchema(input []byte) (*Schema, error) {
	tokens, err := lexSchema(input)
	if err != nil {
		return nil, err
	}

	return &Schema{
		Tokens: tokens,
	}, nil
}