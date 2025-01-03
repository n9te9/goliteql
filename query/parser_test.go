package query_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lkeix/gg-parser/query"
)

func TestQueryParse(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *query.Document
		wantErr  error
	}{
		{
			name: "Parse simple graphql query operation",
			input: []byte(`query MyQuery {
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
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
						Name:          "MyQuery",
						Variables: []*query.Variable{
							{
								Name: []byte("id"),
								Type: &query.FieldType{
									Name:     []byte("Int"),
									Nullable: false,
									IsList:   false,
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
						Name:          "UpdateSettings",
						Variables: []*query.Variable{
							{
								Name: []byte("settings"),
								Type: &query.FieldType{
									Name:     []byte("SettingsInput"),
									Nullable: true,
									IsList:   false,
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
						Name:          "UpdateSettings",
						Variables: []*query.Variable{
							{
								Name: []byte("settings"),
								Type: &query.FieldType{
									Nullable: false,
									IsList:   true,
									ListType: &query.FieldType{
										Name:     []byte("SettingInput"),
										Nullable: true,
										IsList:   false,
									},
								},
								DefaultValue: []byte(`[{theme:"dark",notifications:true,options:{a:1,b:2}}]`),
							},
						},
					},
				},
			},
		},
		{
			name: "Parse query with nullable variables and default value",
			input: []byte(`query MyQuery($id: Int = 42, $name: String = "default") {
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
						Variables: []*query.Variable{
							{
								Name: []byte("id"),
								Type: &query.FieldType{
									Name:     []byte("Int"),
									Nullable: true,
									IsList:   false,
								},
								DefaultValue: []byte("42"),
							},
							{
								Name: []byte("name"),
								Type: &query.FieldType{
									Name:     []byte("String"),
									Nullable: true,
									IsList:   false,
								},
								DefaultValue: []byte(`"default"`),
							},
						},
					},
				},
			},
		}, {
			name: "Parse query with nested list variables",
			input: []byte(`query NestedList($matrix: [[Int!]!]!) {
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "NestedList",
						Variables: []*query.Variable{
							{
								Name: []byte("matrix"),
								Type: &query.FieldType{
									Nullable: false,
									IsList:   true,
									ListType: &query.FieldType{
										Nullable: false,
										IsList:   true,
										ListType: &query.FieldType{
											Name:     []byte("Int"),
											Nullable: false,
											IsList:   false,
										},
									},
								},
								DefaultValue: nil,
							},
						},
					},
				},
			},
		}, {
			name: "Parse query with nullable complex default object values",
			input: []byte(`query DefaultComplex($config: ConfigInput = { retries: 3, timeout: 30.5 }) {
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "DefaultComplex",
						Variables: []*query.Variable{
							{
								Name: []byte("config"),
								Type: &query.FieldType{
									Name:     []byte("ConfigInput"),
									Nullable: true,
									IsList:   false,
								},
								DefaultValue: []byte(`{retries:3,timeout:30.5}`),
							},
						},
					},
				},
			},
		}, {
			name: "Parse query with deeply nested default object values",
			input: []byte(`query DeepDefaults($data: DataInput = { user: { id: 1, settings: { theme: "light" } } }) {
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "DeepDefaults",
						Variables: []*query.Variable{
							{
								Name: []byte("data"),
								Type: &query.FieldType{
									Name:     []byte("DataInput"),
									Nullable: true,
									IsList:   false,
								},
								DefaultValue: []byte(`{user:{id:1,settings:{theme:"light"}}}`),
							},
						},
					},
				},
			},
		}, {
			name: "Parse query with enum default value",
			input: []byte(`query WithEnum($status: Status = ACTIVE) {
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "WithEnum",
						Variables: []*query.Variable{
							{
								Name: []byte("status"),
								Type: &query.FieldType{
									Name:     []byte("Status"),
									Nullable: true,
									IsList:   false,
								},
								DefaultValue: []byte("ACTIVE"),
							},
						},
					},
				},
			},
		}, {
			name: "Parse query with nullable list variables",
			input: []byte(`query NullableList($ids: [ID]) {
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "NullableList",
						Variables: []*query.Variable{
							{
								Name: []byte("ids"),
								Type: &query.FieldType{
									Nullable: true,
									IsList:   true,
									ListType: &query.FieldType{
										Name:     []byte("ID"),
										Nullable: true,
										IsList:   false,
									},
								},
								DefaultValue: nil,
							},
						},
					},
				},
			},
		}, {
			name: "Parse simple query selection",
			input: []byte(`query MyQuery {
				field
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
						Selections: []query.Selection{
							&query.Field{
								Name: []byte("field"),
							},
						},
					},
				},
			},
		}, 
		{
			name: "Parse nested query selection",
			input: []byte(`query MyQuery {
				field {
					subfield
				}
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
						Selections: []query.Selection{
							&query.Field{
								Name: []byte("field"),
								Selections: []query.Selection{
									&query.Field{
										Name: []byte("subfield"),
									},
								},
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
