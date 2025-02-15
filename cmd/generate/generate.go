package main

import (
	"log"
	"os"

	"github.com/lkeix/gg-parser/internal/generator"
)

func main() {
	schemaDirectory := "./internal/golden_files"

	g, err := generator.NewGenerator(schemaDirectory, os.Stdout, os.Stdout)
	if err != nil {
		log.Fatalf("error creating generator: %v", err)
	}

	if err := g.Generate(); err != nil {
		log.Fatalf("error generating code: %v", err)
	}
}