package query

import "fmt"

type OperationType string

const (
	QueryOperation       OperationType = "query"
	MutationOperation    OperationType = "mutation"
	SubscriptionOperation OperationType = "subscription"
)

type FieldType struct {
	Name []byte
	Nullable bool
	IsList bool
	ListType *FieldType
}

type Variable struct {
	Name []byte
	Type *FieldType
	DefaultValue []byte
}

type Operation struct {
	OperationType OperationType
	Name string
	Variables []*Variable
	Selections []Selection
}

type Selection interface {
	isSelection()
}

type Argument struct {
	Name []byte
	Type *FieldType
	DefaultValue []byte
}

type Field struct {
	Name []byte
	Arguments []*Argument
	Selections []Selection
}

func (f *Field) isSelection() {}

type FragmentSpread struct {
	Name []byte
}

func (f *FragmentSpread) isSelection() {}

type InlineFragment struct {
	TypeCondition []byte
	Selections []Selection
}

func (f *InlineFragment) isSelection() {}

type Document struct {
	tokens []*Token
	Operations []*Operation
	Name []byte
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

	cur := 0
	doc := &Document{
		tokens: tokens,
		Operations: make([]*Operation, 0),
	}
	for cur < len(tokens) {
		if tokens[cur].Type.IsOperation() {
			op, newCur, err := p.parseOperation(tokens, cur)
			if err != nil {
				return nil, err
			}

			cur = newCur
			doc.Operations = append(doc.Operations, op)
		}

		if tokens[cur].Type == EOF {
			break
		}
	}

	return doc, nil
}

func (p *Parser) parseOperation(tokens Tokens, cur int) (*Operation, int, error) {
	operationType := OperationType(tokens[cur].Value)
	cur++

	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected name after operation type")
	}

	op := &Operation{
		OperationType: operationType,
		Name: string(tokens[cur].Value),
	}
	cur++

	if tokens[cur].Type == ParenOpen {
		variables, newCur, err := p.parseOperationVariables(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}
		cur = newCur
		op.Variables = variables
	}

	if tokens[cur].Type != CurlyOpen {
		return nil, cur, fmt.Errorf("expected { after operation")
	}
	cur++

	// TODO: parse selections
	if tokens[cur].Type != CurlyClose {
		return nil, cur, fmt.Errorf("expected } after operation")
	}
	cur++

	return op, cur, nil
}

func (p *Parser) parseOperationVariables(tokens Tokens, cur int) ([]*Variable, int, error) {
	variables := make([]*Variable, 0)
	cur++

	var prev *Token
	for tokens[cur].Type != ParenClose {
		variable, newCur, err := p.parseOperationVariable(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}
		variables = append(variables, variable)
		cur = newCur

		if tokens[cur].Type == Comma {
			cur++
			continue
		}

		if prev == tokens[cur] {
			return nil, cur, fmt.Errorf("expected variable after %v", prev)
		}
		prev = tokens[cur]
	}
	cur++

	return variables, cur, nil
}

func (p *Parser) parseOperationVariable(tokens Tokens, cur int) (*Variable, int, error) {
	if tokens[cur].Type != Dollar {
		return nil, cur, fmt.Errorf("expected $ before variable")
	}
	cur++

	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected variable name but got %s", tokens[cur].Value)
	}

	variableName := tokens[cur].Value
	cur++

	if tokens[cur].Type != Colon {
		return nil, cur, fmt.Errorf("expected : after variable name")
	}
	cur++

	variableType, newCur, err := p.parseFieldType(tokens, cur, 0)
	if err != nil {
		return nil, newCur, err
	}
	cur = newCur
	
	var defaultValue []byte
	if tokens[cur].Type == Equal {
		cur++
		if tokens[cur].Type != Value {
			return nil, cur, fmt.Errorf("expected default value")
		}

		defaultValue, cur, err = p.parseDefaultValue(tokens, cur)
		if err != nil {
			return nil, cur, err
		}

		cur = newCur
	}

	return &Variable{
		Name: variableName,
		Type: variableType,
		DefaultValue: defaultValue,
	}, cur, nil
}

func (p *Parser) parseFieldType(tokens Tokens, cur, nestedRank int) (*FieldType, int, error) {
	fieldType := &FieldType{}
	if tokens[cur].Type == BracketOpen {
		newFieldType, cur, err := p.parseFieldType(tokens, cur, nestedRank + 1)
		if err != nil {
			return nil, cur, err
		}

		fieldType.ListType = newFieldType
	}

	if tokens[cur].Type == BracketClose {
		if nestedRank == 0 {
			return nil, cur, fmt.Errorf("unexpected ]")
		}
		cur++
		if tokens[cur].Type == Exclamation {
			fieldType.Nullable = false
			cur++
		}
		return fieldType, cur, nil
	}

	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected type name but got %s", tokens[cur].Value)	
	}

	fieldType.Name = tokens[cur].Value
	cur++

	if tokens[cur].Type == Exclamation {
		fieldType.Nullable = false
		cur++
	}

	return fieldType, cur, nil
}

func (p *Parser) parseDefaultValue(tokens Tokens, cur int) ([]byte, int, error) {
	if tokens[cur].Type != Value && tokens[cur].Type != CurlyOpen && tokens[cur].Type != BracketOpen {
		return nil, cur, fmt.Errorf("expected value")
	}

	if tokens[cur].Type == Value {
		return tokens[cur].Value, cur, nil
	}

	if tokens[cur].Type == CurlyOpen {
		return p.parseObjectValue(tokens, cur)
	}

	if tokens[cur].Type == BracketOpen {
		return p.parseListValue(tokens, cur)
	}

	return nil, cur, fmt.Errorf("unexpected token")
}

func (p *Parser) parseObjectValue(tokens Tokens, cur int) ([]byte, int, error) {
	cur++
	objectValue := make([]byte, 0)
	for tokens[cur].Type != CurlyClose {
		objectValue = append(objectValue, tokens[cur].Value...)
		cur++
	}

	return objectValue, cur, nil
}

func (p *Parser) parseListValue(tokens Tokens, cur int) ([]byte, int, error) {
	cur++
	listValue := make([]byte, 0)
	for tokens[cur].Type != BracketClose {
		listValue = append(listValue, tokens[cur].Value...)
		cur++
	}

	return listValue, cur, nil
}
