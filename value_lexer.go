package goliteql

import (
	"fmt"
)

type ValueLexer struct {
}

func NewValueLexer() *ValueLexer {
	return &ValueLexer{}
}

type ValueLexerTokenType int

const (
	ILLEGAL ValueLexerTokenType = iota
	EOF
	IDENT
	INT
	FLOAT
	STRING
	ID
	BOOL
	NULL
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	COLON
	COMMA
)

type ValueLexerToken struct {
	Type  ValueLexerTokenType
	Value []byte
}

func (vl *ValueLexer) Lex(input []byte) ([]*ValueLexerToken, error) {
	tokens := make([]*ValueLexerToken, 0)

	for i := 0; i < len(input); {
		switch input[i] {
		case ' ', '\t', '\n', '\r':
			i++ // Skip whitespace
			continue
		case '(':
			tokens = append(tokens, &ValueLexerToken{Type: LPAREN, Value: []byte("(")})
			i++
			continue
		case ')':
			tokens = append(tokens, &ValueLexerToken{Type: RPAREN, Value: []byte(")")})
			i++
			continue
		case '{':
			tokens = append(tokens, &ValueLexerToken{Type: LBRACE, Value: []byte("{")})
			i++
			continue
		case '}':
			tokens = append(tokens, &ValueLexerToken{Type: RBRACE, Value: []byte("}")})
			i++
			continue
		case '[':
			tokens = append(tokens, &ValueLexerToken{Type: LBRACKET, Value: []byte("[")})
			i++
			continue
		case ']':
			tokens = append(tokens, &ValueLexerToken{Type: RBRACKET, Value: []byte("]")})
			i++
			continue
		case ':':
			tokens = append(tokens, &ValueLexerToken{Type: COLON, Value: []byte(":")})
			i++
			continue
		case ',':
			tokens = append(tokens, &ValueLexerToken{Type: COMMA, Value: []byte(",")})
			i++
			continue
		default:
			token, nextIndex := vl.lexToken(input, i)
			if token == ILLEGAL {
				return nil, fmt.Errorf("illegal token at position %d", i)
			}
			tokens = append(tokens, &ValueLexerToken{Type: token, Value: input[i:nextIndex]})
			i = nextIndex
			continue
		}
	}

	tokens = append(tokens, &ValueLexerToken{Type: EOF, Value: nil})

	return tokens, nil
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isMinus(b byte) bool {
	return b == '-'
}

func findDot(value []byte) (int, error) {
	ret := -1
	for i, b := range value {
		if ret != -1 && b == '.' {
			return -1, fmt.Errorf("multiple dots found in number: %s", value)
		}

		if b == '.' {
			ret = i
		}
	}

	return ret, nil
}

func isFloat(value []byte, idx int) bool {
	if idx > 0 && idx < len(value)-1 {
		if isDigit(value[idx-1]) && isDigit(value[idx+1]) {
			return true
		}
	}

	return false
}

func (vl *ValueLexer) lexToken(input []byte, start int) (ValueLexerTokenType, int) {
	i := start
	if isMinus(input[start]) {
		i++
	}

	for i < len(input) && (isLetter(input[i]) || isDigit(input[i]) || input[i] == '_' || isFloat(input, i)) {
		i++
	}

	if i > start {
		value := input[start:i]
		switch string(value) {
		case "true", "false":
			return BOOL, i
		case "null":
			return NULL, i
		default:
			if isDigit(value[0]) || isMinus(value[0]) {
				for i := 1; i < len(value); i++ {
					if !isDigit(value[i]) && value[i] != '.' {
						return ILLEGAL, i
					}
				}

				dotIndex, err := findDot(value)
				if err != nil {
					return ILLEGAL, i
				}

				if dotIndex != -1 {
					return FLOAT, i
				}

				return INT, i
			}
			return IDENT, i
		}
	}

	if start < len(input) && input[start] == '"' {
		end := start + 1
		for end < len(input) && input[end] != '"' {
			end++
		}
		if end < len(input) && input[end] == '"' {
			return STRING, end + 1
		}
	}

	return ILLEGAL, start + 1
}
