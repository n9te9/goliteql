package schema

import (
	"unicode"
)

type Type string

const (
	ReservedType Type = "TYPE"
	Identifier Type = "IDENTIFIER"
	Field Type = "FIELD"

	Extend Type = "EXTEND"
	Scalar Type = "SCALAR"
	Enum Type = "ENUM"
	Input Type = "INPUT"
	Interface Type = "INTERFACE"
	Union Type = "UNION"

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
	Exclamation Type = "EXCLAMATION" // '!'
)

type Token struct {
	Type Type
	Value []byte
	Line int
	Column int
}

func newKeywordToken(input []byte, t Type, cur, col, line int) (*Token, int) {
	start := cur
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) {
		cur++
	}

	return &Token{Type: t, Value: input[start:cur], Column: col, Line: line}, cur
}

func newPunctuatorToken(input []byte, t Type, cur, col, line int) (*Token, int) {
	return &Token{Type: t, Value: []byte{input[cur]}, Column: col, Line: line}, cur + 1
}

func newExclamationToken(input []byte, cur, col, line int) (*Token, int) {
	return &Token{Type: Exclamation, Value: []byte{input[cur]}, Column: col, Line: line}, cur + 1
}

func newFieldToken(input []byte, cur, col, line int) (*Token, int) {
	start := cur
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) || unicode.IsDigit(rune(input[cur])) {
		cur++
	}

	return &Token{Type: Field, Value: input[start:cur], Column: col, Line: line}, cur
}

func newIdentifierToken(input []byte, cur, col, line int) (*Token, int) {
	start := cur
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) {
		cur++
	}

	return &Token{Type: Identifier, Value: input[start:cur], Column: col, Line: line}, cur
}

type Tokens []*Token

func (t Tokens) isType() bool {
	if len(t) == 0 {
		return false
	}

	lastToken := t[len(t)-1]
	return lastToken.Type == ReservedType || 
	lastToken.Type == Colon || 
	lastToken.Type == BracketOpen ||
	lastToken.Type == At
}

func (t Tokens) isInput() bool {
	if len(t) == 0 {
		return false
	}

	lastToken := t[len(t)-1]
	return lastToken.Type == Input
}

func (t Tokens) isInterface() bool {
	if len(t) == 0 {
		return false
	}

	lastToken := t[len(t)-1]
	return lastToken.Type == Interface
}

type Lexer struct {}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (l *Lexer) Lex(input []byte) ([]*Token, error) {
	tokens := make(Tokens, 0)
	cur := 0
	col, line := 1, 1

	var token *Token
	for cur < len(input) {
		switch input[cur] {
		case ' ', '\t':
			cur++
			col++
			continue
		case '\n':
			cur++
			line++
			col = 1
			continue
		}

		switch input[cur] {
		case '{', '}', '(', ')', ':', '@', ',', '=', '[', ']':
			if t, ok := punctuators[punctuator(input[cur])]; ok {
				token, cur = newPunctuatorToken(input, t, cur, col, line)
				tokens = append(tokens, token)
				col++
			}
			continue
		case '!':
			token, cur = newExclamationToken(input, cur, col, line)
			tokens = append(tokens, token)
			col++
			continue
		}

		if unicode.IsLetter(rune(input[cur])) {
			end := keywordEnd(input, cur)
			keyword := keyword(input[cur:end])
			if t, ok := keywords[keyword]; ok {
				token, cur = newKeywordToken(input, t, cur, col, line)
				tokens = append(tokens, token)
				col += len(token.Value)
				continue
			}

			if tokens.isType() || 
			tokens.isInput() ||
			tokens.isInterface() {
				token, cur = newIdentifierToken(input, cur, col, line)
				tokens = append(tokens, token)
				col += len(token.Value)
				continue
			}

			token, cur = newFieldToken(input, cur, col, line)
			tokens = append(tokens, token)
			col += len(token.Value)
		}
	}

	tokens = append(tokens, &Token{Type: EOF, Value: nil, Column: col, Line: line})

	return tokens, nil
}

type keyword string

func (k keyword) String() string {
	return string(k)
}

var keywords = map[keyword]Type{
	"type": ReservedType,
	"extend": Extend,
	"scalar": Scalar,
	"enum": Enum,
	"input": Input,
	"interface": Interface,
	"union": Union,
}

type punctuator byte

var punctuators = map[punctuator]Type{
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

func keywordEnd(input []byte, cur int) int {
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) {
		cur++
	}

	return cur
}