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

	selections, newCur, err := p.parseSelections(tokens, cur)
	if err != nil {
		return nil, newCur, err
	}
	cur = newCur
	op.Selections = selections

	if tokens[cur].Type != CurlyClose {
		return nil, cur, fmt.Errorf("expected } after operation")
	}
	cur++

	return op, cur, nil
}

func (p *Parser) parseSelections(tokens Tokens, cur int) ([]Selection, int, error) {
	var selections []Selection = nil
	for tokens[cur].Type != CurlyClose {
		newSelection, newCur, err := p.parseSelection(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}
		selections = append(selections, newSelection)
		cur = newCur
	}

	return selections, cur, nil
}

func (p *Parser) parseSelection(tokens Tokens, cur int) (Selection, int, error) {
	if tokens[cur].Type == Spread {
		cur++
		return p.parseFragment(tokens, cur)
	}

	return p.parseField(tokens, cur)
}

func (p *Parser) parseFragment(tokens Tokens, cur int) (Selection, int, error) {
	if tokens[cur].Type == On {
		cur++
		return p.parseInlineFragment(tokens, cur)
	}

	return p.parseFragmentSpread(tokens, cur)
}

func (p *Parser) parseInlineFragment(tokens Tokens, cur int) (*InlineFragment, int, error) {
	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected type name but got %s", tokens[cur].Value)
	}

	v := tokens[cur].Value
	cur++

	if tokens[cur].Type == CurlyOpen {
		cur++
	} else {
		return nil, cur, fmt.Errorf("expected { after type name")
	}

	selections, newCur, err := p.parseSelections(tokens, cur)
	if err != nil {
		return nil, newCur, err
	}
	cur = newCur

	return &InlineFragment{
		TypeCondition: v,
		Selections: selections,
	}, cur + 1, nil
}

func (p *Parser) parseFragmentSpread(tokens Tokens, cur int) (*FragmentSpread, int, error) {
	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected fragment name but got %s", tokens[cur].Value)
	}

	return &FragmentSpread{
		Name: tokens[cur].Value,
	}, cur + 1, nil
}

func (p *Parser) parseField(tokens Tokens, cur int) (*Field, int, error) {
	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected field name but got %s", tokens[cur].Value)
	}

	field := &Field{
		Name: tokens[cur].Value,
	}
	cur++

	if tokens[cur].Type == ParenOpen {
		arguments, newCur, err := p.parseFieldArguments(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}
		cur = newCur
		field.Arguments = arguments
	}

	if tokens[cur].Type == CurlyOpen {
		cur++

		selections, newCur, err := p.parseSelections(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}
		cur = newCur + 1
		field.Selections = selections
	}

	return field, cur, nil
}

func (p *Parser) parseFieldArguments(tokens Tokens, cur int) ([]*Argument, int, error) {
	arguments := make([]*Argument, 0)
	cur++

	for tokens[cur].Type != ParenClose {
		argument, newCur, err := p.parseFieldArgument(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}
		arguments = append(arguments, argument)
		cur = newCur

		if tokens[cur].Type == Comma {
			cur++
			continue
		}
	}

	cur++

	return arguments, cur, nil
}

func (p *Parser) parseFieldArgument(tokens Tokens, cur int) (*Argument, int, error) {
	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected argument name")
	}

	argument := &Argument{
		Name: tokens[cur].Value,
	}
	cur++

	if tokens[cur].Type != Colon {
		return nil, cur, fmt.Errorf("expected : after argument name")
	}
	cur++

	fieldType, newCur, err := p.parseFieldType(tokens, cur, 0)
	if err != nil {
		return nil, newCur, err
	}
	cur = newCur
	argument.Type = fieldType

	if tokens[cur].Type == Equal {
		cur++
		if tokens[cur].Type != Value {
			return nil, cur, fmt.Errorf("expected value after =")
		}
		argument.DefaultValue, cur, err = p.parseDefaultValue(tokens, cur)
		if err != nil {
			return nil, cur, err
		}

		cur++
	}

	return argument, cur, nil
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
		if tokens[cur].Type != Value &&
		 tokens[cur].Type != CurlyOpen &&
		 tokens[cur].Type != BracketOpen &&
		 tokens[cur].Type != Name {
			return nil, cur, fmt.Errorf("expected default value but got %s(%s)", tokens[cur].Value, tokens[cur].Type)
		}

		defaultValue, cur, err = p.parseDefaultValue(tokens, cur)
		if err != nil {
			return nil, cur, err
		}
	}

	return &Variable{
		Name: variableName,
		Type: variableType,
		DefaultValue: defaultValue,
	}, cur, nil
}

func (p *Parser) parseFieldType(tokens Tokens, cur, nestedRank int) (*FieldType, int, error) {
	fieldType := &FieldType{
		Nullable: true,
	}

	if tokens[cur].Type == BracketOpen {
		newFieldType, newCur, err := p.parseFieldType(tokens, cur + 1, nestedRank + 1)
		if err != nil {
			return nil, cur, err
		}

		fieldType.ListType = newFieldType
		fieldType.IsList = true
		cur = newCur
	}

	if tokens[cur].Type == BracketClose {
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
	if tokens[cur].Type != Value && 
	tokens[cur].Type != CurlyOpen && 
	tokens[cur].Type != BracketOpen &&
	tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected value")
	}

	if tokens[cur].Type == Value || tokens[cur].Type == Name {
		return tokens[cur].Value, cur + 1, nil
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
	objectValue := make([]byte, 0)

	nested := 0
	for {
		objectValue = append(objectValue, tokens[cur].Value...)

		if tokens[cur].Type == CurlyOpen {
			nested++
		}

		if tokens[cur].Type == CurlyClose {
			nested--
		}

		if nested == 0 {
			break
		}

		cur++
	}

	return objectValue, cur + 1, nil
}

func (p *Parser) parseListValue(tokens Tokens, cur int) ([]byte, int, error) {
	listValue := make([]byte, 0)

	nested := 0
	for {
		listValue = append(listValue, tokens[cur].Value...)

		if tokens[cur].Type == BracketOpen {
			nested++
		}

		if tokens[cur].Type == BracketClose {
			nested--
		}

		if nested == 0 {
			break
		}

		cur++
	}

	return listValue, cur + 1, nil
}
