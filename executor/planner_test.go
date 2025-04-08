package executor_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lkeix/gg-executor/executor"
	"github.com/lkeix/gg-executor/query"
)

func TestPlanExecution(t *testing.T) {
	// Test cases for PlanExecution function
	tests := []struct {
		name         string
		input        []query.Selection
		resultTree   *executor.Node
		expectedName []byte
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
			expectedName: []byte("field1"),
			resultTree: &executor.Node{
				Name: []byte("field1"),
				SelectSets: []query.Selection{
					&query.Field{Name: []byte("childField1")},
				},
				Children: []*executor.Node{
					{
						Name:       []byte("childField1"),
						SelectSets: nil,
						Children:   nil,
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
			expectedName: []byte("level1"),
			resultTree: &executor.Node{
				Name: []byte("level1"),
				SelectSets: []query.Selection{
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
				Children: []*executor.Node{
					{
						Name: []byte("level2"),
						SelectSets: []query.Selection{
							&query.Field{
								Name: []byte("level3"),
								Selections: []query.Selection{
									&query.Field{
										Name: []byte("leafField"),
									},
								},
							},
						},
						Children: []*executor.Node{
							{
								Name: []byte("level3"),
								SelectSets: []query.Selection{
									&query.Field{
										Name: []byte("leafField"),
									},
								},
								Children: []*executor.Node{
									{
										Name:       []byte("leafField"),
										SelectSets: nil,
										Children:   nil,
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
			result := executor.PlanExecution(tt.input)
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
