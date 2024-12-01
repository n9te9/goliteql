package schema_test

import (
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := schema.NewLexer()
			parser := schema.NewParser(lexer)
			got, err := parser.Parse(tt.input)
			if err != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(schema.Schema{}, schema.TypeDefinition{})); diff != "" {
				t.Errorf("Parse() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}