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
	Query Type = "QUERY"
	Mutation Type = "MUTATION"
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
	Daller      Type = "DALLAR"       // '$'
	Exclamation Type = "EXCLAMATION" // '!'
	Spread 		Type = "SPREAD"      // '...'
	On 			Type = "ON"
	Fragment Type = "FRAGMENT"
	Value Type = "VALUE"
)

var queryKeywords = map[string]Type{
	"query": Query,
	"mutation": Mutation,
	"subscription": Subscription,
	"on": On,
	"fragment": Fragment,
}

func newNameToken(input []byte, cur, col, line int) (*Token, int) {
	start := cur
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) || unicode.IsDigit(rune(input[cur])) {
		cur++
	}

	if tokenType, ok := queryKeywords[string(input[start:cur])]; ok {
		return &Token{Type: tokenType, Value: input[start:cur], Column: col, Line: line}, cur
	}
	return &Token{Type: Name, Value: input[start:cur], Column: col, Line: line}, cur
}

func newStringValueToken(input []byte, cur, col, line int) (*Token, int, error) {
	start := cur
	cur++
	for cur < len(input) && input[cur] != '"' {
		cur++
	}

	if cur >= len(input) {
		return nil, -1, errors.New("unterminated string")
	}

	return &Token{Type: Value, Value: input[start:cur+1], Column: col, Line: line}, cur + 1, nil
}


func newValueToken(input []byte, cur, col, line int) (*Token, int) {
	start := cur
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) || unicode.IsDigit(rune(input[cur])) {
		cur++
	}

	if tokenType, ok := queryKeywords[string(input[start:cur])]; ok {
		return &Token{Type: tokenType, Value: input[start:cur], Column: col, Line: line}, cur
	}
	return &Token{Type: Name, Value: input[start:cur], Column: col, Line: line}, cur
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
	'$': Daller,
	'!': Exclamation,
}

type Lexer struct {
}

func NewLexer() *Lexer {
	return &Lexer{}
}

type Tokens []*Token

func (t Tokens) isDefaultValue() bool {
	if len(t) == 0 {
		return false
	}

	if t[len(t) - 1].Type == Equal {
		return true
	}

	return false
}

func (t Tokens) isArgument() bool {
	if len(t) == 0 {
		return false
	}

	if t[len(t) - 1].Type == Colon {
		return true
	}

	return false
}

func (l *Lexer) Lex(input []byte) (Tokens, error) {
	tokens := make(Tokens, 0)
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
		case '{', '}', '(', ')', ':', '@', ',', '=', '[', ']', '$', '!':
			tokens = append(tokens, &Token{Type: queryPunctuators[input[cur]], Value: []byte{input[cur]}, Column: col, Line: line})
			cur++
			col++
			continue
		}

		if tokens.isDefaultValue() || tokens.isArgument() {
			if unicode.IsLetter(rune(input[cur])) || unicode.IsDigit(rune(input[cur])) {
				token, cur = newValueToken(input, cur, col, line)
				tokens = append(tokens, token)
				col += len(token.Value)
				continue
			}

			if input[cur] == '"' {
				token, cur, err = newStringValueToken(input, cur, col, line)
				if err != nil {
					return nil, err
				}

				tokens = append(tokens, token)
				col += len(token.Value) + 2
				continue
			}
		}

		if unicode.IsLetter(rune(input[cur])) || input[cur] == '_' || unicode.IsDigit(rune(input[cur])) {
			token, cur = newNameToken(input, cur, col, line)
			tokens = append(tokens, token)
			col += len(token.Value)
			continue
		}

		if input[cur] == '.' {
			if input[cur + 1] == '.' && input[cur + 2] == '.' {
				tokens = append(tokens, &Token{Type: Spread, Value: []byte("..."), Column: col, Line: line})
				cur += 3
				col += 3
				continue
			} else {
				return nil, errors.New("invalid token")
			}
		}
	}

	tokens = append(tokens, newEOFToken(col, line))
	return tokens, nil
}
