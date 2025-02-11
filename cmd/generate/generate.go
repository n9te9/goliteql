package main

import (
	"log"

	"github.com/lkeix/gg-parser/internal/generator"
)

func main() {
	schemaDirectory := "./internal/golden_files"
	outputDirectory := "./outputs"

	g, err := generator.NewGenerator(schemaDirectory, outputDirectory)
	if err != nil {
		log.Fatalf("error creating generator: %v", err)
	}

	if err := g.Generate(); err != nil {
		log.Fatalf("error generating code: %v", err)
	}
}