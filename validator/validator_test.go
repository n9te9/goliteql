package validator_test

import (
	"testing"

	"errors"
	"github.com/lkeix/gg-parser/query"
	"github.com/lkeix/gg-parser/schema"
	"github.com/lkeix/gg-parser/validator"
)

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name       string
		schemaFunc func(parser *schema.Parser) *schema.Schema
		query      []byte
		want       error
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
		{
			name: "Validate query with undefined field",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					users: [User]
				}

				type User {
					id: ID!
					name: String
				}`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}

				return s
			},
			query: []byte(`query {
				users {
					id
					unknownField
				}
			}`),
			want: errors.New("error validating operations: error validating field users: field unknownField is not defined in schema"),
		},
		{
			name: "Validate query with missing required argument",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					user(id: ID!): User
				}

				type User {
					id: ID!
					name: String
				}`)
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
				}
			}`),
			want: errors.New("error validating operations: error validating field user: missing required arguments: [id]"),
		},
		{
			name: "Validate query with type mismatch in argument",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					user(id: ID!): User
				}

				type User {
					id: ID!
					name: String
				}`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}

				return s
			},
			query: []byte(`query {
				user(id: 123) {
					id
					name
				}
			}`),
		},
		{
			name: "Validate query with valid fragment",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					user: User
				}

				type User {
					id: ID!
					name: String
					age: Int
				}`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}

				return s
			},
			query: []byte(`query {
				user {
					...UserFragment
				}
			}

			fragment UserFragment on User {
				id
				name
				age
			}`),
			want: nil,
		},
		{
			name: "Validate query with missing field in nested type",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					users: [User]
				}

				type User {
					id: ID!
					name: String
					posts: [Post]
				}

				type Post {
					id: ID!
					title: String
				}`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}

				return s
			},
			query: []byte(`query {
				users {
					id
					posts {
						id
						unknownField
					}
				}
			}`),
			want: errors.New("error validating operations: error validating field users: error validating field posts: field unknownField is not defined in schema"),
		},
		{
			name: "Validate valid nested query",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					users: [User]
				}

				type User {
					id: ID!
					name: String
					posts: [Post]
				}

				type Post {
					id: ID!
					title: String
				}`)
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
					posts {
						id
						title
					}
				}
			}`),
			want: nil,
		},
		{
			name: "Validate query with invalid fragment",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					user: User
				}

				type User {
					id: ID!
					name: String
					age: Int
				}

				type Post {
					id: ID!
					title: String
				}`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}

				return s
			},
			query: []byte(`query {
				user {
					...PostFragment
				}
			}

			fragment PostFragment on Post {
				id
				title
			}`),
			want: errors.New(`error validating operations: error validating field user: fragment PostFragment is based on type Post, but field is of type User`),
		},
		{
			name: "Validate query with valid fragment spread",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					user: User
				}
		
				type User {
					id: ID!
					name: String
					posts: [Post]
				}
		
				type Post {
					id: ID!
					title: String
				}`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}

				return s
			},
			query: []byte(`query {
				user {
					posts {
						...PostFragment
					}
				}
			}
		
			fragment PostFragment on Post {
				id
				title
			}`),
			want: nil,
		},
		{
			name: "Validate query with valid inline fragment",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					searchResults: [SearchResult]
				}
		
				union SearchResult = User | Post
		
				type User {
					id: ID!
					name: String
				}
		
				type Post {
					id: ID!
					title: String
				}`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}
		
				return s
			},
			query: []byte(`query {
				searchResults {
					...on User {
						id
						name
					}
					...on Post {
						id
						title
					}
				}
			}`),
			want: nil,
		},
		{
			name: "Validate query with empty invalid inline fragment",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					searchResults: [SearchResult]
				}
		
				union SearchResult = User | Post
		
				type User {
					id: ID!
					name: String
				}
		
				type Post {
					id: ID!
					title: String
				}`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}
		
				return s
			},
			query: []byte(`query {
				searchResults {}
			}`),
			want: errors.New("error validating operations: error validating field searchResults: union type SearchResult must have subfields"),
		},
		{
			name: "Validate query with invalid inline fragment type",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					searchResults: [SearchResult]
				}
		
				union SearchResult = User | Post
		
				type User {
					id: ID!
					name: String
				}
		
				type Post {
					id: ID!
					title: String
				}`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}
		
				return s
			},
			query: []byte(`query {
				searchResults {
					...on InvalidType {
						id
					}
				}
			}`),
			want: errors.New("error validating operations: error validating field searchResults: type InvalidType is not defined in schema"),
		},
		{
			name: "Validate query with nested inline fragment",
			schemaFunc: func(parser *schema.Parser) *schema.Schema {
				input := []byte(`type Query {
					searchResults: [SearchResult]
				}
		
				union SearchResult = User | Post
		
				type User {
					id: ID!
					name: String
					posts: [Post]
				}
		
				type Post {
					id: ID!
					title: String
					comments: [Comment]
				}
		
				type Comment {
					id: ID!
					content: String
				}`)
				s, err := parser.Parse(input)
				if err != nil {
					panic(err)
				}
		
				return s
			},
			query: []byte(`query {
				searchResults {
					...on User {
						id
						name
						posts {
							...on Post {
								id
								title
								comments {
									...on Comment {
										id
										content
									}
								}
							}
						}
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
