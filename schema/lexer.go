package schema

import (
	"fmt"
	"strings"
	"unicode"
)

type Type string

const (
	ReservedType Type = "TYPE"
	ReservedSchema Type = "SCHEMA"
	Identifier   Type = "IDENTIFIER"
	Field        Type = "FIELD"

	Extend    Type = "EXTEND"
	Scalar    Type = "SCALAR"
	Enum      Type = "ENUM"
	Input     Type = "INPUT"
	Interface Type = "INTERFACE"
	Union     Type = "UNION"

	Int     Type = "INT"
	String  Type = "STRING"
	Boolean Type = "BOOLEAN"
	Float   Type = "FLOAT"
	Null    Type = "NULL"

	Value Type = "VALUE"

	Query        Type = "QUERY"
	Mutate       Type = "MUTATION"
	Subscription Type = "SUBSCRIPTION"
	EOF          Type = "EOF"

	CurlyOpen    Type = "CURLY_OPEN"    // '{'
	CurlyClose   Type = "CURLY_CLOSE"   // '}'
	ParenOpen    Type = "PAREN_OPEN"    // '('
	ParenClose   Type = "PAREN_CLOSE"   // ')'
	Colon        Type = "COLON"         // ':'
	At           Type = "AT"            // '@'
	Comma        Type = "COMMA"         // ','
	Equal        Type = "EQUAL"         // '='
	BracketOpen  Type = "BRACKET_OPEN"  // '['
	BracketClose Type = "BRACKET_CLOSE" // ']'
	Exclamation  Type = "EXCLAMATION"   // '!'
	Pipe         Type = "PIPE"          // '|'
)

type Token struct {
	Type   Type
	Value  []byte
	Line   int
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

var (
	queryValue = []byte(`Query`)
	mutateValue = []byte(`Mutation`)
	subscriptionValue = []byte(`Subscription`)
)

func newIdentifierToken(input []byte, cur, col, line int) (*Token, int) {
	start := cur
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) {
		cur++
	}

	if string(input[start:cur]) == string(queryValue) {
		return &Token{Type: Query, Value: queryValue, Column: col, Line: line}, cur
	}

	if string(input[start:cur]) == string(mutateValue) {
		return &Token{Type: Mutate, Value: mutateValue, Column: col, Line: line}, cur
	}

	if string(input[start:cur]) == string(subscriptionValue) {
		return &Token{Type: Subscription, Value: subscriptionValue, Column: col, Line: line}, cur
	}

	return &Token{Type: Identifier, Value: input[start:cur], Column: col, Line: line}, cur
}

func newIntToken(input []byte, cur, col, line int) (*Token, int) {
	start := cur
	for cur < len(input) && unicode.IsDigit(rune(input[cur])) {
		cur++
	}

	return &Token{Type: Int, Value: input[start:cur], Column: col, Line: line}, cur
}

func newFloatToken(input []byte, cur, col, line int) (*Token, int) {
	start := cur
	for cur < len(input) && unicode.IsDigit(rune(input[cur])) || input[cur] == '.' {
		cur++
	}

	return &Token{Type: Float, Value: input[start:cur], Column: col, Line: line}, cur
}

func newStringToken(input []byte, cur, col, line int) (*Token, int) {
	start := cur
	cur++
	for cur < len(input) && input[cur] != '"' {
		cur++
	}

	return &Token{Type: String, Value: input[start : cur+1], Column: col, Line: line}, cur + 1
}

func newBooleanToken(input []byte, cur, col, line int) (*Token, int) {
	start := cur
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) {
		cur++
	}

	return &Token{Type: Boolean, Value: input[start:cur], Column: col, Line: line}, cur
}

func newNullToken(cur, col, line int) (*Token, int) {
	return &Token{Type: Null, Value: []byte("null"), Column: col, Line: line}, cur + 4
}

func newValueToken(input []byte, cur, end, col, line int) (*Token, int) {
	v := string(input[cur:end])
	var token *Token
	// expect int or float
	if unicode.IsDigit(rune(input[cur])) {
		if !strings.Contains(v, ".") {
			// expect int
			token, cur = newIntToken(input, cur, col, line)
			return token, cur
		} else {
			// expect float
			token, cur = newFloatToken(input, cur, col, line)
			return token, cur
		}
	}

	// expect string
	if input[cur] == '"' {
		token, cur = newStringToken(input, cur, col, line)
		return token, cur
	}

	// expect boolean
	if v == "true" || v == "false" {
		token, cur = newBooleanToken(input, cur, col, line)
		return token, cur
	}

	// expect null
	if v == "null" || v == "NULL" || v == "Null" {
		token, cur = newNullToken(cur, col, line)
		return token, cur
	}

	return nil, cur
}

func newListTokens(input []byte, cur, col, line int) (Tokens, int) {
	tokens := Tokens{
		&Token{Type: BracketOpen, Value: []byte{'['}, Column: col, Line: line},
	}

	cur++
	col++
	
	var token *Token
	for cur < len(input) && input[cur] != ']' {
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

		if input[cur] == '[' {
			nestedTokens, newCur := newListTokens(input, cur, col, line)
			tokens = append(tokens, nestedTokens...)
			col += newCur - cur
			cur = newCur
			continue
		}

		if input[cur] == ',' {
			token, cur = newPunctuatorToken(input, Comma, cur, col, line)
			tokens = append(tokens, token)
			col++
			continue
		}

		end := keywordEnd(input, cur)
		token, cur = newValueToken(input, cur, end, col, line)

		tokens = append(tokens, token)
		col += len(token.Value)
	}

	tokens = append(tokens, &Token{Type: BracketClose, Value: []byte{']'}, Column: col, Line: line})
	cur++

	return tokens, cur
}

func newDirectiveValueToken(input []byte, start, cur, col, line int) (*Token, int) {
	return &Token{Type: Value, Value: input[start:cur], Column: col, Line: line}, cur
}

func newDirectiveArgumentTokens(input []byte, cur, col, line int) (Tokens, int) {
	tokens := make(Tokens, 0)

	var token *Token
	for cur < len(input) && input[cur] != ')' {
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

		if input[cur] == ':' {
			token, cur = newPunctuatorToken(input, Colon, cur, col, line)
			tokens = append(tokens, token)
			col++
			continue
		}

		end := directiveKeywordEnd(input, cur)
		if tokens.isField() {
			token, cur = newDirectiveValueToken(input, cur, end, col, line)
			tokens = append(tokens, token)
			col += len(token.Value)
			continue
		}

		if input[cur] == ',' {
			token, cur = newPunctuatorToken(input, Comma, cur, col, line)
			tokens = append(tokens, token)
			col++
			continue
		}

		token, cur = newFieldToken(input, cur, col, line)
		tokens = append(tokens, token)
		col += len(token.Value)
	}
	
	tokens = append(tokens, &Token{Type: ParenClose, Value: []byte{')'}, Column: col, Line: line})
	cur++

	return tokens, cur
}

func newEnumTokens(input []byte, cur, col, line int) (Tokens, int, int, int) {
	tokens := make(Tokens, 0)

	var token *Token
	for cur < len(input) && input[cur] != '}' {
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
		case '{', ',':
			token, cur = newPunctuatorToken(input, punctuators[punctuator(input[cur])], cur, col, line)
			tokens = append(tokens, token)
			col++
			continue
		}

		token, cur = newIdentifierToken(input, cur, col, line)
		tokens = append(tokens, token)
		col += len(token.Value)
	}

	return tokens, cur, line, col
}

func newUnionTokens(input []byte, cur, col, line int) (Tokens, int, int, int) {
	tokens := make(Tokens, 0)

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
		case '|', '=':
			token, cur = newPunctuatorToken(input, punctuators[punctuator(input[cur])], cur, col, line)
			tokens = append(tokens, token)
			col++
			continue
		}

		if token != nil && token.Type != Pipe && token.Type != Equal {
			break
		}

		token, cur = newIdentifierToken(input, cur, col, line)
		tokens = append(tokens, token)
		col += len(token.Value)
		continue
	}

	return tokens, cur, line, col
}

type Tokens []*Token

func (t Tokens) isType() bool {
	if len(t) == 0 {
		return false
	}

	lastToken := t[len(t)-1]
	return lastToken.Type == ReservedType
}

func (t Tokens) isField() bool {
	if len(t) == 0 {
		return false
	}

	lastToken := t[len(t)-1]
	return lastToken.Type == Colon ||
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

func (t Tokens) isDefaultArgument() bool {
	if len(t) == 0 {
		return false
	}

	lastToken := t[len(t)-1]
	return lastToken.Type == Equal
}

func (t Tokens) isDirectiveArgument() bool {
	if len(t) < 3 {
		return false
	}

	lastToken := t[len(t)-1]
	secondLastToken := t[len(t)-2]
	thirdLastToken := t[len(t)-3]
	return thirdLastToken.Type == At &&
		secondLastToken.Type == Identifier &&
		lastToken.Type == ParenOpen
}

func (t Tokens) isEnum() bool {
	if len(t) == 0 {
		return false
	}

	lastToken := t[len(t)-1]
	return lastToken.Type == Enum
}

func (t Tokens) isUnion() bool {
	if len(t) == 0 {
		return false
	}

	lastToken := t[len(t)-1]
	return lastToken.Type == Union
}

type Lexer struct{}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (l *Lexer) Lex(input []byte) ([]*Token, error) {
	tokens := make(Tokens, 0)
	cur := 0
	prev := 0
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

		if tokens.isEnum() {
			newTokens, newCur, newLine, newCol := newEnumTokens(input, cur, col, line)
			tokens = append(tokens, newTokens...)
			line = newLine
			col = newCol
			cur = newCur
			continue
		}

		if tokens.isUnion() {
			newTokens, newCur, newLine, newCol := newUnionTokens(input, cur, col, line)
			tokens = append(tokens, newTokens...)
			line = newLine
			col = newCol
			cur = newCur
			continue
		}

		end := keywordEnd(input, cur)
		if tokens.isDefaultArgument() {
			// expect list
			if input[cur] == '[' {
				newTokens, newCur := newListTokens(input, cur, col, line)
				tokens = append(tokens, newTokens...)
				col += newCur - cur
				cur = newCur

				continue
			}

			// expect value
			token, cur = newValueToken(input, cur, end, col, line)
			tokens = append(tokens, token)
			col += len(token.Value)
			continue
		}

		if tokens.isDirectiveArgument() {
			newTokens, newCur := newDirectiveArgumentTokens(input, cur, col, line)
			tokens = append(tokens, newTokens...)
			col += newCur - cur
			cur = newCur

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

		if unicode.IsLetter(rune(input[cur])) || unicode.IsDigit(rune(input[cur])) {
			keyword := keyword(input[cur:end])
			if t, ok := keywords[keyword]; ok {
				token, cur = newKeywordToken(input, t, cur, col, line)
				tokens = append(tokens, token)
				col += len(token.Value)
				continue
			}

			if tokens.isField() ||
				tokens.isType() ||
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

		if cur == prev {
			return nil, fmt.Errorf("unexpected character %q at %d:%d", input[cur], line, col)
		}
		prev = cur
	}

	tokens = append(tokens, &Token{Type: EOF, Value: nil, Column: col, Line: line})

	return tokens, nil
}

type keyword string

func (k keyword) String() string {
	return string(k)
}

var keywords = map[keyword]Type{
	"type":      ReservedType,
	"extend":    Extend,
	"scalar":    Scalar,
	"enum":      Enum,
	"input":     Input,
	"interface": Interface,
	"union":     Union,
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
	'|': Pipe,
}

func directiveKeywordEnd(input []byte, cur int) int {
	bracketOpenCount := 0
	for cur < len(input) {
		if input[cur] == '[' {
			bracketOpenCount++
		}
		if input[cur] == ']' {
			bracketOpenCount--
		}

		cur++
		if input[cur] == ')' || (input[cur] == ',' && bracketOpenCount == 0) {
			break
		}
	}

	return cur
}

func keywordEnd(input []byte, cur int) int {
	for cur < len(input) && (unicode.IsLetter(rune(input[cur])) || unicode.IsDigit(rune(input[cur]))) {
		cur++
		if input[cur] == '"' {
			cur++
			break
		}
	}

	return cur
}
