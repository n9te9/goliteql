package schema_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/n9te9/goliteql/schema"
)

func TestSchema_Merge(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *schema.Schema
		wantErr error
		isSkip  bool
	}{
		{
			name: "Merge extend query definition",
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
				Indexes: &schema.Indexes{
					TypeIndex:        make(map[string]*schema.TypeDefinition),
					OperationIndexes: make(map[schema.OperationType]map[string]*schema.OperationDefinition),
					EnumIndex:        make(map[string]*schema.EnumDefinition),
					UnionIndex:       make(map[string]*schema.UnionDefinition),
					InterfaceIndex:   make(map[string]*schema.InterfaceDefinition),
					InputIndex:       make(map[string]*schema.InputDefinition),
					ScalarIndex:      make(map[string]*schema.ScalarDefinition),
					ExtendIndex:      make(map[string]schema.ExtendDefinition),
				},
				Directives: schema.NewBuildInDirectives(),
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Fields: []*schema.FieldDefinition{
							{
								Name:       []byte("user"),
								Directives: []*schema.Directive{},
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
										Name: []byte("isActive"),
										Type: &schema.FieldType{
											Name:     []byte("Boolean"),
											Nullable: false,
											IsList:   false,
										},
										Default: []byte("false"),
									},
								},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: false,
									IsList:   false,
								},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
							{
								Name:       []byte("users"),
								Arguments:  []*schema.ArgumentDefinition{},
								Directives: []*schema.Directive{},
								Type: &schema.FieldType{
									Name:     nil,
									Nullable: false,
									IsList:   true,
									ListType: &schema.FieldType{
										Name:     []byte("User"),
										Nullable: true,
										IsList:   false,
									},
								},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
						},
					},
				},
			},
		}, {
			name: "Merge extend mutation definition",
			input: []byte(`type Mutation {
				createUser(name: String!, email: String!): User!
			}
				
			extend type Mutation {
				createUser(name: String!, email: String!, isActive: Boolean! = false): User!
				deleteUser(id: ID!): Boolean!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Directives: schema.NewBuildInDirectives(),
				Indexes: &schema.Indexes{
					TypeIndex:        make(map[string]*schema.TypeDefinition),
					OperationIndexes: make(map[schema.OperationType]map[string]*schema.OperationDefinition),
					EnumIndex:        make(map[string]*schema.EnumDefinition),
					UnionIndex:       make(map[string]*schema.UnionDefinition),
					InterfaceIndex:   make(map[string]*schema.InterfaceDefinition),
					InputIndex:       make(map[string]*schema.InputDefinition),
					ScalarIndex:      make(map[string]*schema.ScalarDefinition),
					ExtendIndex:      make(map[string]schema.ExtendDefinition),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.MutationOperation,
						Fields: []*schema.FieldDefinition{
							{
								Name:       []byte("createUser"),
								Directives: []*schema.Directive{},
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("name"),
										Type: &schema.FieldType{
											Name:     []byte("String"),
											Nullable: false,
											IsList:   false,
										},
									},
									{
										Name: []byte("email"),
										Type: &schema.FieldType{
											Name:     []byte("String"),
											Nullable: false,
											IsList:   false,
										},
									},
									{
										Name: []byte("isActive"),
										Type: &schema.FieldType{
											Name:     []byte("Boolean"),
											Nullable: false,
											IsList:   false,
										},
										Default: []byte("false"),
									},
								},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: false,
									IsList:   false,
								},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
							{
								Name:       []byte("deleteUser"),
								Directives: []*schema.Directive{},
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
									Name:     []byte("Boolean"),
									Nullable: false,
									IsList:   false,
								},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
						},
					},
				},
			},
		}, {
			name: "Merge extend subscription definition",
			input: []byte(`type Subscription {
				userCreated: User!
			}
			
			extend type Subscription {
				userCreated: User!
				userDeleted: User!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Directives: schema.NewBuildInDirectives(),
				Indexes: &schema.Indexes{
					TypeIndex:        make(map[string]*schema.TypeDefinition),
					OperationIndexes: make(map[schema.OperationType]map[string]*schema.OperationDefinition),
					EnumIndex:        make(map[string]*schema.EnumDefinition),
					UnionIndex:       make(map[string]*schema.UnionDefinition),
					InterfaceIndex:   make(map[string]*schema.InterfaceDefinition),
					InputIndex:       make(map[string]*schema.InputDefinition),
					ScalarIndex:      make(map[string]*schema.ScalarDefinition),
					ExtendIndex:      make(map[string]schema.ExtendDefinition),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.SubscriptionOperation,
						Fields: []*schema.FieldDefinition{
							{
								Name:       []byte("userCreated"),
								Directives: []*schema.Directive{},
								Arguments:  []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: false,
									IsList:   false,
								},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
							{
								Name:       []byte("userDeleted"),
								Directives: []*schema.Directive{},
								Arguments:  []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name:     []byte("User"),
									Nullable: false,
									IsList:   false,
								},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Merge extend type definition",
			input: []byte(`type User {
				id: ID!
				name: String!
			}
			
			extend type User {
				email: String!
			}`),
			want: &schema.Schema{
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Indexes: &schema.Indexes{
					TypeIndex:        make(map[string]*schema.TypeDefinition),
					OperationIndexes: make(map[schema.OperationType]map[string]*schema.OperationDefinition),
					EnumIndex:        make(map[string]*schema.EnumDefinition),
					UnionIndex:       make(map[string]*schema.UnionDefinition),
					InterfaceIndex:   make(map[string]*schema.InterfaceDefinition),
					InputIndex:       make(map[string]*schema.InputDefinition),
					ScalarIndex:      make(map[string]*schema.ScalarDefinition),
					ExtendIndex:      make(map[string]schema.ExtendDefinition),
				},
				Directives: schema.NewBuildInDirectives(),
				Types: []*schema.TypeDefinition{
					{
						Name: []byte("User"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("email"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
							{
								Name: []byte("id"),
								Type: &schema.FieldType{
									Name:     []byte("ID"),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
							{
								Name: []byte("name"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Merge extend interface definition",
			input: []byte(`interface Node {
				id: ID!
			}

			extend interface Node {
				createdAt: DateTime!
			}`),
			want: &schema.Schema{
				Indexes: &schema.Indexes{
					TypeIndex:        make(map[string]*schema.TypeDefinition),
					OperationIndexes: make(map[schema.OperationType]map[string]*schema.OperationDefinition),
					EnumIndex:        make(map[string]*schema.EnumDefinition),
					UnionIndex:       make(map[string]*schema.UnionDefinition),
					InterfaceIndex:   make(map[string]*schema.InterfaceDefinition),
					InputIndex:       make(map[string]*schema.InputDefinition),
					ScalarIndex:      make(map[string]*schema.ScalarDefinition),
					ExtendIndex:      make(map[string]schema.ExtendDefinition),
				},
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Directives: schema.NewBuildInDirectives(),
				Interfaces: []*schema.InterfaceDefinition{
					{
						Name: []byte("Node"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("createdAt"),
								Type: &schema.FieldType{
									Name:     []byte("DateTime"),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
							{
								Name: []byte("id"),
								Type: &schema.FieldType{
									Name:     []byte("ID"),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Merge extend interface definition with multiple fields",
			input: []byte(`interface Node {
					id: ID!
			}
	
			extend interface Node {
					createdAt: DateTime!
					hogehoge: [String!]!
			}`),
			want: &schema.Schema{
				Indexes: &schema.Indexes{
					TypeIndex:        make(map[string]*schema.TypeDefinition),
					OperationIndexes: make(map[schema.OperationType]map[string]*schema.OperationDefinition),
					EnumIndex:        make(map[string]*schema.EnumDefinition),
					UnionIndex:       make(map[string]*schema.UnionDefinition),
					InterfaceIndex:   make(map[string]*schema.InterfaceDefinition),
					InputIndex:       make(map[string]*schema.InputDefinition),
					ScalarIndex:      make(map[string]*schema.ScalarDefinition),
					ExtendIndex:      make(map[string]schema.ExtendDefinition),
				},
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Directives: schema.NewBuildInDirectives(),
				Interfaces: []*schema.InterfaceDefinition{
					{
						Name: []byte("Node"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("createdAt"),
								Type: &schema.FieldType{
									Name:     []byte("DateTime"),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
							{
								Name: []byte("hogehoge"),
								Type: &schema.FieldType{
									Name:     nil,
									Nullable: false,
									IsList:   true,
									ListType: &schema.FieldType{
										Name:     []byte("String"),
										Nullable: false,
										IsList:   false,
									},
								},
								Directives: []*schema.Directive{},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
							{
								Name: []byte("id"),
								Type: &schema.FieldType{
									Name:     []byte("ID"),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
								Location: &schema.Location{
									Name: []byte("FIELD_DEFINITION"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Merge extend union definition",
			input: []byte(`union SearchResult = User | Post
			
			extend union SearchResult = Comment`),
			want: &schema.Schema{
				Indexes: &schema.Indexes{
					TypeIndex:        make(map[string]*schema.TypeDefinition),
					OperationIndexes: make(map[schema.OperationType]map[string]*schema.OperationDefinition),
					EnumIndex:        make(map[string]*schema.EnumDefinition),
					UnionIndex:       make(map[string]*schema.UnionDefinition),
					InterfaceIndex:   make(map[string]*schema.InterfaceDefinition),
					InputIndex:       make(map[string]*schema.InputDefinition),
					ScalarIndex:      make(map[string]*schema.ScalarDefinition),
					ExtendIndex:      make(map[string]schema.ExtendDefinition),
				},
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Directives: schema.NewBuildInDirectives(),
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
		},
		{
			name: "Merge extend enum definition",
			input: []byte(`enum Role {
				ADMIN
				USER
			}

			extend enum Role {
				EDITOR
			}`),
			want: &schema.Schema{
				Indexes: &schema.Indexes{
					TypeIndex:        make(map[string]*schema.TypeDefinition),
					OperationIndexes: make(map[schema.OperationType]map[string]*schema.OperationDefinition),
					EnumIndex:        make(map[string]*schema.EnumDefinition),
					UnionIndex:       make(map[string]*schema.UnionDefinition),
					InterfaceIndex:   make(map[string]*schema.InterfaceDefinition),
					InputIndex:       make(map[string]*schema.InputDefinition),
					ScalarIndex:      make(map[string]*schema.ScalarDefinition),
					ExtendIndex:      make(map[string]schema.ExtendDefinition),
				},
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Directives: schema.NewBuildInDirectives(),
				Enums: []*schema.EnumDefinition{
					{
						Name: []byte("Role"),
						Values: []*schema.EnumElement{
							{
								Name:  []byte("ADMIN"),
								Value: []byte("ADMIN"),
							},
							{
								Name:  []byte("USER"),
								Value: []byte("USER"),
							},
							{
								Name:  []byte("EDITOR"),
								Value: []byte("EDITOR"),
							},
						},
					},
				},
			},
		}, {
			name: "Merge extend input definition",
			input: []byte(`input CreateUserInput {
				name: String!
				email: String!
			}

			extend input CreateUserInput {
				isActive: Boolean! = false
			}`),
			want: &schema.Schema{
				Indexes: &schema.Indexes{
					TypeIndex:        make(map[string]*schema.TypeDefinition),
					OperationIndexes: make(map[schema.OperationType]map[string]*schema.OperationDefinition),
					EnumIndex:        make(map[string]*schema.EnumDefinition),
					UnionIndex:       make(map[string]*schema.UnionDefinition),
					InterfaceIndex:   make(map[string]*schema.InterfaceDefinition),
					InputIndex:       make(map[string]*schema.InputDefinition),
					ScalarIndex:      make(map[string]*schema.ScalarDefinition),
					ExtendIndex:      make(map[string]schema.ExtendDefinition),
				},
				Definition: &schema.SchemaDefinition{
					Query:        []byte("Query"),
					Mutation:     []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Directives: schema.NewBuildInDirectives(),
				Inputs: []*schema.InputDefinition{
					{
						Name: []byte("CreateUserInput"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("isActive"),
								Type: &schema.FieldType{
									Name:     []byte("Boolean"),
									Nullable: false,
									IsList:   false,
								},
								Default:    []byte("false"),
								Directives: []*schema.Directive{},
								Location: &schema.Location{
									Name: []byte("INPUT_FIELD_DEFINITION"),
								},
							},
							{
								Name: []byte("name"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
								Location: &schema.Location{
									Name: []byte("INPUT_FIELD_DEFINITION"),
								},
							},
							{
								Name: []byte("email"),
								Type: &schema.FieldType{
									Name:     []byte("String"),
									Nullable: false,
									IsList:   false,
								},
								Directives: []*schema.Directive{},
								Location: &schema.Location{
									Name: []byte("INPUT_FIELD_DEFINITION"),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := schema.NewLexer()
			parser := schema.NewParser(lexer)
			s, err := parser.Parse(tt.input)
			if err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
			}

			got, err := s.Merge()
			if err != tt.wantErr {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}

			ttWant := schema.WithBuiltin(tt.want)
			ttWant = schema.WithTypeIntrospection(ttWant)

			if diff := cmp.Diff(got, ttWant, cmpopts.IgnoreFields(schema.Schema{}, "Indexes", "Tokens")); diff != "" {
				t.Errorf("Parse() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
