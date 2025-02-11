package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/lkeix/gg-parser/schema"
)

type Generator struct {
	Schema *schema.Schema
}

var gqlFilePattern = regexp.MustCompile(`.*[\.gql | \.graphql]`)

func NewGenerator(schemaDirectory, outputDirectory string) (*Generator, error) {
	gqlFilePaths := make([]string, 0)

	err := filepath.Walk(schemaDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && gqlFilePattern.MatchString(info.Name()) {
			gqlFilePaths = append(gqlFilePaths, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error get gql file path: %w", err)
	}

	fmt.Println(gqlFilePaths)

	fileContents := make([]byte, 0)

	for _, path := range gqlFilePaths {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
		}

		content, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		content = append(content, []byte("\n")...)

		fileContents = append(fileContents, content...)
	}

	lexer := schema.NewLexer()
	parser := schema.NewParser(lexer)
	s, err := parser.Parse(fileContents)
	if err != nil {
		return nil, fmt.Errorf("error parsing schema: %w", err)
	}

	s, err = s.Merge()
	if err != nil {
		return nil, fmt.Errorf("error merging schema: %w", err)
	}

	g := &Generator{
		Schema: s,
	}

	return g, nil
}

func (g *Generator) Generate() error {
	if err := g.generateModel(); err != nil {
		panic(err)
	}

	return nil
}

func (g *Generator) generateModel() error {
	for _, input := range g.Schema.Inputs {
		fmt.Println(string(input.Name))
	}

	for _, t := range g.Schema.Types {
		fmt.Println(string(t.Name))
	}

	return nil
}
