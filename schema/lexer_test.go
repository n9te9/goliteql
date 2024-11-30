package schema_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lkeix/gg-parser/schema"
)

func TestLexer_Lex(t *testing.T) {
	tests := []struct{
		name string
		input []byte
		expected []*schema.Token
		wantErr error
	}{
		{
			name: "Lex simple schema",
			input: []byte(`type User {
				id: ID!
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("User"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 11, Line: 1},
				{Type: schema.Field, Value: []byte("id"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 7, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 9, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 11, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Lex simple schema has optional field",
			input: []byte(`type User {
				id: ID!
				option: String
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("User"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 11, Line: 1},
				{Type: schema.Field, Value: []byte("id"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 7, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 9, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 11, Line: 2},
				{Type: schema.Field, Value: []byte("option"), Column: 5, Line: 3},
				{Type: schema.Colon, Value: []byte(":"), Column: 11, Line: 3},
				{Type: schema.Identifier, Value: []byte("String"), Column: 13, Line: 3},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 4},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 4},
			},
		},
		{
			name: "Lex simple schema has required list field",
			input: []byte(`type User {
				id: ID!
				tags: [String]!
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("User"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 11, Line: 1},
				{Type: schema.Field, Value: []byte("id"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 7, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 9, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 11, Line: 2},
				{Type: schema.Field, Value: []byte("tags"), Column: 5, Line: 3},
				{Type: schema.Colon, Value: []byte(":"), Column: 9, Line: 3},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 11, Line: 3},
				{Type: schema.Identifier, Value: []byte("String"), Column: 12, Line: 3},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 18, Line: 3},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 19, Line: 3},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 4},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 4},
			},
		},
		{
			name: "Lex simple schema has optional list field",
			input: []byte(`type User {
				id: ID!
				tags: [String]
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("User"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 11, Line: 1},
				{Type: schema.Field, Value: []byte("id"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 7, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 9, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 11, Line: 2},
				{Type: schema.Field, Value: []byte("tags"), Column: 5, Line: 3},
				{Type: schema.Colon, Value: []byte(":"), Column: 9, Line: 3},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 11, Line: 3},
				{Type: schema.Identifier, Value: []byte("String"), Column: 12, Line: 3},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 18, Line: 3},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 4},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 4},
			},
		},
		{
			name: "Lex simple schema has simple key directive",
			input: []byte(`type User {
				id: ID!
				key: @key
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("User"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 11, Line: 1},
				{Type: schema.Field, Value: []byte("id"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 7, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 9, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 11, Line: 2},
				{Type: schema.Field, Value: []byte("key"), Column: 5, Line: 3},
				{Type: schema.Colon, Value: []byte(":"), Column: 8, Line: 3},
				{Type: schema.At, Value: []byte("@"), Column: 10, Line: 3},
				{Type: schema.Identifier, Value: []byte("key"), Column: 11, Line: 3},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 4},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 4},
			},
		},
		{
			name: "Lex simple input schema",
			input: []byte(`input User {
				id: ID!
			}`),
			expected: []*schema.Token{
				{Type: schema.Input, Value: []byte("input"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("User"), Column: 7, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("id"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 7, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 9, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 11, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Lex simple interface schema",
			input: []byte(`interface Node {
				id: ID!
			}`),
			expected: []*schema.Token{
				{Type: schema.Interface, Value: []byte("interface"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Node"), Column: 11, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 16, Line: 1},
				{Type: schema.Field, Value: []byte("id"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 7, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 9, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 11, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Lex simple Query schema",
			input: []byte(`type Query {
				getUser(id: ID!): User
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("getUser"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 12, Line: 2},
				{Type: schema.Field, Value: []byte("id"), Column: 13, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 15, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 17, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 19, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 20, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 21, Line: 2},
				{Type: schema.Identifier, Value: []byte("User"), Column: 23, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Lex simple Mutate schema",
			input: []byte(`type Mutate {
				createUser(name: String!): User
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Mutate"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 13, Line: 1},
				{Type: schema.Field, Value: []byte("createUser"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 15, Line: 2},
				{Type: schema.Field, Value: []byte("name"), Column: 16, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 20, Line: 2},
				{Type: schema.Identifier, Value: []byte("String"), Column: 22, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 28, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 29, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 30, Line: 2},
				{Type: schema.Identifier, Value: []byte("User"), Column: 32, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := schema.NewLexer()
			got, err := lexer.Lex(tt.input)
			if err != tt.wantErr {
				t.Errorf("Lex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Errorf("Lex() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}