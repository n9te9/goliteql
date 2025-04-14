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

	"github.com/lkeix/gg-executor/schema"
)

type Generator struct {
	Schema              *schema.Schema
	queryAST            *ast.File
	mutationAST         *ast.File
	subscriptionAST     *ast.File
	modelAST            *ast.File
	modelPackagePath    string
	resolverPackagePath string

	modelOutput    io.Writer
	queryResolverOutput io.Writer
	queryResolverAST         *ast.File

	mutationResolverOutput io.Writer
	mutationResolverAST         *ast.File

	rootResolverOutput io.Writer
	resolverAST         *ast.File
}

var gqlFilePattern = regexp.MustCompile(`^.+\.gql$|^.+\.graphql$`)

func NewGenerator(schemaDirectory string, modelOutput, queryResolverOutput, mutationResolverOutput, rootResolverOutput io.Writer, modelPackagePath, resolverPackagePath string) (*Generator, error) {
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

	importDecl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"net/http"`,
				},
			},
		},
	}

	g := &Generator{
		Schema:          s,
		queryAST:        &ast.File{},
		mutationAST:     &ast.File{},
		subscriptionAST: &ast.File{},
		modelAST: &ast.File{
			Name: ast.NewIdent(filepath.Base(modelPackagePath)),
		},
		resolverAST: &ast.File{
			Name: ast.NewIdent(filepath.Base(resolverPackagePath)),
		},
		queryResolverAST: &ast.File{
			Name: ast.NewIdent(filepath.Base(resolverPackagePath)),
			Decls: []ast.Decl{
				importDecl,
			},
		},
		mutationResolverAST: &ast.File{
			Name: ast.NewIdent(filepath.Base(resolverPackagePath)),
			Decls: []ast.Decl{
				importDecl,
			},
		},
		modelOutput:         modelOutput,
		modelPackagePath:    modelPackagePath,
		queryResolverOutput: queryResolverOutput,
		mutationResolverOutput: mutationResolverOutput,
		rootResolverOutput:     rootResolverOutput,
		resolverPackagePath: resolverPackagePath,
	}

	return g, nil
}

func (g *Generator) Generate() error {
	// generate resolver code
	if err := g.generateResolver(); err != nil {
		panic(err)
	}

	if err := g.generateModel(); err != nil {
		panic(err)
	}

	return nil
}

func (g *Generator) generateModel() error {
	g.modelAST.Decls = append(g.modelAST.Decls, generateModelImport())

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

		g.modelAST.Decls = append(g.modelAST.Decls, generateInputModelUnmarshalJSON(input))
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

	if op := g.Schema.GetQuery(); op != nil {
		g.modelAST.Decls = append(g.modelAST.Decls, generateSelectionSetInput(op.Fields)...)
	}

	if op := g.Schema.GetMutation(); op != nil {
		g.modelAST.Decls = append(g.modelAST.Decls, generateSelectionSetInput(op.Fields)...)
	}

	if op := g.Schema.GetSubscription(); op != nil {
		g.modelAST.Decls = append(g.modelAST.Decls, generateSelectionSetInput(op.Fields)...)
	}

	format.Node(g.modelOutput, token.NewFileSet(), g.modelAST)

	return nil
}

func (g *Generator) generateResolver() error {
	if isUsedDefinedType(g.Schema.GetQuery()) || isUsedDefinedType(g.Schema.GetMutation()) || isUsedDefinedType(g.Schema.GetSubscription()) {
		importSpecs := []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"io"`,
				},
			},
			&ast.ImportSpec{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: `// remove _, when use model package`,
						},
					},
				},
				Name: ast.NewIdent("_"),
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, g.modelPackagePath),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/lkeix/gg-executor/query"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/lkeix/gg-executor/query/utils"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/lkeix/gg-executor/executor"`,
				},
			},
		}

		importSpecs = append(importSpecs, generateResolverImport().Specs...)

		// generate import statement
		g.resolverAST.Decls = append(g.resolverAST.Decls, &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: importSpecs,
		})
	}

	g.resolverAST.Decls = append(g.resolverAST.Decls, generateResolverInterface(g.Schema.GetQuery(), g.Schema.GetMutation(), g.Schema.GetSubscription()))

	queryFields := make(schema.FieldDefinitions, 0)
	mutationFields := make(schema.FieldDefinitions, 0)
	fields := make(schema.FieldDefinitions, 0)

	if q := g.Schema.GetQuery(); q != nil {
		queryFields = q.Fields
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateQueryExecutor(g.Schema.GetQuery()))
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateWrapResponseWriter(g.Schema.GetQuery(), g.Schema.Indexes.TypeIndex)...)
	}

	if m := g.Schema.GetMutation(); m != nil {
		mutationFields = m.Fields
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateMutationExecutor(g.Schema.GetMutation()))
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateWrapResponseWriter(g.Schema.GetMutation(), g.Schema.Indexes.TypeIndex)...)
	}

	if s := g.Schema.GetSubscription(); s != nil {
		fields = append(fields, s.Fields...)
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateSubscriptionExecutor(g.Schema.GetSubscription()))
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateWrapResponseWriter(g.Schema.GetSubscription(), g.Schema.Indexes.TypeIndex)...)
	}

	if g.Schema.GetQuery() != nil {
		g.queryResolverAST.Decls = append(g.queryResolverAST.Decls, generateInterfaceField(g.Schema.GetQuery()))
	}

	if g.Schema.GetMutation() != nil {
		g.mutationResolverAST.Decls = append(g.mutationResolverAST.Decls, generateInterfaceField(g.Schema.GetMutation()))
	}

	g.resolverAST.Decls = append(g.resolverAST.Decls, generateResolverImplementationStruct()...)
	g.resolverAST.Decls = append(g.resolverAST.Decls, generateResolverImplementation(fields)...)

	g.queryResolverAST.Decls = append(g.queryResolverAST.Decls, generateResolverImplementation(queryFields)...)
	g.mutationResolverAST.Decls = append(g.mutationResolverAST.Decls, generateResolverImplementation(mutationFields)...)

	g.resolverAST.Decls = append(g.resolverAST.Decls, generateResolverServeHTTP(g.Schema.GetQuery(), g.Schema.GetMutation(), g.Schema.GetSubscription()))

	if err := format.Node(g.rootResolverOutput, token.NewFileSet(), g.resolverAST); err != nil {
		return fmt.Errorf("error formatting resolver: %w", err)
	}

	if err := format.Node(g.queryResolverOutput, token.NewFileSet(), g.queryResolverAST); err != nil {
		return fmt.Errorf("error formatting query resolver: %w", err)
	}

	if err := format.Node(g.mutationResolverOutput, token.NewFileSet(), g.mutationResolverAST); err != nil {
		return fmt.Errorf("error formatting mutation resolver: %w", err)
	}

	return nil
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
