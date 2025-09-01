package schema

import (
	"bytes"
	"fmt"
	"reflect"
)

type OperationType string

func (o OperationType) IsQuery() bool {
	return o == QueryOperation
}

func (o OperationType) IsMutation() bool {
	return o == MutationOperation
}

func (o OperationType) IsSubscription() bool {
	return o == SubscriptionOperation
}

const (
	QueryOperation        OperationType = "query"
	MutationOperation     OperationType = "mutation"
	SubscriptionOperation OperationType = "subscription"
)

type CompositeType interface {
	GetFieldByName(name []byte) *FieldDefinition
	TypeName() []byte
}

type TypeDefinition struct {
	Name              []byte
	Fields            FieldDefinitions
	Required          map[*FieldDefinition]struct{}
	PrimitiveTypeName []byte
	Interfaces        [][]byte
	Directives        []*Directive
	Extentions        []*TypeDefinition
}

func (t *TypeDefinition) IsDefinition() bool {
	return true
}

type TypeDefinitions []*TypeDefinition

func (t TypeDefinitions) WithoutMetaDefinition() TypeDefinitions {
	res := make(TypeDefinitions, 0, len(t))

	for _, td := range t {
		if !td.IsIntrospection() {
			res = append(res, td)
		}
	}

	return res
}

func (t *TypeDefinition) IsPrimitive() bool {
	return bytes.Equal(t.Name, []byte("String")) || bytes.Equal(t.Name, []byte("Int")) || bytes.Equal(t.Name, []byte("Float")) || bytes.Equal(t.Name, []byte("Boolean")) || bytes.Equal(t.Name, []byte("ID"))
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

func (t *TypeDefinition) IsIntrospection() bool {
	return bytes.Equal(t.Name, []byte("__Schema")) ||
		bytes.Equal(t.Name, []byte("__Type")) ||
		bytes.Equal(t.Name, []byte("__Field")) ||
		bytes.Equal(t.Name, []byte("__InputValue")) ||
		bytes.Equal(t.Name, []byte("__EnumValue")) ||
		bytes.Equal(t.Name, []byte("__Directive")) ||
		bytes.Equal(t.Name, []byte("__DirectiveLocation")) ||
		bytes.Equal(t.Name, []byte("__TypeKind"))
}

type FieldType struct {
	Name     []byte
	Nullable bool
	IsList   bool
	ListType *FieldType
}

func (f *FieldType) IsObject() bool {
	return !bytes.Equal(f.Name, []byte(""))
}

func (f *FieldType) IsBoolean() bool {
	return bytes.Equal(f.Name, []byte("Boolean"))
}

func (f *FieldType) IsString() bool {
	return bytes.Equal(f.Name, []byte("String"))
}

func (f *FieldType) IsInt() bool {
	return bytes.Equal(f.Name, []byte("Int"))
}

func (f *FieldType) IsFloat() bool {
	return bytes.Equal(f.Name, []byte("Float"))
}

func (f *FieldType) IsID() bool {
	return bytes.Equal(f.Name, []byte("ID"))
}

func (f *FieldType) GetRootType() *FieldType {
	if f.IsList {
		return f.ListType.GetRootType()
	}

	return f
}

func (f *FieldType) IsPrimitive() bool {
	return bytes.Equal(f.Name, []byte("String")) || bytes.Equal(f.Name, []byte("Int")) || bytes.Equal(f.Name, []byte("Float")) || bytes.Equal(f.Name, []byte("Boolean")) || bytes.Equal(f.Name, []byte("ID"))
}

func (f *FieldType) GetNestFieldType(nestCount, currentNestCount int) *FieldType {
	if nestCount == currentNestCount {
		return f
	}

	if f.IsList {
		return f.ListType.GetNestFieldType(nestCount, currentNestCount+1)
	}

	return nil
}

type OperationDefinition struct {
	OperationType OperationType
	Name          []byte
	Fields        FieldDefinitions
	Extentions    []*OperationDefinition
}

func (o *OperationDefinition) IsDefinition() bool {
	return true
}

func (f FieldDefinitions) HasDeprecatedDirective() bool {
	for _, field := range f {
		for _, directive := range field.Directives {
			if bytes.Equal(directive.Name, []byte("deprecated")) {
				return true
			}
		}
	}

	return false
}

func (o *OperationDefinition) GetFieldByName(name []byte) *FieldDefinition {
	for _, field := range o.Fields {
		if bytes.Equal(field.Name, name) {
			return field
		}
	}

	return nil
}

type DefinitionType interface {
	*TypeDefinition | *OperationDefinition | *EnumDefinition | *UnionDefinition | *InterfaceDefinition | *InputDefinition | *ScalarDefinition | *DirectiveDefinition
}

type Indexes struct {
	TypeIndex        map[string]*TypeDefinition
	EnumIndex        map[string]*EnumDefinition
	UnionIndex       map[string]*UnionDefinition
	InterfaceIndex   map[string]*InterfaceDefinition
	InputIndex       map[string]*InputDefinition
	OperationIndexes map[OperationType]map[string]*OperationDefinition
	ScalarIndex      map[string]*ScalarDefinition
	DirectiveIndex   map[string]*DirectiveDefinition
	ExtendIndex      map[string]ExtendDefinition
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
			if bytes.Equal(iface, id.Name) {
				res = append(res, t)
			}
		}
	}

	return res
}

type Schema struct {
	Tokens     Tokens
	Definition *SchemaDefinition
	Operations []*OperationDefinition
	Types      []*TypeDefinition
	Enums      []*EnumDefinition
	Unions     []*UnionDefinition
	Interfaces []*InterfaceDefinition
	Directives DirectiveDefinitions
	Inputs     []*InputDefinition
	Scalars    []*ScalarDefinition
	Extends    []ExtendDefinition

	Indexes *Indexes
}

func NewSchema(tokens Tokens) *Schema {
	operationIndexes := make(map[OperationType]map[string]*OperationDefinition)
	for _, t := range []OperationType{QueryOperation, MutationOperation, SubscriptionOperation} {
		operationIndexes[t] = make(map[string]*OperationDefinition)
	}

	s := &Schema{
		Tokens: tokens,
		Definition: &SchemaDefinition{
			Query:        []byte("Query"),
			Mutation:     []byte("Mutation"),
			Subscription: []byte("Subscription"),
		},
		Indexes: &Indexes{
			TypeIndex:        make(map[string]*TypeDefinition),
			OperationIndexes: operationIndexes,
			EnumIndex:        make(map[string]*EnumDefinition),
			UnionIndex:       make(map[string]*UnionDefinition),
			InterfaceIndex:   make(map[string]*InterfaceDefinition),
			InputIndex:       make(map[string]*InputDefinition),
			ScalarIndex:      make(map[string]*ScalarDefinition),
			ExtendIndex:      make(map[string]ExtendDefinition),
		},
		Directives: NewBuildInDirectives(),
	}

	return s
}

func (s *Schema) extendOperationFields(ops []*OperationDefinition) FieldDefinitions {
	res := make(FieldDefinitions, 0)

	for _, op := range ops {
		res = append(res, op.Fields...)
	}

	return res
}

func (s *Schema) mergeOperation(newSchema *Schema) error {
	for _, t := range s.Operations {
		newOp := new(OperationDefinition)
		newOp.OperationType = t.OperationType
		newOp.Name = t.Name
		newOp.Fields = t.Fields

		extendDefinitions := getOperationDefinitionsFromExtendDefinitions(t.OperationType, s.Extends)
		field := s.extendOperationFields(extendDefinitions)

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

func getOperationDefinitionsFromExtendDefinitions(opType OperationType, extendDefinitions []ExtendDefinition) []*OperationDefinition {
	ret := make([]*OperationDefinition, 0, len(extendDefinitions))
	for _, ext := range extendDefinitions {
		if opDef, ok := ext.(*OperationDefinition); ok && opDef.OperationType == opType {
			ret = append(ret, opDef)
		}
	}

	return ret
}

func (s *Schema) extendTypeDefinitionFields(types []*TypeDefinition) FieldDefinitions {
	res := make(FieldDefinitions, 0)

	for _, t := range types {
		res = append(res, t.Fields...)
	}

	return res
}

func (s *Schema) mergeTypeDefinition(newSchema *Schema) error {
	for _, t := range s.Types {
		newType := new(TypeDefinition)
		newType.Name = t.Name
		newType.Fields = t.Fields
		newType.Interfaces = t.Interfaces
		newType.Directives = t.Directives
		newType.PrimitiveTypeName = t.PrimitiveTypeName

		newFields := make(FieldDefinitions, 0)

		field := s.extendTypeDefinitionFields(getTypeDefinitionsFromExtendDefinitions(s.Extends, string(newType.Name)))

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

func getTypeDefinitionsFromExtendDefinitions(extendDefinitions []ExtendDefinition, name string) []*TypeDefinition {
	ret := make([]*TypeDefinition, 0, len(extendDefinitions))
	for _, ext := range extendDefinitions {
		if typeDef, ok := ext.(*TypeDefinition); ok && string(typeDef.Name) == name {
			ret = append(ret, typeDef)
		}
	}

	return ret
}

func (s *Schema) extendInterfaceDefinition(interfaces []*InterfaceDefinition) FieldDefinitions {
	res := make(FieldDefinitions, 0)

	for _, t := range interfaces {
		res = append(res, t.Fields...)
	}

	return res
}

func (s *Schema) mergeInterfaceDefinition(newSchema *Schema) error {
	for _, t := range s.Interfaces {
		newInterface := new(InterfaceDefinition)
		newInterface.Name = t.Name
		newInterface.Fields = t.Fields
		newInterface.Directives = t.Directives

		newFields := make(FieldDefinitions, 0)

		field := s.extendInterfaceDefinition(getInterfaceDefinitionsFromExtendDefinitions(s.Extends, string(newInterface.Name)))
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

func getInterfaceDefinitionsFromExtendDefinitions(extendDefinitions []ExtendDefinition, name string) []*InterfaceDefinition {
	ret := make([]*InterfaceDefinition, 0, len(extendDefinitions))
	for _, ext := range extendDefinitions {
		if ifaceDef, ok := ext.(*InterfaceDefinition); ok && string(ifaceDef.Name) == name {
			ret = append(ret, ifaceDef)
		}
	}

	return ret
}

func (s *Schema) extendUnionDefinition(unions UnionDefinitions) ([][]byte, error) {
	var res [][]byte

	for _, u := range unions {
		res = append(res, u.Types...)
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

		types, err := s.extendUnionDefinition(getUnionDefinitionFromExtendDefinition(s.Extends, string(newUnion.Name)))
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

	return nil
}

func getUnionDefinitionFromExtendDefinition(extendDefinitions []ExtendDefinition, name string) []*UnionDefinition {
	ret := make([]*UnionDefinition, 0, len(extendDefinitions))
	for _, ext := range extendDefinitions {
		if unionDef, ok := ext.(*UnionDefinition); ok && string(unionDef.Name) == name {
			ret = append(ret, unionDef)
		}
	}
	return ret
}

func (s *Schema) extendEnumDefinition(extentions EnumDefinitions) []*EnumElement {
	res := make([]*EnumElement, 0)

	for _, ext := range extentions {
		res = append(res, ext.Values...)
	}

	return res
}

func (s *Schema) mergeEnumDefinition(newSchema *Schema) error {
	for _, enum := range s.Enums {
		newEnum := new(EnumDefinition)
		newEnum.Name = enum.Name
		newEnum.Directives = enum.Directives
		newEnum.Values = enum.Values
		newEnum.Type = enum.Type

		newSchema.Enums = append(newSchema.Enums, newEnum)
		enumValues := s.extendEnumDefinition(getEnumDefinitionFromExtendDefinitions(s.Extends, string(newEnum.Name)))

		newSchema.Indexes.EnumIndex[string(newEnum.Name)] = newEnum
		newEnum.Values = append(newEnum.Values, enumValues...)
	}

	return nil
}

func getEnumDefinitionFromExtendDefinitions(extendDefinitions []ExtendDefinition, name string) []*EnumDefinition {
	ret := make([]*EnumDefinition, 0, len(extendDefinitions))
	for _, ext := range extendDefinitions {
		if enumDef, ok := ext.(*EnumDefinition); ok && string(enumDef.Name) == name {
			ret = append(ret, enumDef)
		}
	}

	return ret
}

func (s *Schema) extendInputDefinition(exts []*InputDefinition) (FieldDefinitions, error) {
	res := make(FieldDefinitions, 0)

	for _, ext := range exts {
		res = append(res, ext.Fields...)
	}

	return res, nil
}

func (s *Schema) mergeInputDefinition(newSchema *Schema) error {
	for _, input := range s.Inputs {
		newInput := new(InputDefinition)
		newInput.Name = input.Name
		newInput.Fields = input.Fields

		newFields := make(FieldDefinitions, 0)

		field, err := s.extendInputDefinition(getInputDefinitionFromExtendDefinitions(s.Extends, string(newInput.Name)))
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

func getInputDefinitionFromExtendDefinitions(extendDefinitions []ExtendDefinition, name string) []*InputDefinition {
	ret := make([]*InputDefinition, 0, len(extendDefinitions))
	for _, ext := range extendDefinitions {
		if inputDef, ok := ext.(*InputDefinition); ok && string(inputDef.Name) == name {
			ret = append(ret, inputDef)
		}
	}

	return ret
}

func (s *Schema) Merge() (*Schema, error) {
	newSchema := new(Schema)
	newSchema.Definition = s.Definition
	newSchema.Tokens = s.Tokens
	newSchema.Indexes = s.Indexes
	newSchema.Directives = s.Directives
	newSchema.Scalars = s.Scalars

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

	newSchema = WithTypeIntrospection(newSchema)
	newSchema = WithBuiltin(newSchema)

	return newSchema, nil
}

func (s *Schema) GetQuery() *OperationDefinition {
	for _, op := range s.Operations {
		if op.OperationType.IsQuery() {
			return op
		}
	}

	return nil
}

func (s *Schema) GetMutation() *OperationDefinition {
	for _, op := range s.Operations {
		if op.OperationType.IsMutation() {
			return op
		}
	}

	return nil
}

func (s *Schema) GetSubscription() *OperationDefinition {
	for _, op := range s.Operations {
		if op.OperationType.IsSubscription() {
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
	case *ScalarDefinition:
		indexes.ScalarIndex[string(d.Name)] = d
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
	case *ScalarDefinition:
		return any(indexes.ScalarIndex[key]).(T)
	case *DirectiveDefinition:
		return any(indexes.DirectiveIndex[key]).(T)
	}

	return nil
}

type ScalarDefinition struct {
	Name       []byte
	Directives []*Directive
	Extentions []*ScalarDefinition
}

func (s *ScalarDefinition) IsDefinition() bool {
	return true
}

type SchemaDefinition struct {
	Query        []byte
	Mutation     []byte
	Subscription []byte
	Extentions   []*SchemaDefinition

	Directives []*Directive
}

func (s *SchemaDefinition) IsDefinition() bool {
	return true
}

type ExtendDefinition interface {
	IsDefinition() bool
}
