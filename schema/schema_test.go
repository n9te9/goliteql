package schema_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lkeix/gg-parser/schema"
)

func TestSchema_Merge(t *testing.T) {
	tests := []struct{
		name string
		input []byte
		want *schema.Schema
		wantErr error
		isSkip bool
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
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("user"),
								Directives: []*schema.Directive{},
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("id"),
										Type: &schema.FieldType{
											Name: []byte("ID"),
											Nullable: false,
											IsList: false,
										},
									},
									{
										Name: []byte("isActive"),
										Type: &schema.FieldType{
											Name: []byte("Boolean"),
											Nullable: false,
											IsList: false,
										},
										Default: []byte("false"),
									},
								},
								Type: &schema.FieldType{
									Name: []byte("User"),
									Nullable: false,
									IsList: false,
								},
							},
							{
								Name: []byte("users"),
								Arguments: []*schema.ArgumentDefinition{},
								Directives: []*schema.Directive{},
								Type: &schema.FieldType{
									Name: nil,
									Nullable: false,
									IsList: true,
									ListType: &schema.FieldType{
										Name: []byte("User"),
										Nullable: true,
										IsList: false,
									},
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
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.MutationOperation,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("createUser"),
								Directives: []*schema.Directive{},
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("name"),
										Type: &schema.FieldType{
											Name: []byte("String"),
											Nullable: false,
											IsList: false,
										},
									},
									{
										Name: []byte("email"),
										Type: &schema.FieldType{
											Name: []byte("String"),
											Nullable: false,
											IsList: false,
										},
									},
									{
										Name: []byte("isActive"),
										Type: &schema.FieldType{
											Name: []byte("Boolean"),
											Nullable: false,
											IsList: false,
										},
										Default: []byte("false"),
									},
								},
								Type: &schema.FieldType{
									Name: []byte("User"),
									Nullable: false,
									IsList: false,
								},
							},
							{
								Name: []byte("deleteUser"),
								Directives: []*schema.Directive{},
								Arguments: []*schema.ArgumentDefinition{
									{
										Name: []byte("id"),
										Type: &schema.FieldType{
											Name: []byte("ID"),
											Nullable: false,
											IsList: false,
										},
									},
								},
								Type: &schema.FieldType{
									Name: []byte("Boolean"),
									Nullable: false,
									IsList: false,
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
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.SubscriptionOperation,
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("userCreated"),
								Directives: []*schema.Directive{},
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name: []byte("User"),
									Nullable: false,
									IsList: false,
								},
							},
							{
								Name: []byte("userDeleted"),
								Directives: []*schema.Directive{},
								Arguments: []*schema.ArgumentDefinition{},
								Type: &schema.FieldType{
									Name: []byte("User"),
									Nullable: false,
									IsList: false,
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
					Query: []byte("Query"),
					Mutation: []byte("Mutation"),
					Subscription: []byte("Subscription"),
				},
				Types: []*schema.TypeDefinition{
					{
						Name: []byte("User"),
						Fields: []*schema.FieldDefinition{
							{
								Name: []byte("email"),
								Type: &schema.FieldType{
									Name: []byte("String"),
									Nullable: false,
									IsList: false,
								},
								Directives: []*schema.Directive{},
							},
							{
								Name: []byte("id"),
								Type: &schema.FieldType{
									Name: []byte("ID"),
									Nullable: false,
									IsList: false,
								},
								Directives: []*schema.Directive{},
							},
							{
								Name: []byte("name"),
								Type: &schema.FieldType{
									Name: []byte("String"),
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
	}

	ignores := []any{
		schema.Schema{},
		schema.TypeDefinition{},
		schema.InputDefinition{},
		schema.Indexes{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := schema.NewLexer()
			parser := schema.NewParser(lexer)
			schema, err := parser.Parse(tt.input)
			if err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
			}

			got, err := schema.Merge()
			if err != tt.wantErr {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(ignores...)); diff != "" {
				t.Errorf("Parse() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}