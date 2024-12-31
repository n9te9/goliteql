package query_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/lkeix/gg-parser/query"
)

func TestQueryLex(t *testing.T) {
	tests := []struct{
		name string
		input []byte
		expected query.Tokens
		wantErr error
	}{
		{
			name: "Lex simple graphql query",
			input: []byte(`query MyQuery {
				user(id: 123) {
					name
				}
			}`),
			expected: query.Tokens{
				{Type: query.Query, Value: []byte("query"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("MyQuery"), Line: 1, Column: 7},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 15},
				{Type: query.Name, Value: []byte("user"), Line: 2, Column: 5},
				{Type: query.ParenOpen, Value: []byte("("), Line: 2, Column: 9},
				{Type: query.Name, Value: []byte("id"), Line: 2, Column: 10},
				{Type: query.Colon, Value: []byte(":"), Line: 2, Column: 12},
				{Type: query.Name, Value: []byte("123"), Line: 2, Column: 14},
				{Type: query.ParenClose, Value: []byte(")"), Line: 2, Column: 17},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 2, Column: 19},
				{Type: query.Name, Value: []byte("name"), Line: 3, Column: 6},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 4, Column: 5},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 5, Column: 4},
				{Type: query.EOF, Value: nil, Line: 5, Column: 5},
			},
		},
		{
			name: "Lex query with variables",
			input: []byte(`query GetUser($id: ID!, $type: String) {
		user(id: $id, type: $type) {
			name
		}
	}`),
			expected: query.Tokens{
				{Type: query.Query,      Value: []byte("query"),   Line: 1, Column: 1},
				{Type: query.Name,       Value: []byte("GetUser"), Line: 1, Column: 7},
	
				// 引数定義開始
				{Type: query.ParenOpen,  Value: []byte("("),       Line: 1, Column: 14},
				// $id: ID!
				{Type: query.Dollar, 	  Value: []byte("$"),       Line: 1, Column: 15},
				{Type: query.Name,   Value: []byte("id"),      Line: 1, Column: 16},
				{Type: query.Colon,      Value: []byte(":"),       Line: 1, Column: 18},
				{Type: query.Name,       Value: []byte("ID"),      Line: 1, Column: 20},
				{Type: query.Exclamation,       Value: []byte("!"),       Line: 1, Column: 22},
				{Type: query.Comma,      Value: []byte(","),       Line: 1, Column: 23},
				// $type: String
				{Type: query.Dollar, 	  Value: []byte("$"),       Line: 1, Column: 25},
				{Type: query.Name,   Value: []byte("type"),    Line: 1, Column: 26},
				{Type: query.Colon,      Value: []byte(":"),       Line: 1, Column: 30},
				{Type: query.Name,       Value: []byte("String"),  Line: 1, Column: 32},
				{Type: query.ParenClose, Value: []byte(")"),       Line: 1, Column: 38},
	
				// 本体
				{Type: query.CurlyOpen,  Value: []byte("{"),       Line: 1, Column: 40},
				{Type: query.Name,       Value: []byte("user"),    Line: 2, Column: 3},
				{Type: query.ParenOpen,  Value: []byte("("),       Line: 2, Column: 7},
				{Type: query.Name,       Value: []byte("id"),      Line: 2, Column: 8},
				{Type: query.Colon,      Value: []byte(":"),       Line: 2, Column: 10},
				{Type: query.Dollar,     Value: []byte("$"),       Line: 2, Column: 12},
				{Type: query.Name,       Value: []byte("id"),      Line: 2, Column: 13},
				{Type: query.Comma,      Value: []byte(","),       Line: 2, Column: 15},
				{Type: query.Name,       Value: []byte("type"),    Line: 2, Column: 17},
				{Type: query.Colon,      Value: []byte(":"),       Line: 2, Column: 21},
				{Type: query.Dollar,     Value: []byte("$"),       Line: 2, Column: 23},
				{Type: query.Name,       Value: []byte("type"),    Line: 2, Column: 24},
				{Type: query.ParenClose, Value: []byte(")"),       Line: 2, Column: 28},
				{Type: query.CurlyOpen,  Value: []byte("{"),       Line: 2, Column: 30},
	
				{Type: query.Name,       Value: []byte("name"),    Line: 3, Column: 4},
				{Type: query.CurlyClose, Value: []byte("}"),       Line: 4, Column: 3},
				{Type: query.CurlyClose, Value: []byte("}"),       Line: 5, Column: 2},
				{Type: query.EOF,        Value: nil,               Line: 5, Column: 3},
			},
		},
		{
			name: "Lex inline fragment",
			input: []byte(`query MixedTypes {
		user {
			... on Admin {
				adminField
			}
			... on Member {
				memberField
			}
		}
	}`),
			expected: query.Tokens{
				{Type: query.Query,       Value: []byte("query"),   Line: 1, Column: 1},
				{Type: query.Name,        Value: []byte("MixedTypes"), Line: 1, Column: 7},
				{Type: query.CurlyOpen,   Value: []byte("{"),       Line: 1, Column: 18},
	
				{Type: query.Name,        Value: []byte("user"),    Line: 2, Column: 3},
				{Type: query.CurlyOpen,   Value: []byte("{"),       Line: 2, Column: 8},
	
				// ... on Admin
				{Type: query.Spread,      Value: []byte("..."),     Line: 3, Column: 4},
				{Type: query.On,          Value: []byte("on"),      Line: 3, Column: 8},
				{Type: query.Name,        Value: []byte("Admin"),   Line: 3, Column: 11},
				{Type: query.CurlyOpen,   Value: []byte("{"),       Line: 3, Column: 17},
				{Type: query.Name,        Value: []byte("adminField"), Line: 4, Column: 5},
				{Type: query.CurlyClose,  Value: []byte("}"),       Line: 5, Column: 4},
	
				// ... on Member
				{Type: query.Spread,      Value: []byte("..."),     Line: 6, Column: 4},
				{Type: query.On,          Value: []byte("on"),      Line: 6, Column: 8},
				{Type: query.Name,        Value: []byte("Member"),  Line: 6, Column: 11},
				{Type: query.CurlyOpen,   Value: []byte("{"),       Line: 6, Column: 18},
				{Type: query.Name,        Value: []byte("memberField"), Line: 7, Column: 5},
				{Type: query.CurlyClose,  Value: []byte("}"),       Line: 8, Column: 4},
	
				{Type: query.CurlyClose,  Value: []byte("}"),       Line: 9, Column: 3},
				{Type: query.CurlyClose,  Value: []byte("}"),       Line: 10, Column: 2},
				{Type: query.EOF,         Value: nil,               Line: 10, Column: 3},
			},
		},
		{
			name: "Lex query with alias",
			input: []byte(`query AliasQuery {
				user1: user(id: 1) {
					name
				}
				user2: user(id: 2) {
					email
				}
			}`),
			expected: query.Tokens{
				{Type: query.Query, Value: []byte("query"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("AliasQuery"), Line: 1, Column: 7},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 18},
				{Type: query.Name, Value: []byte("user1"), Line: 2, Column: 5},
				{Type: query.Colon, Value: []byte(":"), Line: 2, Column: 10},
				{Type: query.Name, Value: []byte("user"), Line: 2, Column: 12},
				{Type: query.ParenOpen, Value: []byte("("), Line: 2, Column: 16},
				{Type: query.Name, Value: []byte("id"), Line: 2, Column: 17},
				{Type: query.Colon, Value: []byte(":"), Line: 2, Column: 19},
				{Type: query.Name, Value: []byte("1"), Line: 2, Column: 21},
				{Type: query.ParenClose, Value: []byte(")"), Line: 2, Column: 22},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 2, Column: 24},
				{Type: query.Name, Value: []byte("name"), Line: 3, Column: 6},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 4, Column: 5},
				{Type: query.Name, Value: []byte("user2"), Line: 5, Column: 5},
				{Type: query.Colon, Value: []byte(":"), Line: 5, Column: 10},
				{Type: query.Name, Value: []byte("user"), Line: 5, Column: 12},
				{Type: query.ParenOpen, Value: []byte("("), Line: 5, Column: 16},
				{Type: query.Name, Value: []byte("id"), Line: 5, Column: 17},
				{Type: query.Colon, Value: []byte(":"), Line: 5, Column: 19},
				{Type: query.Name, Value: []byte("2"), Line: 5, Column: 21},
				{Type: query.ParenClose, Value: []byte(")"), Line: 5, Column: 22},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 5, Column: 24},
				{Type: query.Name, Value: []byte("email"), Line: 6, Column: 6},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 7, Column: 5},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 8, Column: 4},
				{Type: query.EOF, Value: nil, Line: 8, Column: 5},
			},
		},
		{
			name: "Lex query with directives",
			input: []byte(`query QueryWithDirective($includeEmail: Boolean!) {
				user {
					email @include(if: $includeEmail)
					name
				}
			}`),
			expected: query.Tokens{
				{Type: query.Query, Value: []byte("query"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("QueryWithDirective"), Line: 1, Column: 7},
				{Type: query.ParenOpen, Value: []byte("("), Line: 1, Column: 25},
				{Type: query.Dollar, Value: []byte("$"), Line: 1, Column: 26},
				{Type: query.Name, Value: []byte("includeEmail"), Line: 1, Column: 27},
				{Type: query.Colon, Value: []byte(":"), Line: 1, Column: 39},
				{Type: query.Name, Value: []byte("Boolean"), Line: 1, Column: 41},
				{Type: query.Exclamation, Value: []byte("!"), Line: 1, Column: 48},
				{Type: query.ParenClose, Value: []byte(")"), Line: 1, Column: 49},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 51},
				{Type: query.Name, Value: []byte("user"), Line: 2, Column: 5},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 2, Column: 10},
				{Type: query.Name, Value: []byte("email"), Line: 3, Column: 6},
				{Type: query.At, Value: []byte("@"), Line: 3, Column: 12},
				{Type: query.Name, Value: []byte("include"), Line: 3, Column: 13},
				{Type: query.ParenOpen, Value: []byte("("), Line: 3, Column: 20},
				{Type: query.Name, Value: []byte("if"), Line: 3, Column: 21},
				{Type: query.Colon, Value: []byte(":"), Line: 3, Column: 23},
				{Type: query.Dollar, Value: []byte("$"), Line: 3, Column: 25},
				{Type: query.Name, Value: []byte("includeEmail"), Line: 3, Column: 26},
				{Type: query.ParenClose, Value: []byte(")"), Line: 3, Column: 38},
				{Type: query.Name, Value: []byte("name"), Line: 4, Column: 6},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 5, Column: 5},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 6, Column: 4},
				{Type: query.EOF, Value: nil, Line: 6, Column: 5},
			},
		},{
			name: "Mutation simple example",
			input: []byte(`mutation CreateUser {
				createUser(name: "Alice", email: "alice@example.com") {
					id
					name
				}
			}`),
			expected: query.Tokens{
				{Type: query.Mutation, Value: []byte("mutation"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("CreateUser"), Line: 1, Column: 10},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 21},
				{Type: query.Name, Value: []byte("createUser"), Line: 2, Column: 5},
				{Type: query.ParenOpen, Value: []byte("("), Line: 2, Column: 15},
				{Type: query.Name, Value: []byte("name"), Line: 2, Column: 16},
				{Type: query.Colon, Value: []byte(":"), Line: 2, Column: 20},
				{Type: query.Value, Value: []byte("\"Alice\""), Line: 2, Column: 22},
				{Type: query.Comma, Value: []byte(","), Line: 2, Column: 31},
				{Type: query.Name, Value: []byte("email"), Line: 2, Column: 33},
				{Type: query.Colon, Value: []byte(":"), Line: 2, Column: 38},
				{Type: query.Value, Value: []byte("\"alice@example.com\""), Line: 2, Column: 40},
				{Type: query.ParenClose, Value: []byte(")"), Line: 2, Column: 61},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 2, Column: 63},
				{Type: query.Name, Value: []byte("id"), Line: 3, Column: 6},
				{Type: query.Name, Value: []byte("name"), Line: 4, Column: 6},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 5, Column: 5},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 6, Column: 4},
				{Type: query.EOF, Value: nil, Line: 6, Column: 5},
			},
		},
		{
			name: "Subscription example",
			input: []byte(`subscription OnUserAdded {
				userAdded {
					id
					name
				}
			}`),
			expected: query.Tokens{
				{Type: query.Subscription, Value: []byte("subscription"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("OnUserAdded"), Line: 1, Column: 14},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 26},
				{Type: query.Name, Value: []byte("userAdded"), Line: 2, Column: 5},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 2, Column: 15},
				{Type: query.Name, Value: []byte("id"), Line: 3, Column: 6},
				{Type: query.Name, Value: []byte("name"), Line: 4, Column: 6},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 5, Column: 5},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 6, Column: 4},
				{Type: query.EOF, Value: nil, Line: 6, Column: 5},
			},
		},
		{
			name: "Mutation with variables",
			input: []byte(`mutation CreateUser($name: String!, $email: String!) {
				createUser(name: $name, email: $email) {
					id
					name
				}
			}`),
			expected: []*query.Token{
				{Type: query.Mutation, Value: []byte("mutation"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("CreateUser"), Line: 1, Column: 10},
				{Type: query.ParenOpen, Value: []byte("("), Line: 1, Column: 20},
				{Type: query.Dollar, Value: []byte("$"), Line: 1, Column: 21},
				{Type: query.Name, Value: []byte("name"), Line: 1, Column: 22},
				{Type: query.Colon, Value: []byte(":"), Line: 1, Column: 26},
				{Type: query.Name, Value: []byte("String"), Line: 1, Column: 28},
				{Type: query.Exclamation, Value: []byte("!"), Line: 1, Column: 34},
				{Type: query.Comma, Value: []byte(","), Line: 1, Column: 35},
				{Type: query.Dollar, Value: []byte("$"), Line: 1, Column: 37},
				{Type: query.Name, Value: []byte("email"), Line: 1, Column: 38},
				{Type: query.Colon, Value: []byte(":"), Line: 1, Column: 43},
				{Type: query.Name, Value: []byte("String"), Line: 1, Column: 45},
				{Type: query.Exclamation, Value: []byte("!"), Line: 1, Column: 51},
				{Type: query.ParenClose, Value: []byte(")"), Line: 1, Column: 52},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 54},
				{Type: query.Name, Value: []byte("createUser"), Line: 2, Column: 3},
				{Type: query.ParenOpen, Value: []byte("("), Line: 2, Column: 13},
				{Type: query.Name, Value: []byte("name"), Line: 2, Column: 14},
				{Type: query.Colon, Value: []byte(":"), Line: 2, Column: 18},
				{Type: query.Dollar, Value: []byte("$"), Line: 2, Column: 20},
				{Type: query.Name, Value: []byte("name"), Line: 2, Column: 21},
				{Type: query.Comma, Value: []byte(","), Line: 2, Column: 25},
				{Type: query.Name, Value: []byte("email"), Line: 2, Column: 27},
				{Type: query.Colon, Value: []byte(":"), Line: 2, Column: 32},
				{Type: query.Dollar, Value: []byte("$"), Line: 2, Column: 34},
				{Type: query.Name, Value: []byte("email"), Line: 2, Column: 35},
				{Type: query.ParenClose, Value: []byte(")"), Line: 2, Column: 41},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 2, Column: 43},
				{Type: query.Name, Value: []byte("id"), Line: 3, Column: 4},
				{Type: query.Name, Value: []byte("name"), Line: 4, Column: 4},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 5, Column: 3},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 6, Column: 2},
				{Type: query.EOF, Value: nil, Line: 6, Column: 3},
			},
		},
		{
			name: "Fragment usage",
			input: []byte(`query GetUser {
				user {
					...UserFields
				}
			}
			
			fragment UserFields on User {
				id
				name
			}`),
			expected: []*query.Token{
				{Type: query.Query, Value: []byte("query"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("GetUser"), Line: 1, Column: 7},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 15},
				{Type: query.Name, Value: []byte("user"), Line: 2, Column: 3},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 2, Column: 8},
				{Type: query.Spread, Value: []byte("..."), Line: 3, Column: 4},
				{Type: query.Name, Value: []byte("UserFields"), Line: 3, Column: 7},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 4, Column: 3},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 5, Column: 2},
				{Type: query.Fragment, Value: []byte("fragment"), Line: 7, Column: 1},
				{Type: query.Name, Value: []byte("UserFields"), Line: 7, Column: 10},
				{Type: query.On, Value: []byte("on"), Line: 7, Column: 21},
				{Type: query.Name, Value: []byte("User"), Line: 7, Column: 24},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 7, Column: 29},
				{Type: query.Name, Value: []byte("id"), Line: 8, Column: 3},
				{Type: query.Name, Value: []byte("name"), Line: 9, Column: 3},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 10, Column: 2},
				{Type: query.EOF, Value: nil, Line: 10, Column: 3},
			},
		},
		{
			name: "Deeply nested query",
			input: []byte(`query DeepQuery {
				user {
					posts {
						comments {
							author {
								name
							}
						}
					}
				}
			}`),
			expected: []*query.Token{
				{Type: query.Query, Value: []byte("query"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("DeepQuery"), Line: 1, Column: 7},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 16},
				{Type: query.Name, Value: []byte("user"), Line: 2, Column: 3},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 2, Column: 8},
				{Type: query.Name, Value: []byte("posts"), Line: 3, Column: 4},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 3, Column: 10},
				{Type: query.Name, Value: []byte("comments"), Line: 4, Column: 5},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 4, Column: 14},
				{Type: query.Name, Value: []byte("author"), Line: 5, Column: 6},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 5, Column: 13},
				{Type: query.Name, Value: []byte("name"), Line: 6, Column: 7},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 7, Column: 6},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 8, Column: 5},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 9, Column: 4},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 10, Column: 3},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 11, Column: 2},
				{Type: query.EOF, Value: nil, Line: 11, Column: 3},
			},
		},{
			name: "Mutation with default variable values",
			input: []byte(`mutation UpdateSettings($settings: SettingsInput = { theme: "dark", notifications: true }) {
				updateSettings(settings: $settings) {
					success
				}
			}`),
			expected: []*query.Token{
				{Type: query.Mutation, Value: []byte("mutation"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("UpdateSettings"), Line: 1, Column: 10},
				{Type: query.ParenOpen, Value: []byte("("), Line: 1, Column: 24},
				{Type: query.Dollar, Value: []byte("$"), Line: 1, Column: 25},
				{Type: query.Name, Value: []byte("settings"), Line: 1, Column: 26},
				{Type: query.Colon, Value: []byte(":"), Line: 1, Column: 34},
				{Type: query.Name, Value: []byte("SettingsInput"), Line: 1, Column: 36},
				{Type: query.Equal, Value: []byte("="), Line: 1, Column: 49},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 51},
				{Type: query.Name, Value: []byte("theme"), Line: 1, Column: 53},
				{Type: query.Colon, Value: []byte(":"), Line: 1, Column: 58},
				{Type: query.Value, Value: []byte("\"dark\""), Line: 1, Column: 60},
				{Type: query.Comma, Value: []byte(","), Line: 1, Column: 67},
				{Type: query.Name, Value: []byte("notifications"), Line: 1, Column: 69},
				{Type: query.Colon, Value: []byte(":"), Line: 1, Column: 82},
				{Type: query.Name, Value: []byte("true"), Line: 1, Column: 84},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 1, Column: 88},
				{Type: query.ParenClose, Value: []byte(")"), Line: 1, Column: 89},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 91},
				{Type: query.Name, Value: []byte("updateSettings"), Line: 2, Column: 4},
				{Type: query.ParenOpen, Value: []byte("("), Line: 2, Column: 18},
				{Type: query.Name, Value: []byte("settings"), Line: 2, Column: 19},
				{Type: query.Colon, Value: []byte(":"), Line: 2, Column: 27},
				{Type: query.Dollar, Value: []byte("$"), Line: 2, Column: 29},
				{Type: query.Name, Value: []byte("settings"), Line: 2, Column: 30},
				{Type: query.ParenClose, Value: []byte(")"), Line: 2, Column: 38},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 2, Column: 40},
				{Type: query.Name, Value: []byte("success"), Line: 3, Column: 5},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 4, Column: 4},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 5, Column: 3},
				{Type: query.EOF, Value: nil, Line: 5, Column: 4},
			},
		},
		{
			name: "Query with list arguments",
			input: []byte(`query UsersByIds($ids: [ID!]!) {
				users(ids: $ids) {
					id
					name
				}
			}`),
			expected: []*query.Token{
				{Type: query.Query, Value: []byte("query"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("UsersByIds"), Line: 1, Column: 7},
				{Type: query.ParenOpen, Value: []byte("("), Line: 1, Column: 18},
				{Type: query.Dollar, Value: []byte("$"), Line: 1, Column: 19},
				{Type: query.Name, Value: []byte("ids"), Line: 1, Column: 20},
				{Type: query.Colon, Value: []byte(":"), Line: 1, Column: 23},
				{Type: query.BracketOpen, Value: []byte("["), Line: 1, Column: 25},
				{Type: query.Name, Value: []byte("ID"), Line: 1, Column: 26},
				{Type: query.Exclamation, Value: []byte("!"), Line: 1, Column: 28},
				{Type: query.BracketClose, Value: []byte("]"), Line: 1, Column: 29},
				{Type: query.Exclamation, Value: []byte("!"), Line: 1, Column: 30},
				{Type: query.ParenClose, Value: []byte(")"), Line: 1, Column: 31},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 33},
				{Type: query.Name, Value: []byte("users"), Line: 2, Column: 4},
				{Type: query.ParenOpen, Value: []byte("("), Line: 2, Column: 9},
				{Type: query.Name, Value: []byte("ids"), Line: 2, Column: 10},
				{Type: query.Colon, Value: []byte(":"), Line: 2, Column: 13},
				{Type: query.Dollar, Value: []byte("$"), Line: 2, Column: 15},
				{Type: query.Name, Value: []byte("ids"), Line: 2, Column: 16},
				{Type: query.ParenClose, Value: []byte(")"), Line: 2, Column: 19},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 2, Column: 21},
				{Type: query.Name, Value: []byte("id"), Line: 3, Column: 5},
				{Type: query.Name, Value: []byte("name"), Line: 4, Column: 5},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 5, Column: 4},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 6, Column: 3},
				{Type: query.EOF, Value: nil, Line: 6, Column: 4},
			},
		},
		{
			name: "Query with escaped strings",
			input: []byte(`query EscapedStrings {
				user(name: "Alice\nBob\tCarol\"Quoted\"") {
					id
					name
				}
			}`),
			expected: query.Tokens{
				{Type: query.Query, Value: []byte("query"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("EscapedStrings"), Line: 1, Column: 7},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 22},
				{Type: query.Name, Value: []byte("user"), Line: 2, Column: 3},
				{Type: query.ParenOpen, Value: []byte("("), Line: 2, Column: 7},
				{Type: query.Name, Value: []byte("name"), Line: 2, Column: 8},
				{Type: query.Colon, Value: []byte(":"), Line: 2, Column: 12},
				{Type: query.Value, Value: []byte("\"Alice\\nBob\\tCarol\\\"Quoted\\\"\""), Line: 2, Column: 14},
				{Type: query.ParenClose, Value: []byte(")"), Line: 2, Column: 43},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 2, Column: 45},
				{Type: query.Name, Value: []byte("id"), Line: 3, Column: 4},
				{Type: query.Name, Value: []byte("name"), Line: 4, Column: 4},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 5, Column: 3},
				{Type: query.CurlyClose, Value: []byte("}"), Line: 6, Column: 2},
				{Type: query.EOF, Value: nil, Line: 6, Column: 3},
			},
		},
		
	}

	ignores := cmpopts.IgnoreFields(query.Token{}, "Column")
	
	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T)  {
			lexer := query.NewLexer()
			got, err := lexer.Lex(tt.input)
			if tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && err != nil {
				t.Errorf("Parse() error %v", err)
				return
			}

			if diff := cmp.Diff(got, tt.expected, ignores); diff != "" {
				t.Errorf("Parse() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}