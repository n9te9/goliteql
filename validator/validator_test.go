package validator_test

import (
	"testing"

	"github.com/lkeix/gg-parser/query"
	"github.com/lkeix/gg-parser/schema"
	"github.com/lkeix/gg-parser/validator"
	"errors"
)

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name string
		schemaFunc func(parser *schema.Parser) *schema.Schema
		query []byte
		want error
	}{
		{
			name: "Validate query with missing operation",
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
				user {
					id
					name
					age
				}
			}`),
			want: errors.New("error validating operations: field user is not defined in schema"),
		},
		{
			name: "Validate query with missing operation arguments",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					user(id: ID!): User
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
				user {
					id
					name
					age
				}
			}`),
			want: errors.New("error validating operations: error validating field user: missing required arguments: [id]"),
		},
		{
			name: "Validate query with missing subfields",
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
					posts
				}
			}`),
			want: errors.New("error validating operations: error validating field users: field posts is not defined in schema"),
		},
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
		},{
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
		t.Run(tt.name, func(t *testing.T) {
			lexer := schema.NewLexer()
			s := tt.schemaFunc(schema.NewParser(lexer))
			s, _ = s.Merge()
			s.Preload()
			
			queryLexer := query.NewLexer()
			queryParser := query.NewParser(queryLexer)

			v := validator.NewValidator(s, queryParser)

			err := v.Validate(tt.query)

			if tt.want != nil && err.Error() != tt.want.Error() {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.want)
				return
			}

			if tt.want == nil && err != nil {
				t.Errorf("Parse() error %v", err)
				return
			}

			if tt.want == nil && err == nil {
				return
			}
		})
	}
}