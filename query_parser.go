package ggparser

type Document struct {
	Tokens []*QueryToken
}

func ParseQuery(input []byte) (*Document, error) {
	tokens, err := lexQuery(input)
	if err != nil {
		return nil, err
	}

	return &Document{
		Tokens: tokens,
	}, nil
}