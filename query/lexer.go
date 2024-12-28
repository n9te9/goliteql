package query

import (
	"errors"
	"unicode"
)

type Type string

type Token struct {
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

func newFieldToken(input []byte, cur, col, line int) (*Token, int) {
	start := cur
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) || unicode.IsDigit(rune(input[cur])) {
		cur++
	}

	if tokenType, ok := queryKeywords[string(input[start:cur])]; ok {
		return &Token{Type: tokenType, Value: input[start:cur], Column: col, Line: line}, cur
	}
	return &Token{Type: Name, Value: input[start:cur], Column: col, Line: line}, cur
}

func newIntToken(input []byte, cur, col, line int) (*Token, int) {
	start := cur
	for cur < len(input) && unicode.IsDigit(rune(input[cur])) {
		cur++
	}

	return &Token{Type: Int, Value: input[start:cur], Column: col, Line: line}, cur
}

func newStringToken(input []byte, cur, col, line int) (*Token, int, error) {
	start := cur + 1
	cur++
	for cur < len(input) && input[cur] != '"' {
		cur++
	}

	if cur >= len(input) {
		return nil, -1, errors.New("unterminated string")
	}

	return &Token{Type: String, Value: input[start:cur], Column: col, Line: line}, cur + 1, nil
}

func newEOFToken(col, line int) *Token {
	return &Token{Type: EOF, Value: nil, Column: col, Line: line}
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

func (l *Lexer) Lex(input []byte) ([]*Token, error) {
	tokens := make([]*Token, 0)
	cur := 0
	col, line := 1, 1

	var token *Token
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
			tokens = append(tokens, &Token{Type: queryPunctuators[input[cur]], Value: []byte{input[cur]}, Column: col, Line: line})
			cur++
			col++
			continue
		}

		if unicode.IsLetter(rune(input[cur])) {
			token, cur = newFieldToken(input, cur, col, line)
			tokens = append(tokens, token)
			col += len(token.Value)
			continue
		}

		if unicode.IsDigit(rune(input[cur])) {
			token, cur = newIntToken(input, cur, col, line)
			tokens = append(tokens, token)
			col += len(token.Value)
			continue
		}

		if input[cur] == '"' {
			token, cur, err = newStringToken(input, cur, col, line)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
			col += len(token.Value) + 2
			continue
		}
	}

	tokens = append(tokens, newEOFToken(col, line))
	return tokens, nil
}
