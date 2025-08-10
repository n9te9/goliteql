package query_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/n9te9/goliteql/query"
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
		{
			name: "Parse query with fragment spread",
			input: []byte(`query MyQuery {
				field {
					...FragmentName
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
									&query.FragmentSpread{
										Name: []byte("FragmentName"),
									},
								},
							},
						},
					},
				},
			},
		}, {
			name: "Parse query with inline fragment",
			input: []byte(`query MyQuery {
				field {
					... on TypeName {
						subfield
					}
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
									&query.InlineFragment{
										TypeCondition: []byte("TypeName"),
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
			},
		}, {
			name: "Parse query with fragment spread and inline fragment",
			input: []byte(`query MyQuery {
				field {
					...FragmentName
					... on TypeName {
						subfield
					}
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
									&query.FragmentSpread{
										Name: []byte("FragmentName"),
									},
									&query.InlineFragment{
										TypeCondition: []byte("TypeName"),
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
			},
		},
		{
			name: "Parse fragment spread with a single directive",
			input: []byte(`query MyQuery {
				field {
					...FragmentName @include(if: true)
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
									&query.FragmentSpread{
										Name: []byte("FragmentName"),
										Directives: []*query.Directive{
											{
												Name: []byte("include"),
												Arguments: []*query.DirectiveArgument{
													{
														Name:  []byte("if"),
														Value: []byte("true"),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Parse fragment spread with multiple directives",
			input: []byte(`query MyQuery {
				field {
					...FragmentName @include(if: true)
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
									&query.FragmentSpread{
										Name: []byte("FragmentName"),
										Directives: []*query.Directive{
											{
												Name: []byte("include"),
												Arguments: []*query.DirectiveArgument{
													{
														Name:  []byte("if"),
														Value: []byte("true"),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Parse field with directive",
			input: []byte(`query MyQuery {
				user @include(if: true)
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
						Selections: []query.Selection{
							&query.Field{
								Name: []byte("user"),
								Directives: []*query.Directive{
									{
										Name: []byte("include"),
										Arguments: []*query.DirectiveArgument{
											{
												Name:  []byte("if"),
												Value: []byte("true"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Parse field with directive",
			input: []byte(`query MyQuery {
				user @include(if: true)
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
						Selections: []query.Selection{
							&query.Field{
								Name: []byte("user"),
								Directives: []*query.Directive{
									{
										Name: []byte("include"),
										Arguments: []*query.DirectiveArgument{
											{
												Name:  []byte("if"),
												Value: []byte("true"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Parse field with skip directive",
			input: []byte(`query MyQuery {
				post {
					id
					title
					content
					author @skip(if: true) {
						id
					}
				}
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
						Selections: []query.Selection{
							&query.Field{
								Name: []byte("post"),
								Selections: []query.Selection{
									&query.Field{
										Name:       []byte("id"),
										Directives: nil,
									},
									&query.Field{
										Name:       []byte("title"),
										Directives: nil,
									},
									&query.Field{
										Name:       []byte("content"),
										Directives: nil,
									},
									&query.Field{
										Name: []byte("author"),
										Selections: []query.Selection{
											&query.Field{
												Name: []byte("id"),
											},
										},
										Directives: []*query.Directive{
											{
												Name: []byte("skip"),
												Arguments: []*query.DirectiveArgument{
													{
														Name:  []byte("if"),
														Value: []byte("true"),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Parse inline fragment with directive",
			input: []byte(`query MyQuery {
				user {
					... on User @deprecated(reason: "Use newField") {
						newField
					}
				}
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
						Selections: []query.Selection{
							&query.Field{
								Name: []byte("user"),
								Selections: []query.Selection{
									&query.InlineFragment{
										TypeCondition: []byte("User"),
										Directives: []*query.Directive{
											{
												Name: []byte("deprecated"),
												Arguments: []*query.DirectiveArgument{
													{
														Name:  []byte("reason"),
														Value: []byte(`"Use newField"`),
													},
												},
											},
										},
										Selections: []query.Selection{
											&query.Field{
												Name: []byte("newField"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Parse query operation with simple directive",
			input: []byte(`query MyQuery @deprecated {
				field
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
						Directives: []*query.Directive{
							{
								Name: []byte("deprecated"),
							},
						},
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
			name: "Parse query operation with directive having arguments",
			input: []byte(`query MyQuery @include(if: true) {
				field
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
						Directives: []*query.Directive{
							{
								Name: []byte("include"),
								Arguments: []*query.DirectiveArgument{
									{
										Name:  []byte("if"),
										Value: []byte("true"),
									},
								},
							},
						},
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
			name: "Parse query operation with multiple directives",
			input: []byte(`query MyQuery @include(if: true) @deprecated {
				field
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
						Directives: []*query.Directive{
							{
								Name: []byte("include"),
								Arguments: []*query.DirectiveArgument{
									{
										Name:  []byte("if"),
										Value: []byte("true"),
									},
								},
							},
							{
								Name: []byte("deprecated"),
							},
						},
						Selections: []query.Selection{
							&query.Field{
								Name: []byte("field"),
							},
						},
					},
				},
			},
		}, {
			name: "Parse query operation with directive having complex arguments",
			input: []byte(`query MyQuery @settings(config: { theme: "dark", features: ["a", "b"] }) {
				field
			}`),
			expected: &query.Document{
				Operations: []*query.Operation{
					{
						OperationType: query.QueryOperation,
						Name:          "MyQuery",
						Directives: []*query.Directive{
							{
								Name: []byte("settings"),
								Arguments: []*query.DirectiveArgument{
									{
										Name:  []byte("config"),
										Value: []byte(`{theme:"dark",features:["a","b"]}`),
									},
								},
							},
						},
						Selections: []query.Selection{
							&query.Field{
								Name: []byte("field"),
							},
						},
					},
				},
			},
		}, {
			name: "Parse fragment definition",
			input: []byte(`fragment MyFragment on User {
				id
				name
			}`),
			expected: &query.Document{
				Operations: query.Operations{},
				FragmentDefinitions: query.FragmentDefinitions{
					{
						Name:          []byte("MyFragment"),
						BasedTypeName: []byte("User"),
						Selections: []query.Selection{
							&query.Field{
								Name: []byte("id"),
							},
							&query.Field{
								Name: []byte("name"),
							},
						},
					},
				},
			},
		},
	}

	opts := cmp.FilterPath(func(p cmp.Path) bool {
		return p.Last().String() == ".tokens" || p.Last().String() == ".isVariable"
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
