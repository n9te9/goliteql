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
			Name: ast.NewIdent("generated_model"),
		},
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
						Fields: g.generateModelField(input.Fields),
					},
				},
			},
		})

		fmt.Println(string(input.Name))
	}

	for _, t := range g.Schema.Types {
		fmt.Println(string(t.Name))
	}

	// output to file
	// f, err := os.Create("model.go")
	// if err != nil {
	// 	return fmt.Errorf("error creating file: %w", err)
	// }

	format.Node(os.Stdout, token.NewFileSet(), g.modelAST)

	return nil
}

func (g *Generator) generateModelField(field schema.FieldDefinitions) *ast.FieldList {
	fields := make([]*ast.Field, 0)

	for _, f := range field {
		fieldType := GraphQLType(f.Type.Name)
		var fieldTypeIdent *ast.Ident
		if fieldType.IsPrimitive() {
			fieldTypeIdent = ast.NewIdent(fieldType.golangType())
		}

		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{
				{
					Name: toUpperCase(string(f.Name)),
				},
			},
			Type: fieldTypeIdent,
			Tag: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", string(f.Name)),
			},
		})
	}

	return &ast.FieldList{
		List: fields,
	}
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

func toLowerCase(s string) string {
	return string(s[0]+32) + s[1:]
}