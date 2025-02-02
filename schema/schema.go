package schema

import (
	"bytes"
	"fmt"
	"reflect"
)

type OperationType string

func (o OperationType) isQuery() bool {
	return o == QueryOperation
}

func (o OperationType) isMutation() bool {
	return o == MutationOperation
}

func (o OperationType) isSubscription() bool {
	return o == SubscriptionOperation
}

const (
	QueryOperation OperationType = "query"
	MutationOperation OperationType = "mutation"
	SubscriptionOperation OperationType = "subscription"
)

type CompositeType interface {
	GetFieldByName(name []byte) *FieldDefinition
	TypeName() []byte
}

type TypeDefinition struct {
	Name []byte
	Fields FieldDefinitions
	Required map[*FieldDefinition]struct{}
	tokens Tokens
	Interfaces []*InterfaceDefinition
	Directives []*Directive
	Extentions []*TypeDefinition
}

func (t *TypeDefinition) GetFieldByName(name []byte) *FieldDefinition {
	for _, field := range t.Fields {
		if bytes.Equal(field.Name, name) {
			return field
		}
	}

	return nil
}

func (t *TypeDefinition) TypeName() []byte {
	return t.Name
}

type FieldType struct {
	Name []byte
	Nullable bool
	IsList bool
	ListType *FieldType
}

func (f *FieldType) GetPremitiveType() *FieldType {
	if f.IsList {
		return f.ListType.GetPremitiveType()
	}

	return f
}

type ArgumentDefinition struct {
	Name []byte
	Default []byte
	Type *FieldType
}

type ArgumentDefinitions []*ArgumentDefinition

func (a ArgumentDefinitions) RequiredArguments() map[*ArgumentDefinition]struct{} {
	res := make(map[*ArgumentDefinition]struct{})
	for _, arg := range a {
		if !arg.Type.Nullable {
			res[arg] = struct{}{}
		}
	}

	return res
}

type OperationDefinition struct {
	OperationType OperationType
	Name []byte
	Fields FieldDefinitions
	Extentions []*OperationDefinition
}

func (o *OperationDefinition) GetFieldByName(name []byte) *FieldDefinition {
	for _, field := range o.Fields {
		if bytes.Equal(field.Name, name) {
			return field
		}
	}

	return nil
}

type UnionDefinition struct {
	Name []byte
	Types [][]byte
	Extentions []*UnionDefinition
	Directives []*Directive
}

func (u *UnionDefinition) GetFieldByName(name []byte) *FieldDefinition {
	for _, t := range u.Types {
		if bytes.Equal(t, name) {
			return &FieldDefinition{Name: name, Type: &FieldType{Name: name, Nullable: false}}
		}
	}

	return nil
}

func (u *UnionDefinition) HasType(name string) bool {
	for _, t := range u.Types {
		if string(t) == name {
			return true
		}
	}
	
	return false
}

type UnionDefinitions []*UnionDefinition

func (u UnionDefinitions) Has(name string) bool {
	for _, union := range u {
		if string(union.Name) == name {
			return true
		}
	}
	
	return false
}

func (u *UnionDefinition) TypeName() []byte {
	return u.Name
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

func (i *Indexes) GetTypeDefinition(name string) *TypeDefinition {
	return i.TypeIndex[name]
}

func (i *Indexes) GetInterfaceDefinition(name string) *InterfaceDefinition {
	return i.InterfaceIndex[name]
}

func (i *Indexes) GetUnionDefinition(name string) *UnionDefinition {
	return i.UnionIndex[name]
}

func (i *Indexes) GetImplementedType(id *InterfaceDefinition) []*TypeDefinition {
	res := make([]*TypeDefinition, 0)

	for _, t := range i.TypeIndex {
		for _, iface := range t.Interfaces {
			if bytes.Equal(iface.Name, id.Name) {
				res = append(res, t)
			}
		}
	}

	return res
}

type Schema struct {
	tokens Tokens
	Definition *SchemaDefinition
	Operations []*OperationDefinition
	Types []*TypeDefinition
	Enums []*EnumDefinition
	Unions []*UnionDefinition
	Interfaces []*InterfaceDefinition
	Directives DirectiveDefinitions
	Inputs []*InputDefinition
	Scalars []*ScalarDefinition

	Indexes *Indexes
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
		Indexes: &Indexes{
			TypeIndex: make(map[string]*TypeDefinition),
			OperationIndexes: operationIndexes,
			EnumIndex: make(map[string]*EnumDefinition),
			UnionIndex: make(map[string]*UnionDefinition),
			InterfaceIndex: make(map[string]*InterfaceDefinition),
			InputIndex: make(map[string]*InputDefinition),
		},
		Directives: newBuildInDirectives(),
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

		newSchema.Indexes.OperationIndexes[newOp.OperationType][string(newOp.Name)] = newOp
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
		newSchema.Indexes.TypeIndex[string(newType.Name)] = newType
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

		newSchema.Indexes.InterfaceIndex[string(newInterface.Name)] = newInterface
		newSchema.Interfaces = append(newSchema.Interfaces, newInterface)
	}

	return nil
}

func (s *Schema) digUnionDefinition(name string, unions UnionDefinitions) ([][]byte, error) {
	var res [][]byte

	if !unions.Has(name) {
		return nil, fmt.Errorf("union %s not found", name)
	}

	for _, u := range unions {
		if string(u.Name) == name {
			res = append(res, u.Types...)

			if len(u.Extentions) > 0 {
				types, err := s.digUnionDefinition(name, u.Extentions)
				if err != nil {
					return nil, err
				}

				res = append(res, types...)
			}
		}
	}

	return res, nil
}

func (s *Schema) mergeUnionDefinition(newSchema *Schema) error {
	for _, t := range s.Unions {
		newUnion := new(UnionDefinition)
		newUnion.Name = t.Name
		newUnion.Types = t.Types
		newUnion.Directives = t.Directives

		newSchema.Unions = append(newSchema.Unions, newUnion)

		if len(t.Extentions) > 0 {
			types, err := s.digUnionDefinition(string(newUnion.Name), t.Extentions)
			if err != nil {
				return err
			}

			for _, t := range types {
				if newUnion.HasType(string(t)) {
					return fmt.Errorf("type %s already exists in union %s", t, newUnion.Name)
				}
			}

			newSchema.Indexes.UnionIndex[string(newUnion.Name)] = newUnion
			newUnion.Types = append(newUnion.Types, types...)
		}
	}

	return nil
}

func (s *Schema) digEnumDefinition(name string, extentions EnumDefinitions) ([]*EnumElement, error) {
	if !extentions.Has(name) {
		return nil, fmt.Errorf("enum %s not found", name)
	}

	res := make([]*EnumElement, 0)
	for _, ext := range extentions {
		if string(ext.Name) == name {
			res = append(res, ext.Values...)
			
			if len(ext.Extentions) > 0 {
				elms, err := s.digEnumDefinition(name, ext.Extentions)
				if err != nil {
					return nil, err
				}

				res = append(res, elms...)
			}
		}
	}

	return res, nil
}

func (s *Schema) mergeEnumDefinition(newSchema *Schema) error {
	for _, enum := range s.Enums {
		newEnum := new(EnumDefinition)
		newEnum.Name = enum.Name
		newEnum.Directives = enum.Directives
		newEnum.Values = enum.Values

		newSchema.Enums = append(newSchema.Enums, newEnum)
		if len(enum.Extentions) > 0 {
			enumValues, err := s.digEnumDefinition(string(newEnum.Name), enum.Extentions)
			if err != nil {
				return err
			}
			
			newSchema.Indexes.EnumIndex[string(newEnum.Name)] = newEnum
			newEnum.Values = append(newEnum.Values, enumValues...)
		}
	}

	return nil
}

func (s *Schema) digInputDefinition(name string, exts []*InputDefinition) (FieldDefinitions, error) {
	res := make(FieldDefinitions, 0)

	for _, ext := range exts {
		if string(ext.Name) == name {
			res = append(res, ext.Fields...)

			if len(ext.Extentions) > 0 {
				fields, err := s.digInputDefinition(name, ext.Extentions)
				if err != nil {
					return nil, err
				}

				res = append(res, fields...)
			}
		}
	}

	return res, nil
}

func (s *Schema) mergeInputDefinition(newSchema *Schema) error {
	for _, input := range s.Inputs {
		newInput := new(InputDefinition)
		newInput.Name = input.Name
		newInput.Fields = input.Fields

		newFields := make(FieldDefinitions, 0)

		field, err := s.digInputDefinition(string(newInput.Name), input.Extentions)
		if err != nil {
			return err
		}

		newFields = append(newFields, field...)

		for _, field := range newInput.Fields {
			newField := newInput.Fields.Last(string(field.Name))
			if newField == nil {
				newFields = append(newFields, field)
			} else if !newFields.Has(string(newField.Name)) {
				newFields = append(newFields, newField)
			}
		}
		newInput.Fields = newFields

		newSchema.Indexes.InputIndex[string(newInput.Name)] = newInput
		newSchema.Inputs = append(newSchema.Inputs, newInput)
	}

	return nil
} 

func (s *Schema) Merge() (*Schema, error) {
	newSchema := new(Schema)
	newSchema.Definition = s.Definition
	newSchema.tokens = s.tokens
	newSchema.Indexes = s.Indexes
	
	if err := s.mergeOperation(newSchema); err != nil {
		return nil, err
	}

	if err := s.mergeTypeDefinition(newSchema); err != nil {
		return nil, err
	}

	if err := s.mergeInterfaceDefinition(newSchema); err != nil {
		return nil, err
	}

	if err := s.mergeUnionDefinition(newSchema); err != nil {
		return nil, err
	}

	if err := s.mergeEnumDefinition(newSchema); err != nil {
		return nil, err
	}

	if err := s.mergeInputDefinition(newSchema); err != nil {
		return nil, err
	}

	return newSchema, nil
}

func (s *Schema) GetQuery() *OperationDefinition {
	for _, op := range s.Operations {
		if op.OperationType.isQuery() {
			return op
		}
	}

	return nil
}

func (s *Schema) GetMutation() *OperationDefinition {
	for _, op := range s.Operations {
		if op.OperationType.isMutation() {
			return op
		}
	}

	return nil
}

func (s *Schema) GetSubscription() *OperationDefinition {
	for _, op := range s.Operations {
		if op.OperationType.isSubscription() {
			return op
		}
	}

	return nil
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

	return nil
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