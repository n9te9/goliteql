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

type Directive struct {
	Name []byte
	Arguments []*DirectiveArgument
	Locations [][]byte
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
	Locations [][]byte
}

type InputDefinition struct {
	Name []byte
	Fields []*FieldDefinition
	tokens Tokens
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
}

type SchemaDefinition struct {
	Query []byte
	Mutation []byte
	Subscription []byte
}