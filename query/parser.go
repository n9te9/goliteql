package query

import (
	"bytes"
	"fmt"
)

type OperationType string

const (
	QueryOperation        OperationType = "query"
	MutationOperation     OperationType = "mutation"
	SubscriptionOperation OperationType = "subscription"
)

type FieldType struct {
	Name     []byte
	Nullable bool
	IsList   bool
	ListType *FieldType
}

type Variable struct {
	Name         []byte
	Type         *FieldType
	DefaultValue []byte
}

type Directive struct {
	Name      []byte
	Arguments []*DirectiveArgument
}

type Operation struct {
	OperationType OperationType
	Name          string
	Variables     []*Variable
	Selections    []Selection
	Directives    []*Directive
}

type Operations []*Operation

func (o Operations) GetQuery() *Operation {
	for _, op := range o {
		if op.OperationType == QueryOperation {
			return op
		}
	}

	return nil
}

func (o Operations) GetMutation() *Operation {
	for _, op := range o {
		if op.OperationType == MutationOperation {
			return op
		}
	}

	return nil
}

func (o Operations) GetSubscription() *Operation {
	for _, op := range o {
		if op.OperationType == SubscriptionOperation {
			return op
		}
	}

	return nil
}

type Selection interface {
	isSelection()
	GetSelections() []Selection
}

type Argument struct {
	Name         []byte
	Type         *FieldType
	DefaultValue []byte
}

type DirectiveArgument struct {
	Name       []byte
	Value      []byte
	IsVariable bool
}

type Field struct {
	Name       []byte
	Arguments  []*Argument
	Selections []Selection
	Directives []*Directive
}

func (f *Field) isSelection() {}

func (f *Field) GetSelections() []Selection {
	return f.Selections
}

type FragmentSpread struct {
	Name       []byte
	Directives []*Directive
}

func (f *FragmentSpread) isSelection() {}

func (f *FragmentSpread) GetSelections() []Selection {
	return nil
}

type InlineFragment struct {
	TypeCondition []byte
	Selections    []Selection
	Directives    []*Directive
}

func (f *InlineFragment) isSelection() {}

func (f *InlineFragment) GetSelections() []Selection {
	return f.Selections
}

type Document struct {
	tokens              []*Token
	Operations          Operations
	FragmentDefinitions FragmentDefinitions
	Name                []byte
}

type FragmentDefinitions []*FragmentDefinition

func (f FragmentDefinitions) GetFragment(name []byte) *FragmentDefinition {
	for _, fragment := range f {
		if bytes.Equal(fragment.Name, name) {
			return fragment
		}
	}

	return nil
}

type FragmentDefinition struct {
	Name          []byte
	BasedTypeName []byte
	Selections    []Selection
}

func (f *FragmentDefinition) isSelection() {}

func (f *FragmentDefinition) GetSelections() []Selection {
	return f.Selections
}

type Parser struct {
	Lexer *Lexer
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		Lexer: lexer,
	}
}

func NewParserWithLexer() *Parser {
	return &Parser{
		Lexer: NewLexer(),
	}
}

func (p *Parser) Parse(input []byte) (*Document, error) {
	tokens, err := p.Lexer.Lex(input)
	if err != nil {
		return nil, err
	}

	cur := 0
	doc := &Document{
		tokens:     tokens,
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

		if tokens[cur].Type == Fragment {
			fragmentDefinition, newCur, err := p.parseFragmentDefinition(tokens, cur)
			if err != nil {
				return nil, err
			}

			cur = newCur
			doc.FragmentDefinitions = append(doc.FragmentDefinitions, fragmentDefinition)
		}

		if tokens[cur].Type == EOF {
			break
		}
	}

	return doc, nil
}

func (p *Parser) parseFragmentDefinition(tokens Tokens, cur int) (*FragmentDefinition, int, error) {
	cur++
	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected fragment name but got %s", tokens[cur].Value)
	}

	fragmentName := tokens[cur].Value
	cur++

	if tokens[cur].Type != On {
		return nil, cur, fmt.Errorf("expected on after fragment name")
	}
	cur++

	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected type name after on")
	}

	typeName := tokens[cur].Value
	cur++

	if tokens[cur].Type != CurlyOpen {
		return nil, cur, fmt.Errorf("expected { after type name")
	}
	cur++

	selections, newCur, err := p.parseSelections(tokens, cur)
	if err != nil {
		return nil, newCur, err
	}
	cur = newCur

	if tokens[cur].Type != CurlyClose {
		return nil, cur, fmt.Errorf("expected } after fragment")
	}
	cur++

	return &FragmentDefinition{
		Name:          fragmentName,
		BasedTypeName: typeName,
		Selections:    selections,
	}, cur, nil
}

func (p *Parser) parseOperation(tokens Tokens, cur int) (*Operation, int, error) {
	operationType := OperationType(tokens[cur].Value)
	cur++

	operationName := ""
	if tokens[cur].Type == Name {
		operationName = string(tokens[cur].Value)
		cur++
	}

	op := &Operation{
		OperationType: operationType,
		Name:          operationName,
	}

	if tokens[cur].Type == ParenOpen {
		variables, newCur, err := p.parseOperationVariables(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}
		cur = newCur
		op.Variables = variables
	}

	for tokens[cur].Type == At {
		cur++
		directive, newCur, err := p.parseDirective(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}
		cur = newCur
		op.Directives = append(op.Directives, directive)
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

	var directives []*Directive = nil
	for tokens[cur].Type == At {
		cur++
		directive, newCur, err := p.parseDirective(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}
		cur = newCur
		directives = append(directives, directive)
	}

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
		Selections:    selections,
		Directives:    directives,
	}, cur + 1, nil
}

func (p *Parser) parseFragmentSpread(tokens Tokens, cur int) (*FragmentSpread, int, error) {
	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected fragment name but got %s", tokens[cur].Value)
	}

	v := tokens[cur].Value
	var directives []*Directive = nil
	cur++
	for tokens[cur].Type == At {
		cur++
		directive, newCur, err := p.parseDirective(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}

		cur = newCur
		directives = append(directives, directive)
	}

	return &FragmentSpread{
		Name:       v,
		Directives: directives,
	}, cur, nil
}

func (p *Parser) parseDirective(tokens Tokens, cur int) (*Directive, int, error) {
	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected directive name but got %s", tokens[cur].Value)
	}

	v := tokens[cur].Value

	cur++

	var arguments []*DirectiveArgument = nil
	var err error
	if tokens[cur].Type == ParenOpen {
		arguments, cur, err = p.parseDirectiveArguments(tokens, cur)
		if err != nil {
			return nil, -1, err
		}
	}

	return &Directive{
		Arguments: arguments,
		Name:      v,
	}, cur, nil
}

func (p *Parser) parseDirectiveArguments(tokens Tokens, cur int) ([]*DirectiveArgument, int, error) {
	arguments := make([]*DirectiveArgument, 0)
	cur++

	for tokens[cur].Type != ParenClose {
		argument, newCur, err := p.parseDirectiveArgument(tokens, cur)
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

func (p *Parser) parseDirectiveArgument(tokens Tokens, cur int) (*DirectiveArgument, int, error) {
	if tokens[cur].Type != Name {
		return nil, cur, fmt.Errorf("expected directive argument name but got %s", tokens[cur].Value)
	}

	name := tokens[cur].Value
	cur++

	if tokens[cur].Type != Colon {
		return nil, cur, fmt.Errorf("expected : but got %s", tokens[cur].Value)
	}
	cur++

	var isVariable bool
	if tokens[cur].Type == Dollar {
		isVariable = true
		cur++

		if tokens[cur].Type != Name {
			return nil, cur, fmt.Errorf("expected variable name after $")
		}

		return &DirectiveArgument{
			Name:       name,
			Value:      tokens[cur].Value,
			IsVariable: isVariable,
		}, cur, nil
	}

	if tokens[cur].Type == Value || tokens[cur].Type == Name {
		return &DirectiveArgument{
			Name:       name,
			Value:      tokens[cur].Value,
			IsVariable: isVariable,
		}, cur + 1, nil
	}

	if tokens[cur].Type == CurlyOpen {
		newValue, newCur, err := p.parseObjectValue(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}

		return &DirectiveArgument{
			Name:       name,
			Value:      newValue,
			IsVariable: isVariable,
		}, newCur, nil
	}

	if tokens[cur].Type == BracketOpen {
		newValue, newCur, err := p.parseListValue(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}

		return &DirectiveArgument{
			Name:       name,
			Value:      newValue,
			IsVariable: isVariable,
		}, newCur, nil
	}

	return nil, cur, fmt.Errorf("unexpected token")
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

	for tokens[cur].Type == At {
		cur++
		directive, newCur, err := p.parseDirective(tokens, cur)
		if err != nil {
			return nil, newCur, err
		}
		cur = newCur
		field.Directives = append(field.Directives, directive)
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
	if tokens[cur].Type == Dollar {
		cur++
	}

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
		Name:         variableName,
		Type:         variableType,
		DefaultValue: defaultValue,
	}, cur, nil
}

func (p *Parser) parseFieldType(tokens Tokens, cur, nestedRank int) (*FieldType, int, error) {
	fieldType := &FieldType{
		Nullable: true,
	}

	if tokens[cur].Type == BracketOpen {
		newFieldType, newCur, err := p.parseFieldType(tokens, cur+1, nestedRank+1)
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
		if tokens[cur].Type == BracketOpen {
			listValue, newCur, err := p.parseListValue(tokens, cur)
			if err != nil {
				return nil, newCur, err
			}

			objectValue = append(objectValue, listValue...)
			cur = newCur
			continue
		}

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
		if tokens[cur].Type == CurlyOpen {
			objectValue, newCur, err := p.parseObjectValue(tokens, cur)
			if err != nil {
				return nil, newCur, err
			}

			listValue = append(listValue, objectValue...)
			cur = newCur
			continue
		}

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
