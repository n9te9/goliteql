package query_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
				{Type: query.Daller, 	  Value: []byte("$"),       Line: 1, Column: 15},
				{Type: query.Name,   Value: []byte("id"),      Line: 1, Column: 16},
				{Type: query.Colon,      Value: []byte(":"),       Line: 1, Column: 18},
				{Type: query.Name,       Value: []byte("ID"),      Line: 1, Column: 20},
				{Type: query.Exclamation,       Value: []byte("!"),       Line: 1, Column: 22},
				{Type: query.Comma,      Value: []byte(","),       Line: 1, Column: 23},
				// $type: String
				{Type: query.Daller, 	  Value: []byte("$"),       Line: 1, Column: 25},
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
				{Type: query.Daller,     Value: []byte("$"),       Line: 2, Column: 12},
				{Type: query.Name,       Value: []byte("id"),      Line: 2, Column: 13},
				{Type: query.Comma,      Value: []byte(","),       Line: 2, Column: 15},
				{Type: query.Name,       Value: []byte("type"),    Line: 2, Column: 17},
				{Type: query.Colon,      Value: []byte(":"),       Line: 2, Column: 21},
				{Type: query.Daller,     Value: []byte("$"),       Line: 2, Column: 23},
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
				{Type: query.Daller, Value: []byte("$"), Line: 1, Column: 26},
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
				{Type: query.Daller, Value: []byte("$"), Line: 3, Column: 25},
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
	}

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

			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Errorf("Parse() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}