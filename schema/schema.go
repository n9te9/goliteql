package schema

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
}

type EnumDefinition struct {
	Name []byte
	Values [][]byte
}

type UnionDefinition struct {
	Name []byte
	Types [][]byte
}

type InterfaceDefinition struct {
	Name []byte
	Fields []*FieldDefinition
}

type DirectiveDefinition struct {
	Name []byte
	Arguments []*FieldDefinition
	Locations [][]byte
}

type InputDefinition struct {
	Name []byte
	Fields []*FieldDefinition
	tokens Tokens
}

type Schema struct {
	tokens Tokens
	Operations []*OperationDefinition
	Types []*TypeDefinition
	Enums []*EnumDefinition
	Unions []*UnionDefinition
	Interfaces []*InterfaceDefinition
	Directives []*DirectiveDefinition
	Inputs []*InputDefinition
}