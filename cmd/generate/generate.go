package main

import (
	"log"
	"os"

	"github.com/lkeix/gg-executor/internal/generator"
)

func main() {
	schemaDirectory := "./internal/golden_files/operation_test"
	modelOutputFile, err := os.Create("./internal/generated/operation_test/model/models.go")
	if err != nil {
		panic(err)
	}

	queryResolverOutputFile, err := os.Create("./internal/generated/operation_test/resolver/query.resolver.go")
	if err != nil {
		panic(err)
	}

	mutationResolverOutputFile, err := os.Create("./internal/generated/operation_test/resolver/mutation.resolver.go")
	if err != nil {
		panic(err)
	}

	rootResolverOutputFile, err := os.Create("./internal/generated/operation_test/resolver/resolver.go")
	if err != nil {
		panic(err)
	}

	g, err := generator.NewGenerator(schemaDirectory, modelOutputFile, queryResolverOutputFile, mutationResolverOutputFile, rootResolverOutputFile, "github.com/lkeix/gg-executor/internal/generated/operation_test/model", "github.com/lkeix/gg-executor/internal/generated/operation_test/resolver")
	if err != nil {
		log.Fatalf("error creating generator: %v", err)
	}

	if err := g.Generate(); err != nil {
		log.Fatalf("error generating code: %v", err)
	}
}
