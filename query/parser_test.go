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
			name: "Parse simple graphql query operation",
			input: []byte(`query MyQuery {
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name: "MyQuery",
					},
				},
			},
		}, {
			name: "Parse simple graphql query operation with variables",
			input: []byte(`query MyQuery($id: Int!) {
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name: "MyQuery",
						Variables: []*query.Variable{
							{
								Name: []byte("id"),
								Type: &query.FieldType{
									Name: []byte("Int"),
									Nullable: false,
									IsList: false,
								},
								DefaultValue: nil,
							},
						},
					},
				},
			},
		},
		{
			name: "Query with default complex object values",
			input: []byte(`query UpdateSettings($settings: SettingsInput = { theme: "dark", notifications: true }) {
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name: "UpdateSettings",
						Variables: []*query.Variable{
							{
								Name: []byte("settings"),
								Type: &query.FieldType{
									Name: []byte("SettingsInput"),
									Nullable: true,
									IsList: false,
								},
								DefaultValue: []byte(`{theme:"dark",notifications:true}`),
							},
						},
					},
				},
			},
		},
		{
			name: "Query with default complex object list values",
			input: []byte(`query UpdateSettings($settings: [SettingInput]! = [{ theme: "dark", notifications: true, options: { a: 1, b: 2 } }]) {
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name: "UpdateSettings",
						Variables: []*query.Variable{
							{
								Name: []byte("settings"),
								Type : &query.FieldType{
									Nullable: false,
									IsList: true,
									ListType: &query.FieldType{
										Name: []byte("SettingInput"),
										Nullable: true,
										IsList: false,
									},
								},
								DefaultValue: []byte(`[{theme:"dark",notifications:true,options:{a:1,b:2}}]`),
							},
						},
					},
				},
			},
		},
	}

	opts := cmp.FilterPath(func(p cmp.Path) bool {
		return p.Last().String() == ".tokens"
	}, cmp.Ignore())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := query.NewLexer()
			parser := query.NewParser(lexer)
			got, err := parser.Parse(tt.input)
			if tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && err != nil {
				t.Errorf("Parse() error %v", err)
				return
			}

			if diff := cmp.Diff(got, tt.expected, opts); diff != "" {
				t.Errorf("Parse() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}