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
	JSONBytes() ([]byte, error)
}

type ValueParserLiteral struct {
	Value     []byte
	TokenType ValueLexerTokenType
	IsField   bool
}

func (v *ValueParserLiteral) Type() ValueParserExprType {
	return ValueParserExprTypeLiteral
}

func (v *ValueParserLiteral) IsString() bool {
	return v.TokenType == STRING
}

func (v *ValueParserLiteral) StringValue() string {
	return string(v.Value)
}

func (v *ValueParserLiteral) IDValue() string {
	return string(v.Value)
}

func (v *ValueParserLiteral) IsInt() bool {
	return v.TokenType == INT
}

func (v *ValueParserLiteral) IntValue() int {
	var value int
	fmt.Sscanf(string(v.Value), "%d", &value)
	return value
}

func (v *ValueParserLiteral) IsFloat() bool {
	return v.TokenType == FLOAT
}

func (v *ValueParserLiteral) FloatValue() float64 {
	var value float64
	fmt.Sscanf(string(v.Value), "%f", &value)
	return value
}

func (v *ValueParserLiteral) IsBool() bool {
	return v.TokenType == BOOL
}

func (v *ValueParserLiteral) BoolValue() bool {
	return string(v.Value) == "true"
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

	expr, _, err := vp.parseValue(tokens, false)
	if err != nil {
		return nil, fmt.Errorf("error parsing value: %w", err)
	}

	return expr, nil
}

func (vp *ValueParser) parseValue(tokens []*ValueLexerToken, isField bool) (ValueParserExpr, int, error) {
	if len(tokens) == 0 {
		return nil, -1, nil
	}

	switch tokens[0].Type {
	case LBRACE:
		return vp.parseObject(tokens)
	case LBRACKET:
		return vp.parseArray(tokens)
	default:
		fmt.Println(isField)
		return vp.parseLiteral(tokens, isField)
	}
}

func (vp *ValueParser) parseLiteral(tokens []*ValueLexerToken, isField bool) (ValueParserExpr, int, error) {
	if len(tokens) == 0 {
		return nil, -1, fmt.Errorf("expected value token")
	}

	if tokens[0].Type == IDENT {
		return nil, -1, fmt.Errorf("expected Scalar token, got IDENT: %s", tokens[0].Value)
	}

	return &ValueParserLiteral{
		Value:     tokens[0].Value,
		TokenType: tokens[0].Type,
		IsField:   isField,
	}, 1, nil
}

func (vpl *ValueParserLiteral) JSONBytes() ([]byte, error) {
	switch vpl.TokenType {
	case STRING:
		if !vpl.IsField {
			return []byte(fmt.Sprintf(`"%s"`, vpl.IDValue())), nil
		}
		return []byte(vpl.StringValue()), nil
	case INT, FLOAT:
		return vpl.Value, nil
	case BOOL:
		return []byte(fmt.Sprintf("%t", vpl.BoolValue())), nil
	case NULL:
		return []byte("null"), nil
	default:
		return nil, fmt.Errorf("unsupported literal type: %v", vpl.TokenType)
	}
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

		value, nextIndex, err := vp.parseValue(tokens[i:], true)
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

func (vpo *ValueParserObject) JSONBytes() ([]byte, error) {
	result := "{"
	first := true

	for key, value := range vpo.Fields {
		if !first {
			result += ","
		}
		first = false

		jsonValue, err := value.JSONBytes()
		if err != nil {
			return nil, err
		}

		result += fmt.Sprintf(`"%s":%s`, key, jsonValue)
	}

	result += "}"
	return []byte(result), nil
}

func (vp *ValueParser) parseArray(tokens []*ValueLexerToken) (ValueParserExpr, int, error) {
	items := make([]ValueParserExpr, 0)
	i := 1

	for i < len(tokens) && tokens[i].Type != RBRACKET {
		value, nextIndex, err := vp.parseValue(tokens[i:], false)
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

func (vpa *ValueParserArray) JSONBytes() ([]byte, error) {
	result := "["
	first := true

	for _, item := range vpa.Items {
		if !first {
			result += ","
		}
		first = false

		jsonValue, err := item.JSONBytes()
		if err != nil {
			return nil, err
		}

		result += string(jsonValue)
	}

	result += "]"
	return []byte(result), nil
}
