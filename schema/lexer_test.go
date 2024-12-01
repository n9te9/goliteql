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
			name: "Lex field with multiple arguments",
			input: []byte(`type Query {
				getUsers(age: Int!, active: Boolean): [User]!
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("getUsers"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 13, Line: 2},
				{Type: schema.Field, Value: []byte("age"), Column: 14, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 17, Line: 2},
				{Type: schema.Identifier, Value: []byte("Int"), Column: 19, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 22, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Column: 23, Line: 2},
				{Type: schema.Field, Value: []byte("active"), Column: 25, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 31, Line: 2},
				{Type: schema.Identifier, Value: []byte("Boolean"), Column: 33, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 40, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 41, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 43, Line: 2},
				{Type: schema.Identifier, Value: []byte("User"), Column: 44, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 48, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 49, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Lex field with arguments and default values",
			input: []byte(`type Query {
				getUsers(age: Int = 18, active: Boolean = true): [User]!
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("getUsers"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 13, Line: 2},
				{Type: schema.Field, Value: []byte("age"), Column: 14, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 17, Line: 2},
				{Type: schema.Identifier, Value: []byte("Int"), Column: 19, Line: 2},
				{Type: schema.Equal, Value: []byte("="), Column: 23, Line: 2},
				{Type: schema.Int, Value: []byte("18"), Column: 25, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Column: 27, Line: 2},
				{Type: schema.Field, Value: []byte("active"), Column: 29, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 35, Line: 2},
				{Type: schema.Identifier, Value: []byte("Boolean"), Column: 37, Line: 2},
				{Type: schema.Equal, Value: []byte("="), Column: 45, Line: 2},
				{Type: schema.Boolean, Value: []byte("true"), Column: 47, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 51, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 52, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 54, Line: 2},
				{Type: schema.Identifier, Value: []byte("User"), Column: 55, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 59, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 60, Line: 2},
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
		{
			name: "Lex simple Subscription schema",
			input: []byte(`type Subscription {
				messageAdded(channelId: ID!): Message
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Subscription"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 19, Line: 1},
				{Type: schema.Field, Value: []byte("messageAdded"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 17, Line: 2},
				{Type: schema.Field, Value: []byte("channelId"), Column: 18, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 27, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 29, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 31, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 32, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 33, Line: 2},
				{Type: schema.Identifier, Value: []byte("Message"), Column: 35, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Lex field with list default value",
			input: []byte(`type Query {
				getUsers(ids: [ID] = ["id1", "id2", "id3"]): [User]!
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("getUsers"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 13, Line: 2},
				{Type: schema.Field, Value: []byte("ids"), Column: 14, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 17, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 19, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 20, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 22, Line: 2},
				{Type: schema.Equal, Value: []byte("="), Column: 24, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 26, Line: 2},
				{Type: schema.String, Value: []byte(`"id1"`), Column: 27, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Column: 32, Line: 2},
				{Type: schema.String, Value: []byte(`"id2"`), Column: 34, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Column: 39, Line: 2},
				{Type: schema.String, Value: []byte(`"id3"`), Column: 41, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 46, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 47, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 48, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 50, Line: 2},
				{Type: schema.Identifier, Value: []byte("User"), Column: 51, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 55, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 56, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Type with list of non-nullable IDs",
			input: []byte(`type User {
				ids: [ID!]!
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("User"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 11, Line: 1},
				{Type: schema.Field, Value: []byte("ids"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 8, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 10, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 11, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 13, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 14, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"), Column: 15, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Field with default list value",
			input: []byte(`type Query {
				getIds(ids: [ID] = ["id1", "id2"])
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("getIds"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 11, Line: 2},
				{Type: schema.Field, Value: []byte("ids"), Column: 12, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 15, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 17, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 18, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 20, Line: 2},
				{Type: schema.Equal, Value: []byte("="), Column: 22, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 24, Line: 2},
				{Type: schema.String, Value: []byte(`"id1"`), Column: 25, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Column: 30, Line: 2},
				{Type: schema.String, Value: []byte(`"id2"`), Column: 32, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 37, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 38, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Field with nested lists",
			input: []byte(`type Query {
				getNestedLists(ids: [[ID]])
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("getNestedLists"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 19, Line: 2},
				{Type: schema.Field, Value: []byte("ids"), Column: 20, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 23, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 25, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 26, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 27, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 29, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 30, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 31, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Directive with arguments",
			input: []byte(`type Query {
				users: [User] @include(if: true)
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("users"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 10, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 12, Line: 2},
				{Type: schema.Identifier, Value: []byte("User"), Column: 13, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 17, Line: 2},
				{Type: schema.At, Value: []byte("@"), Column: 19, Line: 2},
				{Type: schema.Identifier, Value: []byte("include"), Column: 20, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 27, Line: 2},
				{Type: schema.Field, Value: []byte("if"), Column: 28, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 30, Line: 2},
				{Type: schema.Boolean, Value: []byte("true"), Column: 32, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 36, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Directive with multiple arguments",
			input: []byte(`type Query {
				users: [User] @deprecated(reason: "use new field", version: 2)
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("users"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 10, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 12, Line: 2},
				{Type: schema.Identifier, Value: []byte("User"), Column: 13, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 17, Line: 2},
				{Type: schema.At, Value: []byte("@"), Column: 19, Line: 2},
				{Type: schema.Identifier, Value: []byte("deprecated"), Column: 20, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 30, Line: 2},
				{Type: schema.Field, Value: []byte("reason"), Column: 31, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 37, Line: 2},
				{Type: schema.String, Value: []byte(`"use new field"`), Column: 39, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Column: 54, Line: 2},
				{Type: schema.Field, Value: []byte("version"), Column: 56, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 63, Line: 2},
				{Type: schema.Int, Value: []byte("2"), Column: 65, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 66, Line: 2},
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