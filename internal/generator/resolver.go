package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func newQueryIdent(query *schema.OperationDefinition) *ast.Ident {
	return ast.NewIdent(newQueryName(query))
}

func newQueryName(query *schema.OperationDefinition) string {
	queryName := "Query"
	if query != nil {
		if len(query.Name) != 0 {
			queryName = string(query.Name)
		}
	}

	return queryName + "Resolver"
}

func newMutationIdent(mutation *schema.OperationDefinition) *ast.Ident {
	return ast.NewIdent(newMutationName(mutation))
}

func newMutationName(mutation *schema.OperationDefinition) string {
	mutationName := "Mutation"
	if mutation != nil {
		if len(mutation.Name) != 0 {
			mutationName = string(mutation.Name)
		}
	}

	return mutationName + "Resolver"
}

func newSubscriptionIdent(subscription *schema.OperationDefinition) *ast.Ident {
	return ast.NewIdent(newSubscriptionName(subscription))
}

func newSubscriptionName(subscription *schema.OperationDefinition) string {
	subscriptionName := "Subscription"
	if subscription != nil {
		if len(subscription.Name) != 0 {
			subscriptionName = string(subscription.Name)
		}
	}

	return subscriptionName + "Resolver"
}

func generateResolverImport() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"net/http"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"fmt"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"strings"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"encoding/json"`,
				},
			},
		},
	}
}

func generateResolverInterface(query, mutation, subscription *schema.OperationDefinition) *ast.GenDecl {
	generateField := func(query, mutation, subscription *schema.OperationDefinition) []*ast.Field {
		fields := make([]*ast.Field, 0, 3)
		if query != nil {
			fields = append(fields, &ast.Field{
				Type: newQueryIdent(query),
			})
		}

		if mutation != nil {
			fields = append(fields, &ast.Field{
				Type: newMutationIdent(mutation),
			})
		}

		if subscription != nil {
			fields = append(fields, &ast.Field{
				Type: newSubscriptionIdent(subscription),
			})
		}

		return fields
	}

	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "Resolver",
				},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: generateField(query, mutation, subscription),
					},
				},
			},
		},
	}
}

func generateResolverServeHTTP(query, mutation, subscription *schema.OperationDefinition) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("ServeHTTP"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
				},
			},
		},
		Type: &ast.FuncType{
			Params:  generateServeHTTPArgs(),
			Results: &ast.FieldList{},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "// *********** AUTO GENERATED CODE ***********",
				},
				{
					Text: "// *********** DON'T EDIT ***********",
				},
			},
		},
		Body: generateServeHTTPBody(query, mutation, subscription),
	}
}

func generateTypeExprFromFieldType(fieldType *schema.FieldType) ast.Expr {
	if fieldType.IsList {
		return &ast.ArrayType{
			Elt: generateTypeExprFromFieldType(fieldType.ListType),
		}
	}

	graphQLType := GraphQLType(fieldType.Name)

	var baseTypeExpr ast.Expr = ast.NewIdent(graphQLType.golangType())
	if !graphQLType.IsPrimitive() {
		baseTypeExpr = &ast.SelectorExpr{
			X:   ast.NewIdent("model"),
			Sel: ast.NewIdent(graphQLType.golangType()),
		}
	}

	if fieldType.Nullable {
		return &ast.StarExpr{
			X: baseTypeExpr,
		}
	}

	return baseTypeExpr
}

func generateResponseGraphQLResponseStruct(operationName string, fieldDefinition *schema.FieldDefinition) ast.Decl {
	structName := fmt.Sprintf("%sGraphQLResponse", operationName)

	return &ast.GenDecl{
		Tok: token.TYPE,
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: fmt.Sprintf("// use %s for %s resposne", structName, operationName),
				},
			},
		},
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(structName),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("Data"),
								},
								Type: generateTypeExprFromFieldType(fieldDefinition.Type),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("Errors"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("executor"),
									Sel: ast.NewIdent("GraphQLError"),
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateQueryExecutor(query *schema.OperationDefinition) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("queryExecutor"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
				},
			},
		},
		Type: &ast.FuncType{
			Params:  generateOperationExecutorArgs(),
			Results: &ast.FieldList{},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "// *********** AUTO GENERATED CODE ***********",
				},
				{
					Text: "// *********** DON'T EDIT ***********",
				},
			},
		},
		Body: generateExecutorBody(query, "query"),
	}
}

func generateMutationExecutor(mutation *schema.OperationDefinition) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("mutationExecutor"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
				},
			},
		},
		Type: &ast.FuncType{
			Params:  generateOperationExecutorArgs(),
			Results: &ast.FieldList{},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "// *********** AUTO GENERATED CODE ***********",
				},
				{
					Text: "// *********** DON'T EDIT ***********",
				},
			},
		},
		Body: generateExecutorBody(mutation, "mutation"),
	}
}

func generateSubscriptionExecutor(subscription *schema.OperationDefinition) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("subscriptionExecutor"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
				},
			},
		},
		Type: &ast.FuncType{
			Params:  generateOperationExecutorArgs(),
			Results: &ast.FieldList{},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "// *********** AUTO GENERATED CODE ***********",
				},
				{
					Text: "// *********** DON'T EDIT ***********",
				},
			},
		},
		Body: generateExecutorBody(subscription, "subscription"),
	}
}

func generateWrapResponseWriter(op *schema.OperationDefinition) []ast.Decl {
	res := make([]ast.Decl, 0, len(op.Fields))

	for _, field := range op.Fields {
		res = append(res, generateWrapResponseWriterStruct(field))
		res = append(res, generateWrapResponseWriterFunc(field))
		res = append(res, generateWrapResponseWriterWrite(string(field.Name), field))
	}

	return res
}

func generateWrapResponseWriterStruct(field *schema.FieldDefinition) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("Wrap" + string(field.Name) + "ResponseWriter"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("http"),
									Sel: ast.NewIdent("ResponseWriter"),
								},
							},
							{
								Names: []*ast.Ident{ast.NewIdent("selections")},
								Type: &ast.ArrayType{
									Len: nil,
									Elt: &ast.SelectorExpr{
										X:   ast.NewIdent("query"),
										Sel: ast.NewIdent("Selection"),
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("variables"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("json"),
									Sel: ast.NewIdent("RawMessage"),
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateWrapResponseWriterFunc(field *schema.FieldDefinition) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("new" + string(field.Name) + "Writer"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("w")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("http"),
							Sel: ast.NewIdent("ResponseWriter"),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("selections")},
						Type: &ast.ArrayType{
							Len: nil,
							Elt: &ast.SelectorExpr{
								X:   ast.NewIdent("query"),
								Sel: ast.NewIdent("Selection"),
							},
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("variables"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("json"),
							Sel: ast.NewIdent("RawMessage"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.Ident{
								Name: "Wrap" + string(field.Name) + "ResponseWriter",
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent("Wrap" + string(field.Name) + "ResponseWriter"),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("ResponseWriter"),
										Value: ast.NewIdent("w"),
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("selections"),
										Value: ast.NewIdent("selections"),
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("variables"),
										Value: ast.NewIdent("variables"),
									},
								},
							},
						},
					},
				},
			},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "// *********** AUTO GENERATED CODE ***********",
				},
				{
					Text: "// *********** DON'T EDIT ***********",
				},
			},
		},
	}
}

func extractWillDeclTypeDefinition(typeIndex map[string]*schema.TypeDefinition, fieldType *schema.FieldType) *schema.TypeDefinition {
	if fieldType.IsList {
		return extractWillDeclTypeDefinition(typeIndex, fieldType.ListType)
	}

	typeDefinition := typeIndex[string(fieldType.Name)]
	if typeDefinition == nil {
		return nil
	}

	return typeDefinition
}

func generateResponseStructForWrapResponseWriter(typeIndex map[string]*schema.TypeDefinition, operation *schema.OperationDefinition) []ast.Decl {
	if operation == nil {
		return nil
	}

	decls := make([]ast.Decl, 0)

	for _, field := range operation.Fields {
		if field == nil {
			continue
		}

		t := extractWillDeclTypeDefinition(typeIndex, field.Type)

		decls = append(decls, generateResponseStructDeclsForWrapResponseWriter(string(field.Name), field, t, typeIndex)...)
	}

	return decls
}

func generateResponseStructArrayType(nestCount int, responseStructName, prefix string) ast.Expr {
	if nestCount > 0 {
		return &ast.ArrayType{
			Elt: generateResponseStructArrayType(nestCount-1, responseStructName, prefix),
		}
	}

	if prefix == "" {
		return &ast.ArrayType{
			Elt: &ast.StarExpr{
				X: ast.NewIdent(responseStructName),
			},
		}
	}

	return &ast.ArrayType{
		Elt: &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(prefix),
				Sel: ast.NewIdent(responseStructName),
			},
		},
	}
}

// TODO: refactor
func generateWrapResponseWriterNestedTypeInitializerForFieldType(fieldType *schema.FieldType, responseStructName string, maxNestCount, nestCount int) []ast.Stmt {
	nestedItr := ""
	respNestedItr := ""

	arrayType := generateResponseStructArrayType(maxNestCount-nestCount-1, responseStructName, "")
	if nestCount >= 0 {
		for i := 0; i < nestCount+1; i++ {
			nestedItr += fmt.Sprintf("[k%d]", i)
		}

		for i := 0; i < nestCount; i++ {
			respNestedItr += fmt.Sprintf("[k%d]", i)
		}
	}

	stmts := make([]ast.Stmt, 0)
	if maxNestCount != nestCount {
		stmts = append(stmts, &ast.AssignStmt{
			Tok: token.ASSIGN,
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("resp%s", respNestedItr)),
			},
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("append"),
					Args: []ast.Expr{
						ast.NewIdent(fmt.Sprintf("resp%s", respNestedItr)),
						&ast.CallExpr{
							Fun: ast.NewIdent("make"),
							Args: []ast.Expr{
								arrayType,
								&ast.CallExpr{
									Fun: ast.NewIdent("len"),
									Args: []ast.Expr{
										ast.NewIdent("baseResp" + nestedItr),
									},
								},
							},
						},
					},
				},
			},
		})
	}

	if fieldType.IsList {
		stmts = append(stmts, generateWrapResponseWriterNestedTypeInitializerForFieldType(fieldType.ListType, responseStructName, maxNestCount, nestCount+1)...)
		body := &ast.BlockStmt{
			List: stmts,
		}

		return []ast.Stmt{
			&ast.RangeStmt{
				Key:  ast.NewIdent(fmt.Sprintf("k%d", nestCount)),
				Tok:  token.DEFINE,
				X:    ast.NewIdent("baseResp"),
				Body: body,
			},
		}
	}

	if !fieldType.IsList {
		return stmts
	}

	panic("unknown type")
}

func generateWrapResponseWriterNestedTypeInitializer(responseStructName string, field *schema.FieldDefinition) []ast.Stmt {
	nestCount := getFieldSliceNestLevel(field.Type)
	if nestCount < 2 {
		return []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.BasicLit{},
			},
		}
	}

	return generateWrapResponseWriterNestedTypeInitializerForFieldType(field.Type.ListType, responseStructName, nestCount-1, 0)
}

func generateApplySkipDirective() ast.Stmt {
	return &ast.IfStmt{
		Cond: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("executor"),
				Sel: ast.NewIdent("IsSkipped"),
			},
			Args: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("sel"),
					Sel: ast.NewIdent("Directives"),
				},
				&ast.SelectorExpr{
					X:   ast.NewIdent("w"),
					Sel: ast.NewIdent("variables"),
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.BasicLit{
						Kind:  token.CONTINUE,
						Value: "continue",
					},
				},
			},
		},
	}
}

func generateApplyIncludeDirective() ast.Stmt {
	return &ast.IfStmt{
		Cond: &ast.UnaryExpr{
			Op: token.NOT,
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("executor"),
					Sel: ast.NewIdent("IsIncluded"),
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("sel"),
						Sel: ast.NewIdent("Directives"),
					},
					&ast.SelectorExpr{
						X:   ast.NewIdent("w"),
						Sel: ast.NewIdent("variables"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.BasicLit{
						Kind:  token.CONTINUE,
						Value: "continue",
					},
				},
			},
		},
	}
}

func generateWrapResponseWriterReponseFieldWalkerValidationStmts(fieldType *schema.FieldType, typeDefinition *schema.TypeDefinition, nestCount int) []ast.Stmt {
	xName := "baseResp"

	if fieldType.IsList {
		if nestCount > 0 {
			xName = fmt.Sprintf("v%d", nestCount-1)
		}

		valueName := ast.NewIdent(fmt.Sprintf("v%d", nestCount))
		if !fieldType.ListType.IsList {
			valueName = ast.NewIdent("_")
		}

		return []ast.Stmt{
			&ast.RangeStmt{
				Key:   ast.NewIdent(fmt.Sprintf("k%d", nestCount)),
				Tok:   token.DEFINE,
				Value: valueName,
				X:     ast.NewIdent(xName),
				Body: &ast.BlockStmt{
					List: generateWrapResponseWriterReponseFieldWalkerValidationStmts(fieldType.ListType, typeDefinition, nestCount+1),
				},
			},
		}
	}

	nestItr := ""
	for i := 0; i < nestCount; i++ {
		nestItr += fmt.Sprintf("[k%d]", i)
	}

	validationTargetName := "baseResp" + nestItr
	responseTargetName := "resp" + nestItr

	stmts := make([]ast.Stmt, 0)
	for _, field := range typeDefinition.Fields {
		fieldName := string(field.Name)
		if isLowerCase(string(field.Name)) {
			fieldName = toUpperCase(fieldName)
		}

		var assignExpr ast.Expr = &ast.UnaryExpr{
			Op: token.AND,
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(validationTargetName),
				Sel: ast.NewIdent(fieldName),
			},
		}

		if field.Type.Nullable {
			assignExpr = &ast.SelectorExpr{
				X:   ast.NewIdent(validationTargetName),
				Sel: ast.NewIdent(fieldName),
			}
		}

		if field.IsPrimitive() {
			stmts = append(stmts, &ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X: &ast.CallExpr{
						Fun: ast.NewIdent("string"),
						Args: []ast.Expr{
							ast.NewIdent("sel.Name"),
						},
					},
					Op: token.EQL,
					Y: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("\"%s\"", string(field.Name)),
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						generateApplySkipDirective(),
						generateApplyIncludeDirective(),
						&ast.AssignStmt{
							Tok: token.ASSIGN,
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent(responseTargetName),
									Sel: ast.NewIdent(fieldName),
								},
							},
							Rhs: []ast.Expr{
								assignExpr,
							},
						},
					},
				},
			})

			continue
		}

		var argExpr ast.Expr = &ast.SelectorExpr{
			X:   ast.NewIdent(validationTargetName),
			Sel: ast.NewIdent(fieldName),
		}

		if !field.Type.Nullable && !field.Type.IsList {
			argExpr = &ast.SelectorExpr{
				X:   ast.NewIdent(validationTargetName),
				Sel: ast.NewIdent(fieldName),
			}
		}

		if field.Type.IsList {
			argExpr = &ast.SelectorExpr{
				X:   ast.NewIdent(validationTargetName),
				Sel: ast.NewIdent(fieldName),
			}
		}

		if field.Type.IsList && field.Type.Nullable {
			argExpr = &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent(validationTargetName),
					Sel: ast.NewIdent(fieldName),
				},
			}
		}

		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: ast.NewIdent("string"),
					Args: []ast.Expr{
						ast.NewIdent("sel.Name"),
					},
				},
				Op: token.EQL,
				Y: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s\"", string(field.Name)),
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					generateApplySkipDirective(),
					generateApplyIncludeDirective(),
					&ast.AssignStmt{
						Tok: token.ASSIGN,
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent(responseTargetName),
								Sel: ast.NewIdent(fieldName),
							},
						},
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("w"),
									Sel: ast.NewIdent("walk" + string(field.Name)),
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("sel"),
										Sel: ast.NewIdent("GetSelections()"),
									},
									argExpr,
								},
							},
						},
					},
				},
			},
		})
	}

	return stmts
}

func generateWrapResponseWriterResponseFieldWalkerStmts(field *schema.FieldDefinition, typeDefinition *schema.TypeDefinition) ast.Stmt {
	stmts := make([]ast.Stmt, 0)
	stmts = append(stmts, &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{
			ast.NewIdent("sel"),
		},
		Rhs: []ast.Expr{
			&ast.TypeAssertExpr{
				X: ast.NewIdent("selection"),
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("query"),
						Sel: ast.NewIdent("Field"),
					},
				},
			},
		},
	})

	stmts = append(stmts, generateWrapResponseWriterReponseFieldWalkerValidationStmts(field.Type, typeDefinition, 0)...)

	return &ast.RangeStmt{
		Key:   ast.NewIdent("_"),
		Tok:   token.DEFINE,
		Value: ast.NewIdent("selection"),
		X:     ast.NewIdent("selections"),
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}
}

func getFieldSliceNestLevel(fieldType *schema.FieldType) int {
	nestLevel := 0

	for fieldType.IsList {
		fieldType = fieldType.ListType
		nestLevel++
	}

	return nestLevel
}

func generateWrapResponseWriterResponseInitializeStmt(responseStructName string, field *schema.FieldDefinition) ast.Stmt {
	if field.Type.IsList {
		return &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{
				ast.NewIdent("resp"),
			},
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("make"),
					Args: []ast.Expr{
						generateWrapResponseWriterResponseFieldArgType(responseStructName, field.Type),
						&ast.CallExpr{
							Fun: ast.NewIdent("len"),
							Args: []ast.Expr{
								ast.NewIdent("baseResp"),
							},
						},
					},
				},
			},
		}
	}

	return &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: []ast.Expr{
			ast.NewIdent("resp"),
		},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: ast.NewIdent("new"),
				Args: []ast.Expr{
					ast.NewIdent(responseStructName),
				},
			},
		},
	}
}

func generateWrapResponseWriterResponseFieldWalkerBody(responseStructName string, field *schema.FieldDefinition, typeDefinition *schema.TypeDefinition) *ast.BlockStmt {
	stmts := []ast.Stmt{
		generateWrapResponseWriterResponseInitializeStmt(responseStructName, field),
		&ast.ExprStmt{X: &ast.BasicLit{}},
	}
	stmts = append(stmts, generateWrapResponseWriterNestedTypeInitializer(responseStructName, field)...)
	stmts = append(stmts, generateWrapResponseWriterResponseFieldWalkerStmts(field, typeDefinition))
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("resp"),
		},
	})

	return &ast.BlockStmt{
		List: stmts,
	}
}

func generateWrapResponseWriterResponseFieldArgType(structName string, fieldType *schema.FieldType) ast.Expr {
	if fieldType.IsList {
		return &ast.ArrayType{
			Elt: generateWrapResponseWriterResponseFieldArgType(structName, fieldType.ListType),
		}
	}

	if fieldType.Nullable {
		return &ast.StarExpr{
			X: &ast.Ident{
				Name: structName,
			},
		}
	}

	return &ast.Ident{
		Name: structName,
	}
}

func generateWrapResponseWriterResponseFieldWalker(operationName string, field *schema.FieldDefinition, typeDefinition *schema.TypeDefinition) ast.Decl {
	methodSufix := string(field.Name)

	responseStructName := operationName + string(typeDefinition.Name) + "Response"

	var resultType ast.Expr = &ast.StarExpr{
		X: ast.NewIdent(responseStructName),
	}

	if field.Type.IsList {
		resultType = generateWrapResponseWriterResponseFieldArgType(responseStructName, field.Type)
	}

	return &ast.FuncDecl{
		Name: ast.NewIdent("walk" + methodSufix),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("w")},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: "Wrap" + string(operationName) + "ResponseWriter",
						},
					},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("selections")},
						Type: &ast.ArrayType{
							Elt: &ast.SelectorExpr{
								X:   ast.NewIdent("query"),
								Sel: ast.NewIdent("Selection"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("baseResp")},
						Type:  generateTypeExprFromFieldType(field.Type),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: resultType,
					},
				},
			},
		},
		Body: generateWrapResponseWriterResponseFieldWalkerBody(responseStructName, field, typeDefinition),
	}
}

func generateResponseStructDeclsForWrapResponseWriter(rootFieldName string, field *schema.FieldDefinition, typeDefinition *schema.TypeDefinition, index map[string]*schema.TypeDefinition) []ast.Decl {
	fields := make([]*ast.Field, 0, len(typeDefinition.Fields))
	operationName := rootFieldName
	structPrefix := operationName + string(typeDefinition.Name)

	ret := make([]ast.Decl, 0)
	for _, field := range typeDefinition.Fields {
		graphqlType := GraphQLType(field.Type.Name)

		var typeExpr ast.Expr
		typeExpr = &ast.StarExpr{
			X: ast.NewIdent(string(graphqlType.golangType())),
		}

		targetField := field
		if !graphqlType.IsPrimitive() && field.Type.IsList {
			ft := field.Type
			for ft.IsList {
				ft = ft.ListType
			}
			td := index[string(ft.Name)]

			typeExpr = generateWrapResponseWriterResponseFieldArgType(rootFieldName+string(td.Name)+"Response", field.Type)

			graphqlType = GraphQLType(td.TypeName())
			if graphqlType == "" {
				panic("unknown type")
			}

			ret = append(ret, generateResponseStructDeclsForWrapResponseWriter(rootFieldName, field, td, index)...)
		}

		if !graphqlType.IsPrimitive() && !field.Type.IsList {
			graphqlType = GraphQLType(field.Type.Name)
			if isLowerCase(string(graphqlType.golangType())) {
				graphqlType = GraphQLType(toUpperCase(string(graphqlType.golangType())))
			}

			typeExpr = &ast.StarExpr{
				X: ast.NewIdent(rootFieldName + graphqlType.golangType() + "Response"),
			}
		}

		fieldName := FieldName(targetField.Name)
		tag := fmt.Sprintf("`json:\"%s,omitempty\"`", string(targetField.Name))
		// TODO: implment validation
		// if !targetField.Type.Nullable {
		// 	tag = fmt.Sprintf("`json:\"%s\"`", string(targetField.Name))
		// }
		fields = append(fields, &ast.Field{
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: tag,
			},
			Names: []*ast.Ident{ast.NewIdent(fieldName.ExportedGolangFieldName())},
			Type:  typeExpr,
		})
	}

	structDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(structPrefix + "Response"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: fields,
					},
				},
			},
		},
	}

	typeFieldName := FieldName(typeDefinition.TypeName())
	fieldName := FieldName(string(rootFieldName) + string(typeFieldName))

	typeExpr := generateWrapResponseWriterResponseTypeExprBasedTypeName(rootFieldName+typeFieldName.ExportedGolangFieldName()+"Response", field.Type)
	if !field.Type.IsList && !field.IsPrimitive() {
		typeExpr = &ast.StarExpr{
			X: typeExpr,
		}
	}

	wrapStructDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("wrap" + fieldName.ExportedGolangFieldName() + "Response"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent(fieldName.ExportedGolangFieldName() + "Response")},
								Tag: &ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf("`json:\"%s\"`", string(rootFieldName)),
								},
								Type: typeExpr,
							},
						},
					},
				},
			},
		},
	}

	ret = append(ret, structDecl)
	ret = append(ret, wrapStructDecl)
	ret = append(ret, generateWrapResponseWriterResponseFieldWalker(rootFieldName, field, typeDefinition))

	for _, field := range typeDefinition.Fields {
		if field.IsPrimitive() {
			continue
		}

		typeDefinition := index[string(field.Type.Name)]
		if typeDefinition == nil {
			continue
		}

		ret = append(ret, generateResponseStructDeclsForWrapResponseWriter(rootFieldName, field, typeDefinition, index)...)
	}

	return ret
}

func generateWrapResponseWriterResponseTypeExprBasedTypeName(basedTypeName string, fieldType *schema.FieldType) ast.Expr {
	graphqlType := GraphQLType(fieldType.Name)

	if fieldType.IsList {
		return &ast.ArrayType{
			Elt: generateWrapResponseWriterResponseTypeExprBasedTypeName(basedTypeName, fieldType.ListType),
		}
	}

	if graphqlType == "" {
		return &ast.Ident{
			Name: "interface{}",
		}
	}

	return &ast.Ident{
		Name: basedTypeName,
	}
}

func extractRootGraphQLType(fieldType *schema.FieldType) GraphQLType {
	if fieldType.IsList {
		return extractRootGraphQLType(fieldType.ListType)
	}

	return GraphQLType(fieldType.Name)
}

func generateWrapResponseWriterWrite(rootFieldName string, field *schema.FieldDefinition) *ast.FuncDecl {
	graphqlType := GraphQLType(field.Type.Name)

	if field.Type.IsList {
		ft := field.Type
		for ft.IsList {
			ft = ft.ListType
		}
		graphqlType = GraphQLType(ft.Name)
	}

	if graphqlType == "" {
		return &ast.FuncDecl{
			Name: ast.NewIdent("Write"),
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("w")},
						Type: &ast.StarExpr{
							X: &ast.Ident{
								Name: "Wrap" + string(field.Name) + "ResponseWriter",
							},
						},
					},
				},
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("b")},
							Type: &ast.ArrayType{
								Elt: &ast.Ident{
									Name: "byte",
								},
								Len: nil,
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.Ident{
								Name: "int",
							},
						},
						{
							Type: &ast.Ident{
								Name: "error",
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("0"),
							ast.NewIdent("nil"),
						},
					},
				},
			},
		}
	}

	var argExpr ast.Expr = &ast.SelectorExpr{
		X:   ast.NewIdent("resp"),
		Sel: ast.NewIdent("Data"),
	}

	return &ast.FuncDecl{
		Name: ast.NewIdent("Write"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("w")},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: "Wrap" + string(field.Name) + "ResponseWriter",
						},
					},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("b")},
						Type: &ast.ArrayType{
							Elt: &ast.Ident{
								Name: "byte",
							},
							Len: nil,
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{
							Name: "int",
						},
					},
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "// *********** AUTO GENERATED CODE ***********",
				},
				{
					Text: "// *********** DON'T EDIT ***********",
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("resp"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent(fmt.Sprintf("%sGraphQLResponse", rootFieldName)),
								Elts: []ast.Expr{},
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.BasicLit{},
				},
				&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("json"),
								Sel: ast.NewIdent("Unmarshal(b, &resp)"),
							},
						},
					},
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("w"),
												Sel: ast.NewIdent("ResponseWriter"),
											},
											Sel: ast.NewIdent("Write"),
										},
										Args: []ast.Expr{
											ast.NewIdent(`[]byte(fmt.Sprintf("failed to Unmarshal\nuse \"executor.GraphQLResponse\" for response: %s", err.Error()))`),
										},
									},
								},
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.BasicLit{},
				},

				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						ast.NewIdent("selected"),
					},
					Rhs: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent("wrap" + FieldName(rootFieldName).ExportedGolangFieldName() + extractRootGraphQLType(field.Type).golangType() + "Response"),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: ast.NewIdent(FieldName(rootFieldName).ExportedGolangFieldName() + extractRootGraphQLType(field.Type).golangType() + "Response"),
										Value: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("w"),
												Sel: ast.NewIdent("walk" + string(rootFieldName)),
											},
											Args: []ast.Expr{
												&ast.SelectorExpr{
													X:   ast.NewIdent("w"),
													Sel: ast.NewIdent("selections"),
												},
												argExpr,
											},
										},
									},
								},
							},
						},
					},
				},

				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("selectedResp"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("make"),
							Args: []ast.Expr{
								&ast.MapType{
									Key:   &ast.Ident{Name: "string"},
									Value: ast.NewIdent("any"),
								},
							},
						},
					},
				},

				&ast.AssignStmt{
					Tok: token.ASSIGN,
					Lhs: []ast.Expr{
						&ast.IndexExpr{
							X:     ast.NewIdent("selectedResp"),
							Index: ast.NewIdent("\"data\""),
						},
					},
					Rhs: []ast.Expr{
						ast.NewIdent("selected"),
					},
				},

				&ast.AssignStmt{
					Tok: token.ASSIGN,
					Lhs: []ast.Expr{
						&ast.IndexExpr{
							X:     ast.NewIdent("selectedResp"),
							Index: ast.NewIdent("\"errors\""),
						},
					},
					Rhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("resp"),
							Sel: ast.NewIdent("Errors"),
						},
					},
				},

				&ast.ExprStmt{
					X: &ast.BasicLit{},
				},

				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						ast.NewIdent("respByte"),
						ast.NewIdent("err"),
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("json"),
								Sel: ast.NewIdent("Marshal"),
							},
							Args: []ast.Expr{
								ast.NewIdent("selectedResp"),
							},
						},
					},
				},

				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("w"),
											Sel: ast.NewIdent("ResponseWriter"),
										},
										Sel: ast.NewIdent("WriteHeader"),
									},
									Args: []ast.Expr{
										ast.NewIdent("http.StatusInternalServerError"),
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("w"),
												Sel: ast.NewIdent("ResponseWriter"),
											},
											Sel: ast.NewIdent("Write"),
										},
										Args: []ast.Expr{
											ast.NewIdent("[]byte(err.Error())"),
										},
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("w"),
									Sel: ast.NewIdent("ResponseWriter"),
								},
								Sel: ast.NewIdent("Write"),
							},
							Args: []ast.Expr{
								ast.NewIdent("respByte"),
							},
						},
					},
				},
			},
		},
	}
}

func generateExecutorBody(op *schema.OperationDefinition, operationType string) *ast.BlockStmt {
	body := []ast.Stmt{}

	if op == nil {
		return &ast.BlockStmt{
			List: body,
		}
	}

	var methodName, executorName string
	if operationType == "query" {
		methodName = "GetQuery"
		executorName = "queryExecutor"
	}

	if operationType == "mutation" {
		methodName = "GetMutation"
		executorName = "mutationExecutor"
	}

	if operationType == "subscription" {
		methodName = "GetSubscription"
		executorName = "subscriptionExecutor"
	}

	bodyStmt := make([]ast.Stmt, 0)
	for _, field := range op.Fields {
		caseBody := make([]ast.Stmt, 0)
		fieldName := fmt.Sprintf("\"%s\"", field.Name)
		caseBody = append(caseBody, generateBodyForArgument(methodName, string(fieldName))...)
		caseBody = append(caseBody, &ast.AssignStmt{
			Tok: token.ASSIGN,
			Lhs: []ast.Expr{
				ast.NewIdent("w"),
			},
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("new" + string(field.Name) + "Writer"),
					Args: []ast.Expr{
						ast.NewIdent("w"),
						&ast.SelectorExpr{
							X:   ast.NewIdent("node"),
							Sel: ast.NewIdent("SelectSets"),
						},
						ast.NewIdent("variables"),
					},
				},
			},
		})
		caseBody = append(caseBody,
			&ast.ExprStmt{X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("r"),
					Sel: ast.NewIdent(toUpperCase(string(field.Name))),
				},
				Args: []ast.Expr{
					ast.NewIdent("w"),
					ast.NewIdent("req"),
				},
			}},
			&ast.RangeStmt{
				Value: ast.NewIdent("child"),
				Tok:   token.DEFINE,
				Key:   ast.NewIdent("_"),
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("node"),
					Sel: ast.NewIdent("Children"),
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("child.SelectSets"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("r"),
											Sel: ast.NewIdent(executorName),
										},
										Args: []ast.Expr{
											ast.NewIdent("w"),
											ast.NewIdent("req"),
											ast.NewIdent("child"),
											ast.NewIdent("parsedQuery"),
											ast.NewIdent("variables"),
										},
									},
									},
								},
							},
						},
					},
				},
			},
		)

		bodyStmt = append(bodyStmt, &ast.CaseClause{
			List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: fieldName}},
			Body: caseBody,
		})
	}

	if operationType == "query" {
		schemaCase := &ast.CaseClause{
			List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "\"__schema\""}},
			Body: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent("__schema"),
						},
						Args: []ast.Expr{
							ast.NewIdent("w"),
							ast.NewIdent("req"),
							ast.NewIdent("node"),
							ast.NewIdent("parsedQuery"),
							ast.NewIdent("variables"),
						},
					},
				},
			},
		}
		bodyStmt = append(bodyStmt, schemaCase)
	}

	body = append(body, &ast.SwitchStmt{
		Tag: ast.NewIdent("string(node.Name)"),
		Body: &ast.BlockStmt{
			List: bodyStmt,
		},
	})

	return &ast.BlockStmt{
		List: body,
	}
}

func generateBodyForArgument(operationMethodName, fieldName string) []ast.Stmt {
	return []ast.Stmt{
		&ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{
				ast.NewIdent("args"),
			},
			Rhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("utils"),
					Sel: ast.NewIdent(fmt.Sprintf("ExtractSelectorArgs(parsedQuery.Operations.%s(), %s)", operationMethodName, fieldName)),
				},
			},
		},
		&ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{
				ast.NewIdent("body"),
				ast.NewIdent("err"),
			},
			Rhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("utils"),
					Sel: ast.NewIdent("ConvRequestBodyFromVariables(variables, args)"),
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("err"),
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.CallExpr{
						Fun: ast.NewIdent("http.Error"),
						Args: []ast.Expr{
							ast.NewIdent("w"),
							&ast.BasicLit{Kind: token.STRING, Value: "\"Unknown arguments\""},
							ast.NewIdent("http.StatusUnprocessableEntity"),
						},
					}},
					&ast.ReturnStmt{},
				},
			},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("req.Body")},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("io"),
						Sel: ast.NewIdent("NopCloser"),
					},
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("strings"),
								Sel: ast.NewIdent("NewReader"),
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("string"),
									Args: []ast.Expr{
										ast.NewIdent("body"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateServeHTTPBody(query, mutation, subscription *schema.OperationDefinition) *ast.BlockStmt {
	querySwitchCases := []ast.Stmt{}
	// req.Body = io.NopCloser(strings.NewReader(string(request.Variables)))

	if query != nil {
		querySwitchCases = append(querySwitchCases, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("r"),
					Sel: ast.NewIdent("queryExecutor"),
				},
				Args: []ast.Expr{
					ast.NewIdent("w"),
					ast.NewIdent("req"),
					ast.NewIdent("node"),
					ast.NewIdent("parsedQuery"),
					ast.NewIdent("variables"),
				},
			},
		})
	}

	mutationSwitchCases := []ast.Stmt{}
	if mutation != nil {
		mutationSwitchCases = append(mutationSwitchCases, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("r"),
					Sel: ast.NewIdent("mutationExecutor"),
				},
				Args: []ast.Expr{
					ast.NewIdent("w"),
					ast.NewIdent("req"),
					ast.NewIdent("node"),
					ast.NewIdent("parsedQuery"),
					ast.NewIdent("variables"),
				},
			},
		})
	}

	subscriptionSwitchCases := []ast.Stmt{}
	if subscription != nil {
		subscriptionSwitchCases = append(subscriptionSwitchCases, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("r"),
					Sel: ast.NewIdent("subscriptionExecutor"),
				},
				Args: []ast.Expr{
					ast.NewIdent("w"),
					ast.NewIdent("req"),
					ast.NewIdent("node"),
					ast.NewIdent("parsedQuery"),
					ast.NewIdent("variables"),
				},
			},
		})
	}

	return &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("detectOperationType")},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.FuncLit{
						Type: generateDetectOperationType().Type,
						Body: generateDetectOperationType().Body,
					},
				},
			},

			&ast.ExprStmt{X: &ast.BasicLit{}},
			&ast.DeclStmt{Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent("request")},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{Names: []*ast.Ident{ast.NewIdent("OperationName")}, Type: ast.NewIdent("string")},
									{Names: []*ast.Ident{ast.NewIdent("Query")}, Type: ast.NewIdent("string")},
									{Names: []*ast.Ident{ast.NewIdent("Variables")}, Type: ast.NewIdent("json.RawMessage")},
								},
							},
						},
					},
				},
			}},

			&ast.ExprStmt{X: &ast.BasicLit{}},
			&ast.IfStmt{
				Init: &ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.BasicLit{
							Value: "err",
							Kind:  token.ASSIGN,
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("json"),
							Sel: ast.NewIdent("NewDecoder(req.Body).Decode(&request)"),
						},
					},
				},
				Cond: &ast.BinaryExpr{
					X:  ast.NewIdent("err"),
					Op: token.NEQ,
					Y:  ast.NewIdent("nil"),
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{X: &ast.CallExpr{
							Fun: ast.NewIdent("http.Error"),
							Args: []ast.Expr{
								ast.NewIdent("w"),
								&ast.BasicLit{Kind: token.STRING, Value: "\"Invalid JSON\""},
								ast.NewIdent("http.StatusUnprocessableEntity"),
							},
						}},
						&ast.ReturnStmt{},
					},
				},
			},

			&ast.ExprStmt{X: &ast.BasicLit{}},

			&ast.ExprStmt{
				X: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `// replacing req.Body is in order to use variables instinctly in each resolvers from model package`,
				},
			},

			&ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("operationType")},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("detectOperationType"),
						Args: []ast.Expr{
							ast.NewIdent("request.Query"),
						},
					},
				},
			},

			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("parsedQuery"),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("r.parser"),
						Sel: ast.NewIdent("Parse([]byte(request.Query))"),
					},
				},
			},

			&ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X:  ast.NewIdent("err"),
					Op: token.NEQ,
					Y:  ast.NewIdent("nil"),
				},

				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{X: &ast.CallExpr{
							Fun: ast.NewIdent("http.Error"),
							Args: []ast.Expr{
								ast.NewIdent("w"),
								&ast.BasicLit{Kind: token.STRING, Value: "\"failed to parse query\""},
								ast.NewIdent("http.StatusInternalServerError"),
							},
						}},
						&ast.ReturnStmt{},
					},
				},
			},

			&ast.ExprStmt{X: &ast.BasicLit{}},

			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("variables"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("request"),
						Sel: ast.NewIdent("Variables"),
					},
				},
			},

			&ast.SwitchStmt{
				Tag: ast.NewIdent("operationType"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.CaseClause{
							List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "\"query\""}},
							Body: append([]ast.Stmt{
								&ast.AssignStmt{
									Tok: token.DEFINE,
									Lhs: []ast.Expr{
										ast.NewIdent("rootSelectionSet"),
									},
									Rhs: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("utils"),
											Sel: ast.NewIdent("ExtractExecuteSelector(parsedQuery.Operations.GetQuery(), request.OperationName)"),
										},
									},
								},

								&ast.AssignStmt{
									Tok: token.DEFINE,
									Lhs: []ast.Expr{
										ast.NewIdent("node"),
									},
									Rhs: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("executor"),
											Sel: ast.NewIdent("PlanExecution(rootSelectionSet)"),
										},
									},
								},
							}, querySwitchCases...),
						},

						&ast.CaseClause{
							List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "\"mutation\""}},
							Body: append([]ast.Stmt{
								&ast.AssignStmt{
									Tok: token.DEFINE,
									Lhs: []ast.Expr{
										ast.NewIdent("rootSelectionSet"),
									},
									Rhs: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("utils"),
											Sel: ast.NewIdent("ExtractExecuteSelector(parsedQuery.Operations.GetMutation(), request.OperationName)"),
										},
									},
								},

								&ast.AssignStmt{
									Tok: token.DEFINE,
									Lhs: []ast.Expr{
										ast.NewIdent("node"),
									},
									Rhs: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("executor"),
											Sel: ast.NewIdent("PlanExecution(rootSelectionSet)"),
										},
									},
								},
							}, mutationSwitchCases...),
						},

						&ast.CaseClause{
							List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "\"subscription\""}},
							Body: []ast.Stmt{
								&ast.AssignStmt{
									Tok: token.DEFINE,
									Lhs: []ast.Expr{
										ast.NewIdent("operationName"),
									},
									Rhs: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("utils"),
											Sel: ast.NewIdent("ExtractSelectorName(parsedQuery.Operations.GetSubscription(), request.OperationName)"),
										},
									},
								},
								&ast.SwitchStmt{
									Tag: ast.NewIdent("operationName"),
									Body: &ast.BlockStmt{
										List: subscriptionSwitchCases,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateServeHTTPArgs() *ast.FieldList {
	return &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{
					{
						Name: "w",
					},
				},
				Type: &ast.ParenExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("http"),
						Sel: ast.NewIdent("ResponseWriter"),
					},
				},
			},
			{
				Names: []*ast.Ident{
					{
						Name: "req",
					},
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("http"),
						Sel: ast.NewIdent("Request"),
					},
				},
			},
		},
	}
}

func generateOperationExecutorArgs() *ast.FieldList {
	ExecutorArgs := generateServeHTTPArgs()

	additionalFields := []*ast.Field{
		{
			Names: []*ast.Ident{
				{
					Name: "node",
				},
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("executor"),
					Sel: ast.NewIdent("Node"),
				},
			},
		},
		{
			Names: []*ast.Ident{
				{
					Name: "parsedQuery",
				},
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("query"),
					Sel: ast.NewIdent("Document"),
				},
			},
		},
		{
			Names: []*ast.Ident{
				{
					Name: "variables",
				},
			},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("json"),
				Sel: ast.NewIdent("RawMessage"),
			},
		},
	}

	ExecutorArgs.List = append(ExecutorArgs.List, additionalFields...)
	return ExecutorArgs
}

func generateDetectOperationType() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("detectOperationType"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("query")},
						Type:  ast.NewIdent("string"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// query = strings.TrimSpace(query)
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("query")},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("strings"),
								Sel: ast.NewIdent("TrimSpace"),
							},
							Args: []ast.Expr{ast.NewIdent("query")},
						},
					},
				},
				// if strings.HasPrefix(query, "query") { return "query" }
				generatePrefixCheck("query"),
				// if strings.HasPrefix(query, "mutation") { return "mutation" }
				generatePrefixCheck("mutation"),
				// if strings.HasPrefix(query, "subscription") { return "subscription" }
				generatePrefixCheck("subscription"),
				// return ""
				&ast.ReturnStmt{
					Results: []ast.Expr{&ast.BasicLit{
						Kind:  token.STRING,
						Value: `""`,
					}},
				},
			},
		},
	}
}

func generatePrefixCheck(operation string) *ast.IfStmt {
	return &ast.IfStmt{
		Cond: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("strings"),
				Sel: ast.NewIdent("HasPrefix"),
			},
			Args: []ast.Expr{
				ast.NewIdent("query"),
				&ast.BasicLit{Kind: token.STRING, Value: `"` + operation + `"`},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BasicLit{Kind: token.STRING, Value: `"` + operation + `"`},
					},
				},
			},
		},
	}
}

func generateInterfaceField(operation *schema.OperationDefinition) *ast.GenDecl {
	generateField := func(field schema.FieldDefinitions) *ast.FieldList {
		fields := make([]*ast.Field, 0, len(field))

		for _, f := range field {
			fields = append(fields, &ast.Field{
				Names: []*ast.Ident{
					{
						Name: toUpperCase(string(f.Name)),
					},
				},
				Type: &ast.FuncType{
					Params:  generateServeHTTPArgs(),
					Results: &ast.FieldList{},
				},
			})
		}

		return &ast.FieldList{
			List: fields,
		}
	}

	var ident *ast.Ident
	if operation.OperationType.IsQuery() {
		ident = newQueryIdent(operation)
	}

	if operation.OperationType.IsMutation() {
		ident = newMutationIdent(operation)
	}

	if operation.OperationType.IsSubscription() {
		ident = newSubscriptionIdent(operation)
	}

	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ident,
				Type: &ast.InterfaceType{
					Methods: generateField(operation.Fields),
				},
			},
		},
	}
}

func isUsedDefinedType(operation *schema.OperationDefinition) bool {
	if operation != nil {
		for _, field := range operation.Fields {
			if !GraphQLType(field.Type.Name).IsPrimitive() {
				return true
			}

			for _, arg := range field.Arguments {
				if !GraphQLType(arg.Type.Name).IsPrimitive() {
					return true
				}
			}
		}
	}

	return false
}

func generateResolverImplementationStruct() []ast.Decl {
	return []ast.Decl{
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{
						Name: "resolver",
					},
					Type: &ast.StructType{
						Fields: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										ast.NewIdent("parser"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("query"),
											Sel: ast.NewIdent("Parser"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		&ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{
						{
							Name: "_",
						},
					},
					Type: ast.NewIdent("Resolver"),
					Values: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.Ident{
								Name: "resolver{}",
							},
						},
					},
				},
			},
		},
		&ast.FuncDecl{
			Name: ast.NewIdent("NewResolver"),
			Type: &ast.FuncType{
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.StarExpr{
								X: ast.NewIdent("resolver"),
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CompositeLit{
								Type: ast.NewIdent("&resolver"),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: ast.NewIdent("parser"),
										Value: &ast.SelectorExpr{
											Sel: ast.NewIdent("NewParserWithLexer()"),
											X:   ast.NewIdent("query"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateResolverImplementation(fields schema.FieldDefinitions) []ast.Decl {
	decls := make([]ast.Decl, 0, len(fields))

	recv := func(t *schema.FieldType) string {
		if t.IsList {
			return fmt.Sprintf("[]%s", t.ListType.Name)
		}

		graphQLType := GraphQLType(t.Name)
		return graphQLType.golangType()
	}

	for _, f := range fields {
		returnsStr := recv(f.Type)

		if f.Type != nil {
			decls = append(decls, generateResponseGraphQLResponseStruct(string(f.Name), f))
		}

		decls = append(decls, &ast.FuncDecl{
			Doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{
						Text: fmt.Sprintf("// Read request body for %sArgs", toUpperCase(string(f.Name))),
					},
					{
						Text: fmt.Sprintf("// Write response body for %s", returnsStr),
					},
				},
			},
			Name: ast.NewIdent(toUpperCase(string(f.Name))),
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "r",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.Ident{
								Name: "resolver",
							},
						},
					},
				},
			},
			Type: &ast.FuncType{
				Params:  generateServeHTTPArgs(),
				Results: &ast.FieldList{},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{},
			},
		})
	}

	return decls
}

var fieldsIntrospectionFieldDefinition = &schema.FieldDefinition{
	Name: []byte("__fields"),
	Arguments: schema.ArgumentDefinitions{
		{
			Name: []byte("includeDeprecated"),
			Type: &schema.FieldType{
				Name:     []byte("Boolean"),
				IsList:   false,
				Nullable: true,
				ListType: nil,
			},
			Default: []byte("false"),
		},
	},
	Type: &schema.FieldType{
		Name:     nil,
		IsList:   true,
		Nullable: true,
		ListType: &schema.FieldType{
			Name:     []byte("__Field"),
			IsList:   false,
			Nullable: false,
			ListType: nil,
		},
	},
}

func generateExtractOperationArgumentsDecl(field *schema.FieldDefinition, indexes *schema.Indexes) ast.Decl {

	bodyStmts := make([]ast.Stmt, 0)
	bodyStmts = append(bodyStmts, generateDeclareStmts(field.Arguments)...)

	errReturnStmt := generateArgumentReturnStmt(field.Arguments)
	errReturnStmt.Results = append(errReturnStmt.Results, ast.NewIdent("err"))
	bodyStmts = append(bodyStmts, generateDefaultValueAssignmentStmts(field.Arguments, indexes, errReturnStmt)...)

	okReturnStmt := generateArgumentReturnStmt(field.Arguments)
	okReturnStmt.Results = append(okReturnStmt.Results, ast.NewIdent("nil"))
	bodyStmts = append(bodyStmts, okReturnStmt)

	results := generateArgumentsParams(field.Arguments)
	results.List = append(results.List, &ast.Field{
		Type: ast.NewIdent("error"),
	})

	return &ast.FuncDecl{
		Name: ast.NewIdent("extract" + string(field.Name) + "Args"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "r",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: "resolver",
						},
					},
				},
			},
		},
		Type: &ast.FuncType{
			Params:  generateNodeWalkerArgs(),
			Results: results,
		},
		Body: &ast.BlockStmt{
			List: bodyStmts,
		},
	}
}

func generateArgumentsParams(args schema.ArgumentDefinitions) *ast.FieldList {
	fields := make([]*ast.Field, 0, len(args))

	for _, arg := range args {
		fields = append(fields, &ast.Field{
			Type: generateTypeExprFromFieldType(arg.Type),
		})
	}

	return &ast.FieldList{
		List: fields,
	}
}

func generateDeclareStmts(args schema.ArgumentDefinitions) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(args))

	stmts = append(stmts, &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok:   token.VAR,
			Specs: generateVarSpecs(args),
		},
	})

	return stmts
}

func generateDefaultValueAssignmentStmts(args schema.ArgumentDefinitions, indexes *schema.Indexes, errReturnStmt ast.Stmt) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(args))

	for i, arg := range args {
		if arg.Default != nil {
			stmts = append(stmts, generateAssignDefaultValueStmt(arg, indexes))
			stmts = append(stmts, generateResolveDefaultValueStmt(arg, i, errReturnStmt)...)
		}
	}

	return stmts
}

func generateResolveDefaultValueStmt(arg *schema.ArgumentDefinition, index int, errReturnStmt ast.Stmt) []ast.Stmt {
	argExpr := &ast.SelectorExpr{
		X: &ast.IndexExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("node"),
				Sel: ast.NewIdent("Arguments"),
			},
			Index: ast.NewIdent(fmt.Sprintf("%d", index)),
		},
		Sel: ast.NewIdent("Value"),
	}

	ifStmt := &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  argExpr,
			Op: token.NEQ,
			Y:  ast.NewIdent("nil"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						ast.NewIdent("ast"),
						ast.NewIdent("err"),
					},
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:  &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: ast.NewIdent("r"),
										Sel: ast.NewIdent("parser"),
									},
									Sel: ast.NewIdent("ValueParser"),
								},
								Sel: ast.NewIdent("Parse"),
							},
							Args: []ast.Expr{
								argExpr,
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							errReturnStmt,
						},
					},
				},
			},
		},
	}

	ret := make([]ast.Stmt, 0)
	ret = append(ret, ifStmt)
	
	return ret
}

func generateAssignDefaultValueStmt(arg *schema.ArgumentDefinition, indexes *schema.Indexes) ast.Stmt {
	return &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			ast.NewIdent(string(arg.Name)),
		},
		Rhs: []ast.Expr{
			generateDefaultValueExpr(arg, indexes),
		},
	}
}

func generateDefaultValueExpr(arg *schema.ArgumentDefinition, indexes *schema.Indexes) ast.Expr {
	if arg.Type.IsBoolean() {
		if arg.Type.Nullable {
			return generateBoolPointerAST(string(arg.Default))
		} else {
			return &ast.BasicLit{
				Kind:  token.STRING,
				Value: string(arg.Default),
			}
		}
	}

	if arg.Type.IsString() {
		if arg.Type.Nullable {
			return generateStringPointerAST(string(arg.Default))
		} else {
			return &ast.BasicLit{
				Kind:  token.STRING,
				Value: string(arg.Default),
			}
		}
	}

	if arg.Type.IsInt() {
		if arg.Type.Nullable {
			return generateIntPointerAST(string(arg.Default))
		} else {
			return &ast.BasicLit{
				Kind:  token.STRING,
				Value: string(arg.Default),
			}
		}
	}

	if arg.Type.IsFloat() {
		if arg.Type.Nullable {
			return generateFloatPointerAST(string(arg.Default))
		} else {
			return &ast.BasicLit{
				Kind:  token.STRING,
				Value: string(arg.Default),
			}
		}
	}

	// TODO: implement other types like ID, Enum, etc.
	return nil
}

func generateIntPointerAST(value string) ast.Expr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.IndexExpr{
			X: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: ast.NewIdent("int"),
				},
				Elts: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s"`, value),
					},
				},
			},
			Index: &ast.BasicLit{
				Kind:  token.INT,
				Value: "0",
			},
		},
	}
}

func generateFloatPointerAST(value string) ast.Expr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.IndexExpr{
			X: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: ast.NewIdent("float64"),
				},
				Elts: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.FLOAT,
						Value: value,
					},
				},
			},
			Index: &ast.BasicLit{
				Kind:  token.INT,
				Value: "0",
			},
		},
	}
}

func generateVarSpecs(args schema.ArgumentDefinitions) []ast.Spec {
	specs := make([]ast.Spec, 0)
	for _, arg := range args {
		specs = append(specs, &ast.ValueSpec{
			Names: []*ast.Ident{
				ast.NewIdent(string(arg.Name)),
			},
			Type: generateTypeExprFromFieldType(arg.Type),
		})
	}

	return specs
}

func generateArgumentReturnStmt(args schema.ArgumentDefinitions) *ast.ReturnStmt {
	results := make([]ast.Expr, 0, len(args))

	for _, arg := range args {
		results = append(results, ast.NewIdent(string(arg.Name)))
	}

	return &ast.ReturnStmt{
		Results: results,
	}
}
