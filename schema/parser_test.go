package schema_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lkeix/gg-parser/schema"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		isSkip  bool
		want    *schema.Schema
		wantErr error
	}{
		{
			name: "Parse standard schema definition",
			input: []byte(`schema {
				query: Query
				mutation: Mutation
				subscription: Subscription
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
			},
		},
		{
			name: "Parse custom schema definition",
			input: []byte(`schema {
				query: RootQuery
				mutation: RootMutation
				subscription: RootSubscription
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("RootQuery"),
					Mutation: []byte("RootMutation"),
					Subscription: []byte("RootSubscription"),
				},
			},
		},
		{
			name: "Parse custom lack optional schema definition",
			input: []byte(`schema {
				query: RootQuery
				mutation: RootMutation
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("RootQuery"),
					Mutation: []byte("RootMutation"),
				},
			},
		},
		{
			name: "simple type with scalar fields",
			input: []byte(`type User {
				id: ID!
				name: String!
				age: Int
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Types: []*schema.TypeDefinition{
					{
						Name: []byte("User"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("id"),
								Type: &schema.FieldType{
									Name:     []byte("ID"),
									Nullable: false,
									IsList:   false,
								},
							},
							{
								Name: []byte("name"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
							},
							{
								Name: []byte("age"),
								Type: &schema.FieldType{
									Name:     []byte("Int"),
									Nullable: true,
									IsList:   false,
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		}, {
			name: "type with list fields",
			input: []byte(`type User {
				friends: [User!]!
				posts: [Post]
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Types: []*schema.TypeDefinition{
					{
						Name: []byte("User"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("friends"),
								Type: &schema.FieldType{
									Name:     nil,
									Nullable: false,
									IsList:   true,
									ListType: &schema.FieldType{
										Name:     []byte("User"),
										Nullable: false,
										IsList:   false,
									},
								},
							},
							{
								Name: []byte("posts"),
								Type: &schema.FieldType{
									Name:     nil,
									Nullable: true,
									IsList:   true,
									ListType: &schema.FieldType{
										Name:     []byte("Post"),
										Nullable: true,
										IsList:   false,
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		}, {
			name: "type with deeply nested list",
			input: []byte(`type Data {
				matrix: [[[Int!]!]!]!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Types: []*schema.TypeDefinition{
					{
						Name: []byte("Data"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("matrix"),
								Type: &schema.FieldType{
									IsList:   true,
									Nullable: false,
									ListType: &schema.FieldType{
										IsList:   true,
										Nullable: false,
										ListType: &schema.FieldType{
											IsList:   true,
											Nullable: false,
											ListType: &schema.FieldType{
												Name:     []byte("Int"),
												Nullable: false,
												IsList:   false,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "type with nullable nested list",
			input: []byte(`type Example {
				data: [[String]]
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Types: []*schema.TypeDefinition{
					{
						Name: []byte("Example"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("data"),
								Type: &schema.FieldType{
									IsList:   true,
									Nullable: true,
									ListType: &schema.FieldType{
										IsList:   true,
										Nullable: true,
										ListType: &schema.FieldType{
											Name:     []byte("String"),
											Nullable: true,
											IsList:   false,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name:  "simple union type",
			input: []byte(`union SearchResult = User | Post`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Unions: []*schema.UnionDefinition{
					{
						Name: []byte("SearchResult"),
						Types: [][]byte{
							[]byte("User"),
							[]byte("Post"),
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "multi-line union type",
			input: []byte(`union SearchResult = User
				| Post
				| Comment`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Unions: []*schema.UnionDefinition{
					{
						Name: []byte("SearchResult"),
						Types: [][]byte{
							[]byte("User"),
							[]byte("Post"),
							[]byte("Comment"),
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name:    "invalid union type",
			input:   []byte(`union SearchResult = User |`),
			want:    nil,
			wantErr: errors.New("unexpected end of input"),
		},
		{
			name:    "empty union type",
			input:   []byte(`union SearchResult =`),
			want:    nil,
			wantErr: errors.New("unexpected end of input"),
		},
		{
			name: "simple input type",
			input: []byte(`input Filter {
				field: String!
				value: String!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Inputs: []*schema.InputDefinition{
					{
						Name: []byte("Filter"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("field"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
							},
							{
								Name: []byte("value"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "input type with a simple default value",
			input: []byte(`input Filter {
				field: String! = "name"
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Inputs: []*schema.InputDefinition{
					{
						Name: []byte("Filter"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("field"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
								Default: []byte(`"name"`),
							},
						},
					},
				},
			},
		},
		{
			name: "input type with default values",
			input: []byte(`input Filter {
				field: String! = "name"
				value: String! = "John Doe"
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Inputs: []*schema.InputDefinition{
					{
						Name: []byte("Filter"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("field"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
								Default: []byte(`"name"`),
							},
							{
								Name: []byte("value"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
								Default: []byte(`"John Doe"`),
							},
						},
					},
				},
			},
		},
		{
			name: "simple interface type",
			input: []byte(`interface Node {
				id: ID!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Interfaces: []*schema.InterfaceDefinition{
					{
						Name: []byte("Node"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("id"),
								Type: &schema.FieldType{
									Name:     []byte("ID"),
									Nullable: false,
									IsList:   false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "interface with nested lists",
			input: []byte(`interface Nested {
				items: [[Item!]!]!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Interfaces: []*schema.InterfaceDefinition{
					{
						Name: []byte("Nested"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("items"),
								Type: &schema.FieldType{
									IsList: true,
									ListType: &schema.FieldType{
										IsList: true,
										ListType: &schema.FieldType{
											Name:     []byte("Item"),
											Nullable: false,
										},
										Nullable: false,
									},
									Nullable: false,
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name:  "interface without fields",
			input: []byte(`interface Empty {}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Interfaces: []*schema.InterfaceDefinition{
					{
						Name:   []byte("Empty"),
						Fields: []*schema.FieldDefinition{},
					},
				},
			},
			wantErr: nil,
		},
		{
			name:    "invalid interface (missing opening curly brace)",
			input:   []byte(`interface Node id: ID!}`),
			want:    nil,
			wantErr: errors.New("expected '{' but got id"),
		},
		{
			name: "invalid interface (missing closing curly brace)",
			input: []byte(`interface Node {
				id: ID!`),
			want:    nil,
			wantErr: errors.New("unexpected end of input"),
		},
		{
			name: "simple non-argument Query operation",
			input: []byte(`type Query {
				user: User
				users: [User]
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name:      []byte("user"),
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
							{
								Name:      []byte("users"),
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     nil,
									Nullable: true,
									IsList:   true,
									ListType: &schema.FieldType{
										Name:     []byte("User"),
										Nullable: true,
									},
								},
								Directives: []*schema.Directive{},
							},
						},
					},
				},
			},
		},
		{
			name: "simple Query operation",
			input: []byte(`type Query {
				user(id: ID!): User
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("user"),
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("id"),
										Type: &schema.FieldType{
											Name:     []byte("ID"),
											Nullable: false,
											IsList:   false,
										},
									},
								},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
						},
					},
				},
			},
		},
		{
			name: "simple default argument Query operation",
			input: []byte(`type Query {
				users(offset: INT = 1): [User]
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("users"),
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("offset"),
										Type: &schema.FieldType{
											Name:     []byte("INT"),
											Nullable: true,
											IsList:   false,
										},
										Default: []byte("1"),
									},
								},
								Type: &schema.FieldType{
									Name:     nil,
									Nullable: true,
									IsList:   true,
									ListType: &schema.FieldType{
										Name:     []byte("User"),
										Nullable: true,
									},
								},
								Directives: []*schema.Directive{},
							},
						},
					},
				},
			},
		},
		// TODO: support nested arguments with default values
		{
			name: "deep nested argument Query operation",
			input: []byte(`type Query {
				getUser(filter: [[FilterInput!]!]!): User
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("getUser"),
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("filter"),
										Type: &schema.FieldType{
											IsList:   true,
											Nullable: false,
											ListType: &schema.FieldType{
												IsList:   true,
												Nullable: false,
												ListType: &schema.FieldType{
													Name:     []byte("FilterInput"),
													Nullable: false,
												},
											},
										},
									},
								},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
						},
					},
				},
			},
		},
		{
			name: "deep nested argument with default value Query operation",
			input: []byte(`type Query {
				getUser(filter: [[FilterInput!]!]! = [[{field: "name", value: "John Doe"}]]): User
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("getUser"),
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("filter"),
										Type: &schema.FieldType{
											IsList:   true,
											Nullable: false,
											ListType: &schema.FieldType{
												IsList:   true,
												Nullable: false,
												ListType: &schema.FieldType{
													Name:     []byte("FilterInput"),
													Nullable: false,
												},
											},
										},
										Default: []byte(`[[{field: "name", value: "John Doe"}]]`),
									},
								},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
						},
					},
				},
			},
		},
		{
			name: "simple Mutation operation",
			input: []byte(`type Mutation {
				createUser(mail: String!, name: String!): Boolean
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.MutationOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("createUser"),
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("mail"),
										Type: &schema.FieldType{
											Name:     []byte("String"),
											Nullable: false,
											IsList:   false,
										},
									},
									{
										Name: []byte("name"),
										Type: &schema.FieldType{
											Name:     []byte("String"),
											Nullable: false,
											IsList:   false,
										},
									},
								},
								Type: &schema.FieldType{
									Name:     []byte("Boolean"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
						},
					},
				},
			},
		},
		{
			name: "non-argument Mutation operation",
			input: []byte(`type Mutation {
				addAnonymousUser: Boolean
				deleteUsers: Boolean
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.MutationOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name:      []byte("addAnonymousUser"),
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     []byte("Boolean"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
							{
								Name:      []byte("deleteUsers"),
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     []byte("Boolean"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
						},
					},
				},
			},
		},
		{
			name: "default argument Mutation operation",
			input: []byte(`type Mutation {
				updateUser(id: ID!, name: String = "John Doe"): Boolean
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.MutationOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("updateUser"),
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("id"),
										Type: &schema.FieldType{
											Name:     []byte("ID"),
											Nullable: false,
											IsList:   false,
										},
									},
									{
										Name: []byte("name"),
										Type: &schema.FieldType{
											Name:     []byte("String"),
											Nullable: true,
											IsList:   false,
										},
										Default: []byte(`"John Doe"`),
									},
								},
								Type: &schema.FieldType{
									Name:     []byte("Boolean"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
						},
					},
				},
			},
		},
		{
			name: "non-argument Subscription operation",
			input: []byte(`type Subscription {
				newUser: User
				newUsers: [User]
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.SubscriptionOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name:      []byte("newUser"),
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
							{
								Name:      []byte("newUsers"),
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     nil,
									Nullable: true,
									IsList:   true,
									ListType: &schema.FieldType{
										Name:     []byte("User"),
										Nullable: true,
									},
								},
								Directives: []*schema.Directive{},
							},
						},
					},
				},
			},
		},
		{
			name: "simple Subscription operation",
			input: []byte(`type Subscription {
				userFollowed(followerId: ID!, followeeId: ID!): Notification
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.SubscriptionOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("userFollowed"),
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("followerId"),
										Type: &schema.FieldType{
											Name:     []byte("ID"),
											Nullable: false,
											IsList:   false,
										},
									},
									{
										Name: []byte("followeeId"),
										Type: &schema.FieldType{
											Name:     []byte("ID"),
											Nullable: false,
											IsList:   false,
										},
									},
								},
								Type: &schema.FieldType{
									Name:     []byte("Notification"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
						},
					},
				},
			},
		},
		{
			name: "simple default argument Subscription operation",
			input: []byte(`type Subscription {
				userFollowed(followerId: ID!, followeeId: ID!, notify: Boolean = true): Notification
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.SubscriptionOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("userFollowed"),
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("followerId"),
										Type: &schema.FieldType{
											Name:     []byte("ID"),
											Nullable: false,
											IsList:   false,
										},
									},
									{
										Name: []byte("followeeId"),
										Type: &schema.FieldType{
											Name:     []byte("ID"),
											Nullable: false,
											IsList:   false,
										},
									},
									{
										Name: []byte("notify"),
										Type: &schema.FieldType{
											Name:     []byte("Boolean"),
											Nullable: true,
											IsList:   false,
										},
										Default: []byte("true"),
									},
								},
								Type: &schema.FieldType{
									Name:     []byte("Notification"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
						},
					},
				},
			},
		},
		// directives
		{
			name: "field with a simple directive",
			input: []byte(`type Query {
					user: User @deprecated
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name:      []byte("user"),
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{
									{
										Name:      []byte("deprecated"),
										Arguments: nil,
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "field with a directive that has arguments",
			input: []byte(`type Query {
					user: User @include(if: true)
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name:      []byte("user"),
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{
									{
										Name: []byte("include"),
										Arguments: []*schema.DirectiveArgument{
											{
												Name: []byte("if"),
												Value: []byte("true"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "field with multiple directives",
			input: []byte(`type Query {
			user: User @deprecated @include(if: true)
	}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name:      []byte("user"),
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{
									{
										Name:      []byte("deprecated"),
										Arguments: nil,
									},
									{
										Name: []byte("include"),
										Arguments: []*schema.DirectiveArgument{
											{
												Name: []byte("if"),
												Value: []byte("true"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "simple enum type",
			input: []byte(`enum Role {
				ADMIN
				USER
				GUEST
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Enums: []*schema.EnumDefinition{
					{
						Name: []byte("Role"),
						Values: [][]byte{
							[]byte("ADMIN"),
							[]byte("USER"),
							[]byte("GUEST"),
						},
					},
				},
			},
		},
		{
			name: "Parse extend schema definition",
			input: []byte(`schema {
				query: Query
				mutation: Mutation
				subscription: Subscription
			}
			
			extend schema {
			  query: RootQuery
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
					Extentions: []*schema.SchemaDefinition{
						{
							Query: []byte("RootQuery"),
						},
					},
				},
			},
		},
		{
			name: "Parse extend type definition",
			input: []byte(`type User {
				id: ID!
			}
			
			extend type User {
			  created_at: String
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Types: []*schema.TypeDefinition{
					{
						Name: []byte("User"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("id"),
								Type: &schema.FieldType{
									Name: []byte("ID"),
									Nullable: false,
									IsList: false,
								},
							},
						},
						Extentions: []*schema.TypeDefinition{
							{
								Name: []byte("User"),
								Fields: []*schema.FieldDefinition{
									{
										Name: []byte("created_at"),
										Type: &schema.FieldType{
											Name: []byte("String"),
											Nullable: true,
											IsList: false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Parse extend query definition",
			input: []byte(`type Query {
				user(id: ID!): User!
			}
				
			extend type Query {
				user(id: ID!, isActive: Boolean! = false): User!
				users(): [User]!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Name: nil,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("user"),
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte(`id`),
										Default: nil,
										Type: &schema.FieldType{
											Name: []byte(`ID`),
											Nullable: false,
											IsList: false,
										},
									},
								},
								Type: &schema.FieldType{
									Name: []byte(`User`),
									Nullable: false,
									IsList: false,
								},
								Directives: []*schema.Directive{},
							},
						},
						Extentions: []*schema.OperationDefinition{
							{
								OperationType: schema.QueryOperation,
								Name: nil,
								Fields: []*schema.FieldDefinition{
									{
										Name: []byte("user"),
										Arguments: []*schema.ArgumentDefinition{
											{
												Name: []byte(`id`),
												Default: nil,
												Type: &schema.FieldType{
													Name: []byte(`ID`),
													Nullable: false,
													IsList: false,
												},
											},
											{
												Name: []byte(`isActive`),
												Default: []byte(`false`),
												Type: &schema.FieldType{
													Name: []byte(`Boolean`),
													Nullable: false,
													IsList: false,
												},
											},
										},
										Type: &schema.FieldType{
											Name: []byte(`User`),
											Nullable: false,
											IsList: false,
										},
										Directives: []*schema.Directive{},
									},
									{
										Name: []byte("users"),
										Arguments: []*schema.ArgumentDefinition{
										},
										Type: &schema.FieldType{
											Name: nil,
											Nullable: false,
											IsList: true,
											ListType: &schema.FieldType{
												Name: []byte(`User`),
												Nullable: true,
												IsList: false,
											},
										},
										Directives: []*schema.Directive{},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Parse extend mutation definition",
			input: []byte(`type Mutation {
				createUser(name: String!): User!
			}
				
			extend type Mutation {
				createUser(name: String!, isActive: Boolean! = true): User!
				updateUser(id: ID!, name: String): User!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.MutationOperation,
						Name: nil,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("createUser"),
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte(`name`),
										Default: nil,
										Type: &schema.FieldType{
											Name: []byte(`String`),
											Nullable: false,
											IsList: false,
										},
									},
								},
								Type: &schema.FieldType{
									Name: []byte(`User`),
									Nullable: false,
									IsList: false,
								},
								Directives: []*schema.Directive{},
							},
						},
						Extentions: []*schema.OperationDefinition{
							{
								OperationType: schema.MutationOperation,
								Name: nil,
								Fields: []*schema.FieldDefinition{
									{
										Name: []byte("createUser"),
										Arguments: []*schema.ArgumentDefinition{
											{
												Name: []byte(`name`),
												Default: nil,
												Type: &schema.FieldType{
													Name: []byte(`String`),
													Nullable: false,
													IsList: false,
												},
											},
											{
												Name: []byte(`isActive`),
												Default: []byte(`true`),
												Type: &schema.FieldType{
													Name: []byte(`Boolean`),
													Nullable: false,
													IsList: false,
												},
											},
										},
										Type: &schema.FieldType{
											Name: []byte(`User`),
											Nullable: false,
											IsList: false,
										},
										Directives: []*schema.Directive{},
									},
									{
										Name: []byte("updateUser"),
										Arguments: []*schema.ArgumentDefinition{
											{
												Name: []byte(`id`),
												Default: nil,
												Type: &schema.FieldType{
													Name: []byte(`ID`),
													Nullable: false,
													IsList: false,
												},
											},
											{
												Name: []byte(`name`),
												Default: nil,
												Type: &schema.FieldType{
													Name: []byte(`String`),
													Nullable: true,
													IsList: false,
												},
											},
										},
										Type: &schema.FieldType{
											Name: []byte(`User`),
											Nullable: false,
											IsList: false,
										},
										Directives: []*schema.Directive{},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Parse extend subscription definition",
			input: []byte(`type Subscription {
				userCreated: User!
			}
				
			extend type Subscription {
				userCreated(isActive: Boolean! = false): User!
				userUpdated(id: ID!): User!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.SubscriptionOperation,
						Name: nil,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("userCreated"),
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name: []byte(`User`),
									Nullable: false,
									IsList: false,
								},
								Directives: []*schema.Directive{},
							},
						},
						Extentions: []*schema.OperationDefinition{
							{
								OperationType: schema.SubscriptionOperation,
								Name: nil,
								Fields: []*schema.FieldDefinition{
									{
										Name: []byte("userCreated"),
										Arguments: []*schema.ArgumentDefinition{
											{
												Name: []byte(`isActive`),
												Default: []byte(`false`),
												Type: &schema.FieldType{
													Name: []byte(`Boolean`),
													Nullable: false,
													IsList: false,
												},
											},
										},
										Type: &schema.FieldType{
											Name: []byte(`User`),
											Nullable: false,
											IsList: false,
										},
										Directives: []*schema.Directive{},
									},
									{
										Name: []byte("userUpdated"),
										Arguments: []*schema.ArgumentDefinition{
											{
												Name: []byte(`id`),
												Default: nil,
												Type: &schema.FieldType{
													Name: []byte(`ID`),
													Nullable: false,
													IsList: false,
												},
											},
										},
										Type: &schema.FieldType{
											Name: []byte(`User`),
											Nullable: false,
											IsList: false,
										},
										Directives: []*schema.Directive{},
									},
								},
							},
						},
					},
				},
			},
		},				
	}

	ignores := []any{
		schema.Schema{},
		schema.TypeDefinition{},
		schema.InputDefinition{},
		schema.Indexes{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isSkip {
				t.Skip()
			}
			lexer := schema.NewLexer()
			parser := schema.NewParser(lexer)
			got, err := parser.Parse(tt.input)
			if tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && err != nil {
				t.Errorf("Parse() error %v", err)
			}

			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(ignores...)); diff != "" {
				t.Errorf("Parse() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
