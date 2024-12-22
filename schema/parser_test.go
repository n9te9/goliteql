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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("RootQuery"),
					Mutation:     []byte("RootMutation"),
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
					Query:    []byte("RootQuery"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
							},
							{
								Name: []byte("name"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
							{
								Name: []byte("age"),
								Type: &schema.FieldType{
									Name:     []byte("Int"),
									Nullable: true,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
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
								Directives: []*schema.Directive{},
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
							},
							{
								Name: []byte("value"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
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
			name: "input type with a simple default value",
			input: []byte(`input Filter {
				field: String! = "name"
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
							},
							{
								Name: []byte("value"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
								Default: []byte(`"John Doe"`),
								Directives: []*schema.Directive{},
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
							},
						},
						Extentions: []*schema.TypeDefinition{
							{
								Name: []byte("User"),
								Fields: []*schema.FieldDefinition{
									{
										Name: []byte("created_at"),
										Type: &schema.FieldType{
											Name:     []byte("String"),
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
										Name:    []byte(`id`),
										Default: nil,
										Type: &schema.FieldType{
											Name:     []byte(`ID`),
											Nullable: false,
											IsList:   false,
										},
									},
								},
								Type: &schema.FieldType{
									Name:     []byte(`User`),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
						},
						Extentions: []*schema.OperationDefinition{
							{
								OperationType: schema.QueryOperation,
								Name:          nil,
								Fields: []*schema.FieldDefinition{
									{
										Name: []byte("user"),
										Arguments: []*schema.ArgumentDefinition{
											{
												Name:    []byte(`id`),
												Default: nil,
												Type: &schema.FieldType{
													Name:     []byte(`ID`),
													Nullable: false,
													IsList:   false,
												},
											},
											{
												Name:    []byte(`isActive`),
												Default: []byte(`false`),
												Type: &schema.FieldType{
													Name:     []byte(`Boolean`),
													Nullable: false,
													IsList:   false,
												},
											},
										},
										Type: &schema.FieldType{
											Name:     []byte(`User`),
											Nullable: false,
											IsList:   false,
										},
										Directives: []*schema.Directive{},
									},
									{
										Name:      []byte("users"),
										Arguments: []*schema.ArgumentDefinition{},
										Type: &schema.FieldType{
											Name:     nil,
											Nullable: false,
											IsList:   true,
											ListType: &schema.FieldType{
												Name:     []byte(`User`),
												Nullable: true,
												IsList:   false,
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
										Name:    []byte(`name`),
										Default: nil,
										Type: &schema.FieldType{
											Name:     []byte(`String`),
											Nullable: false,
											IsList:   false,
										},
									},
								},
								Type: &schema.FieldType{
									Name:     []byte(`User`),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
						},
						Extentions: []*schema.OperationDefinition{
							{
								OperationType: schema.MutationOperation,
								Name:          nil,
								Fields: []*schema.FieldDefinition{
									{
										Name: []byte("createUser"),
										Arguments: []*schema.ArgumentDefinition{
											{
												Name:    []byte(`name`),
												Default: nil,
												Type: &schema.FieldType{
													Name:     []byte(`String`),
													Nullable: false,
													IsList:   false,
												},
											},
											{
												Name:    []byte(`isActive`),
												Default: []byte(`true`),
												Type: &schema.FieldType{
													Name:     []byte(`Boolean`),
													Nullable: false,
													IsList:   false,
												},
											},
										},
										Type: &schema.FieldType{
											Name:     []byte(`User`),
											Nullable: false,
											IsList:   false,
										},
										Directives: []*schema.Directive{},
									},
									{
										Name: []byte("updateUser"),
										Arguments: []*schema.ArgumentDefinition{
											{
												Name:    []byte(`id`),
												Default: nil,
												Type: &schema.FieldType{
													Name:     []byte(`ID`),
													Nullable: false,
													IsList:   false,
												},
											},
											{
												Name:    []byte(`name`),
												Default: nil,
												Type: &schema.FieldType{
													Name:     []byte(`String`),
													Nullable: true,
													IsList:   false,
												},
											},
										},
										Type: &schema.FieldType{
											Name:     []byte(`User`),
											Nullable: false,
											IsList:   false,
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
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.SubscriptionOperation,
						Name:          nil,
						Fields: []*schema.FieldDefinition{
							{
								Name:      []byte("userCreated"),
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     []byte(`User`),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
							},
						},
						Extentions: []*schema.OperationDefinition{
							{
								OperationType: schema.SubscriptionOperation,
								Name:          nil,
								Fields: []*schema.FieldDefinition{
									{
										Name: []byte("userCreated"),
										Arguments: []*schema.ArgumentDefinition{
											{
												Name:    []byte(`isActive`),
												Default: []byte(`false`),
												Type: &schema.FieldType{
													Name:     []byte(`Boolean`),
													Nullable: false,
													IsList:   false,
												},
											},
										},
										Type: &schema.FieldType{
											Name:     []byte(`User`),
											Nullable: false,
											IsList:   false,
										},
										Directives: []*schema.Directive{},
									},
									{
										Name: []byte("userUpdated"),
										Arguments: []*schema.ArgumentDefinition{
											{
												Name:    []byte(`id`),
												Default: nil,
												Type: &schema.FieldType{
													Name:     []byte(`ID`),
													Nullable: false,
													IsList:   false,
												},
											},
										},
										Type: &schema.FieldType{
											Name:     []byte(`User`),
											Nullable: false,
											IsList:   false,
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
			name: "Parse extend enum definition",
			input: []byte(`enum Role {
					ADMIN
					USER
			}

			extend enum Role {
					GUEST
					SUPERADMIN
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Enums: []*schema.EnumDefinition{
					{
						Name: []byte("Role"),
						Values: []*schema.EnumElement{
							{
								Name: []byte("ADMIN"),
								Value: []byte("ADMIN"),
								Directives: nil,
							},
							{
								Name: []byte("USER"),
								Value: []byte("USER"),
								Directives: nil,
							},
						},
						Extentions: []*schema.EnumDefinition{
							{
								Name: []byte("Role"),
								Values: []*schema.EnumElement{
									{
										Name: []byte("GUEST"),
										Value: []byte("GUEST"),
										Directives: nil,
									},
									{
										Name: []byte("SUPERADMIN"),
										Value: []byte("SUPERADMIN"),
										Directives: nil,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Parse extend input definition",
			input: []byte(`input UserInput {
					id: ID!
					name: String!
			}

			extend input UserInput {
					age: Int
					email: String
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Inputs: []*schema.InputDefinition{
					{
						Name: []byte("UserInput"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("id"),
								Type: &schema.FieldType{
									Name:     []byte("ID"),
									Nullable: false,
								},
								Directives: []*schema.Directive{},
							},
							{
								Name: []byte("name"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
								},
								Directives: []*schema.Directive{},
							},
						},
						Extentions: []*schema.InputDefinition{
							{
								Name: []byte("UserInput"),
								Fields: []*schema.FieldDefinition{
									{
										Name: []byte("age"),
										Type: &schema.FieldType{
											Name:     []byte("Int"),
											Nullable: true,
										},
										Directives: []*schema.Directive{},
									},
									{
										Name: []byte("email"),
										Type: &schema.FieldType{
											Name:     []byte("String"),
											Nullable: true,
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
			name: "Parse extend interface definition",
			input: []byte(`interface Node {
					id: ID!
			}

			extend interface Node {
					createdAt: String
					updatedAt: String
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								},
								Directives: []*schema.Directive{},
							},
						},
						Extentions: []*schema.InterfaceDefinition{
							{
								Name: []byte("Node"),
								Fields: []*schema.FieldDefinition{
									{
										Name: []byte("createdAt"),
										Type: &schema.FieldType{
											Name:     []byte("String"),
											Nullable: true,
										},
										Directives: []*schema.Directive{},
									},
									{
										Name: []byte("updatedAt"),
										Type: &schema.FieldType{
											Name:     []byte("String"),
											Nullable: true,
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
			name: "Parse extend union definition",
			input: []byte(`union SearchResult = User | Post

			extend union SearchResult = Comment | Page`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Unions: []*schema.UnionDefinition{
					{
						Name: []byte("SearchResult"),
						Types: [][]byte{
							[]byte("User"),
							[]byte("Post"),
						},
						Extentions: []*schema.UnionDefinition{
							{
								Name: []byte("SearchResult"),
								Types: [][]byte{
									[]byte("Comment"),
									[]byte("Page"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Directive on field (@deprecated)",
			input: []byte(`
				type User {
					name: String @deprecated(reason: "Use fullName instead")
				}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Types: []*schema.TypeDefinition{
					{
						Name: []byte("User"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("name"),
								Type: &schema.FieldType{Name: []byte("String"), Nullable: true},
								Directives: []*schema.Directive{
									{
										Name: []byte("deprecated"),
										Arguments: []*schema.DirectiveArgument{
											{Name: []byte("reason"), Value: []byte("\"Use fullName instead\"")},
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
			name: "Directive on input field",
			input: []byte(`
				directive @length(max: Int) on FIELD_DEFINITION
				
				input Filter {
					query: String @length(max: 50)
				}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Directives: []*schema.DirectiveDefinition{
					{
						Name: []byte("length"),
						Arguments: []*schema.ArgumentDefinition{
							{
								Name: []byte("max"), Type: &schema.FieldType{Name: []byte("Int"), Nullable: true},
							},
						},
						Locations: []*schema.Location{
							{Name: []byte("FIELD_DEFINITION")},
						},
					},
				},
				Inputs: []*schema.InputDefinition{
					{
						Name: []byte("Filter"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("query"),
								Type: &schema.FieldType{Name: []byte("String"), Nullable: true},
								Directives: []*schema.Directive{
									{
										Name: []byte("length"),
										Arguments: []*schema.DirectiveArgument{
											{Name: []byte("max"), Value: []byte("50")},
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
			name: "Multiple directives on field",
			input: []byte(`
				type Query {
					user: User @deprecated @auth(role: "USER")
				}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("user"),
								Type: &schema.FieldType{Name: []byte("User"), Nullable: true},
								Directives: []*schema.Directive{
									{Name: []byte("deprecated")},
									{
										Name: []byte("auth"),
										Arguments: []*schema.DirectiveArgument{
											{Name: []byte("role"), Value: []byte("\"USER\"")},
										},
									},
								},
								Arguments: []*schema.ArgumentDefinition{},
							},
						},
					},
				},
			},
		},
		{
			name: "Simple directive definition (no args, single location)",
			input: []byte(`
				directive @deprecated on FIELD_DEFINITION
			`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Directives: []*schema.DirectiveDefinition{
					{
						Name: []byte("deprecated"),
						Arguments: []*schema.ArgumentDefinition{},
						Locations: []*schema.Location{
							{Name: []byte("FIELD_DEFINITION")},
						},
					},
				},
			},
		},
		{
			name: "Directive definition with arguments, repeatable, multiple locations",
			input: []byte(`
				directive @auth(
					role: String = "USER",
					enabled: Boolean!
				) repeatable on FIELD_DEFINITION | OBJECT
			`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Directives: []*schema.DirectiveDefinition{
					{
						Name:       []byte("auth"),
						Repeatable: true,
						Arguments: []*schema.ArgumentDefinition{
							{
								Name: []byte("role"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: true,
									IsList:   false,
								},
								Default: []byte(`"USER"`),
							},
							{
								Name: []byte("enabled"),
								Type: &schema.FieldType{
									Name:     []byte("Boolean"),
									Nullable: false,
									IsList:   false,
								},
							},
						},
						Locations: []*schema.Location{
							{Name: []byte("FIELD_DEFINITION")},
							{Name: []byte("OBJECT")},
						},
					},
				},
			},
		},
		{
			name: "Directive usage on enum value",
			input: []byte(`enum Direction {
				NORTH
				EAST @deprecated(reason: "No longer used")
				SOUTH
				WEST
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Enums: []*schema.EnumDefinition{
					{
						Name: []byte("Direction"),
						Values: []*schema.EnumElement{
							{
								Name: []byte("NORTH"),
								Value: []byte("NORTH"),
								Directives: nil,
							},
							{
								Name: []byte("EAST"),
								Value: []byte("EAST"),
								Directives: []*schema.Directive{
									{
										Name: []byte("deprecated"),
										Arguments: []*schema.DirectiveArgument{
											{Name: []byte("reason"), Value: []byte("\"No longer used\"")},
										},
									},
								},
							},
							{
								Name: []byte("SOUTH"),
								Value: []byte("SOUTH"),
								Directives: nil,
							},
							{
								Name: []byte("WEST"),
								Value: []byte("WEST"),
								Directives: nil,
							},
						},
					},
				},
			},
		},
		{
			name: "Directive usage on interface",
			input: []byte(`interface Node @auth(role: "ADMIN") {
				id: ID!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
							},
						},
						Directives: []*schema.Directive{
							{
								Name: []byte("auth"),
								Arguments: []*schema.DirectiveArgument{
									{
										Name:  []byte("role"),
										Value: []byte(`"ADMIN"`),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Directive usage on union",
			input: []byte(`union Entity @deprecated(reason: "Use another union") = User | Post`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Unions: []*schema.UnionDefinition{
					{
						Name: []byte("Entity"),
						Types: [][]byte{
							[]byte("User"),
							[]byte("Post"),
						},
						Directives: []*schema.Directive{
							{
								Name: []byte("deprecated"),
								Arguments: []*schema.DirectiveArgument{
									{
										Name:  []byte("reason"),
										Value: []byte(`"Use another union"`),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Directive usage on schema",
			input: []byte(`schema @example {
				query: RootQuery
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query: []byte("RootQuery"),
					Directives: []*schema.Directive{
						{
							Name: []byte("example"),
							Arguments: nil,
						},
					},
				},
			},
		},
		{
			name: "Multiple directives on interface",
			input: []byte(`interface Node @auth(role: "ADMIN") @deprecated(reason: "Will be removed") {
				id: ID!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
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
								Directives: []*schema.Directive{},
							},
						},
						Directives: []*schema.Directive{
							{
								Name: []byte("auth"),
								Arguments: []*schema.DirectiveArgument{
									{Name: []byte("role"), Value: []byte(`"ADMIN"`)},
								},
							},
							{
								Name: []byte("deprecated"),
								Arguments: []*schema.DirectiveArgument{
									{Name: []byte("reason"), Value: []byte(`"Will be removed"`)},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Enum with directives on itself and values",
			input: []byte(`enum Role @upperCase {
			ADMIN @deprecated
			USER
		}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Enums: []*schema.EnumDefinition{
					{
						Name: []byte("Role"),
						Directives: []*schema.Directive{
							{Name: []byte("upperCase")},
						},
						Values: []*schema.EnumElement{
							{
								Name: []byte("ADMIN"),
								Value: []byte("ADMIN"),
								Directives: []*schema.Directive{
									{Name: []byte("deprecated")},
								},
							},
							{
								Name: []byte("USER"),
								Value: []byte("USER"),
								Directives: nil,
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
				return
			}

			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(ignores...)); diff != "" {
				t.Errorf("Parse() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
