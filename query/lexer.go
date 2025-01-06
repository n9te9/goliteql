package query

import (
	"errors"
	"fmt"
	"slices"
	"unicode"
)

type Type string

func (t Type) IsOperation() bool {
	return t == Query || t == Mutation || t == Subscription
}

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
	Dollar      Type = "DALLAR"       // '$'
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

func newBlockStringValueToken(input []byte, cur, col, line int) (*Token, int, int, int, error) {
	start := cur
	cur += 3

	tokenStartLine := line
	for cur + 2 < len(input) {
		if input[cur] == '"' && input[cur + 1] == '"' && input[cur + 2] == '"' {
			break
		}
		cur++

		if input[cur] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}

	if cur + 2 > len(input) {
		return nil, -1, -1, -1, fmt.Errorf("unterminated string at line %d, column %d", tokenStartLine, col)
	}
	cur += 3

	return &Token{Type: Value, Value: input[start:cur], Column: col, Line: tokenStartLine}, cur, line, col, nil
}

func newStringValueToken(input []byte, cur, col, line int) (*Token, int, int, int, error) {
	if cur + 3 < len(input) && input[cur] == '"' && input[cur + 1] == '"' && input[cur + 2] == '"' {
		return newBlockStringValueToken(input, cur, col, line)
	}

	start := cur
	cur++
	escape := false
	for (cur < len(input) && input[cur] != '"') || escape {
		if input[cur] == '\\' {
			escape = true
		} else {
			escape = false
		}
		cur++
		col++
	}

	if cur >= len(input) {
		return nil, -1, -1, -1, fmt.Errorf("unterminated string at line %d, column %d", line, col)
	}

	return &Token{Type: Value, Value: input[start:cur+1], Column: col, Line: line}, cur + 1, line, col, nil
}


func newValueToken(input []byte, cur, col, line int) (*Token, int, int, int) {
	start := cur
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) || unicode.IsDigit(rune(input[cur])) || input[cur] == '.' {
		cur++
		col++
	}

	if tokenType, ok := queryKeywords[string(input[start:cur])]; ok {
		return &Token{Type: tokenType, Value: input[start:cur], Column: col, Line: line}, cur, line, col
	}
	return &Token{Type: Name, Value: input[start:cur], Column: col, Line: line}, cur, line, col
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
	'$': Dollar,
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

type Types []Type

func (t Types) isArgument() bool {
	if len(t) == 0 {
		return false
	}

	if slices.Contains(t, ParenOpen) {
		return true
	}

	return false
}

func (l *Lexer) Lex(input []byte) (Tokens, error) {
	tokens := make(Tokens, 0)
	cur := 0
	col, line := 1, 1

	var token, prev *Token
	var err error
	stack := make(Types, 0)
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
		case '{', '(', '[':
			stack = append(stack, queryPunctuators[input[cur]])
			tokens = append(tokens, &Token{Type: queryPunctuators[input[cur]], Value: []byte{input[cur]}, Column: col, Line: line})
			cur++
			col++
			continue
		case '}', ')', ']':
			if len(stack) == 0 {
				return nil, errors.New("invalid token")
			}

			if stack[len(stack) - 1] == CurlyOpen && input[cur] != '}' {
				return nil, errors.New("invalid token")
			}

			if stack[len(stack) - 1] == ParenOpen && input[cur] != ')' {
				return nil, errors.New("invalid token")
			}

			if stack[len(stack) - 1] == BracketOpen && input[cur] != ']' {
				return nil, errors.New("invalid token")
			}

			stack = stack[:len(stack) - 1]
			tokens = append(tokens, &Token{Type: queryPunctuators[input[cur]], Value: []byte{input[cur]}, Column: col, Line: line})
			cur++
			col++
			continue
		}

		switch input[cur] {
		case ':', '@', ',', '=', '$', '!':
			tokens = append(tokens, &Token{Type: queryPunctuators[input[cur]], Value: []byte{input[cur]}, Column: col, Line: line})
			cur++
			col++
			continue
		}

		if tokens.isDefaultValue() || tokens.isArgument() || stack.isArgument() {
			if unicode.IsLetter(rune(input[cur])) || unicode.IsDigit(rune(input[cur])) {
				token, cur, line, col = newValueToken(input, cur, col, line)
				tokens = append(tokens, token)
				col += len(token.Value)
				continue
			}

			if input[cur] == '"' {
				token, cur, line, col, err = newStringValueToken(input, cur, col, line)
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

		if token == prev {
			return nil, fmt.Errorf("invalid token at line %d, column %d", line, col)
		}
		prev = token
	}

	tokens = append(tokens, newEOFToken(col, line))
	return tokens, nil
}
