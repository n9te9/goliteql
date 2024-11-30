package ggparser_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	ggparser "github.com/lkeix/gg-parser"
)

func TestQueryParse(t *testing.T) {
	tests := []struct{
		name string
		input []byte
		expected *ggparser.Document
		wantErr error
	}{
		{
			name: "Parse simple graphql query",
			input: []byte(`query MyQuery {
				user(id: 123) {
					name
				}
			}`),
			expected: &ggparser.Document{
				Tokens: []*ggparser.QueryToken{
					{QueryType: ggparser.QueryTypeQuery, Value: []byte("query"), Column: 1, Line: 1},
					{QueryType: ggparser.QueryTypeName, Value: []byte("MyQuery"), Column: 7, Line: 1},
					{QueryType: ggparser.QueryTypeCurlyOpen, Value: []byte("{"), Column: 15, Line: 1},
					{QueryType: ggparser.QueryTypeName, Value: []byte("user"), Column: 5, Line: 2},
					{QueryType: ggparser.QueryTypeParenOpen, Value: []byte("("), Column: 9, Line: 2},
					{QueryType: ggparser.QueryTypeName, Value: []byte("id"), Column: 10, Line: 2},
					{QueryType: ggparser.QueryTypeColon, Value: []byte(":"), Column: 12, Line: 2},
					{QueryType: ggparser.QueryTypeInt, Value: []byte("123"), Column: 14, Line: 2},
					{QueryType: ggparser.QueryTypeParenClose, Value: []byte(")"), Column: 17, Line: 2},
					{QueryType: ggparser.QueryTypeCurlyOpen, Value: []byte("{"), Column: 19, Line: 2},
					{QueryType: ggparser.QueryTypeName, Value: []byte("name"), Column: 6, Line: 3},
					{QueryType: ggparser.QueryTypeCurlyClose, Value: []byte("}"), Column: 5, Line: 4},
					{QueryType: ggparser.QueryTypeCurlyClose, Value:[]byte("}"), Column: 4, Line: 5},
					{QueryType: ggparser.QueryTypeEOF, Value: nil, Column: 5, Line: 5},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ggparser.ParseQuery(tt.input)
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