package generator_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lkeix/gg-parser/internal/generator"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name string
		schemaDirectory string
		output *bytes.Buffer
		expected error
		expectGoFilePath string
	}{
		{
			name: "Generate code",
			schemaDirectory:  "../golden_files/model_test",
			output: bytes.NewBuffer(nil),
			expected: nil,
			expectGoFilePath: "../golden_files/model_test/model.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(filepath.Abs(tt.schemaDirectory))
			generator, err := generator.NewGenerator(tt.schemaDirectory, tt.output)
			if err != nil {
				t.Fatalf("error creating generator: %v", err)
			}

			err = generator.Generate()
			if err != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, err)
			}

			expectedContent, err := os.ReadFile(tt.expectGoFilePath)
			if err != nil {
				t.Fatalf("error reading file: %v", err)
			}

			if cmp.Diff(expectedContent, tt.output.Bytes()) != "" {
				t.Fatalf("expected \n%s, got \n%s", expectedContent, tt.output.Bytes())
			}
		})
	}
}