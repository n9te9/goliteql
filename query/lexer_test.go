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