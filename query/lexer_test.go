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
		expected []*query.Token
		wantErr error
	}{
		{
			name: "Lex simple graphql query",
			input: []byte(`query MyQuery {
				user(id: 123) {
					name
				}
			}`),
			expected: []*query.Token{
				{Type: query.Query, Value: []byte("query"), Line: 1, Column: 1},
				{Type: query.Name, Value: []byte("MyQuery"), Line: 1, Column: 7},
				{Type: query.CurlyOpen, Value: []byte("{"), Line: 1, Column: 15},
				{Type: query.Name, Value: []byte("user"), Line: 2, Column: 5},
				{Type: query.ParenOpen, Value: []byte("("), Line: 2, Column: 9},
				{Type: query.Name, Value: []byte("id"), Line: 2, Column: 10},
				{Type: query.Colon, Value: []byte(":"), Line: 2, Column: 12},
				{Type: query.Int, Value: []byte("123"), Line: 2, Column: 14},
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
			expected: []*query.Token{
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
	// 	{
	// 		name: "Lex inline fragment",
	// 		input: []byte(`query MixedTypes {
	// 	user {
	// 		... on Admin {
	// 			adminField
	// 		}
	// 		... on Member {
	// 			memberField
	// 		}
	// 	}
	// }`),
	// 		expected: []*query.Token{
	// 			{Type: query.Query,       Value: []byte("query"),   Line: 1, Column: 1},
	// 			{Type: query.Name,        Value: []byte("MixedTypes"), Line: 1, Column: 7},
	// 			{Type: query.CurlyOpen,   Value: []byte("{"),       Line: 1, Column: 18},
	
	// 			{Type: query.Name,        Value: []byte("user"),    Line: 2, Column: 3},
	// 			{Type: query.CurlyOpen,   Value: []byte("{"),       Line: 2, Column: 8},
	
	// 			// ... on Admin
	// 			{Type: query.Spread,      Value: []byte("..."),     Line: 3, Column: 5},
	// 			{Type: query.On,          Value: []byte("on"),      Line: 3, Column: 9},
	// 			{Type: query.Name,        Value: []byte("Admin"),   Line: 3, Column: 12},
	// 			{Type: query.CurlyOpen,   Value: []byte("{"),       Line: 3, Column: 18},
	// 			{Type: query.Name,        Value: []byte("adminField"), Line: 4, Column: 7},
	// 			{Type: query.CurlyClose,  Value: []byte("}"),       Line: 5, Column: 5},
	
	// 			// ... on Member
	// 			{Type: query.Spread,      Value: []byte("..."),     Line: 6, Column: 5},
	// 			{Type: query.On,          Value: []byte("on"),      Line: 6, Column: 9},
	// 			{Type: query.Name,        Value: []byte("Member"),  Line: 6, Column: 12},
	// 			{Type: query.CurlyOpen,   Value: []byte("{"),       Line: 6, Column: 19},
	// 			{Type: query.Name,        Value: []byte("memberField"), Line: 7, Column: 7},
	// 			{Type: query.CurlyClose,  Value: []byte("}"),       Line: 8, Column: 5},
	
	// 			{Type: query.CurlyClose,  Value: []byte("}"),       Line: 9, Column: 3},
	// 			{Type: query.CurlyClose,  Value: []byte("}"),       Line: 10, Column: 1},
	// 			{Type: query.EOF,         Value: nil,               Line: 10, Column: 2},
	// 		},
	// 	},
	// 	{
	// 		name: "Lex mutation",
	// 		input: []byte(`mutation UpdateUser($id: ID!, $name: String) {
	// 	updateUser(id: $id, name: $name) {
	// 		success
	// 		user {
	// 			id
	// 			name
	// 		}
	// 	}
	// }`),
	// 		expected: []*query.Token{
	// 			{Type: query.Mutation,    Value: []byte("mutation"),Line: 1, Column: 1},
	// 			{Type: query.Name,        Value: []byte("UpdateUser"), Line: 1, Column: 10},
	
	// 			{Type: query.ParenOpen,   Value: []byte("("),       Line: 1, Column: 20},
	// 			// $id: ID!
	// 			{Type: query.Variable,    Value: []byte("id"),      Line: 1, Column: 21},
	// 			{Type: query.Colon,       Value: []byte(":"),       Line: 1, Column: 24},
	// 			{Type: query.Name,        Value: []byte("ID"),      Line: 1, Column: 26},
	// 			{Type: query.Bang,        Value: []byte("!"),       Line: 1, Column: 28},
	// 			{Type: query.Comma,       Value: []byte(","),       Line: 1, Column: 29},
	// 			// $name: String
	// 			{Type: query.Variable,    Value: []byte("name"),    Line: 1, Column: 31},
	// 			{Type: query.Colon,       Value: []byte(":"),       Line: 1, Column: 36},
	// 			{Type: query.Name,        Value: []byte("String"),  Line: 1, Column: 38},
	// 			{Type: query.ParenClose,  Value: []byte(")"),       Line: 1, Column: 44},
	
	// 			{Type: query.CurlyOpen,   Value: []byte("{"),       Line: 1, Column: 46},
	// 			// updateUser(id: $id, name: $name) {
	// 			{Type: query.Name,        Value: []byte("updateUser"), Line: 2, Column: 3},
	// 			{Type: query.ParenOpen,   Value: []byte("("),       Line: 2, Column: 13},
	// 			{Type: query.Name,        Value: []byte("id"),      Line: 2, Column: 14},
	// 			{Type: query.Colon,       Value: []byte(":"),       Line: 2, Column: 16},
	// 			{Type: query.Variable,    Value: []byte("id"),      Line: 2, Column: 18},
	// 			{Type: query.Comma,       Value: []byte(","),       Line: 2, Column: 20},
	// 			{Type: query.Name,        Value: []byte("name"),    Line: 2, Column: 22},
	// 			{Type: query.Colon,       Value: []byte(":"),       Line: 2, Column: 26},
	// 			{Type: query.Variable,    Value: []byte("name"),    Line: 2, Column: 28},
	// 			{Type: query.ParenClose,  Value: []byte(")"),       Line: 2, Column: 32},
	// 			{Type: query.CurlyOpen,   Value: []byte("{"),       Line: 2, Column: 34},
	
	// 			{Type: query.Name,        Value: []byte("success"), Line: 3, Column: 5},
	// 			{Type: query.Name,        Value: []byte("user"),    Line: 4, Column: 5},
	// 			{Type: query.CurlyOpen,   Value: []byte("{"),       Line: 4, Column: 10},
	// 			{Type: query.Name,        Value: []byte("id"),      Line: 5, Column: 7},
	// 			{Type: query.Name,        Value: []byte("name"),    Line: 6, Column: 7},
	// 			{Type: query.CurlyClose,  Value: []byte("}"),       Line: 7, Column: 5},
	
	// 			{Type: query.CurlyClose,  Value: []byte("}"),       Line: 8, Column: 3},
	// 			{Type: query.CurlyClose,  Value: []byte("}"),       Line: 9, Column: 1},
	// 			{Type: query.EOF,         Value: nil,               Line: 9, Column: 2},
	// 		},
	// 	},
	// 	{
	// 		name: "Lex query with directive",
	// 		input: []byte(`query UserList @deprecated(reason: "Use new endpoint") {
	// 	users {
	// 		id
	// 		name
	// 	}
	// }`),
	// 		expected: []*query.Token{
	// 			{Type: query.Query,        Value: []byte("query"),       Line: 1, Column: 1},
	// 			{Type: query.Name,         Value: []byte("UserList"),    Line: 1, Column: 7},
	
	// 			// ディレクティブ @deprecated(reason: "Use new endpoint")
	// 			{Type: query.At,           Value: []byte("@"),           Line: 1, Column: 16},
	// 			{Type: query.Name,         Value: []byte("deprecated"),  Line: 1, Column: 17},
	// 			{Type: query.ParenOpen,    Value: []byte("("),           Line: 1, Column: 27},
	// 			{Type: query.Name,         Value: []byte("reason"),      Line: 1, Column: 28},
	// 			{Type: query.Colon,        Value: []byte(":"),           Line: 1, Column: 34},
	// 			{Type: query.String,       Value: []byte(`"Use new endpoint"`), Line: 1, Column: 36},
	// 			{Type: query.ParenClose,   Value: []byte(")"),           Line: 1, Column: 54},
	
	// 			{Type: query.CurlyOpen,    Value: []byte("{"),           Line: 1, Column: 56},
	// 			{Type: query.Name,         Value: []byte("users"),       Line: 2, Column: 3},
	// 			{Type: query.CurlyOpen,    Value: []byte("{"),           Line: 2, Column: 9},
	
	// 			{Type: query.Name,         Value: []byte("id"),          Line: 3, Column: 5},
	// 			{Type: query.Name,         Value: []byte("name"),        Line: 4, Column: 5},
	
	// 			{Type: query.CurlyClose,   Value: []byte("}"),           Line: 5, Column: 3},
	// 			{Type: query.CurlyClose,   Value: []byte("}"),           Line: 6, Column: 1},
	// 			{Type: query.EOF,          Value: nil,                   Line: 6, Column: 2},
	// 		},
	// 	},
	// 	{
	// 		name: "Lex query with named fragment",
	// 		input: []byte(`query MyQuery {
	// 	user(id: 123) {
	// 		...UserFields
	// 	}
	// }
	
	// fragment UserFields on User {
	// 	id
	// 	name
	// }`),
	// 		expected: []*query.Token{
	// 			// query MyQuery { ... }
	// 			{Type: query.Query,       Value: []byte("query"),     Line: 1, Column: 1},
	// 			{Type: query.Name,        Value: []byte("MyQuery"),   Line: 1, Column: 7},
	// 			{Type: query.CurlyOpen,   Value: []byte("{"),         Line: 1, Column: 15},
	
	// 			{Type: query.Name,        Value: []byte("user"),      Line: 2, Column: 3},
	// 			{Type: query.ParenOpen,   Value: []byte("("),         Line: 2, Column: 7},
	// 			{Type: query.Name,        Value: []byte("id"),        Line: 2, Column: 8},
	// 			{Type: query.Colon,       Value: []byte(":"),         Line: 2, Column: 10},
	// 			{Type: query.Int,         Value: []byte("123"),       Line: 2, Column: 12},
	// 			{Type: query.ParenClose,  Value: []byte(")"),         Line: 2, Column: 15},
	// 			{Type: query.CurlyOpen,   Value: []byte("{"),         Line: 2, Column: 17},
	
	// 			// ...UserFields
	// 			{Type: query.Spread,      Value: []byte("..."),       Line: 3, Column: 5},
	// 			{Type: query.Name,        Value: []byte("UserFields"),Line: 3, Column: 8},
	
	// 			{Type: query.CurlyClose,  Value: []byte("}"),         Line: 4, Column: 3},
	// 			{Type: query.CurlyClose,  Value: []byte("}"),         Line: 5, Column: 1},
	
	// 			// fragment UserFields on User { ... }
	// 			{Type: query.Fragment,    Value: []byte("fragment"),  Line: 7, Column: 1},
	// 			{Type: query.Name,        Value: []byte("UserFields"),Line: 7, Column: 10},
	// 			{Type: query.On,          Value: []byte("on"),        Line: 7, Column: 20},
	// 			{Type: query.Name,        Value: []byte("User"),      Line: 7, Column: 23},
	// 			{Type: query.CurlyOpen,   Value: []byte("{"),         Line: 7, Column: 28},
	
	// 			{Type: query.Name,        Value: []byte("id"),        Line: 8, Column: 3},
	// 			{Type: query.Name,        Value: []byte("name"),      Line: 9, Column: 3},
	
	// 			{Type: query.CurlyClose,  Value: []byte("}"),         Line: 10, Column: 1},
	// 			{Type: query.EOF,         Value: nil,                 Line: 10, Column: 2},
	// 		},
	// 	},
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