package schema

import (
	"fmt"
	"reflect"
)

type OperationType string

const (
	QueryOperation OperationType = "query"
	MutationOperation OperationType = "mutation"
	SubscriptionOperation OperationType = "subscription"
)

type FieldDefinitions []*FieldDefinition

func (f FieldDefinitions) Last(name string) *FieldDefinition {
	var res *FieldDefinition
	for _, field := range f {
		if string(field.Name) == name {
			res = field
		}
	}

	return res
}

func (f FieldDefinitions) Has(name string) bool {
	for _, field := range f {
		if string(field.Name) == name {
			return true
		}
	}
	
	return false
}

type TypeDefinition struct {
	Name []byte
	Fields FieldDefinitions
	tokens Tokens
	Interfaces []*InterfaceDefinition
	Directives []*Directive
	Extentions []*TypeDefinition
}

type FieldType struct {
	Name []byte
	Nullable bool
	IsList bool
	ListType *FieldType
}

type FieldDefinition struct {
	Name []byte
	Arguments []*ArgumentDefinition
	Type *FieldType
	Directives []*Directive
	Default []byte
}

type ArgumentDefinition struct {
	Name []byte
	Default []byte
	Type *FieldType
}

type OperationDefinition struct {
	OperationType OperationType
	Name []byte
	Fields FieldDefinitions
	Extentions []*OperationDefinition
}

type EnumElement struct {
	Name []byte
	Value []byte
	Directives []*Directive
}

type EnumDefinition struct {
	Name []byte
	Values []*EnumElement
	Extentions []*EnumDefinition
	Directives []*Directive
}

type UnionDefinition struct {
	Name []byte
	Types [][]byte
	Extentions []*UnionDefinition
	Directives []*Directive
}

type InterfaceDefinition struct {
	Name []byte
	Fields FieldDefinitions
	Extentions []*InterfaceDefinition
	Directives []*Directive
}

type Location struct {
	Name []byte
	
}

type Directive struct {
	Name []byte
	Arguments []*DirectiveArgument
	Locations []*Location
}

type DirectiveArgument struct {
	Name []byte
	Value []byte
}

type DirectiveDefinition struct {
	Name []byte
	Description []byte
	Arguments []*ArgumentDefinition
	Repeatable bool
	Locations []*Location
}

type InputDefinition struct {
	Name []byte
	Fields FieldDefinitions
	tokens Tokens
	Extentions []*InputDefinition
}

type DefinitionType interface {
	*TypeDefinition | *OperationDefinition | *EnumDefinition | *UnionDefinition | *InterfaceDefinition | *InputDefinition
}

type Indexes struct {
	TypeIndex map[string]*TypeDefinition
	EnumIndex map[string]*EnumDefinition
	UnionIndex map[string]*UnionDefinition
	InterfaceIndex map[string]*InterfaceDefinition
	InputIndex map[string]*InputDefinition
	OperationIndexes map[OperationType]map[string]*OperationDefinition
}

type Schema struct {
	tokens Tokens
	Definition *SchemaDefinition
	Operations []*OperationDefinition
	Types []*TypeDefinition
	Enums []*EnumDefinition
	Unions []*UnionDefinition
	Interfaces []*InterfaceDefinition
	Directives []*DirectiveDefinition
	Inputs []*InputDefinition
	Scalars []*ScalarDefinition

	// indexes is used when extend
	indexes *Indexes
}

func NewSchema(tokens Tokens) *Schema {
	operationIndexes := make(map[OperationType]map[string]*OperationDefinition)

	for _, t := range([]OperationType{QueryOperation, MutationOperation, SubscriptionOperation}) {
		operationIndexes[t] = make(map[string]*OperationDefinition)
	}

	return &Schema{
		tokens: tokens,
		Definition: &SchemaDefinition{
			Query:        []byte("Query"),
			Mutation:     []byte("Mutation"),
			Subscription: []byte("Subscription"),
		},
		indexes: &Indexes{
			TypeIndex: make(map[string]*TypeDefinition),
			OperationIndexes: operationIndexes,
			EnumIndex: make(map[string]*EnumDefinition),
			UnionIndex: make(map[string]*UnionDefinition),
			InterfaceIndex: make(map[string]*InterfaceDefinition),
			InputIndex: make(map[string]*InputDefinition),
		},
	}
}

func (s *Schema) digOperation(name string, ops []*OperationDefinition) (FieldDefinitions, error) {
	res := make(FieldDefinitions, 0)

	for _, op := range ops {
		if string(op.Name) != name {
			return nil, fmt.Errorf("operation %s not found", name)
		}

		res = append(res, op.Fields...)
		if len(op.Extentions) > 0 {
			field, err := s.digOperation(name, op.Extentions)
			if err != nil {
				return nil, err
			}
			res = append(res, field...)
		}
	}

	return res, nil
}

func (s *Schema) mergeOperation(newSchema *Schema) error {
	for _, t := range s.Operations {
		newOp := new(OperationDefinition)
		newOp.OperationType = t.OperationType
		newOp.Name = t.Name
		newOp.Fields = t.Fields

		field, err := s.digOperation(string(newOp.Name), t.Extentions)
		if err != nil {
			return err
		}

		newOp.Fields = append(newOp.Fields, field...)

		newFields := make(FieldDefinitions, 0)
		for _, field := range newOp.Fields {
			newField := newOp.Fields.Last(string(field.Name))
			if newField == nil {
				newFields = append(newFields, field)
			} else if !newFields.Has(string(newField.Name)) {
				newFields = append(newFields, newField)
			}
		}
		newOp.Fields = newFields

		newSchema.Operations = append(newSchema.Operations, newOp)
	}

	return nil
}

func (s *Schema) digTypeDefinition(name string, types []*TypeDefinition) (FieldDefinitions, error) {
	res := make(FieldDefinitions, 0)

	for _, t := range types {
		if string(t.Name) != name {
			return nil, fmt.Errorf("type %s not found", name)
		}

		res = append(res, t.Fields...)
		if len(t.Extentions) > 0 {
			field, err := s.digTypeDefinition(name, t.Extentions)
			if err != nil {
				return nil, err
			}
			res = append(res, field...)
		}
	}

	return res, nil
}

func (s *Schema) mergeTypeDefinition(newSchema *Schema) error {
	for _, t := range s.Types {
		newType := new(TypeDefinition)
		newType.Name = t.Name
		newType.Fields = t.Fields
		newType.Interfaces = t.Interfaces
		newType.Directives = t.Directives

		newFields := make(FieldDefinitions, 0)

		field, err := s.digTypeDefinition(string(newType.Name), t.Extentions)
		if err != nil {
			return err
		}

		newFields = append(newFields, field...)

		for _, field := range newType.Fields {
			newField := newType.Fields.Last(string(field.Name))
			if newField == nil {
				newFields = append(newFields, field)
			} else if !newFields.Has(string(newField.Name)) {
				newFields = append(newFields, newField)
			}
		}
		newType.Fields = newFields

		newSchema.Types = append(newSchema.Types, newType)
	}

	return nil
}

func (s *Schema) digInterfaceDefinition(name string, interfaces []*InterfaceDefinition) (FieldDefinitions, error) {
	res := make(FieldDefinitions, 0)

	for _, t := range interfaces {
		if string(t.Name) != name {
			return nil, fmt.Errorf("interface %s not found", name)
		}

		res = append(res, t.Fields...)
		if len(t.Extentions) > 0 {
			field, err := s.digInterfaceDefinition(name, t.Extentions)
			if err != nil {
				return nil, err
			}
			res = append(res, field...)
		}
	}

	return res, nil
}

func (s *Schema) mergeInterfaceDefinition(newSchema *Schema) error {
	for _, t := range s.Interfaces {
		newInterface := new(InterfaceDefinition)
		newInterface.Name = t.Name
		newInterface.Fields = t.Fields
		newInterface.Directives = t.Directives

		newFields := make(FieldDefinitions, 0)

		field, err := s.digInterfaceDefinition(string(newInterface.Name), t.Extentions)
		if err != nil {
			return err
		}

		newFields = append(newFields, field...)

		for _, field := range newInterface.Fields {
			newField := newInterface.Fields.Last(string(field.Name))
			if newField == nil {
				newFields = append(newFields, field)
			} else if !newFields.Has(string(newField.Name)) {
				newFields = append(newFields, newField)
			}
		}
		newInterface.Fields = newFields

		newSchema.Interfaces = append(newSchema.Interfaces, newInterface)
	}

	return nil
}

func (s *Schema) Merge() (*Schema, error) {
	newSchema := new(Schema)
	newSchema.Definition = s.Definition
	newSchema.tokens = s.tokens
	newSchema.indexes = s.indexes
	
	if err := s.mergeOperation(newSchema); err != nil {
		return nil, err
	}

	if err := s.mergeTypeDefinition(newSchema); err != nil {
		return nil, err
	}

	if err := s.mergeInterfaceDefinition(newSchema); err != nil {
		return nil, err
	}

	return newSchema, nil
}

func add[T DefinitionType](indexes *Indexes, definition T) (*Indexes, error) {
	switch d := any(definition).(type) {
	case *OperationDefinition:
		indexes.OperationIndexes[d.OperationType][string(d.Name)] = d
	case *TypeDefinition:
		indexes.TypeIndex[string(d.Name)] = d
	case *EnumDefinition:
		indexes.EnumIndex[string(d.Name)] = d
	case *UnionDefinition:
		indexes.UnionIndex[string(d.Name)] = d
	case *InterfaceDefinition:
		indexes.InterfaceIndex[string(d.Name)] = d
	case *InputDefinition:
		indexes.InputIndex[string(d.Name)] = d
	default:
		return nil, fmt.Errorf("definition type %v is unsupported in index", reflect.TypeOf(d))
	}

	return indexes, nil
}

func get[T DefinitionType](indexes *Indexes, key string, t T) T {
	switch d := any(t).(type) {
	case *OperationDefinition:
		return any(indexes.OperationIndexes[d.OperationType][key]).(T)
	case *TypeDefinition:
		return any(indexes.TypeIndex[key]).(T)
	case *EnumDefinition:
		return any(indexes.EnumIndex[key]).(T)
	case *UnionDefinition:
		return any(indexes.UnionIndex[key]).(T)
	case *InterfaceDefinition:
		return any(indexes.InterfaceIndex[key]).(T)
	case *InputDefinition:
		return any(indexes.InputIndex[key]).(T)
	}

	return any(nil).(T)
}

type ScalarDefinition struct {
	Name []byte
	Directives []*Directive
}

type SchemaDefinition struct {
	Query []byte
	Mutation []byte
	Subscription []byte
	Extentions []*SchemaDefinition

	Directives []*Directive
}