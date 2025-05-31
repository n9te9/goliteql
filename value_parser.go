package goliteql

import (
	"fmt"
)

type ValueParser struct {
	lexer *ValueLexer
}

func NewValueParser(lexer *ValueLexer) *ValueParser {
	return &ValueParser{
		lexer: lexer,
	}
}

type ValueParserExprType int

const (
	ValueParserExprTypeLiteral ValueParserExprType = iota
	ValueParserExprTypeObject
	ValueParserExprTypeArray
)

type ValueParserExpr interface {
	Type() ValueParserExprType
}

type ValueParserLiteral struct {
	Value     []byte
	TokenType ValueLexerTokenType
}

func (v *ValueParserLiteral) Type() ValueParserExprType {
	return ValueParserExprTypeLiteral
}

func (v *ValueParserLiteral) IsString() bool {
	return v.TokenType == STRING
}

func (v *ValueParserLiteral) IsInt() bool {
	return v.TokenType == INT
}

func (v *ValueParserLiteral) IsFloat() bool {
	return v.TokenType == FLOAT
}

func (v *ValueParserLiteral) IsBool() bool {
	return v.TokenType == BOOL
}

func (v *ValueParserLiteral) IsNull() bool {
	return v.TokenType == NULL
}

type ValueParserObject struct {
	Fields map[string]ValueParserExpr
}

func (v *ValueParserObject) Type() ValueParserExprType {
	return ValueParserExprTypeObject
}

type ValueParserArray struct {
	Items []ValueParserExpr
}

func (v *ValueParserArray) Type() ValueParserExprType {
	return ValueParserExprTypeArray
}

func (vp *ValueParser) Parse(input []byte) (ValueParserExpr, error) {
	tokens, err := vp.lexer.Lex(input)
	if err != nil {
		return nil, err
	}

	expr, _, err := vp.parseValue(tokens)
	if err != nil {
		return nil, fmt.Errorf("error parsing value: %w", err)
	}

	return expr, nil
}

func (vp *ValueParser) parseValue(tokens []*ValueLexerToken) (ValueParserExpr, int, error) {
	if len(tokens) == 0 {
		return nil, -1, nil
	}

	switch tokens[0].Type {
	case LBRACE:
		return vp.parseObject(tokens)
	case LBRACKET:
		return vp.parseArray(tokens)
	default:
		return vp.parseLiteral(tokens)
	}
}

func (vp *ValueParser) parseLiteral(tokens []*ValueLexerToken) (ValueParserExpr, int, error) {
	if len(tokens) == 0 {
		return nil, -1, fmt.Errorf("expected value token")
	}

	if tokens[0].Type == IDENT {
		return nil, -1, fmt.Errorf("expected Scalar token, got IDENT: %s", tokens[0].Value)
	}

	return &ValueParserLiteral{
		Value:     tokens[0].Value,
		TokenType: tokens[0].Type,
	}, 1, nil
}

func (vp *ValueParser) parseObject(tokens []*ValueLexerToken) (ValueParserExpr, int, error) {
	fields := make(map[string]ValueParserExpr)
	i := 1

	for i < len(tokens) && tokens[i].Type != RBRACE {
		if tokens[i].Type != IDENT {
			return nil, -1, fmt.Errorf("expected IDENT token, got %s", tokens[i].Value)
		}

		key := string(tokens[i].Value)
		i++

		if i >= len(tokens) || tokens[i].Type != COLON {
			return nil, -1, fmt.Errorf("expected COLON token after key %s, got %s", key, tokens[i].Value)
		}
		i++

		value, nextIndex, err := vp.parseValue(tokens[i:])
		if err != nil {
			return nil, -1, err
		}

		fields[key] = value
		i += nextIndex

		if i < len(tokens) && tokens[i].Type == COMMA {
			i++
		}
	}

	return &ValueParserObject{
		Fields: fields,
	}, i + 1, nil
}

func (vp *ValueParser) parseArray(tokens []*ValueLexerToken) (ValueParserExpr, int, error) {
	items := make([]ValueParserExpr, 0)
	i := 1

	for i < len(tokens) && tokens[i].Type != RBRACKET {
		value, nextIndex, err := vp.parseValue(tokens[i:])
		if err != nil {
			return nil, -1, err
		}

		items = append(items, value)
		i += nextIndex

		if i < len(tokens) && tokens[i].Type == COMMA {
			i++
		}
	}

	return &ValueParserArray{
		Items: items,
	}, i + 1, nil
}
