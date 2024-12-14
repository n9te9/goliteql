package schema

import (
	"fmt"

)

type Parser struct {
	Lexer *Lexer
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		Lexer: lexer,
	}
}

func (p *Parser) Parse(input []byte) (*Schema, error) {
	tokens, err := p.Lexer.Lex(input)
	if err != nil {
		return nil, err
	}

	schema := NewSchema(tokens)


	cur := 0
	for cur < len(tokens) {
		switch tokens[cur].Type {
		case Extend:
			cur++
			schema, cur, err = p.parseExtendDefinition(schema, tokens, cur)
			if err != nil {
				return nil, err
			}
		case ReservedSchema:
			cur++
			definition, newCur, err := p.parseSchemaDefinition(tokens, cur)
			if err != nil {
				return nil, err
			}
			cur = newCur
			schema.Definition = definition
		case ReservedType:
			t := tokens[cur].Type
			cur++
			if cur >= len(tokens) {
				return nil, fmt.Errorf("unexpected end of input")
			}

			if t == ReservedType && tokens[cur].Type == Identifier {
				typeDefinition, newCur, err := p.parseTypeDefinition(tokens, cur)
				if err != nil {
					return nil, err
				}
				cur = newCur
				schema.Types = append(schema.Types, typeDefinition)
				schema.indexes, err = add(schema.indexes, typeDefinition)
				if err != nil {
					return nil, err
				}

				continue
			}

			if tokens[cur].Type == Query || tokens[cur].Type == Mutate || tokens[cur].Type == Subscription {
				operationDefinition, newCur, err := p.parseOperationDefinition(tokens, cur)
				if err != nil {
					return nil, err
				}
				cur = newCur

				schema.Operations = append(schema.Operations, operationDefinition)
				schema.indexes, err = add(schema.indexes, operationDefinition)
				if err != nil {
					return nil, err
				}
				continue
			}

			return nil, fmt.Errorf("unexpected token %s", string(tokens[cur].Value))
		case Input:
			cur++
			if tokens[cur].Type == Identifier {
				inputDefinition, newCur, err := p.parseInputDefinition(tokens, cur)
				if err != nil {
					return nil, err
				}
				cur = newCur
				schema.Inputs = append(schema.Inputs, inputDefinition)
				schema.indexes, err = add(schema.indexes, inputDefinition)
				if err != nil {
					return nil, err
				}
				continue
			}
		case Enum:
			enumDefinition, newCur, err := p.parseEnumDefinition(tokens, cur)
			if err != nil {
				return nil, err
			}
			cur = newCur
			schema.Enums = append(schema.Enums, enumDefinition)
			schema.indexes, err = add(schema.indexes, enumDefinition)
			if err != nil {
				return nil, err
			}
			continue
		case Interface:
			interfaceDefinition, newCur, err := p.parseInterfaceDefinition(tokens, cur)
			if err != nil {
				return nil, err
			}
			cur = newCur
			schema.Interfaces = append(schema.Interfaces, interfaceDefinition)
		case Union:
			unionDefinition, newCur, err := p.parseUnionDefinition(tokens, cur)
			if err != nil {
				return nil, err
			}
			cur = newCur
			schema.Unions = append(schema.Unions, unionDefinition)
			schema.indexes, err = add(schema.indexes, unionDefinition)
			if err != nil {
				return nil, err
			}

		case EOF:
			return schema, nil
		}
	}

	return nil, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseExtendDefinition(schema *Schema, tokens Tokens, cur int) (*Schema, int, error) {
	switch tokens[cur].Type {
	case ReservedSchema:
		cur++
		definition, newCur, err := p.parseSchemaDefinition(tokens, cur)
		if err != nil {
			return nil, 0, err
		}

		cur = newCur
		schema.Definition.Extentions = append(schema.Definition.Extentions, definition)
	case ReservedType:
		cur++
		if tokens[cur].Type == Identifier {
			typeDefinition, newCur, err := p.parseTypeDefinition(tokens, cur)
			if err != nil {
				return nil, 0, err
			}
			cur = newCur
			t := get(schema.indexes, string(typeDefinition.Name), typeDefinition)
			if t == nil {
				return nil, 0, fmt.Errorf("%s is not defined", typeDefinition.Name)
			}

			t.Extentions = append(t.Extentions, typeDefinition)
		}
	case Interface:

	case Union:

	case Enum:

	case Input:

	case Scalar:
		// TODO: Support
	}

	return schema, cur, nil
}

func (p *Parser) parseSchemaDefinition(tokens Tokens, cur int) (*SchemaDefinition, int, error) {
	definition := new(SchemaDefinition)
	if tokens[cur].Type != CurlyOpen {
		return nil, 0, fmt.Errorf("expected '{' but got %s", string(tokens[cur].Value))
	}
	cur++

	for cur < len(tokens) {
		if tokens[cur].Type != Field {
			return nil, 0, fmt.Errorf("expected field but got %s", string(tokens[cur].Value))
		}

		v := string(tokens[cur].Value)
		if v != "query" && v != "mutation" && v != "subscription" {
			return nil, 0, fmt.Errorf("expected query, mutation or subscription but got %s", v)
		}
		cur++

		if tokens[cur].Type != Colon {
			return nil, 0, fmt.Errorf("expected ':' but got %s", string(tokens[cur].Value))
		}
		cur++

		switch v {
		case "query":
			if tokens[cur].Type == Identifier || tokens[cur].Type == Query {
				definition.Query = tokens[cur].Value
				cur++
			}
		case "mutation":
			if tokens[cur].Type == Identifier || tokens[cur].Type == Mutate {
				definition.Mutation = tokens[cur].Value
				cur++
			}
		case "subscription":
			if tokens[cur].Type == Identifier || tokens[cur].Type == Subscription {
				definition.Subscription = tokens[cur].Value
				cur++
			}
		default:
			return nil, 0, fmt.Errorf("unexpected token %s", string(tokens[cur].Value))
		}

		if tokens[cur].Type == CurlyClose {
			break
		}
	}

	if tokens[cur].Type != CurlyClose {
		return nil, 0, fmt.Errorf("expected '}' but got %s", string(tokens[cur].Value))
	}
	cur++

	return definition, cur, nil
}

func (p *Parser) parseTypeDefinition(tokens Tokens, cur int) (*TypeDefinition, int, error) {
	start := cur
	definition := &TypeDefinition{
		Fields: make([]*FieldDefinition, 0),
		Name:   tokens[cur].Value,
	}

	cur++
	if tokens[cur].Type != CurlyOpen {
		return nil, 0, fmt.Errorf("expected '{' but got %s", string(tokens[cur].Value))
	}

	cur++
	for cur < len(tokens) {
		switch tokens[cur].Type {
		case Field:
			fieldDefinitions, newCur, err := p.parseFieldDefinitions(tokens, cur)
			if err != nil {
				return nil, 0, err
			}
			definition.Fields = append(definition.Fields, fieldDefinitions...)
			cur = newCur
		case CurlyClose:
			definition.tokens = tokens[start:cur]
			cur++
			return definition, cur, nil
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseInputDefinition(tokens Tokens, cur int) (*InputDefinition, int, error) {
	start := cur
	definition := &InputDefinition{
		Fields: make([]*FieldDefinition, 0),
		Name:   tokens[cur].Value,
	}

	cur++
	if tokens[cur].Type != CurlyOpen {
		return nil, 0, fmt.Errorf("expected '{' but got %s", string(tokens[cur].Value))
	}

	cur++
	for cur < len(tokens) {
		switch tokens[cur].Type {
		case Field:
			fieldDefinitions, newCur, err := p.parseFieldDefinitions(tokens, cur)
			if err != nil {
				return nil, 0, err
			}
			definition.Fields = append(definition.Fields, fieldDefinitions...)
			cur = newCur
		case CurlyClose:
			definition.tokens = tokens[start:cur]
			cur++
			return definition, cur, nil
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseEnumDefinition(tokens Tokens, cur int) (*EnumDefinition, int, error) {
	cur++
	if tokens[cur].Type != Identifier {
		return nil, 0, fmt.Errorf("expected identifier but got %s", string(tokens[cur].Value))
	}

	enumDefinition := &EnumDefinition{
		Name: tokens[cur].Value,
	}
	cur++

	if tokens[cur].Type != CurlyOpen {
		return nil, 0, fmt.Errorf("expected '{' but got %s", string(tokens[cur].Value))
	}

	cur++
	for cur < len(tokens) {
		switch tokens[cur].Type {
		case Identifier:
			enumDefinition.Values = append(enumDefinition.Values, tokens[cur].Value)
			cur++
		case CurlyClose:
			cur++
			return enumDefinition, cur, nil
		default:
			return nil, 0, fmt.Errorf("unexpected token %s", string(tokens[cur].Value))
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseOperationDefinition(tokens Tokens, cur int) (*OperationDefinition, int, error) {
	var operationType OperationType
	switch tokens[cur].Type {
	case Query:
		operationType = QueryOperation
	case Mutate:
		operationType = MutationOperation
	case Subscription:
		operationType = SubscriptionOperation
	default:
		return nil, 0, fmt.Errorf("unexpected token %s", string(tokens[cur].Value))
	}
	cur++

	if tokens[cur].Type != CurlyOpen {
		return nil, 0, fmt.Errorf("expected identifier but got %s", string(tokens[cur].Value))
	}

	operationDefinition := &OperationDefinition{
		OperationType: operationType,
		Fields:        make([]*FieldDefinition, 0),
	}
	cur++

	for cur < len(tokens) {
		switch tokens[cur].Type {
		case Field:
			fieldDefinitions, newCur, err := p.parseOperationFields(tokens, cur)
			if err != nil {
				return nil, 0, err
			}
			operationDefinition.Fields = append(operationDefinition.Fields, fieldDefinitions...)
			cur = newCur
		case CurlyClose:
			cur++
			return operationDefinition, cur, nil
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseOperationFields(tokens Tokens, cur int) ([]*FieldDefinition, int, error) {
	definitions := make([]*FieldDefinition, 0)

	for cur < len(tokens) {
		switch tokens[cur].Type {
		case Field:
			fieldDefinition, newCur, err := p.parseOperationField(tokens, cur)
			if err != nil {
				return nil, 0, err
			}
			definitions = append(definitions, fieldDefinition)
			cur = newCur
			continue
		case CurlyClose:
			return definitions, cur, nil
		case EOF:
			return nil, 0, fmt.Errorf("unexpected end of input")
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseOperationField(tokens Tokens, cur int) (*FieldDefinition, int, error) {
	definition := &FieldDefinition{
		Name:      tokens[cur].Value,
		Arguments: make([]*ArgumentDefinition, 0),
		Type:      nil,
	}
	cur++

	if tokens[cur].Type == ParenOpen {
		args, newCur, err := p.parseArguments(tokens, cur)
		if err != nil {
			return nil, 0, err
		}
		definition.Arguments = args
		cur = newCur
	}

	if tokens[cur].Type == Colon {
		cur++

		fieldType, newCur, err := p.parseFieldType(tokens, cur)
		if err != nil {
			return nil, 0, err
		}
		definition.Type = fieldType
		cur = newCur

		directiveDefinitions, newCur, err := p.parseDirectives(tokens, cur)
		if err != nil {
			return nil, 0, err
		}
		definition.Directives = directiveDefinitions
		cur = newCur
	}

	return definition, cur, nil
}

func (p *Parser) parseDirectives(tokens Tokens, cur int) ([]*Directive, int, error) {
	definitions := make([]*Directive, 0)

	var err error
	for cur < len(tokens) {
		switch tokens[cur].Type {
		case At:
			cur++
			var definition *Directive
			if cur < len(tokens) && tokens[cur].Type == Identifier {
				definition = &Directive{
					Name: tokens[cur].Value,
				}
				cur++
			} else {
				return nil, 0, fmt.Errorf("expected identifier but got %s", string(tokens[cur].Value))
			}

			if tokens[cur].Type == ParenOpen {
				definition.Arguments, cur, err = p.parseDirectiveArguments(tokens, cur)
				if err != nil {
					return nil, 0, err
				}
			}

			definitions = append(definitions, definition)
		default:
			return definitions, cur, nil
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseDirectiveArguments(tokens Tokens, cur int) ([]*DirectiveArgument, int, error) {
	args := make([]*DirectiveArgument, 0)
	for cur < len(tokens) {
		switch tokens[cur].Type {
		case ParenOpen, Comma:
			cur++
			continue
		case Field:
			arg, newCur, err := p.parseDirectiveArgument(tokens, cur)
			if err != nil {
				return nil, 0, err
			}
			args = append(args, arg)
			cur = newCur
		case Colon:
			return args, cur, nil
		case ParenClose:
			cur++
			return args, cur, nil
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseDirectiveArgument(tokens Tokens, cur int) (*DirectiveArgument, int, error) {
	arg := &DirectiveArgument{
		Name: tokens[cur].Value,
	}
	cur++

	if tokens[cur].Type != Colon {
		return nil, 0, fmt.Errorf("expected ':' but got %s", string(tokens[cur].Value))
	}
	cur++

	switch tokens[cur].Type {
	case Value:
		arg.Value = tokens[cur].Value
		cur++
	default:
		return nil, 0, fmt.Errorf("unexpected token %s", string(tokens[cur].Value))
	}

	return arg, cur, nil
}

func (p *Parser) parseArguments(tokens Tokens, cur int) ([]*ArgumentDefinition, int, error) {
	args := make([]*ArgumentDefinition, 0)
	for cur < len(tokens) {
		switch tokens[cur].Type {
		case ParenOpen, Comma:
			cur++
			continue
		case Field:
			arg, newCur, err := p.parseArgument(tokens, cur)
			if err != nil {
				return nil, 0, err
			}
			args = append(args, arg)
			cur = newCur
		case Colon:
			return args, cur, nil
		case ParenClose:
			cur++
			return args, cur, nil
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseArgument(tokens Tokens, cur int) (*ArgumentDefinition, int, error) {
	arg := &ArgumentDefinition{
		Name: tokens[cur].Value,
	}
	cur++

	if tokens[cur].Type != Colon {
		return nil, 0, fmt.Errorf("expected ':' but got %s", string(tokens[cur].Value))
	}
	cur++

	fieldType, newCur, err := p.parseFieldType(tokens, cur)
	if err != nil {
		return nil, 0, err
	}
	arg.Type = fieldType
	cur = newCur

	if tokens[cur].Type == Equal {
		cur++
		switch tokens[cur].Type {
		case Value:
			arg.Default = tokens[cur].Value
			cur++
		default:
			return nil, 0, fmt.Errorf("unexpected token %s", string(tokens[cur].Value))
		}
	}

	return arg, cur, nil
}

func (p *Parser) parseInterfaceDefinition(tokens Tokens, cur int) (*InterfaceDefinition, int, error) {
	cur++
	if tokens[cur].Type != Identifier {
		return nil, 0, fmt.Errorf("expected identifier but got %s", string(tokens[cur].Type))
	}

	interfaceDefinition := &InterfaceDefinition{
		Name:   tokens[cur].Value,
		Fields: make([]*FieldDefinition, 0),
	}
	cur++

	if tokens[cur].Type != CurlyOpen {
		return nil, 0, fmt.Errorf("expected '{' but got %s", string(tokens[cur].Value))
	}

	cur++
	for cur < len(tokens) {
		switch tokens[cur].Type {
		case Field:
			fieldDefinitions, newCur, err := p.parseFieldDefinitions(tokens, cur)
			if err != nil {
				return nil, 0, err
			}
			interfaceDefinition.Fields = append(interfaceDefinition.Fields, fieldDefinitions...)
			cur = newCur
		case CurlyClose:
			cur++
			return interfaceDefinition, cur, nil
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseFieldDefinitions(tokens Tokens, cur int) ([]*FieldDefinition, int, error) {
	definitions := make([]*FieldDefinition, 0)

	for cur < len(tokens) {
		switch tokens[cur].Type {
		case CurlyOpen, ParenOpen:
			cur++
		case Field:
			fieldDefinition, newCur, err := p.parseFieldDefinition(tokens, cur)
			if err != nil {
				return nil, 0, err
			}
			definitions = append(definitions, fieldDefinition)
			cur = newCur
		case CurlyClose, ParenClose:
			return definitions, cur, nil
		case EOF:
			return nil, 0, fmt.Errorf("unexpected end of input")
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseFieldDefinition(tokens Tokens, cur int) (*FieldDefinition, int, error) {
	definition := &FieldDefinition{
		Name: tokens[cur].Value,
	}

	cur++
	if tokens[cur].Type != Colon {
		return nil, 0, fmt.Errorf("expected ':' but got %s", string(tokens[cur].Value))
	}

	cur++
	switch tokens[cur].Type {
	case Identifier, BracketOpen:
		fieldType, newCur, err := p.parseFieldType(tokens, cur)
		if err != nil {
			return nil, 0, err
		}
		cur = newCur
		definition.Type = fieldType
	default:
		return nil, 0, fmt.Errorf("expected identifier or '[' but got %s", string(tokens[cur].Value))
	}

	if tokens[cur].Type == Equal {
		cur++
		switch tokens[cur].Type {
		case Value:
			definition.Default = tokens[cur].Value
			cur++
		default:
			return nil, 0, fmt.Errorf("unexpected token %s", string(tokens[cur].Value))
		}
	}

	return definition, cur, nil
}

func (p *Parser) parseFieldType(tokens Tokens, cur int) (*FieldType, int, error) {
	fieldType := &FieldType{
		Nullable: true,
	}

	if tokens[cur].Type == Identifier {
		fieldType.Name = tokens[cur].Value
		cur++
	}

	// for nested list types
	if tokens[cur].Type == BracketOpen {
		listType, newCur, err := p.parseFieldType(tokens, cur+1)
		if err != nil {
			return nil, 0, err
		}
		cur = newCur
		fieldType.ListType = listType
		fieldType.IsList = true
	}

	if tokens[cur].Type == Exclamation {
		fieldType.Nullable = false
		cur++
	}

	if tokens[cur].Type == BracketClose {
		cur++
	}

	return fieldType, cur, nil
}

func (p *Parser) parseUnionDefinition(tokens Tokens, cur int) (*UnionDefinition, int, error) {
	cur++
	if tokens[cur].Type != Identifier {
		return nil, 0, fmt.Errorf("expected identifier but got %s", string(tokens[cur].Value))
	}

	unionDefinition := &UnionDefinition{
		Name: tokens[cur].Value,
	}
	cur++

	if tokens[cur].Type != Equal {
		return nil, 0, fmt.Errorf("expected '=' but got %s", string(tokens[cur].Value))
	}
	prev := tokens[cur]
	cur++

	for cur < len(tokens) {
		switch tokens[cur].Type {
		case Pipe:
			if prev.Type == Equal || prev.Type == Identifier {
				prev = tokens[cur]
				cur++
			} else {
				return nil, 0, fmt.Errorf("unexpected token %s", string(tokens[cur].Value))
			}
		case Identifier:
			if prev.Type == Equal || prev.Type == Pipe {
				unionDefinition.Types = append(unionDefinition.Types, tokens[cur].Value)
				prev = tokens[cur]
				cur++
			}
		case EOF:
			if prev.Type != Identifier {
				return nil, 0, fmt.Errorf("unexpected end of input")
			}

			return unionDefinition, cur, nil
		case ReservedType, Union, Enum, Interface, Input, Extend, ReservedSchema:
			return unionDefinition, cur, nil
		default:
			return nil, 0, fmt.Errorf("unexpected token %s", string(tokens[cur].Value))
		}
	}

	return nil, 0, fmt.Errorf("unexpected end of input")
}
