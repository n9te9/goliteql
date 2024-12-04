package schema

import "fmt"

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

	schema := &Schema{
		tokens: tokens,
	}

	cur := 0
	for cur < len(tokens) {
		switch tokens[cur].Type {
		case ReservedType:
			typeDefinition, newCur, err := p.parseTypeDefinition(tokens, cur)
			if err != nil {
				return nil, err
			}
			cur = newCur
			schema.Types = append(schema.Types, typeDefinition)
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
		case EOF:
			return schema, nil
		}
	}

	return nil, fmt.Errorf("unexpected end of input")
}

func (p *Parser) parseTypeDefinition(tokens Tokens, cur int) (*TypeDefinition, int, error) {
	start := cur
	definition := &TypeDefinition{
		Fields: make([]*FieldDefinition, 0),
	}

	cur++
	if tokens[cur].Type != Identifier {
		return nil, 0, fmt.Errorf("expected identifier but got %s", string(tokens[cur].Value))
	}
	definition.Name = tokens[cur].Value

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

func (p *Parser) parseInterfaceDefinition(tokens Tokens, cur int) (*InterfaceDefinition, int, error) {
	cur++
	if tokens[cur].Type != Identifier {
		return nil, 0, fmt.Errorf("expected identifier but got %s", string(tokens[cur].Type))
	}

	interfaceDefinition := &InterfaceDefinition{
		Name: tokens[cur].Value,
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
		case CurlyOpen:
			cur++
		case Field:
			fieldDefinition, newCur, err := p.parseFieldDefinition(tokens, cur)
			if err != nil {
				return nil, 0, err
			}
			definitions = append(definitions, fieldDefinition)
			cur = newCur
		case CurlyClose:
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
		listType, newCur, err := p.parseFieldType(tokens, cur + 1)
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
			if prev.Type == Equal ||  prev.Type == Pipe {
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

