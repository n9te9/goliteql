package query

import (
	"errors"
	"unicode"
)

type Type string

type QueryToken struct {
	Type Type
	Value []byte
	Line int
	Column int
}

const (
	Name Type = "NAME"
	Int Type = "INT"
	String Type = "STRING"
	Query Type = "QUERY"
	Mutate Type = "MUTATION"
	Subscription Type = "SUBSCRIPTION"
	EOF Type = "EOF"

	CurlyOpen    Type = "CURLY_OPEN"    // '{'
	CurlyClose   Type = "CURLY_CLOSE"   // '}'
	ParenOpen    Type = "PAREN_OPEN"    // '('
	ParenClose   Type = "PAREN_CLOSE"   // ')'
	Colon        Type = "COLON"         // ':'
	At           Type = "AT"           // '@'
	Comma        Type = "COMMA"        // ','
	Equal        Type = "EQUAL"        // '='
	BracketOpen  Type = "BRACKET_OPEN" // '['
	BracketClose Type = "BRACKET_CLOSE" // ']'
)

var queryKeywords = map[string]Type{
	"query": Query,
	"mutation": Mutate,
	"subscription": Subscription,
}

func newFieldQueryToken(input []byte, cur, col, line int) (*QueryToken, int) {
	start := cur
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) || unicode.IsDigit(rune(input[cur])) {
		cur++
	}

	if tokenType, ok := queryKeywords[string(input[start:cur])]; ok {
		return &QueryToken{Type: tokenType, Value: input[start:cur], Column: col, Line: line}, cur
	}
	return &QueryToken{Type: Name, Value: input[start:cur], Column: col, Line: line}, cur
}

func newIntQueryToken(input []byte, cur, col, line int) (*QueryToken, int) {
	start := cur
	for cur < len(input) && unicode.IsDigit(rune(input[cur])) {
		cur++
	}

	return &QueryToken{Type: Int, Value: input[start:cur], Column: col, Line: line}, cur
}

func newStringQueryToken(input []byte, cur, col, line int) (*QueryToken, int, error) {
	start := cur + 1
	cur++
	for cur < len(input) && input[cur] != '"' {
		cur++
	}

	if cur >= len(input) {
		return nil, -1, errors.New("unterminated string")
	}

	return &QueryToken{Type: String, Value: input[start:cur], Column: col, Line: line}, cur + 1, nil
}

func newEOFQueryToken(col, line int) *QueryToken {
	return &QueryToken{Type: EOF, Value: nil, Column: col, Line: line}
}

var queryPunctuators = map[byte]Type{
	'{': CurlyOpen,
	'}': CurlyClose,
	'(': ParenOpen,
	')': ParenClose,
	':': Colon,
	'@': At,
	',': Comma,
	'=': Equal,
	'[': BracketOpen,
	']': BracketClose,
}

type Lexer struct {
}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (l *Lexer) Lex(input []byte) ([]*QueryToken, error) {
	tokens := make([]*QueryToken, 0)
	cur := 0
	col, line := 1, 1

	var token *QueryToken
	var err error
	for cur < len(input) {
		switch input[cur] {
		case ' ', '\t':
			col++
			cur++
			continue
		case '\n':
			line++
			col = 1
			cur++
			continue
		}

		switch input[cur] {
		case '{', '}', '(', ')', ':', '@', ',', '=':
			tokens = append(tokens, &QueryToken{Type: queryPunctuators[input[cur]], Value: []byte{input[cur]}, Column: col, Line: line})
			cur++
			col++
			continue
		}

		if unicode.IsLetter(rune(input[cur])) {
			token, cur = newFieldQueryToken(input, cur, col, line)
			tokens = append(tokens, token)
			col += len(token.Value)
			continue
		}

		if unicode.IsDigit(rune(input[cur])) {
			token, cur = newIntQueryToken(input, cur, col, line)
			tokens = append(tokens, token)
			col += len(token.Value)
			continue
		}

		if input[cur] == '"' {
			token, cur, err = newStringQueryToken(input, cur, col, line)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
			col += len(token.Value) + 2
			continue
		}
	}

	tokens = append(tokens, newEOFQueryToken(col, line))
	return tokens, nil
}
