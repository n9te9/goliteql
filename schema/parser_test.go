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
		name string
		input []byte
		want *schema.Schema
		wantErr error
	}{
		{
			name: "simple type with scalar fields",
			input: []byte(`type User {
				id: ID!
				name: String!
				age: Int
			}`),
			want: &schema.Schema{
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
										Name: []byte("User"),
										Nullable: false,
										IsList: false,
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
										Name: []byte("Post"),
										Nullable: true,
										IsList: false,
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
			name: "simple union type",
			input: []byte(`union SearchResult = User | Post`),
			want: &schema.Schema{
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
			name: "invalid union type",
			input: []byte(`union SearchResult = User |`),
			want: nil,
			wantErr: errors.New("unexpected end of input"),
		},
		{
			name: "empty union type",
			input: []byte(`union SearchResult =`),
			want: nil,
			wantErr: errors.New("unexpected end of input"),
		},
		{
			name: "simple interface type",
			input: []byte(`interface Node {
				id: ID!
			}`),
			want: &schema.Schema{
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
			name: "interface without fields",
			input: []byte(`interface Empty {}`),
			want: &schema.Schema{
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
			name: "invalid interface (missing opening curly brace)",
			input: []byte(`interface Node id: ID!}`),
			want:  nil,
			wantErr: errors.New("expected '{' but got id"),
		},
		{
			name: "invalid interface (missing closing curly brace)",
			input: []byte(`interface Node {
				id: ID!`),
			want: nil,
			wantErr: errors.New("unexpected end of input"),
		},
		{
			name: "simple Query operation",
			input: []byte(`type Query {
				user(id: ID!): User
			}`),
			want: &schema.Schema{
				Operations: []*schema.OperationDefinition{
					{
						OperationType: schema.QueryOperation,
						Name: nil,
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
			got, err := parser.Parse(tt.input)
			if tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(schema.Schema{}, schema.TypeDefinition{})); diff != "" {
				t.Errorf("Parse() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}