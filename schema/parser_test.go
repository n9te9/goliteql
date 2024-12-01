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
			name: "parse simple type schema",
			input: []byte(`type User {
				id: ID!
				name: String!
			}`),
			want: &schema.Schema{
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
							{
								Name: []byte("name"),
								Type: &schema.FieldType{
									Name: []byte("String"),
									Nullable: false,
									IsList: false,
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

			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(schema.Schema{}, schema.TypeDefinition{}, schema.FieldDefinition{})); diff != "" {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}