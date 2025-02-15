package generator

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/lkeix/gg-parser/schema"
)

type Generator struct {
	Schema *schema.Schema
	queryAST *ast.File
	mutationAST *ast.File
	subscriptionAST *ast.File
	modelAST *ast.File

	output io.Writer
}

var gqlFilePattern = regexp.MustCompile(`^.+\.gql$|^.+\.graphql$`)

func NewGenerator(schemaDirectory string, modelOutput io.Writer) (*Generator, error) {
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
		queryAST: &ast.File{},
		mutationAST: &ast.File{},
		subscriptionAST: &ast.File{},
		modelAST: &ast.File{
			Name: ast.NewIdent("generated"),
		},
		output: modelOutput,
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
		g.modelAST.Decls = append(g.modelAST.Decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{
						Name: string(input.Name),
					},
					Type: &ast.StructType{
						Fields: generateModelField(input.Fields),
					},
				},
			},
		})
	}

	for _, t := range g.Schema.Types {
		g.modelAST.Decls = append(g.modelAST.Decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{
						Name: string(t.Name),
					},
					Type: &ast.StructType{
						Fields: generateModelField(t.Fields),
					},
				},
			},
		})
	}

	format.Node(g.output, token.NewFileSet(), g.modelAST)

	return nil
}

func golangType(fieldType *schema.FieldType, graphQLType GraphQLType) *ast.Ident {
	if fieldType.IsList {
		return ast.NewIdent("[]" + golangType(fieldType.ListType, GraphQLType(fieldType.ListType.Name)).Name)
	}

	if graphQLType.IsPrimitive() {
		if fieldType.Nullable {
			return ast.NewIdent("*" + graphQLType.golangType())
		}

		return ast.NewIdent(graphQLType.golangType())
	}

	if fieldType.Nullable {
		return ast.NewIdent("*" + graphQLType.golangType())
	}

	return ast.NewIdent(graphQLType.golangType())
}

type GraphQLType string

func (g GraphQLType) IsPrimitive() bool {
	switch g {
	case "Int", "Float", "String", "Boolean", "ID":
		return true
	default:
		return false
	}
}

func (g GraphQLType) golangType() string {
	switch g {
	case "Int":
		return "int"
	case "Float":
		return "float64"
	case "String":
		return "string"
	case "Boolean":
		return "bool"
	case "ID":
		return "string"
	default:
		return string(g)
	}
}

func toUpperCase(s string) string {
	return string(s[0]-32) + s[1:]
}
