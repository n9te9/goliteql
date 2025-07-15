package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/n9te9/goliteql/schema"
	"golang.org/x/tools/imports"
)

type Generator struct {
	Schema                     *schema.Schema
	queryAST                   *ast.File
	mutationAST                *ast.File
	subscriptionAST            *ast.File
	modelAST                   *ast.File
	modelPackagePath           string
	resolverPackagePath        string
	rootResolverOutputFilePath string

	modelOutput                 io.Writer
	queryResolverOutput         io.Writer
	queryResolverAST            *ast.File
	queryResolverOutputFilePath string

	mutationResolverOutput         io.Writer
	mutationResolverAST            *ast.File
	mutationResolverOutputFilePath string

	rootResolverOutput io.Writer
	resolverAST        *ast.File

	enumOutput io.Writer
	enumAST    *ast.File
}

type Config struct {
	SchemaDirectory            string `yaml:"schema_directory"`
	ModelOutputFile            string `yaml:"model_output_file"`
	QueryResolverOutputFile    string `yaml:"query_resolver_output_file"`
	MutationResolverOutputFile string `yaml:"mutation_resolver_output_file"`
	RootResolverOutputFile     string `yaml:"root_resolver_output_file"`
	EnumOutputFile             string `yaml:"enum_output_file"`
	ModelPackageName           string `yaml:"model_package_name"`
	ResolverPackageName        string `yaml:"resolver_package_name"`
}

var gqlFilePattern = regexp.MustCompile(`^.+\.gql$|^.+\.graphql$`)

func createDirectories(conf *Config) {
	if err := os.MkdirAll(conf.SchemaDirectory, 0755); err != nil && !os.IsExist(err) {
		log.Fatalf("error creating schema directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(conf.ModelOutputFile), 0755); err != nil {
		log.Fatalf("error creating model output directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(conf.QueryResolverOutputFile), 0755); err != nil {
		log.Fatalf("error creating query resolver output directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(conf.MutationResolverOutputFile), 0755); err != nil {
		log.Fatalf("error creating mutation resolver output directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(conf.RootResolverOutputFile), 0755); err != nil {
		log.Fatalf("error creating root resolver output directory: %v", err)
	}
}

func createFile(filePath string) (*os.File, error) {
	file, err := os.Create(filePath)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("error creating file %s: %w", filePath, err)
	}

	return file, nil
}

func NewGenerator(config *Config) (*Generator, error) {
	createDirectories(config)

	schemaDirectory := config.SchemaDirectory
	modelPackagePath := config.ModelPackageName
	resolverPackagePath := config.ResolverPackageName

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

	var modelOutput, queryResolverOutput, mutationResolverOutput, rootResolverOutput, enumOutput io.Writer
	if len(extractUserEnumDefinitions(s.Enums)) > 0 {
		enumOutput, err = createFile(config.EnumOutputFile)
		if err != nil {
			return nil, fmt.Errorf("error creating enum output file: %w", err)
		}
	}

	if s.Definition.Query != nil {
		queryResolverOutput, err = createFile(config.QueryResolverOutputFile)
		if err != nil {
			return nil, fmt.Errorf("error creating query resolver output file: %w", err)
		}
	}

	if s.Definition.Mutation != nil {
		mutationResolverOutput, err = createFile(config.MutationResolverOutputFile)
		if err != nil {
			return nil, fmt.Errorf("error creating mutation resolver output file: %w", err)
		}
	}

	if s.Definition.Subscription != nil {
		rootResolverOutput, err = createFile(config.RootResolverOutputFile)
		if err != nil {
			return nil, fmt.Errorf("error creating root resolver output file: %w", err)
		}
	}
	modelOutput, err = createFile(config.ModelOutputFile)
	if err != nil {
		return nil, fmt.Errorf("error creating model output file: %w", err)
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
		enumAST: &ast.File{
			Name: ast.NewIdent(filepath.Base(modelPackagePath)),
		},
		queryResolverAST: &ast.File{
			Name:  ast.NewIdent(filepath.Base(resolverPackagePath)),
			Decls: []ast.Decl{},
		},
		mutationResolverAST: &ast.File{
			Name:  ast.NewIdent(filepath.Base(resolverPackagePath)),
			Decls: []ast.Decl{},
		},
		modelOutput:                    modelOutput,
		modelPackagePath:               modelPackagePath,
		queryResolverOutput:            queryResolverOutput,
		mutationResolverOutput:         mutationResolverOutput,
		rootResolverOutput:             rootResolverOutput,
		enumOutput:                     enumOutput,
		resolverPackagePath:            resolverPackagePath,
		queryResolverOutputFilePath:    config.QueryResolverOutputFile,
		mutationResolverOutputFilePath: config.MutationResolverOutputFile,
		rootResolverOutputFilePath:     config.RootResolverOutputFile,
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
	// g.enumAST.Decls = append(g.enumAST.Decls, generateEnumImport())

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
		if t.IsIntrospection() {
			continue
		}

		fields := generateModelField(t.Fields)

		if len(t.Interfaces) > 0 {
			fields.List = append(fields.List, &ast.Field{
				Type: &ast.BasicLit{},
			})
		}

		for _, iface := range t.Interfaces {
			if iface.Name[0] >= 'a' && iface.Name[0] <= 'z' {
				panic(fmt.Sprintf("interface name %s should start with uppercase letter", iface.Name))
			}

			fields.List = append(fields.List, &ast.Field{
				Names: []*ast.Ident{},
				Type:  ast.NewIdent(string(iface.Name)),
			})
		}

		decl := &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{
						Name: string(t.Name),
					},
					Type: &ast.StructType{
						Fields: fields,
					},
				},
			},
		}
		g.modelAST.Decls = append(g.modelAST.Decls, decl)

		if t.PrimitiveTypeName != nil {
			g.modelAST.Decls = append(g.modelAST.Decls, &ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: string(t.Name),
						},
						Type: &ast.Ident{
							Name: string(t.PrimitiveTypeName),
						},
					},
				},
			})
		}
	}

	g.modelAST.Decls = append(g.modelAST.Decls, generateInterfaceTypeDecls(g.Schema.Interfaces)...)

	userEnums := extractUserEnumDefinitions(g.Schema.Enums)
	g.enumAST.Decls = append(g.enumAST.Decls, generateEnumModelAST(userEnums)...)

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

	if g.enumOutput != nil {
		if err := format.Node(g.enumOutput, token.NewFileSet(), g.enumAST); err != nil {
			return fmt.Errorf("error formatting enum: %w", err)
		}
	}

	return nil
}

func generateOperationImport(operation *schema.OperationDefinition, modelPackagePath string) []ast.Spec {
	specs := make([]ast.Spec, 0)

	specs = append(specs,
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"context"`,
			},
		})

	if useDefinedType(operation) {
		specs = append(specs,
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, modelPackagePath),
				},
			})
	}

	return specs
}

func extractUserEnumDefinitions(enums []*schema.EnumDefinition) []*schema.EnumDefinition {
	ret := make([]*schema.EnumDefinition, 0, len(enums))
	for _, e := range enums {
		if !e.IsIntrospection() {
			ret = append(ret, e)
		}
	}

	return ret
}

func extractIntrospectionEnumDefinitions(enums []*schema.EnumDefinition) []*schema.EnumDefinition {
	ret := make([]*schema.EnumDefinition, 0, len(enums))
	for _, e := range enums {
		if e.IsIntrospection() {
			ret = append(ret, e)
		}
	}

	return ret
}

func (g *Generator) generateResolver() error {
	if isUsedDefinedType(g.Schema.GetQuery()) || isUsedDefinedType(g.Schema.GetMutation()) || isUsedDefinedType(g.Schema.GetSubscription()) {
		importSpecs := []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, g.modelPackagePath),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/n9te9/goliteql/query"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/n9te9/goliteql/query/utils"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/n9te9/goliteql/executor"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/n9te9/goliteql"`,
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

	modelPrefix := filepath.Base(g.modelPackagePath)

	if q := g.Schema.GetQuery(); q != nil {
		queryFields = q.Fields
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateQueryExecutor(q))
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateApplyQueryResponseFuncDecls(q, g.Schema.Indexes, 0, modelPrefix)...)
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateOperationArgumentDecls(modelPrefix, q, g.Schema.Indexes)...)
	}

	if m := g.Schema.GetMutation(); m != nil {
		mutationFields = m.Fields
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateMutationExecutor(m))
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateApplyQueryResponseFuncDecls(m, g.Schema.Indexes, 0, modelPrefix)...)
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateOperationArgumentDecls(modelPrefix, m, g.Schema.Indexes)...)
	}

	if s := g.Schema.GetSubscription(); s != nil {
		fields = append(fields, s.Fields...)
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateSubscriptionExecutor(g.Schema.GetSubscription()))
		g.resolverAST.Decls = append(g.resolverAST.Decls, generateWrapResponseWriter(g.Schema.GetSubscription())...)
	}

	if g.Schema.GetQuery() != nil {
		g.queryResolverAST.Decls = append(g.queryResolverAST.Decls, &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: generateOperationImport(g.Schema.GetQuery(), g.modelPackagePath),
		})
		g.queryResolverAST.Decls = append(g.queryResolverAST.Decls, generateInterfaceField(modelPrefix, g.Schema.GetQuery(), g.Schema.Indexes))
	}

	if g.Schema.GetMutation() != nil {
		g.mutationResolverAST.Decls = append(g.mutationResolverAST.Decls, &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: generateOperationImport(g.Schema.GetMutation(), g.modelPackagePath),
		})
		g.mutationResolverAST.Decls = append(g.mutationResolverAST.Decls, generateInterfaceField(modelPrefix, g.Schema.GetMutation(), g.Schema.Indexes))
	}

	g.resolverAST.Decls = append(g.resolverAST.Decls, generateResolverImplementationStruct()...)
	g.resolverAST.Decls = append(g.resolverAST.Decls, generateResolverImplementation(modelPrefix, fields, g.Schema.Indexes)...)

	g.queryResolverAST.Decls = append(g.queryResolverAST.Decls, generateResolverImplementation(modelPrefix, queryFields, g.Schema.Indexes)...)
	g.mutationResolverAST.Decls = append(g.mutationResolverAST.Decls, generateResolverImplementation(modelPrefix, mutationFields, g.Schema.Indexes)...)

	g.resolverAST.Decls = append(g.resolverAST.Decls, generateResolverServeHTTP(g.Schema.GetQuery(), g.Schema.GetMutation(), g.Schema.GetSubscription()))
	g.resolverAST.Decls = append(g.resolverAST.Decls, generateOperationResponseStructDecls(g.Schema)...)

	// Introspection generation
	// g.resolverAST.Decls = append(g.resolverAST.Decls, g.generateIntrospection(g.modelPackagePath)...)

	var rootResolverBuffer bytes.Buffer
	if err := format.Node(&rootResolverBuffer, token.NewFileSet(), g.resolverAST); err != nil {
		return fmt.Errorf("error formatting resolver: %w", err)
	}

	var queryResolverBuffer bytes.Buffer
	if err := format.Node(&queryResolverBuffer, token.NewFileSet(), g.queryResolverAST); err != nil {
		return fmt.Errorf("error formatting query resolver: %w", err)
	}

	var mutationResolverBuffer bytes.Buffer
	if err := format.Node(&mutationResolverBuffer, token.NewFileSet(), g.mutationResolverAST); err != nil {
		return fmt.Errorf("error formatting mutation resolver: %w", err)
	}

	fixed, err := imports.Process(g.rootResolverOutputFilePath, rootResolverBuffer.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("error processing root resolver imports: %w", err)
	}
	if _, err := g.rootResolverOutput.Write(fixed); err != nil {
		return fmt.Errorf("error writing root resolver output: %w", err)
	}
	fixed, err = imports.Process(g.queryResolverOutputFilePath, queryResolverBuffer.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("error processing query resolver imports: %w", err)
	}
	if _, err := g.queryResolverOutput.Write(fixed); err != nil {
		return fmt.Errorf("error writing query resolver output: %w", err)
	}
	fixed, err = imports.Process(g.mutationResolverOutputFilePath, mutationResolverBuffer.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("error processing mutation resolver imports: %w", err)
	}
	if _, err := g.mutationResolverOutput.Write(fixed); err != nil {
		return fmt.Errorf("error writing mutation resolver output: %w", err)
	}

	return nil
}

func (g *Generator) generateIntrospection(modelPrefix string) []ast.Decl {
	ret := make([]ast.Decl, 0)

	ret = append(ret, generateIntrospectionModelAST(g.Schema.Types)...)
	ret = append(ret, generateIntrospectionSchemaQueryAST(g.Schema))
	ret = append(ret, generateIntrospectionTypesFuncDecl(g.Schema))
	ret = append(ret, generateIntrospectionSchemaResponseModelAST())
	ret = append(ret, generateIntrospectionSchemaResponseDataModelAST())
	ret = append(ret, generateIntrospectionTypeResponseDataModelAST())
	ret = append(ret, generateIntrospectionTypeResponseModelAST())
	ret = append(ret, generateIntrospectionQueryTypeMethodAST(g.Schema))
	ret = append(ret, generateIntrospectionMutationTypeMethodAST(g.Schema))
	ret = append(ret, generateIntrospectionTypeMethodDecls(g.Schema)...)
	ret = append(ret, generateIntrospectionFieldTypeTypeOfDecls(g.Schema)...)
	ret = append(ret, generateIntrospectionTypeFieldsDecls(g.Schema.Types)...)
	ret = append(ret, generateIntrospectionInterfaceFieldsDecls(g.Schema.Interfaces)...)
	ret = append(ret, generateIntrospectionInterfaceTypeFuncDecls(g.Schema.Interfaces, g.Schema.Indexes)...)
	ret = append(ret, generateIntrospectionTypeFieldsFuncDecls(g.Schema.Types, g.Schema.Indexes)...)
	ret = append(ret, generateIntrospectionTypeResolverDeclsFromInterfaces(g.Schema.Interfaces, g.Schema.Indexes)...)
	ret = append(ret, generateExtractOperationArgumentsDecl(modelPrefix, fieldsIntrospectionFieldDefinition, g.Schema.Indexes))
	ret = append(ret, generateIntrospectionTypeFuncDecl(g.Schema))
	ret = append(ret, generateIntrospectionTypeFuncDecls(g.Schema.Types)...)
	ret = append(ret, generateIntrospectionInputFuncDecls(g.Schema.Inputs)...)
	ret = append(ret, generateIntrospectionInputFieldsDecls(g.Schema.Inputs)...)
	ret = append(ret, generateIntrospectionInputFieldsFuncDecls(g.Schema.Inputs, g.Schema.Indexes)...)
	ret = append(ret, generateIntrospectionScalarFuncDecls(g.Schema.Scalars)...)
	ret = append(ret, generateEnumModelAST(extractIntrospectionEnumDefinitions(g.Schema.Enums))...)
	ret = append(ret, generateIntrospectionEnumFuncDecls(g.Schema.Enums)...)
	ret = append(ret, generateIntrospectionEnumValuesFuncDecl(g.Schema.Enums)...)
	ret = append(ret, generateIntrospectionOperationFuncDecls(g.Schema)...)
	ret = append(ret, generateIntrospectionUnionTypeFuncDecls(g.Schema.Unions, g.Schema.Indexes)...)

	if q := g.Schema.GetQuery(); q != nil {
		ret = append(ret, generateIntrospectionFieldsFuncsAST(string(g.Schema.Definition.Query), q.Fields)...)
		ret = append(ret, generateIntrospectionOperationFieldFuncDecls(q)...)
	}

	if m := g.Schema.GetMutation(); m != nil {
		ret = append(ret, generateIntrospectionFieldsFuncsAST(string(g.Schema.Definition.Mutation), m.Fields)...)
		ret = append(ret, generateIntrospectionOperationFieldFuncDecls(m)...)
	}

	return ret
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

func isLowerCase(s string) bool {
	return s[0] >= 'a' && s[0] <= 'z'
}

type FieldName string

func (f FieldName) ExportedGolangFieldName() string {
	// check uppercase
	if f[0] >= 'A' && f[0] <= 'Z' {
		return string(f)
	}

	// check lowercase
	if f[0] >= 'a' && f[0] <= 'z' {
		return string(f[0]-32) + string(f[1:])
	}

	panic(fmt.Sprintf("invalid field name: %s", f))
}

func useDefinedType(operationDefinition *schema.OperationDefinition) bool {
	if operationDefinition == nil {
		return false
	}

	for _, field := range operationDefinition.Fields {
		for _, arg := range field.Arguments {
			if useDefinedTypeRecursive(arg.Type) {
				return true
			}
		}

		if useDefinedTypeRecursive(field.Type) {
			return true
		}
	}

	return false
}

func useDefinedTypeRecursive(fieldType *schema.FieldType) bool {
	if fieldType.IsList {
		return useDefinedTypeRecursive(fieldType.ListType)
	}

	return !fieldType.IsPrimitive()
}
