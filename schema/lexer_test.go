package schema_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lkeix/gg-parser/schema"
)

func TestLexer_Lex(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []*schema.Token
		wantErr  error
	}{
		{
			name: "Lex standard schema definition",
			input: []byte(`schema {
				query: Query
				mutation: Mutation
				subscription: Subscription
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedSchema, Value: []byte("schema"), Column: 1, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 8, Line: 1},
				{Type: schema.Field, Value: []byte("query"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 10, Line: 2},
				{Type: schema.Query, Value: []byte("Query"), Column: 12, Line: 2},
				{Type: schema.Field, Value: []byte("mutation"), Column: 5, Line: 3},
				{Type: schema.Colon, Value: []byte(":"), Column: 13, Line: 3},
				{Type: schema.Mutate, Value: []byte("Mutation"), Column: 15, Line: 3},
				{Type: schema.Field, Value: []byte("subscription"), Column: 5, Line: 4},
				{Type: schema.Colon, Value: []byte(":"), Column: 17, Line: 4},
				{Type: schema.Subscription, Value: []byte("Subscription"), Column: 19, Line: 4},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 5},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 5},
			},
		},
		{
			name: "Lex custom schema definition",
			input: []byte(`schema {
				query: RootQuery
				mutation: RootMutation
				subscription: RootSubscription
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedSchema, Value: []byte("schema"), Column: 1, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 8, Line: 1},
				{Type: schema.Field, Value: []byte("query"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 10, Line: 2},
				{Type: schema.Identifier, Value: []byte("RootQuery"), Column: 12, Line: 2},
				{Type: schema.Field, Value: []byte("mutation"), Column: 5, Line: 3},
				{Type: schema.Colon, Value: []byte(":"), Column: 13, Line: 3},
				{Type: schema.Identifier, Value: []byte("RootMutation"), Column: 15, Line: 3},
				{Type: schema.Field, Value: []byte("subscription"), Column: 5, Line: 4},
				{Type: schema.Colon, Value: []byte(":"), Column: 17, Line: 4},
				{Type: schema.Identifier, Value: []byte("RootSubscription"), Column: 19, Line: 4},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 5},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 5},
			},
		},
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
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
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
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
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
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("getUsers"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 13, Line: 2},
				{Type: schema.Field, Value: []byte("age"), Column: 14, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 17, Line: 2},
				{Type: schema.Identifier, Value: []byte("Int"), Column: 19, Line: 2},
				{Type: schema.Equal, Value: []byte("="), Column: 23, Line: 2},
				{Type: schema.Value, Value: []byte("18"), Column: 25, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Column: 27, Line: 2},
				{Type: schema.Field, Value: []byte("active"), Column: 29, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 35, Line: 2},
				{Type: schema.Identifier, Value: []byte("Boolean"), Column: 37, Line: 2},
				{Type: schema.Equal, Value: []byte("="), Column: 45, Line: 2},
				{Type: schema.Value, Value: []byte("true"), Column: 47, Line: 2},
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
			input: []byte(`type Mutation {
				createUser(name: String!): User
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Mutate, Value: []byte("Mutation"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 15, Line: 1},
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
				{Type: schema.Subscription, Value: []byte("Subscription"), Column: 6, Line: 1},
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
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("getUsers"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 13, Line: 2},
				{Type: schema.Field, Value: []byte("ids"), Column: 14, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 17, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 19, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 20, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 22, Line: 2},
				{Type: schema.Equal, Value: []byte("="), Column: 24, Line: 2},
				{Type: schema.Value, Value: []byte(`["id1", "id2", "id3"]`), Column: 26, Line: 2},
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
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("getIds"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 11, Line: 2},
				{Type: schema.Field, Value: []byte("ids"), Column: 12, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 15, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 17, Line: 2},
				{Type: schema.Identifier, Value: []byte("ID"), Column: 18, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 20, Line: 2},
				{Type: schema.Equal, Value: []byte("="), Column: 22, Line: 2},
				{Type: schema.Value, Value: []byte("[\"id1\", \"id2\"]"), Column: 24, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 38, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Field with default nested list value",
			input: []byte(`type Query {
				getIds(ids: [[String]] = [["id1", "id2"], ["id3", "id4", NULL]], hoge: String = "hoge")
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("getIds"), Column: 5, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 11, Line: 2},
				{Type: schema.Field, Value: []byte("ids"), Column: 12, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 15, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 17, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 18, Line: 2},
				{Type: schema.Identifier, Value: []byte("String"), Column: 19, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 25, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 26, Line: 2},
				{Type: schema.Equal, Value: []byte("="), Column: 28, Line: 2},
				{Type: schema.Value, Value: []byte(`[["id1", "id2"], ["id3", "id4", NULL]]`), Column: 30, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Column: 68, Line: 2},
				{Type: schema.Field, Value: []byte("hoge"), Column: 70, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 74, Line: 2},
				{Type: schema.Identifier, Value: []byte("String"), Column: 76, Line: 2},
				{Type: schema.Equal, Value: []byte("="), Column: 83, Line: 2},
				{Type: schema.Value, Value: []byte(`"hoge"`), Column: 85, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 91, Line: 2},
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
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
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
			name: "Directive with boolean argument",
			input: []byte(`type Query {
				users: [User] @include(if: true)
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
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
				{Type: schema.Value, Value: []byte("true"), Column: 32, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 36, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Directive with string argument",
			input: []byte(`type Query {
				users: [User] @deprecated(reason: "hogehoge")
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
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
				{Type: schema.Value, Value: []byte("\"hogehoge\""), Column: 39, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 49, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Directive with list argument",
			input: []byte(`type Query {
				users: [User] @hogehoge(reason: ["hogehoge"])
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("users"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 10, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 12, Line: 2},
				{Type: schema.Identifier, Value: []byte("User"), Column: 13, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 17, Line: 2},
				{Type: schema.At, Value: []byte("@"), Column: 19, Line: 2},
				{Type: schema.Identifier, Value: []byte("hogehoge"), Column: 20, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 28, Line: 2},
				{Type: schema.Field, Value: []byte("reason"), Column: 29, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 35, Line: 2},
				{Type: schema.Value, Value: []byte("[\"hogehoge\"]"), Column: 37, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 49, Line: 2},
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
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
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
				{Type: schema.Value, Value: []byte(`"use new field"`), Column: 39, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Column: 54, Line: 2},
				{Type: schema.Field, Value: []byte("version"), Column: 56, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 63, Line: 2},
				{Type: schema.Value, Value: []byte("2"), Column: 65, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 66, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Directive with complex argument",
			input: []byte(`type Query {
				users: [User] @hogehoge(reason: [["hogehoge"], ["fugafuga"]], version: 2, isActive: true)
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"), Column: 1, Line: 1},
				{Type: schema.Query, Value: []byte("Query"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 12, Line: 1},
				{Type: schema.Field, Value: []byte("users"), Column: 5, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 10, Line: 2},
				{Type: schema.BracketOpen, Value: []byte("["), Column: 12, Line: 2},
				{Type: schema.Identifier, Value: []byte("User"), Column: 13, Line: 2},
				{Type: schema.BracketClose, Value: []byte("]"), Column: 17, Line: 2},
				{Type: schema.At, Value: []byte("@"), Column: 19, Line: 2},
				{Type: schema.Identifier, Value: []byte("hogehoge"), Column: 20, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 28, Line: 2},
				{Type: schema.Field, Value: []byte("reason"), Column: 29, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 35, Line: 2},
				{Type: schema.Value, Value: []byte(`[["hogehoge"], ["fugafuga"]]`), Column: 37, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Column: 65, Line: 2},
				{Type: schema.Field, Value: []byte("version"), Column: 67, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 74, Line: 2},
				{Type: schema.Value, Value: []byte("2"), Column: 76, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Column: 77, Line: 2},
				{Type: schema.Field, Value: []byte("isActive"), Column: 79, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 87, Line: 2},
				{Type: schema.Value, Value: []byte("true"), Column: 89, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 93, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 3},
			},
		},
		{
			name: "Lex enum type",
			input: []byte(`enum Status {
					ACTIVE
					INACTIVE
					PENDING
			}`),
			expected: []*schema.Token{
				{Type: schema.Enum, Value: []byte("enum"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("Status"), Column: 6, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 13, Line: 1},
				{Type: schema.Identifier, Value: []byte("ACTIVE"), Column: 6, Line: 2},
				{Type: schema.Identifier, Value: []byte("INACTIVE"), Column: 6, Line: 3},
				{Type: schema.Identifier, Value: []byte("PENDING"), Column: 6, Line: 4},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 4, Line: 5},
				{Type: schema.EOF, Value: nil, Column: 5, Line: 5},
			},
		},
		{
			name:  "Lex union type1",
			input: []byte(`union SearchResult = User | Post`),
			expected: []*schema.Token{
				{Type: schema.Union, Value: []byte("union"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("SearchResult"), Column: 7, Line: 1},
				{Type: schema.Equal, Value: []byte("="), Column: 20, Line: 1},
				{Type: schema.Identifier, Value: []byte("User"), Column: 22, Line: 1},
				{Type: schema.Pipe, Value: []byte("|"), Column: 27, Line: 1},
				{Type: schema.Identifier, Value: []byte("Post"), Column: 29, Line: 1},
				{Type: schema.EOF, Value: nil, Column: 33, Line: 1},
			},
		},
		{
			name: "Lex union type2",
			input: []byte(`union SearchResult = User |
				Post |
				Comment`),
			expected: []*schema.Token{
				{Type: schema.Union, Value: []byte("union"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("SearchResult"), Column: 7, Line: 1},
				{Type: schema.Equal, Value: []byte("="), Column: 20, Line: 1},
				{Type: schema.Identifier, Value: []byte("User"), Column: 22, Line: 1},
				{Type: schema.Pipe, Value: []byte("|"), Column: 27, Line: 1},
				{Type: schema.Identifier, Value: []byte("Post"), Column: 5, Line: 2},
				{Type: schema.Pipe, Value: []byte("|"), Column: 10, Line: 2},
				{Type: schema.Identifier, Value: []byte("Comment"), Column: 5, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 12, Line: 3},
			},
		},
		{
			name: "Lex union type3",
			input: []byte(`union SearchResult = 
				| User
				| Post`),
			expected: []*schema.Token{
				{Type: schema.Union, Value: []byte("union"), Column: 1, Line: 1},
				{Type: schema.Identifier, Value: []byte("SearchResult"), Column: 7, Line: 1},
				{Type: schema.Equal, Value: []byte("="), Column: 20, Line: 1},
				{Type: schema.Pipe, Value: []byte("|"), Column: 5, Line: 2},
				{Type: schema.Identifier, Value: []byte("User"), Column: 7, Line: 2},
				{Type: schema.Pipe, Value: []byte("|"), Column: 5, Line: 3},
				{Type: schema.Identifier, Value: []byte("Post"), Column: 7, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 11, Line: 3},
			},
		},
		{
			name:  "Lex directive definition",
			input: []byte(`directive @example(arg: String) on FIELD | OBJECT`),
			expected: []*schema.Token{
				{Type: schema.ReservedDirective, Value: []byte("directive"), Column: 1, Line: 1},
				{Type: schema.At, Value: []byte("@"), Column: 11, Line: 1},
				{Type: schema.Identifier, Value: []byte("example"), Column: 12, Line: 1},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 19, Line: 1},
				{Type: schema.Field, Value: []byte("arg"), Column: 20, Line: 1},
				{Type: schema.Colon, Value: []byte(":"), Column: 23, Line: 1},
				{Type: schema.Identifier, Value: []byte("String"), Column: 25, Line: 1},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 31, Line: 1},
				{Type: schema.On, Value: []byte("on"), Column: 33, Line: 1},
				{Type: schema.DirectiveLocation, Value: []byte("FIELD"), Column: 36, Line: 1},
				{Type: schema.Pipe, Value: []byte("|"), Column: 42, Line: 1},
				{Type: schema.DirectiveLocation, Value: []byte("OBJECT"), Column: 44, Line: 1},
				{Type: schema.EOF, Value: nil, Column: 50, Line: 1},
			},
		},
		{
			name: "Simple directive definition (no args, single location)",
			input: []byte(`
				directive @deprecated on FIELD_DEFINITION
			`),
			expected: []*schema.Token{
				{Type: schema.ReservedDirective, Value: []byte("directive"), Column: 5, Line: 2},
				{Type: schema.At, Value: []byte("@"), Column: 15, Line: 2},
				{Type: schema.Identifier, Value: []byte("deprecated"), Column: 16, Line: 2},
				{Type: schema.On, Value: []byte("on"), Column: 27, Line: 2},
				{Type: schema.DirectiveLocation, Value: []byte("FIELD_DEFINITION"), Column: 30, Line: 2},
				{Type: schema.EOF, Value: nil, Column: 50, Line: 2},
			},
		},
		{
			name: "Directive definition with arguments, repeatable, multiple locations",
			input: []byte(`directive @auth(
					role: String = "USER",
					enabled: Boolean!
				) repeatable on FIELD_DEFINITION | OBJECT
			`),
			expected: []*schema.Token{
				{Type: schema.ReservedDirective, Value: []byte("directive"), Column: 1, Line: 1},
				{Type: schema.At, Value: []byte("@"), Column: 11, Line: 1},
				{Type: schema.Identifier, Value: []byte("auth"), Column: 12, Line: 1},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 16, Line: 1},
				{Type: schema.Field, Value: []byte("role"), Column: 6, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 10, Line: 2},
				{Type: schema.Identifier, Value: []byte("String"), Column: 12, Line: 2},
				{Type: schema.Equal, Value: []byte("="), Column: 19, Line: 2},
				{Type: schema.Value, Value: []byte(`"USER"`), Column: 21, Line: 2},
				{Type: schema.Comma, Value: []byte(","), Line: 2, Column: 27},
				{Type: schema.Field, Value: []byte("enabled"), Line: 3, Column: 6},
				{Type: schema.Colon, Value: []byte(":"), Line: 3, Column: 13},
				{Type: schema.Identifier, Value: []byte("Boolean"), Line: 3, Column: 15},
				{Type: schema.Exclamation, Value: []byte("!"), Line: 3, Column: 22},
				{Type: schema.ParenClose, Value: []byte(")"), Line: 4, Column: 5},
				{Type: schema.Repeatable, Value: []byte("repeatable"), Line: 2, Column: 58},
				{Type: schema.On, Value: []byte("on"), Line: 2, Column: 69},
				{Type: schema.DirectiveLocation, Value: []byte("FIELD_DEFINITION"), Line: 2, Column: 72},
				{Type: schema.Pipe, Value: []byte("|"), Line: 2, Column: 89},
				{Type: schema.DirectiveLocation, Value: []byte("OBJECT"), Line: 2, Column: 91},
				{Type: schema.EOF, Value: nil, Line: 2, Column: 101},
			},
		},
		{
			name: "Directive usage on field (valid usage)",
			input: []byte(`directive @deprecated(reason: String) on FIELD_DEFINITION

				type User {
					name: String @deprecated(reason: "Use fullName instead")
				}
			`),
			expected: []*schema.Token{
				{Type: schema.ReservedDirective, Value: []byte("directive"), Column: 1, Line: 1},
				{Type: schema.At, Value: []byte("@"), Column: 11, Line: 1},
				{Type: schema.Identifier, Value: []byte("deprecated"), Column: 12, Line: 1},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 22, Line: 1},
				{Type: schema.Field, Value: []byte("reason"), Column: 23, Line: 1},
				{Type: schema.Colon, Value: []byte(":"), Column: 29, Line: 1},
				{Type: schema.Identifier, Value: []byte("String"), Column: 31, Line: 1},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 37, Line: 1},
				{Type: schema.On, Value: []byte("on"), Column: 39, Line: 1},
				{Type: schema.DirectiveLocation, Value: []byte("FIELD_DEFINITION"), Column: 42, Line: 1},
				{Type: schema.ReservedType, Value: []byte("type"), Column: 64, Line: 1},
				{Type: schema.Identifier, Value: []byte("User"), Column: 69, Line: 1},
				{Type: schema.CurlyOpen, Value: []byte("{"), Column: 74, Line: 1},
				{Type: schema.Field, Value: []byte("name"), Column: 6, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 10, Line: 2},
				{Type: schema.Identifier, Value: []byte("String"), Column: 12, Line: 2},
				{Type: schema.At, Value: []byte("@"), Column: 19, Line: 2},
				{Type: schema.Identifier, Value: []byte("deprecated"), Column: 20, Line: 2},
				{Type: schema.ParenOpen, Value: []byte("("), Column: 30, Line: 2},
				{Type: schema.Field, Value: []byte("reason"), Column: 31, Line: 2},
				{Type: schema.Colon, Value: []byte(":"), Column: 37, Line: 2},
				{Type: schema.Value, Value: []byte(`"Use fullName instead"`), Column: 39, Line: 2},
				{Type: schema.ParenClose, Value: []byte(")"), Column: 61, Line: 2},
				{Type: schema.CurlyClose, Value: []byte("}"), Column: 5, Line: 3},
				{Type: schema.EOF, Value: nil, Column: 4, Line: 4},
			},
		},
		{
			name: "Directive definition with multiple arguments + usage with partial defaults",
			input: []byte(`
				directive @complex(
					flag: Boolean! = true,
					level: Int!,
					name: String
				) on OBJECT | FIELD_DEFINITION

				type Query @complex(level: 5) {
					test: String
				}
			`),
			expected: []*schema.Token{
				{Type: schema.ReservedDirective, Value: []byte("directive"), Column: 5,  Line: 2},
				{Type: schema.At,                Value: []byte("@"),         Column: 15, Line: 2},
				{Type: schema.Identifier,        Value: []byte("complex"),   Column: 16, Line: 2},
				{Type: schema.ParenOpen,         Value: []byte("("),         Column: 23, Line: 2},
				{Type: schema.Field,            Value: []byte("flag"),    Column: 6,  Line: 3},
				{Type: schema.Colon,            Value: []byte(":"),       Column: 10, Line: 3},
				{Type: schema.Identifier,       Value: []byte("Boolean"), Column: 12, Line: 3},
				{Type: schema.Exclamation,      Value: []byte("!"),       Column: 19, Line: 3},
				{Type: schema.Equal,            Value: []byte("="),       Column: 21, Line: 3},
				{Type: schema.Value,            Value: []byte("true"),    Column: 23, Line: 3},
				{Type: schema.Comma,            Value: []byte(","),       Column: 27, Line: 3},
				{Type: schema.Field,       Value: []byte("level"), Column: 6,  Line: 4},
				{Type: schema.Colon,       Value: []byte(":"),     Column: 11, Line: 4},
				{Type: schema.Identifier,  Value: []byte("Int"),   Column: 13, Line: 4},
				{Type: schema.Exclamation, Value: []byte("!"),     Column: 16, Line: 4},
				{Type: schema.Comma,       Value: []byte(","),     Column: 17, Line: 4},
				{Type: schema.Field,      Value: []byte("name"),   Column: 6,  Line: 5},
				{Type: schema.Colon,      Value: []byte(":"),      Column: 10, Line: 5},
				{Type: schema.Identifier, Value: []byte("String"), Column: 12, Line: 5},
				{Type: schema.ParenClose, Value: []byte(")"),      Column: 5,  Line: 6},
				{Type: schema.On,         Value: []byte("on"),     Column: 71, Line: 3},
				{Type: schema.DirectiveLocation, Value: []byte("OBJECT"),           Column: 74, Line: 3},
				{Type: schema.Pipe,              Value: []byte("|"),                Column: 81, Line: 3},
				{Type: schema.DirectiveLocation, Value: []byte("FIELD_DEFINITION"), Column: 83, Line: 3},
				{Type: schema.ReservedType,      Value: []byte("type"),             Column: 105, Line: 3},
				{Type: schema.Query,             Value: []byte("Query"),            Column: 110, Line: 3},
				{Type: schema.At,                Value: []byte("@"),                Column: 116, Line: 3},
				{Type: schema.Identifier,        Value: []byte("complex"),          Column: 117, Line: 3},
				{Type: schema.ParenOpen,         Value: []byte("("),                Column: 124, Line: 3},
				{Type: schema.Field,             Value: []byte("level"),            Column: 125, Line: 3},
				{Type: schema.Colon,             Value: []byte(":"),                Column: 130, Line: 3},
				{Type: schema.Value,             Value: []byte("5"),                Column: 132, Line: 3},
				{Type: schema.ParenClose,        Value: []byte(")"),                Column: 133, Line: 3},
				{Type: schema.CurlyOpen,         Value: []byte("{"),                Column: 135, Line: 3},
				{Type: schema.Field,             Value: []byte("test"),             Column: 6,   Line: 4},
				{Type: schema.Colon,             Value: []byte(":"),                Column: 10,  Line: 4},
				{Type: schema.Identifier,        Value: []byte("String"),           Column: 12,  Line: 4},
				{Type: schema.CurlyClose,        Value: []byte("}"),                Column: 5,   Line: 5},
				{Type: schema.EOF,               Value: nil,                        Column: 4,   Line: 6},
			},			
		},
		{
			name: "Repeatable directive used multiple times on one field",
			input: []byte(`
				directive @tag(label: String!) repeatable on FIELD_DEFINITION

				type Query {
					myField: String 
						@tag(label: "first") 
						@tag(label: "second")
				}
			`),
			expected: []*schema.Token{
				{Type: schema.ReservedDirective, Value: []byte("directive"),       Column: 5,  Line: 2},
				{Type: schema.At,                Value: []byte("@"),               Column: 15, Line: 2},
				{Type: schema.Identifier,        Value: []byte("tag"),             Column: 16, Line: 2},
				{Type: schema.ParenOpen,         Value: []byte("("),               Column: 19, Line: 2},
				{Type: schema.Field,       Value: []byte("label"),    Column: 20, Line: 2},
				{Type: schema.Colon,       Value: []byte(":"),        Column: 25, Line: 2},
				{Type: schema.Identifier,  Value: []byte("String"),   Column: 27, Line: 2},
				{Type: schema.Exclamation, Value: []byte("!"),        Column: 33, Line: 2},
				{Type: schema.ParenClose,  Value: []byte(")"),        Column: 34, Line: 2},
				{Type: schema.Repeatable,  Value: []byte("repeatable"), Column: 36, Line: 2},
				{Type: schema.On,          Value: []byte("on"),       Column: 47, Line: 2},
				{Type: schema.DirectiveLocation, Value: []byte("FIELD_DEFINITION"), Column: 50, Line: 2},
				{Type: schema.ReservedType,      Value: []byte("type"),            Column: 72, Line: 2},
				{Type: schema.Query,             Value: []byte("Query"),           Column: 77, Line: 2},
				{Type: schema.CurlyOpen,         Value: []byte("{"),               Column: 83, Line: 2},
				{Type: schema.Field,       Value: []byte("myField"),  Column: 6,  Line: 3},
				{Type: schema.Colon,       Value: []byte(":"),        Column: 13, Line: 3},
				{Type: schema.Identifier,  Value: []byte("String"),   Column: 15, Line: 3},
				{Type: schema.At,          Value: []byte("@"),        Column: 7,  Line: 4},
				{Type: schema.Identifier,  Value: []byte("tag"),      Column: 8,  Line: 4},
				{Type: schema.ParenOpen,   Value: []byte("("),        Column: 11, Line: 4},
				{Type: schema.Field,       Value: []byte("label"),    Column: 12, Line: 4},
				{Type: schema.Colon,       Value: []byte(":"),        Column: 17, Line: 4},
				{Type: schema.Value,       Value: []byte(`"first"`),  Column: 19, Line: 4},
				{Type: schema.ParenClose,  Value: []byte(")"),        Column: 26, Line: 4},
				{Type: schema.At,          Value: []byte("@"),        Column: 7,  Line: 5},
				{Type: schema.Identifier,  Value: []byte("tag"),      Column: 8,  Line: 5},
				{Type: schema.ParenOpen,   Value: []byte("("),        Column: 11, Line: 5},
				{Type: schema.Field,       Value: []byte("label"),    Column: 12, Line: 5},
				{Type: schema.Colon,       Value: []byte(":"),        Column: 17, Line: 5},
				{Type: schema.Value,       Value: []byte(`"second"`), Column: 19, Line: 5},
				{Type: schema.ParenClose,  Value: []byte(")"),        Column: 27, Line: 5},
				{Type: schema.CurlyClose,  Value: []byte("}"),        Column: 5,  Line: 6},
				{Type: schema.EOF,         Value: nil,                Column: 4,  Line: 7},
			},
		},
		{
			name: "Lex extend schema definition",
			input: []byte(`schema {
				query: Query
				mutation: Mutation
			}
		
			extend schema {
				subscription: MySubscription
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedSchema,  Value: []byte("schema"),       Line: 1, Column: 1},
				{Type: schema.CurlyOpen,       Value: []byte("{"),            Line: 1, Column: 8},
				{Type: schema.Field,           Value: []byte("query"),        Line: 2, Column: 5},
				{Type: schema.Colon,           Value: []byte(":"),            Line: 2, Column: 10},
				{Type: schema.Query,           Value: []byte("Query"),        Line: 2, Column: 12},
				{Type: schema.Field,           Value: []byte("mutation"),     Line: 3, Column: 5},
				{Type: schema.Colon,           Value: []byte(":"),            Line: 3, Column: 13},
				{Type: schema.Mutate,          Value: []byte("Mutation"),     Line: 3, Column: 15},
				{Type: schema.CurlyClose,      Value: []byte("}"),            Line: 4, Column: 4},
				{Type: schema.Extend,          Value: []byte("extend"),       Line: 6, Column: 4},
				{Type: schema.ReservedSchema,  Value: []byte("schema"),       Line: 6, Column: 11},
				{Type: schema.CurlyOpen,       Value: []byte("{"),            Line: 6, Column: 18},
				{Type: schema.Field,           Value: []byte("subscription"), Line: 7, Column: 5},
				{Type: schema.Colon,           Value: []byte(":"),            Line: 7, Column: 17},
				{Type: schema.Identifier,      Value: []byte("MySubscription"),Line:7, Column: 19},
				{Type: schema.CurlyClose,      Value: []byte("}"),            Line: 8, Column: 4},
				{Type: schema.EOF,             Value: nil,                    Line: 8, Column: 5},
			},
		},
		{
			name: "Lex extend type definition",
			input: []byte(`type User {
				id: ID!
			}
		
			extend type User {
				createdAt: String
			}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType,  Value: []byte("type"),    Line: 1, Column: 1},
				{Type: schema.Identifier,    Value: []byte("User"),    Line: 1, Column: 6},
				{Type: schema.CurlyOpen,     Value: []byte("{"),       Line: 1, Column: 11},
				{Type: schema.Field,         Value: []byte("id"),      Line: 2, Column: 5},
				{Type: schema.Colon,         Value: []byte(":"),       Line: 2, Column: 7},
				{Type: schema.Identifier,    Value: []byte("ID"),      Line: 2, Column: 9},
				{Type: schema.Exclamation,   Value: []byte("!"),       Line: 2, Column: 11},
				{Type: schema.CurlyClose,    Value: []byte("}"),       Line: 3, Column: 4},
				{Type: schema.Extend,        Value: []byte("extend"),  Line: 5, Column: 4},
				{Type: schema.ReservedType,  Value: []byte("type"),    Line: 5, Column: 11},
				{Type: schema.Identifier,    Value: []byte("User"),    Line: 5, Column: 16},
				{Type: schema.CurlyOpen,     Value: []byte("{"),       Line: 5, Column: 21},
				{Type: schema.Field,         Value: []byte("createdAt"),Line:6, Column: 5},
				{Type: schema.Colon,         Value: []byte(":"),       Line: 6, Column: 14},
				{Type: schema.Identifier,    Value: []byte("String"),  Line: 6, Column: 16},
				{Type: schema.CurlyClose,    Value: []byte("}"),       Line: 7, Column: 4},
				{Type: schema.EOF,           Value: nil,               Line: 7, Column: 5},
			},
		},
		{
			name: "Lex extend interface definition",
			input: []byte(`interface Node {
				id: ID!
			}
		
			extend interface Node {
				updatedAt: String
			}`),
			expected: []*schema.Token{
				{Type: schema.Interface,  Value: []byte("interface"), Line: 1, Column: 1},
				{Type: schema.Identifier, Value: []byte("Node"),      Line: 1, Column: 11},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line: 1, Column: 16},
				{Type: schema.Field,      Value: []byte("id"),        Line: 2, Column: 5},
				{Type: schema.Colon,      Value: []byte(":"),         Line: 2, Column: 7},
				{Type: schema.Identifier, Value: []byte("ID"),        Line: 2, Column: 9},
				{Type: schema.Exclamation,Value: []byte("!"),         Line: 2, Column: 11},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line: 3, Column: 4},
				{Type: schema.Extend,     Value: []byte("extend"),    Line: 5, Column: 4},
				{Type: schema.Interface,  Value: []byte("interface"), Line: 5, Column: 11},
				{Type: schema.Identifier, Value: []byte("Node"),      Line: 5, Column: 21},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line: 5, Column: 26},
				{Type: schema.Field,      Value: []byte("updatedAt"), Line: 6, Column: 5},
				{Type: schema.Colon,      Value: []byte(":"),         Line: 6, Column: 14},
				{Type: schema.Identifier, Value: []byte("String"),    Line: 6, Column: 16},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line: 7, Column: 4},
				{Type: schema.EOF,        Value: nil,                 Line: 7, Column: 5},
			},
		},
		{
			name: "Lex extend schema definition",
			input: []byte(`schema {
			query: Query
			mutation: Mutation
		}
		
		extend schema {
			subscription: MySubscription
		}`),
			expected: []*schema.Token{
				{Type: schema.ReservedSchema, Value: []byte("schema"),       Line: 1, Column: 1},
				{Type: schema.CurlyOpen,      Value: []byte("{"),            Line: 1, Column: 8},
				{Type: schema.Field,          Value: []byte("query"),        Line: 2, Column: 4},
				{Type: schema.Colon,          Value: []byte(":"),            Line: 2, Column: 9},
				{Type: schema.Query,          Value: []byte("Query"),        Line: 2, Column: 11},
				{Type: schema.Field,          Value: []byte("mutation"),     Line: 3, Column: 4},
				{Type: schema.Colon,          Value: []byte(":"),            Line: 3, Column: 12},
				{Type: schema.Mutate,         Value: []byte("Mutation"),     Line: 3, Column: 14},
				{Type: schema.CurlyClose,     Value: []byte("}"),            Line: 4, Column: 3},
				{Type: schema.Extend,         Value: []byte("extend"),       Line: 6, Column: 3},
				{Type: schema.ReservedSchema, Value: []byte("schema"),       Line: 6, Column: 10},
				{Type: schema.CurlyOpen,      Value: []byte("{"),            Line: 6, Column: 17},
				{Type: schema.Field,          Value: []byte("subscription"), Line: 7, Column: 4},
				{Type: schema.Colon,          Value: []byte(":"),            Line: 7, Column: 16},
				{Type: schema.Identifier,     Value: []byte("MySubscription"),Line:7,Column: 18},
				{Type: schema.CurlyClose,     Value: []byte("}"),            Line: 8, Column: 3},
				{Type: schema.EOF,            Value: nil,                    Line: 8, Column: 4},
			},
		},
		{
			name: "Lex extend type definition",
			input: []byte(`type User {
			id: ID!
		}
		
		extend type User {
			createdAt: String
		}`),
			expected: []*schema.Token{
				{Type: schema.ReservedType, Value: []byte("type"),  Line: 1, Column: 1},
				{Type: schema.Identifier,   Value: []byte("User"),  Line: 1, Column: 6},
				{Type: schema.CurlyOpen,    Value: []byte("{"),     Line: 1, Column: 11},
				{Type: schema.Field,        Value: []byte("id"),    Line: 2, Column: 4},
				{Type: schema.Colon,        Value: []byte(":"),     Line: 2, Column: 6},
				{Type: schema.Identifier,   Value: []byte("ID"),    Line: 2, Column: 8},
				{Type: schema.Exclamation,  Value: []byte("!"),     Line: 2, Column: 10},
				{Type: schema.CurlyClose,   Value: []byte("}"),     Line: 3, Column: 3},
				{Type: schema.Extend,       Value: []byte("extend"),Line: 5, Column: 3},
				{Type: schema.ReservedType, Value: []byte("type"),  Line: 5, Column: 10},
				{Type: schema.Identifier,   Value: []byte("User"),  Line: 5, Column: 15},
				{Type: schema.CurlyOpen,    Value: []byte("{"),     Line: 5, Column: 20},
				{Type: schema.Field,        Value: []byte("createdAt"),Line:6, Column:4},
				{Type: schema.Colon,        Value: []byte(":"),        Line:6, Column:13},
				{Type: schema.Identifier,   Value: []byte("String"),   Line:6, Column:15},
				{Type: schema.CurlyClose,   Value: []byte("}"),        Line: 7, Column: 3},
				{Type: schema.EOF,          Value: nil,                Line: 7, Column: 4},
			},
		},
		{
			name: "Lex extend interface definition",
			input: []byte(`interface Node {
			id: ID!
		}
		
		extend interface Node {
			updatedAt: String
		}`),
			expected: []*schema.Token{
				{Type: schema.Interface,  Value: []byte("interface"), Line:1, Column:1},
				{Type: schema.Identifier, Value: []byte("Node"),      Line:1, Column:11},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:1, Column:16},
				{Type: schema.Field,      Value: []byte("id"),        Line:2, Column:4},
				{Type: schema.Colon,      Value: []byte(":"),         Line:2, Column:6},
				{Type: schema.Identifier, Value: []byte("ID"),        Line:2, Column:8},
				{Type: schema.Exclamation,Value: []byte("!"),         Line:2, Column:10},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:3, Column:3},
				{Type: schema.Extend,     Value: []byte("extend"),    Line:5, Column:3},
				{Type: schema.Interface,  Value: []byte("interface"), Line:5, Column:10},
				{Type: schema.Identifier, Value: []byte("Node"),      Line:5, Column:20},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:5, Column:25},
				{Type: schema.Field,      Value: []byte("updatedAt"), Line:6, Column:4},
				{Type: schema.Colon,      Value: []byte(":"),         Line:6, Column:13},
				{Type: schema.Identifier, Value: []byte("String"),    Line:6, Column:15},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:7, Column:3},
				{Type: schema.EOF,        Value: nil,                 Line:7, Column: 4},
			},
		},
		{
			name: "Lex extend union definition",
			input: []byte(`union SearchResult = User | Post
		
		extend union SearchResult = Comment | Page
		`),
			expected: []*schema.Token{
				{Type: schema.Union,      Value: []byte("union"),        Line:1, Column:1},
				{Type: schema.Identifier, Value: []byte("SearchResult"), Line:1, Column:7},
				{Type: schema.Equal,      Value: []byte("="),            Line:1, Column:20},
				{Type: schema.Identifier, Value: []byte("User"),         Line:1, Column:22},
				{Type: schema.Pipe,       Value: []byte("|"),            Line:1, Column:27},
				{Type: schema.Identifier, Value: []byte("Post"),         Line:1, Column:29},
				{Type: schema.Extend,     Value: []byte("extend"),       Line:3, Column:3},
				{Type: schema.Union,      Value: []byte("union"),        Line:3, Column:10},
				{Type: schema.Identifier, Value: []byte("SearchResult"), Line:3, Column:16},
				{Type: schema.Equal,      Value: []byte("="),            Line:3, Column:29},
				{Type: schema.Identifier, Value: []byte("Comment"),      Line:3, Column:31},
				{Type: schema.Pipe,       Value: []byte("|"),            Line:3, Column:39},
				{Type: schema.Identifier, Value: []byte("Page"),         Line:3, Column:41},
				{Type: schema.EOF,        Value: nil,                    Line:4, Column:3},
			},
		},
		{
			name: "Lex extend enum definition",
			input: []byte(`enum Role {
			ADMIN
			USER
		}
		
		extend enum Role {
			GUEST
			SUPERADMIN
		}`),
			expected: []*schema.Token{
				{Type: schema.Enum,       Value: []byte("enum"),  Line:1, Column:1},
				{Type: schema.Identifier, Value: []byte("Role"),  Line:1, Column:6},
				{Type: schema.CurlyOpen,  Value: []byte("{"),     Line:1, Column:11},
				{Type: schema.Identifier, Value: []byte("ADMIN"), Line:2, Column:4},
				{Type: schema.Identifier, Value: []byte("USER"),  Line:3, Column:4},
				{Type: schema.CurlyClose, Value: []byte("}"),     Line:4, Column:3},
				{Type: schema.Extend,     Value: []byte("extend"),Line:6, Column:3},
				{Type: schema.Enum,       Value: []byte("enum"),  Line:6, Column:10},
				{Type: schema.Identifier, Value: []byte("Role"),  Line:6, Column:15},
				{Type: schema.CurlyOpen,  Value: []byte("{"),     Line:6, Column:20},
				{Type: schema.Identifier, Value: []byte("GUEST"), Line:7, Column:4},
				{Type: schema.Identifier, Value: []byte("SUPERADMIN"),Line:8,Column:4},
				{Type: schema.CurlyClose, Value: []byte("}"),     Line:9, Column:3},
				{Type: schema.EOF,        Value: nil,             Line:9, Column:4},
			},
		},
		{
			name: "Lex extend input definition",
			input: []byte(`input UserInput {
			id: ID!
			name: String
		}
		
		extend input UserInput {
			age: Int
			email: String
		}`),
			expected: []*schema.Token{
				{Type: schema.Input,      Value: []byte("input"),     Line:1, Column:1},
				{Type: schema.Identifier, Value: []byte("UserInput"), Line:1, Column:7},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:1, Column:17},
				{Type: schema.Field,      Value: []byte("id"),        Line:2, Column:4},
				{Type: schema.Colon,      Value: []byte(":"),         Line:2, Column:6},
				{Type: schema.Identifier, Value: []byte("ID"),        Line:2, Column:8},
				{Type: schema.Exclamation,Value: []byte("!"),         Line:2, Column:10},
				{Type: schema.Field,      Value: []byte("name"),      Line:3, Column:4},
				{Type: schema.Colon,      Value: []byte(":"),         Line:3, Column:8},
				{Type: schema.Identifier, Value: []byte("String"),    Line:3, Column:10},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:4, Column:3},
				{Type: schema.Extend,     Value: []byte("extend"),    Line:6, Column:3},
				{Type: schema.Input,      Value: []byte("input"),     Line:6, Column:10},
				{Type: schema.Identifier, Value: []byte("UserInput"), Line:6, Column:16},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:6, Column:26},
				{Type: schema.Field,      Value: []byte("age"),       Line:7, Column:4},
				{Type: schema.Colon,      Value: []byte(":"),         Line:7, Column:7},
				{Type: schema.Identifier, Value: []byte("Int"),       Line:7, Column:9},
				{Type: schema.Field,      Value: []byte("email"),     Line:8, Column:4},
				{Type: schema.Colon,      Value: []byte(":"),         Line:8, Column:9},
				{Type: schema.Identifier, Value: []byte("String"),    Line:8, Column:11},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:9, Column:3},
				{Type: schema.EOF,        Value: nil,                 Line:9, Column:4},
			},
		},
		{
			name: "Lex directive usage on enum value",
			input: []byte(`enum Direction {
			NORTH
			EAST @deprecated(reason: "No longer used")
			SOUTH
			WEST
		}`),
			expected: []*schema.Token{
				{Type: schema.Enum,            Value: []byte("enum"),   Line:1, Column:1},
				{Type: schema.Identifier,      Value: []byte("Direction"), Line:1, Column:6},
				{Type: schema.CurlyOpen,       Value: []byte("{"),      Line:1, Column:16},
				{Type: schema.Identifier,      Value: []byte("NORTH"),  Line:2, Column:4},
				{Type: schema.Identifier,      Value: []byte("EAST"),   Line:3, Column:4},
				{Type: schema.At,              Value: []byte("@"),      Line:3, Column:9},
				{Type: schema.Identifier,      Value: []byte("deprecated"), Line:3,Column:9},
				{Type: schema.ParenOpen,       Value: []byte("("),      Line:3, Column:19},
				{Type: schema.Field,           Value: []byte("reason"), Line:3, Column:20},
				{Type: schema.Colon,           Value: []byte(":"),      Line:3, Column:26},
				{Type: schema.Value,           Value: []byte(`"No longer used"`), Line:3, Column:28},
				{Type: schema.ParenClose,      Value: []byte(")"),      Line:3, Column:44},
				{Type: schema.Identifier,      Value: []byte("SOUTH"),  Line:4, Column:4},
				{Type: schema.Identifier,      Value: []byte("WEST"),   Line:5, Column:4},
				{Type: schema.CurlyClose,      Value: []byte("}"),      Line:6, Column:3},
				{Type: schema.EOF,             Value: nil,              Line:6, Column:4},
			},
		},
		{
			name: "Lex directive usage on interface",
			input: []byte(`interface Node @auth(role: "ADMIN") {
			id: ID!
		}`),
			expected: []*schema.Token{
				{Type: schema.Interface, Value: []byte("interface"), Line:1, Column:1},
				{Type: schema.Identifier,Value: []byte("Node"),      Line:1, Column:11},
				{Type: schema.At,        Value: []byte("@"),         Line:1, Column:16},
				{Type: schema.Identifier,Value: []byte("auth"),      Line:1, Column:17},
				{Type: schema.ParenOpen, Value: []byte("("),         Line:1, Column:21},
				{Type: schema.Field,     Value: []byte("role"),      Line:1, Column:22},
				{Type: schema.Colon,     Value: []byte(":"),         Line:1, Column:26},
				{Type: schema.Value,     Value: []byte(`"ADMIN"`),   Line:1, Column:28},
				{Type: schema.ParenClose,Value: []byte(")"),         Line:1, Column:35},
				{Type: schema.CurlyOpen, Value: []byte("{"),         Line:1, Column:37},
				{Type: schema.Field,     Value: []byte("id"),        Line:2, Column:4},
				{Type: schema.Colon,     Value: []byte(":"),         Line:2, Column:6},
				{Type: schema.Identifier,Value: []byte("ID"),        Line:2, Column:8},
				{Type: schema.Exclamation,Value:[]byte("!"),         Line:2, Column:10},
				{Type: schema.CurlyClose,Value: []byte("}"),         Line:3, Column:3},
				{Type: schema.EOF,       Value: nil,                 Line:3, Column:4},
			},
		},
		{
			name: "Lex directive usage on schema",
			input: []byte(`schema @example {
			query: Query
		}`),
			expected: []*schema.Token{
				{Type: schema.ReservedSchema, Value: []byte("schema"), Line:1, Column:1},
				{Type: schema.At,             Value: []byte("@"),      Line:1, Column:8},
				{Type: schema.Identifier,     Value: []byte("example"),Line:1, Column:9},
				{Type: schema.CurlyOpen,      Value: []byte("{"),      Line:1, Column:17},
				{Type: schema.Field,          Value: []byte("query"),  Line:2, Column:4},
				{Type: schema.Colon,          Value: []byte(":"),      Line:2, Column:9},
				{Type: schema.Query,          Value: []byte("Query"),  Line:2, Column:11},
				{Type: schema.CurlyClose,     Value: []byte("}"),      Line:3, Column:3},
				{Type: schema.EOF,            Value: nil,              Line:3, Column:4},
			},
		},
		{
			name: "Lex directive usage on union",
			input: []byte(`union Entity @deprecated(reason: "Use another union") = User | Post`),
			expected: []*schema.Token{
				{Type: schema.Union,       Value: []byte("union"),   Line:1, Column:1},
				{Type: schema.Identifier,  Value: []byte("Entity"),  Line:1, Column:7},
				{Type: schema.At,          Value: []byte("@"),       Line:1, Column:14},
				{Type: schema.Identifier,  Value: []byte("deprecated"), Line:1, Column:15},
				{Type: schema.ParenOpen,   Value: []byte("("),       Line:1, Column:25},
				{Type: schema.Field,       Value: []byte("reason"),  Line:1, Column:26},
				{Type: schema.Colon,       Value: []byte(":"),       Line:1, Column:32},
				{Type: schema.Value,       Value: []byte(`"Use another union"`), Line:1, Column:34},
				{Type: schema.ParenClose,  Value: []byte(")"),       Line:1, Column:53},
				{Type: schema.Equal,       Value: []byte("="),       Line:1, Column:55},
				{Type: schema.Identifier,  Value: []byte("User"),    Line:1, Column:57},
				{Type: schema.Pipe,        Value: []byte("|"),       Line:1, Column:62},
				{Type: schema.Identifier,  Value: []byte("Post"),    Line:1, Column:64},
				{Type: schema.EOF,         Value: nil,               Line:1, Column:68},
			},
		},
		{
			name: "Lex simple scalar definition",
			input: []byte(`scalar DateTime`),
			expected: []*schema.Token{
				{Type: schema.Scalar,      Value: []byte("scalar"),   Line:1, Column:1},
				{Type: schema.Identifier,  Value: []byte("DateTime"), Line:1, Column:8},
				{Type: schema.EOF,         Value: nil,                 Line:1, Column:16},
			},
		},{
			name: "Parse scalar with directive",
			input: []byte(`scalar URL @specifiedBy(url: "https://example.com/url-spec")`),
			expected: []*schema.Token{
				{Type: schema.Scalar,      Value: []byte("scalar"),   Line:1, Column:1},
				{Type: schema.Identifier,  Value: []byte("URL"),      Line:1, Column:8},
				{Type: schema.At,          Value: []byte("@"),         Line:1, Column:12},
				{Type: schema.Identifier,  Value: []byte("specifiedBy"), Line:1, Column:13},
				{Type: schema.ParenOpen,   Value: []byte("("),         Line:1, Column:24},
				{Type: schema.Field,       Value: []byte("url"),       Line:1, Column:25},
				{Type: schema.Colon,       Value: []byte(":"),         Line:1, Column:28},
				{Type: schema.Value,       Value: []byte(`"https://example.com/url-spec"`), Line:1, Column:30},
				{Type: schema.ParenClose,  Value: []byte(")"),         Line:1, Column:60},
				{Type: schema.EOF,         Value: nil,                 Line:1, Column:61},
			},
		},
		{
			name: "Parse scalar with multiple directives",
			input: []byte(`scalar JSON 
				@specifiedBy(url: "https://example.com/json-spec") 
				@deprecated(reason: "Prefer using JSON2")`),
			expected: []*schema.Token{
				{Type: schema.Scalar,      Value: []byte("scalar"),   Line:1, Column:1},
				{Type: schema.Identifier,  Value: []byte("JSON"),     Line:1, Column:8},
				{Type: schema.At,          Value: []byte("@"),         Line:2, Column:5},
				{Type: schema.Identifier,  Value: []byte("specifiedBy"), Line:2, Column:6},
				{Type: schema.ParenOpen,   Value: []byte("("),         Line:2, Column:17},
				{Type: schema.Field,       Value: []byte("url"),       Line:2, Column:18},
				{Type: schema.Colon,       Value: []byte(":"),         Line:2, Column:21},
				{Type: schema.Value,       Value: []byte(`"https://example.com/json-spec"`), Line:2, Column:23},
				{Type: schema.ParenClose,  Value: []byte(")"),         Line:2, Column:54},
				{Type: schema.At,          Value: []byte("@"),         Line:3, Column:5},
				{Type: schema.Identifier,  Value: []byte("deprecated"), Line:3, Column:6},
				{Type: schema.ParenOpen,   Value: []byte("("),         Line:3, Column:16},
				{Type: schema.Field,       Value: []byte("reason"),    Line:3, Column:17},
				{Type: schema.Colon,       Value: []byte(":"),         Line:3, Column:23},
				{Type: schema.Value,       Value: []byte(`"Prefer using JSON2"`), Line:3, Column:25},
				{Type: schema.ParenClose,  Value: []byte(")"),         Line:3, Column:45},
				{Type: schema.EOF,         Value: nil,                 Line:3, Column:46},
			},
		},
		{
			name: "Lex that implements a single interface",
			input: []byte(`
				interface Node {
					id: ID!
				}

				type User implements Node {
					id: ID!
					name: String
				}
			`),
			expected: []*schema.Token{
				{Type: schema.Interface,  Value: []byte("interface"), Line:2, Column:5},
				{Type: schema.Identifier, Value: []byte("Node"),      Line:2, Column:15},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:2, Column:20},
				{Type: schema.Field,      Value: []byte("id"),        Line:3, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:3, Column:8},
				{Type: schema.Identifier, Value: []byte("ID"),        Line:3, Column:10},
				{Type: schema.Exclamation,Value: []byte("!"),         Line:3, Column:12},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:4, Column:5},
				{Type: schema.ReservedType,Value: []byte("type"),     Line:6, Column:5},
				{Type: schema.Identifier, Value: []byte("User"),      Line:6, Column:10},
				{Type: schema.Implements, Value: []byte("implements"),Line:6, Column:15},
				{Type: schema.Identifier, Value: []byte("Node"),      Line:6, Column:26},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:6, Column:31},
				{Type: schema.Field,      Value: []byte("id"),        Line:7, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:7, Column:8},
				{Type: schema.Identifier, Value: []byte("ID"),        Line:7, Column:10},
				{Type: schema.Exclamation,Value: []byte("!"),         Line:7, Column:12},
				{Type: schema.Field,      Value: []byte("name"),      Line:8, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:8, Column:10},
				{Type: schema.Identifier, Value: []byte("String"),    Line:8, Column:12},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:9, Column:5},
				{Type: schema.EOF,        Value: nil,                 Line:10, Column:4},
			},
		},
		{
			name: "Parse type that implements multiple interfaces",
			input: []byte(`interface Node {
					id: ID!
				}
				interface Timestamp {
					createdAt: String
					updatedAt: String
				}

				type User implements Node & Timestamp {
					id: ID!
					name: String
					createdAt: String
					updatedAt: String
				}
			`),
			expected: []*schema.Token{
				{Type: schema.Interface,  Value: []byte("interface"), Line:1, Column:1},
				{Type: schema.Identifier, Value: []byte("Node"),      Line:1, Column:11},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:1, Column:16},
				{Type: schema.Field,      Value: []byte("id"),        Line:2, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:2, Column:8},
				{Type: schema.Identifier, Value: []byte("ID"),        Line:2, Column:10},
				{Type: schema.Exclamation,Value: []byte("!"),         Line:2, Column:12},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:3, Column:5},
				{Type: schema.Interface,  Value: []byte("interface"), Line:4, Column:5},
				{Type: schema.Identifier, Value: []byte("Timestamp"), Line:4, Column:15},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:4, Column:25},
				{Type: schema.Field,      Value: []byte("createdAt"), Line:5, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:5, Column:15},
				{Type: schema.Identifier, Value: []byte("String"),    Line:5, Column:17},
				{Type: schema.Field,      Value: []byte("updatedAt"), Line:6, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:6, Column:15},
				{Type: schema.Identifier, Value: []byte("String"),    Line:6, Column:17},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:7, Column:5},
				{Type: schema.ReservedType,Value: []byte("type"),     Line:9, Column:5},
				{Type: schema.Identifier, Value: []byte("User"),      Line:9, Column:10},
				{Type: schema.Implements, Value: []byte("implements"),Line:9, Column:15},
				{Type: schema.Identifier, Value: []byte("Node"),      Line:9, Column:26},
				{Type: schema.And,  Value: []byte("&"),         Line:9, Column:31},
				{Type: schema.Identifier, Value: []byte("Timestamp"), Line:9, Column:33},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:9, Column:43},
				{Type: schema.Field,      Value: []byte("id"),        Line:10, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:10, Column:8},
				{Type: schema.Identifier, Value: []byte("ID"),        Line:10, Column:10},
				{Type: schema.Exclamation,Value: []byte("!"),         Line:10, Column:12},
				{Type: schema.Field,      Value: []byte("name"),      Line:11, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:11, Column:10},
				{Type: schema.Identifier, Value: []byte("String"),    Line:11, Column:12},
				{Type: schema.Field,      Value: []byte("createdAt"), Line:12, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:12, Column:15},
				{Type: schema.Identifier, Value: []byte("String"),    Line:12, Column:17},
				{Type: schema.Field,      Value: []byte("updatedAt"), Line:13, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:13, Column:15},
				{Type: schema.Identifier, Value: []byte("String"),    Line:13, Column:17},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:14, Column:5},
				{Type: schema.EOF,        Value: nil,                 Line:15, Column:4},
			},
		},
		{
			name: "Parse type implements interface with directive on type",
			input: []byte(`interface Node {
					id: ID!
				}
		
				type User implements Node @key(fields: "id") {
					id: ID!
					name: String
				}
			`),
			expected: []*schema.Token{
				{Type: schema.Interface,  Value: []byte("interface"), Line:1, Column:1},
				{Type: schema.Identifier, Value: []byte("Node"),      Line:1, Column:11},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:1, Column:16},
				{Type: schema.Field,      Value: []byte("id"),        Line:2, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:2, Column:8},
				{Type: schema.Identifier, Value: []byte("ID"),        Line:2, Column:10},
				{Type: schema.Exclamation,Value: []byte("!"),         Line:2, Column:12},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:3, Column:5},
				{Type: schema.ReservedType,Value: []byte("type"),     Line:5, Column:5},
				{Type: schema.Identifier, Value: []byte("User"),      Line:5, Column:10},
				{Type: schema.Implements, Value: []byte("implements"),Line:5, Column:15},
				{Type: schema.Identifier, Value: []byte("Node"),      Line:5, Column:26},
				{Type: schema.At,          Value: []byte("@"),         Line:5, Column:31},
				{Type: schema.Identifier, Value: []byte("key"), 		 Line:5, Column:32},
				{Type: schema.ParenOpen,   Value: []byte("("),         Line:5, Column:35},
				{Type: schema.Field,       Value: []byte("fields"),    Line:5, Column:36},
				{Type: schema.Colon,       Value: []byte(":"),         Line:5, Column:42},
				{Type: schema.Value,       Value: []byte(`"id"`),      Line:5, Column:44},
				{Type: schema.ParenClose,  Value: []byte(")"),         Line:5, Column:48},
				{Type: schema.CurlyOpen,   Value: []byte("{"),         Line:5, Column:50},
				{Type: schema.Field,       Value: []byte("id"),        Line:6, Column:6},
				{Type: schema.Colon,       Value: []byte(":"),         Line:6, Column:8},
				{Type: schema.Identifier, Value: []byte("ID"),        Line:6, Column:10},
				{Type: schema.Exclamation,Value: []byte("!"),         Line:6, Column:12},
				{Type: schema.Field,       Value: []byte("name"),      Line:7, Column:6},
				{Type: schema.Colon,       Value: []byte(":"),         Line:7, Column:10},
				{Type: schema.Identifier, Value: []byte("String"),    Line:7, Column:12},
				{Type: schema.CurlyClose,  Value: []byte("}"),         Line:8, Column:5},
				{Type: schema.EOF,         Value: nil,                 Line:9, Column:4},
			},
		},
		{
			name: "Parse type implements multiple interfaces with directive",
			input: []byte(`interface Node {
					id: ID!
				}
				interface Timestamp {
					createdAt: String
					updatedAt: String
				}
		
				type User implements Node & Timestamp @anotherDirective {
					id: ID!
					name: String
					createdAt: String
					updatedAt: String
				}
			`),
			expected: []*schema.Token{
				{Type: schema.Interface,  Value: []byte("interface"), Line:1, Column:1},
				{Type: schema.Identifier, Value: []byte("Node"),      Line:1, Column:11},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:1, Column:16},
				{Type: schema.Field,      Value: []byte("id"),        Line:2, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:2, Column:8},
				{Type: schema.Identifier, Value: []byte("ID"),        Line:2, Column:10},
				{Type: schema.Exclamation,Value: []byte("!"),         Line:2, Column:12},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:3, Column:5},
				{Type: schema.Interface,  Value: []byte("interface"), Line:4, Column:5},
				{Type: schema.Identifier, Value: []byte("Timestamp"), Line:4, Column:15},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:4, Column:25},
				{Type: schema.Field,      Value: []byte("createdAt"), Line:5, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:5, Column:15},
				{Type: schema.Identifier, Value: []byte("String"),    Line:5, Column:17},
				{Type: schema.Field,      Value: []byte("updatedAt"), Line:6, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:6, Column:15},
				{Type: schema.Identifier, Value: []byte("String"),    Line:6, Column:17},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:7, Column:5},
				{Type: schema.ReservedType,Value: []byte("type"),     Line:9, Column:5},
				{Type: schema.Identifier, Value: []byte("User"),      Line:9, Column:10},
				{Type: schema.Implements, Value: []byte("implements"),Line:9, Column:15},
				{Type: schema.Identifier, Value: []byte("Node"),      Line:9, Column:26},
				{Type: schema.And,  Value: []byte("&"),         Line:9, Column:31},
				{Type: schema.Identifier, Value: []byte("Timestamp"), Line:9, Column:33},
				{Type: schema.At,          Value: []byte("@"),         Line:9, Column:43},
				{Type: schema.Identifier, Value: []byte("anotherDirective"), Line:9, Column:44},
				{Type: schema.CurlyOpen,  Value: []byte("{"),         Line:9, Column:61},
				{Type: schema.Field,      Value: []byte("id"),        Line:10, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:10, Column:8},
				{Type: schema.Identifier, Value: []byte("ID"),        Line:10, Column:10},
				{Type: schema.Exclamation,Value: []byte("!"),         Line:10, Column:12},
				{Type: schema.Field,      Value: []byte("name"),      Line:11, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:11, Column:10},
				{Type: schema.Identifier, Value: []byte("String"),    Line:11, Column:12},
				{Type: schema.Field,      Value: []byte("createdAt"), Line:12, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:12, Column:15},
				{Type: schema.Identifier, Value: []byte("String"),    Line:12, Column:17},
				{Type: schema.Field,      Value: []byte("updatedAt"), Line:13, Column:6},
				{Type: schema.Colon,      Value: []byte(":"),         Line:13, Column:15},
				{Type: schema.Identifier, Value: []byte("String"),    Line:13, Column:17},
				{Type: schema.CurlyClose, Value: []byte("}"),         Line:14, Column:5},
				{Type: schema.EOF,         Value: nil,                 Line:15, Column:4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := schema.NewLexer()
			got, err := lexer.Lex(tt.input)
			if tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && err != nil {
				t.Errorf("Parse() error %v", err)
				return
			}

			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Errorf("Lex() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
