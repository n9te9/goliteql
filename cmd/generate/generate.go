package main

import (
	"log"
	"os"

	"github.com/lkeix/gg-parser/internal/generator"
)

func main() {
	schemaDirectory := "./internal/golden_files/operation_test"
	modelOutputFile, _ := os.Create("./internal/model.go")
	resolverOutputFile, _ := os.Create("./internal/resolver.go")

	g, err := generator.NewGenerator(schemaDirectory, modelOutputFile, resolverOutputFile, "github.com/lkeix/gg-parser/internal/generated/model", "github.com/lkeix/gg-parser/internal/generated/resolver")
	if err != nil {
		log.Fatalf("error creating generator: %v", err)
	}

	if err := g.Generate(); err != nil {
		log.Fatalf("error generating code: %v", err)
	}
}
