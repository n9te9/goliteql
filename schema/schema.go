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

type TypeDefinition struct {
	Name []byte
	Fields []*FieldDefinition
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
	Fields []*FieldDefinition
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
	Fields []*FieldDefinition
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
	Fields []*FieldDefinition
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

func (s *Schema) Merge() (*Schema, error) {
	return nil, nil
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