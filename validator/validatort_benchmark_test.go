package validator_test

import (
	"testing"

	"github.com/n9te9/goliteql/query"
	"github.com/n9te9/goliteql/schema"
	"github.com/n9te9/goliteql/validator"
)

func BenchmarkValidator_Validate(b *testing.B) {
	tests := []struct {
		name       string
		schemaFunc func(parser *schema.Parser) *schema.Schema
		query      []byte
		want       error
	}{
		{
			name: "Validate simple query",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					users: [User]
				}

				type User {
					id: ID!
					name: String
					age: Int
				}

				directive @deprecated(reason: String) on FIELD_DEFINITION`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}

				return s
			},
			query: []byte(`query {
				users {
					id
					name
					age
				}
			}`),
			want: nil,
		}, {
			name: "Validate nested query",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					users: [User]
				}

				type User {
					id: ID!
					name: String
					age: Int
					posts: [Post]
				}

				type Post {
					id: ID!
					title: String
				}

				directive @deprecated(reason: String) on FIELD_DEFINITION`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}

				return s
			},
			query: []byte(`query {
				users {
					id
					name
					age
					posts {
						id
						title
					}
				}
			}`),
			want: nil,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			lexer := schema.NewLexer()
			s := tt.schemaFunc(schema.NewParser(lexer))
			s, _ = s.Merge()

			queryLexer := query.NewLexer()
			queryParser := query.NewParser(queryLexer)

			v := validator.NewValidator(s, queryParser)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := v.Validate(tt.query)

				if tt.want != nil && err.Error() != tt.want.Error() {
					b.Errorf("Parse() error = %v, wantErr %v", err, tt.want)
					return
				}

				if tt.want == nil && err != nil {
					b.Errorf("Parse() error %v", err)
					return
				}

				if tt.want == nil && err == nil {
				}
			}
		})
	}
}
