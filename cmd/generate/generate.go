package main

import (
	"log"
	"os"

	"github.com/lkeix/gg-executor/internal/generator"
)

func main() {
	schemaDirectory := "./internal/golden_files/operation_test"
	modelOutputFile, err := os.Create("./internal/generated/model/model.go")
	if err != nil {
		panic(err)
	}

	resolverOutputFile, err := os.Create("./internal/generated/resolver/resolver.go")
	if err != nil {
		panic(err)
	}

	g, err := generator.NewGenerator(schemaDirectory, modelOutputFile, resolverOutputFile, "github.com/lkeix/gg-executor/internal/generated/model", "github.com/lkeix/gg-executor/internal/generated/resolver")
	if err != nil {
		log.Fatalf("error creating generator: %v", err)
	}

	if err := g.Generate(); err != nil {
		log.Fatalf("error generating code: %v", err)
	}
}
