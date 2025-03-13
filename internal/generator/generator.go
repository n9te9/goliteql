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
	resolverAST *ast.File
	modelPackagePath string
	resolverPackagePath string

	modelOutput io.Writer
	resolverOutput io.Writer
}

var gqlFilePattern = regexp.MustCompile(`^.+\.gql$|^.+\.graphql$`)

func NewGenerator(schemaDirectory string, modelOutput, resolverOutput io.Writer, modelPackagePath, resolverPackagePath string) (*Generator, error) {
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
			Name: ast.NewIdent(filepath.Base(modelPackagePath)),
		},
		resolverAST: &ast.File{
			Name: ast.NewIdent(filepath.Base(resolverPackagePath)),
		},
		modelOutput: modelOutput,
		modelPackagePath: modelPackagePath,
		resolverOutput: resolverOutput,
		resolverPackagePath: resolverPackagePath,
	}

	return g, nil
}

func (g *Generator) Generate() error {
	if err := g.generateModel(); err != nil {
		panic(err)
	}

	// generate resolver code
	if err := g.generateResolver(); err != nil {
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

	format.Node(g.modelOutput, token.NewFileSet(), g.modelAST)

	return nil
}

func (g *Generator) generateResolver() error {
	if isUsedDefinedType(g.Schema.GetQuery()) || isUsedDefinedType(g.Schema.GetMutation()) || isUsedDefinedType(g.Schema.GetSubscription()) {
		// generate import statement
		g.resolverAST.Decls = append(g.resolverAST.Decls, &ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Name: ast.NewIdent(filepath.Base(g.modelPackagePath)),
					Path: &ast.BasicLit{
						Kind: token.STRING,
						Value: fmt.Sprintf(`"%s"`, g.modelPackagePath),
					},
				},
			},
		})
	}

	g.resolverAST.Decls = append(g.resolverAST.Decls, 
		generateResolverStruct(g.Schema.GetQuery(), g.Schema.GetMutation(), g.Schema.GetSubscription()),
		generateResolverServeHTTP(g.Schema.GetQuery(), g.Schema.GetMutation(), g.Schema.GetSubscription()),
		generateInterfaceField(g.Schema.GetQuery(), g.modelPackagePath),
		generateInterfaceField(g.Schema.GetMutation(), g.modelPackagePath))

	if err := format.Node(g.resolverOutput, token.NewFileSet(), g.resolverAST); err != nil {
		return fmt.Errorf("error formatting resolver: %w", err)
	}

	return nil
}

func golangType(fieldType *schema.FieldType, graphQLType GraphQLType, modelPackagePath string) *ast.Ident {
	if fieldType.IsList {
		return ast.NewIdent("[]" + golangType(fieldType.ListType, GraphQLType(fieldType.ListType.Name), modelPackagePath).Name)
	}

	if graphQLType.IsPrimitive() {
		if fieldType.Nullable {
			return ast.NewIdent("*" + graphQLType.golangType())
		}

		return ast.NewIdent(graphQLType.golangType())
	}

	modelPackagePrefix := filepath.Base(modelPackagePath)
	if fieldType.Nullable {
		return ast.NewIdent("*" + modelPackagePrefix + "." + graphQLType.golangType())
	}

	return ast.NewIdent(modelPackagePrefix + "." + graphQLType.golangType())
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