package query

type Document struct {
	Tokens []*Token
}

type Parser struct {
	Lexer *Lexer
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		Lexer: lexer,
	}
}

func (p *Parser) Parse(input []byte) (*Document, error) {
	tokens, err := p.Lexer.Lex(input)
	if err != nil {
		return nil, err
	}

	return &Document{
		Tokens: tokens,
	}, nil
}
