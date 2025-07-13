package executor_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/n9te9/goliteql/executor"
	"github.com/n9te9/goliteql/query"
)

func TestPlanExecution(t *testing.T) {
	// Test cases for PlanExecution function
	tests := []struct {
		name                string
		input               []query.Selection
		resultTree          []*executor.Node
		fragmentDefinitions query.FragmentDefinitions
	}{
		{
			name: "Test case 1",
			input: []query.Selection{
				&query.Field{
					Name: []byte("field1"),
					Selections: []query.Selection{
						&query.Field{Name: []byte("childField1")},
					},
				},
			},
			resultTree: []*executor.Node{
				{
					Name: []byte("field1"),
					Children: []*executor.Node{
						{
							Name:     []byte("childField1"),
							Children: []*executor.Node{},
						},
					},
				},
			},
		},
		{
			name: "Nested query depth 3",
			input: []query.Selection{
				&query.Field{
					Name: []byte("level1"),
					Selections: []query.Selection{
						&query.Field{
							Name: []byte("level2"),
							Selections: []query.Selection{
								&query.Field{
									Name: []byte("level3"),
									Selections: []query.Selection{
										&query.Field{
											Name: []byte("leafField"),
										},
									},
								},
							},
						},
					},
				},
			},
			resultTree: []*executor.Node{
				{
					Name: []byte("level1"),
					Children: []*executor.Node{
						{
							Name: []byte("level2"),
							Children: []*executor.Node{
								{
									Name: []byte("level3"),
									Children: []*executor.Node{
										{
											Name:     []byte("leafField"),
											Children: []*executor.Node{},
										},
									},
								},
							},
						},
					},
				},
			},
		}, {
			name: "Plan inline fragment",
			input: []query.Selection{
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
			resultTree: []*executor.Node{
				{
					Name: []byte("field"),
					Children: []*executor.Node{
						{
							Type: "TypeName",
							Children: []*executor.Node{
								{
									Name:     []byte("subfield"),
									Children: []*executor.Node{},
								},
							},
						},
					},
				},
			},
		}, {
			name: "Plan depth 3 inline fragment",
			input: []query.Selection{
				&query.Field{
					Name: []byte("field"),
					Selections: []query.Selection{
						&query.InlineFragment{
							TypeCondition: []byte("TypeName1"),
							Selections: []query.Selection{
								&query.InlineFragment{
									TypeCondition: []byte("TypeName2"),
									Selections: []query.Selection{
										&query.InlineFragment{
											TypeCondition: []byte("TypeName3"),
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
			resultTree: []*executor.Node{
				{
					Name: []byte("field"),
					Children: []*executor.Node{
						{
							Type: "TypeName1",
							Children: []*executor.Node{
								{
									Type: "TypeName2",
									Children: []*executor.Node{
										{
											Type: "TypeName3",
											Children: []*executor.Node{
												{
													Name:     []byte("subfield"),
													Children: []*executor.Node{},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.PlanExecution(tt.input, tt.fragmentDefinitions)
			if result == nil {
				t.Errorf("PlanExecution() returned nil")
				return
			}

			if diff := cmp.Diff(tt.resultTree, result); diff != "" {
				t.Errorf("PlanExecution() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
