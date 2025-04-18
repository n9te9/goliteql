package generator_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/n9te9/goliteql/internal/generator"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name                string
		schemaDirectory     string
		modelOutput         *bytes.Buffer
		modelPackagePath    string
		resolverOutput      *bytes.Buffer
		resolverPackagePath string
		expected            error
		expectGoFilePath    string
	}{
		{
			name:                "Generate Model code",
			schemaDirectory:     "../golden_files/model_test",
			modelOutput:         bytes.NewBuffer(nil),
			modelPackagePath:    "github.com/n9te9/goliteql/internal/generated/model",
			resolverOutput:      bytes.NewBuffer(nil),
			resolverPackagePath: "github.com/n9te9/goliteql/internal/generated/resolver",
			expected:            nil,
			expectGoFilePath:    "../golden_files/model_test/model.go",
		},
		{
			name:                "Generate Input code",
			schemaDirectory:     "../golden_files/operation_test",
			modelOutput:         bytes.NewBuffer(nil),
			modelPackagePath:    "github.com/n9te9/goliteql/internal/generated/model",
			resolverOutput:      bytes.NewBuffer(nil),
			resolverPackagePath: "github.com/n9te9/goliteql/internal/generated/resolver",
			expected:            nil,
			expectGoFilePath:    "../golden_files/operation_test/model.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(filepath.Abs(tt.schemaDirectory))
			generator, err := generator.NewGenerator(tt.schemaDirectory, tt.modelOutput, tt.resolverOutput, tt.resolverOutput, tt.resolverOutput, tt.modelPackagePath, tt.resolverPackagePath)
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

			if cmp.Diff(expectedContent, tt.modelOutput.Bytes()) != "" {
				t.Fatalf("expected \n%s, got \n%s", expectedContent, tt.modelOutput.Bytes())
			}
		})
	}
}
