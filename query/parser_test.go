package query_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lkeix/gg-parser/query"
)

func TestQueryParse(t *testing.T) {
	tests := []struct{
		name string
		input []byte
		expected *query.Document
		wantErr error
	}{
		{
			name: "Parse simple graphql query",
			input: []byte(`query MyQuery {
				user(id: 123) {
					name
				}
			}`),
			expected: &query.Document{
				Tokens: []*query.Token{
					{Type: query.Query, Value: []byte("query"), Column: 1, Line: 1},
					{Type: query.Name, Value: []byte("MyQuery"), Column: 7, Line: 1},
					{Type: query.CurlyOpen, Value: []byte("{"), Column: 15, Line: 1},
					{Type: query.Name, Value: []byte("user"), Column: 5, Line: 2},
					{Type: query.ParenOpen, Value: []byte("("), Column: 9, Line: 2},
					{Type: query.Name, Value: []byte("id"), Column: 10, Line: 2},
					{Type: query.Colon, Value: []byte(":"), Column: 12, Line: 2},
					{Type: query.Name, Value: []byte("123"), Column: 14, Line: 2},
					{Type: query.ParenClose, Value: []byte(")"), Column: 17, Line: 2},
					{Type: query.CurlyOpen, Value: []byte("{"), Column: 19, Line: 2},
					{Type: query.Name, Value: []byte("name"), Column: 6, Line: 3},
					{Type: query.CurlyClose, Value: []byte("}"), Column: 5, Line: 4},
					{Type: query.CurlyClose, Value:[]byte("}"), Column: 4, Line: 5},
					{Type: query.EOF, Value: nil, Column: 5, Line: 5},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := query.NewLexer()
			parser := query.NewParser(lexer)
			got, err := parser.Parse(tt.input)
			if err != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Errorf("Parse() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}