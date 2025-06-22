package executor_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/n9te9/goliteql/executor"
)

func TestNullable_MarshalJSON(t *testing.T) {
	type TestStruct struct {
		Nullable executor.Nullable `json:"nullable,omitempty"`
	}

	tests := []struct {
		name       string
		testStruct TestStruct
		expected   string
		wantErr    error
	}{
		{
			name: "nullable is nil should return empty JSON",
			testStruct: TestStruct{
				Nullable: nil,
			},
			expected: "{}",
			wantErr:  nil,
		},
		{
			name: "nullable data is nil should return field is null JSON",
			testStruct: TestStruct{
				Nullable: executor.NewNullable(nil),
			},
			expected: `{"nullable":null}`,
			wantErr:  nil,
		},
		{
			name: "nullable data is not nil should return field with data JSON",
			testStruct: TestStruct{
				Nullable: executor.NewNullable(map[string]string{"key1": "value1"}),
			},
			expected: `{"nullable":{"key1":"value1"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.testStruct)
			if err != nil {
				if err != tt.wantErr {
					t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if cmp.Diff(string(b), tt.expected) != "" {
				t.Errorf("MarshalJSON() = %s, want %s", string(b), tt.expected)
			}
		})
	}
}
