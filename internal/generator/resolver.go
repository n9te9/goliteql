package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/internal/generator/introspection"
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
					Value: `"strconv"`,
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
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"sync"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"time"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"context"`,
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

func generateTypeExprFromFieldType(typePrefix string, fieldType *schema.FieldType) ast.Expr {
	if fieldType.IsList {
		return &ast.ArrayType{
			Elt: generateTypeExprFromFieldType(typePrefix, fieldType.ListType),
		}
	}

	graphQLType := GraphQLType(fieldType.Name)

	var baseTypeExpr ast.Expr = ast.NewIdent(graphQLType.golangType())
	if !graphQLType.IsPrimitive() {
		if typePrefix != "" {
			baseTypeExpr = &ast.SelectorExpr{
				X:   ast.NewIdent(typePrefix),
				Sel: ast.NewIdent(graphQLType.golangType()),
			}
		} else {
			baseTypeExpr = ast.NewIdent(graphQLType.golangType())
		}
	}

	if fieldType.Nullable {
		return &ast.StarExpr{
			X: baseTypeExpr,
		}
	}

	return baseTypeExpr
}

func generateApplyQueryResponseCaseStmts(typeDefinition *schema.TypeDefinition, nestCount int) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if typeDefinition == nil {
		return ret
	}

	for _, field := range typeDefinition.Fields {
		var assignExpr ast.Expr = ast.NewIdent(fmt.Sprintf("ret%d", nestCount-1))
		if nestCount-1 < 0 {
			assignExpr = ast.NewIdent("ret")
		}

		lhs := []ast.Expr{
			&ast.SelectorExpr{
				X:   assignExpr,
				Sel: ast.NewIdent(toUpperCase(string(field.Name))),
			},
		}

		valueName := "resolverRet"
		if nestCount > 0 {
			valueName = fmt.Sprintf("v%d", nestCount-1)
		}

		var rh ast.Expr = &ast.SelectorExpr{
			X:   ast.NewIdent(valueName),
			Sel: ast.NewIdent(toUpperCase(string(field.Name))),
		}

		if !field.Type.Nullable {
			if field.Type.IsBoolean() {
				rh = generateIntPointerExpr(rh)
			} else if field.Type.IsString() || field.Type.IsID() {
				rh = generateStringPointerExpr(rh)
			} else if field.Type.IsInt() {
				rh = generateIntPointerExpr(rh)
			} else if field.Type.IsFloat() {
				rh = generateFloatPointerExpr(rh)
			} else {
				rh = generateObjectPointerExpr(string(field.Type.Name), rh)
			}
		}

		ret = append(ret, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, string(field.Name)),
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: lhs,
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						rh,
					},
				},
			},
		})
	}

	return ret
}

func generateApplyLoopStmtForQueryResponse(fieldDefinition *schema.FieldDefinition, fieldType *schema.FieldType, indexes *schema.Indexes, nestCount int) ast.Stmt {
	if fieldType.IsList {

		rootType := fieldType.GetRootType()
		copiedRootType := &schema.FieldType{
			Name:     []byte(fmt.Sprintf("%sResponse", rootType.Name)),
			IsList:   false,
			ListType: nil,
			Nullable: rootType.Nullable,
		}
		return &ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Value: ast.NewIdent(fmt.Sprintf("v%d", nestCount)),
			X:     ast.NewIdent("resolverRet"),
			Tok:   token.DEFINE,
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("ret%d", nestCount)),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent("new"),
								Args: []ast.Expr{
									generateTypeExprFromFieldType("", copiedRootType),
								},
							},
						},
					},
					generateApplyLoopStmtForQueryResponse(fieldDefinition, fieldType.ListType, indexes, nestCount+1),
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("ret"),
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent("append"),
								Args: []ast.Expr{
									ast.NewIdent("ret"),
									ast.NewIdent(fmt.Sprintf("ret%d", nestCount)),
								},
							},
						},
					},
				},
			},
		}
	}

	return generateApplySwitchStmtForQueryResponse(fieldDefinition, indexes, nestCount)
}

func generateApplySwitchStmtForQueryResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount int) ast.Stmt {
	rootType := field.Type.GetRootType()
	typeDefinition := indexes.TypeIndex[string(rootType.Name)]

	return &ast.RangeStmt{
		Key:   ast.NewIdent("_"),
		Value: ast.NewIdent("child"),
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("node"),
			Sel: ast.NewIdent("Children"),
		},
		Tok: token.DEFINE,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.SwitchStmt{
					Tag: &ast.CallExpr{
						Fun: ast.NewIdent("string"),
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("child"),
								Sel: ast.NewIdent("Name"),
							},
						},
					},
					Body: &ast.BlockStmt{
						List: generateApplyQueryResponseCaseStmts(typeDefinition, nestCount),
					},
				},
			},
		},
	}
}

func generateIndexExprFromNestCount(nestCount int) ast.Expr {
	if nestCount < 0 {
		return ast.NewIdent("ret")
	}

	return &ast.IndexExpr{
		X:     ast.NewIdent("ret"),
		Index: ast.NewIdent(fmt.Sprintf("k%d", nestCount)),
	}
}

func generateApplyQueryResponseFuncDecls(operationDefinition *schema.OperationDefinition, indexes *schema.Indexes, nestCount int, typePrefix string) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, field := range operationDefinition.Fields {
		ret = append(ret, &ast.FuncDecl{
			Name: ast.NewIdent("apply" + string(field.Name) + "QueryResponse"),
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("r")},
						Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					generateRetAssignStmt("", generateFieldTypeForResponse(field.Type)),
					generateApplyLoopStmtForQueryResponse(field, field.Type, indexes, nestCount),
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("ret"),
							ast.NewIdent("nil"),
						},
					},
				},
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("resolverRet"),
							},
							Type: generateTypeExprFromFieldType(typePrefix, field.Type),
						},
						{
							Names: []*ast.Ident{
								ast.NewIdent("node"),
							},
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("executor"),
									Sel: ast.NewIdent("Node"),
								},
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: generateTypeExprFromFieldType("", generateFieldTypeForResponse(field.Type)),
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
		})
	}

	return ret
}

func generateRetAssignStmt(typePrefix string, fieldType *schema.FieldType) ast.Stmt {
	if fieldType.IsList {
		return &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("ret"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("make"),
					Args: []ast.Expr{
						generateTypeExprFromFieldType(typePrefix, fieldType),
						ast.NewIdent("0"),
						&ast.CallExpr{
							Fun: ast.NewIdent("len"),
							Args: []ast.Expr{
								ast.NewIdent("resolverRet"),
							},
						},
					},
				},
			},
		}
	}

	if fieldType.IsPrimitive() {
		if fieldType.Nullable {
			return &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("ret"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("new"),
						Args: []ast.Expr{
							ast.NewIdent(GraphQLType(fieldType.Name).golangType()),
						},
					},
				},
			}
		}

		return &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ret"),
						},
						Type: ast.NewIdent(GraphQLType(fieldType.Name).golangType()),
					},
				},
			},
		}
	}

	rhs := []ast.Expr{
		&ast.CompositeLit{
			Type: ast.NewIdent(string(fieldType.Name)),
			Elts: []ast.Expr{},
		},
	}

	if fieldType.Nullable {
		rhs = []ast.Expr{
			&ast.UnaryExpr{
				Op: token.AND,
				X: &ast.CompositeLit{
					Type: ast.NewIdent(string(fieldType.Name)),
					Elts: []ast.Expr{},
				},
			},
		}
	}

	return &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("ret"),
		},
		Tok: token.DEFINE,
		Rhs: rhs,
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
			Params: generateNodeWalkerArgs(),
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("any"),
					},
					{
						Type: ast.NewIdent("error"),
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
			Params: generateNodeWalkerArgs(),
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("any"),
					},
					{
						Type: ast.NewIdent("error"),
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
			Params:  generateNodeWalkerArgs(),
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

func generateArgumentsAssignStmt(fieldName string, args schema.ArgumentDefinitions) ast.Stmt {
	lhs := make([]ast.Expr, 0, len(args)+1)
	for _, arg := range args {
		lhs = append(lhs, ast.NewIdent(string(arg.Name)))
	}
	lhs = append(lhs, ast.NewIdent("err"))

	rhs := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent("r"),
			Sel: ast.NewIdent(fmt.Sprintf("extract%sArgs", fieldName)),
		},
		Args: []ast.Expr{
			ast.NewIdent("node"),
			ast.NewIdent("variables"),
		},
	}

	return &ast.AssignStmt{
		Tok: token.DEFINE,
		Lhs: lhs,
		Rhs: []ast.Expr{rhs},
	}
}

func generateExecutorBody(op *schema.OperationDefinition, operationType string) *ast.BlockStmt {
	body := []ast.Stmt{}

	if op == nil {
		return &ast.BlockStmt{
			List: body,
		}
	}

	bodyStmt := make([]ast.Stmt, 0)
	for _, field := range op.Fields {
		caseBody := make([]ast.Stmt, 0)
		fieldName := fmt.Sprintf("\"%s\"", field.Name)
		if len(field.Arguments) > 0 {
			caseBody = append(caseBody,
				generateArgumentsAssignStmt(string(field.Name), field.Arguments),
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{},
					},
				})
		}
		caseBody = append(caseBody,
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("resolverRet"),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent(toUpperCase(string(field.Name))),
						},
						Args: generateFieldArguments(field.Arguments),
					},
				},
			},
			generateReturnErrorHandlingStmt([]ast.Expr{
				ast.NewIdent("nil"),
			}),
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("ret"),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent(fmt.Sprintf("apply%sQueryResponse", field.Name)),
						},
						Args: []ast.Expr{
							ast.NewIdent("resolverRet"),
							ast.NewIdent("node"),
						},
					},
				},
			},
			generateReturnErrorHandlingStmt([]ast.Expr{
				ast.NewIdent("nil"),
			}),
			&ast.ReturnStmt{
				Results: []ast.Expr{
					ast.NewIdent("ret"),
					ast.NewIdent("nil"),
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
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("ret"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent("__schema"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("node"),
								ast.NewIdent("variables"),
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
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("ret"),
						ast.NewIdent("nil"),
					},
				},
			},
		}

		typeCase := &ast.CaseClause{
			List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "\"__type\""}},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("ret"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent("__type"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("node"),
								ast.NewIdent("variables"),
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
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("ret"),
						ast.NewIdent("nil"),
					},
				},
			},
		}
		bodyStmt = append(bodyStmt, schemaCase, typeCase)
	}

	stmts := []ast.Stmt{}
	stmts = append(stmts, bodyStmt...)

	body = append(body, &ast.SwitchStmt{
		Tag: ast.NewIdent("string(node.Name)"),
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
			ast.NewIdent("nil"),
		},
	})

	return &ast.BlockStmt{
		List: body,
	}
}

func generateFieldArguments(arguments schema.ArgumentDefinitions) []ast.Expr {
	args := make([]ast.Expr, 0, len(arguments)+1)
	args = append(args, ast.NewIdent("ctx"))

	for _, arg := range arguments {
		args = append(args, ast.NewIdent(string(arg.Name)))
	}

	return args
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
		querySwitchCases = append(querySwitchCases,
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					ast.NewIdent("data"),
				},
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
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					ast.NewIdent("graphqlErrors"),
				},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("make"),
						Args: []ast.Expr{
							&ast.ArrayType{
								Elt: ast.NewIdent("error"),
							},
							ast.NewIdent("0"),
						},
					},
				},
			},
			&ast.RangeStmt{
				Key:   ast.NewIdent("_"),
				Tok:   token.DEFINE,
				Value: ast.NewIdent("node"),
				X:     ast.NewIdent("nodes"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{
								ast.NewIdent("ret"),
								ast.NewIdent("err"),
							},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("r"),
										Sel: ast.NewIdent("queryExecutor"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("req"),
												Sel: ast.NewIdent("Context"),
											},
										},
										ast.NewIdent("node"),
										ast.NewIdent("variables"),
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
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											ast.NewIdent("graphqlErrors"),
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: ast.NewIdent("append"),
												Args: []ast.Expr{
													ast.NewIdent("graphqlErrors"),
													ast.NewIdent("err"),
												},
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Tok: token.ASSIGN,
							Lhs: []ast.Expr{
								&ast.IndexExpr{
									X: ast.NewIdent("data"),
									Index: &ast.CallExpr{
										Fun: ast.NewIdent("string"),
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("node"),
												Sel: ast.NewIdent("Name"),
											},
										},
									},
								},
							},
							Rhs: []ast.Expr{
								ast.NewIdent("ret"),
							},
						},
					},
				},
			}, generateResponseWrite())
	}

	mutationSwitchCases := []ast.Stmt{}
	if mutation != nil {
		mutationSwitchCases = append(mutationSwitchCases, &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{
				ast.NewIdent("data"),
			},
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
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					ast.NewIdent("graphqlErrors"),
				},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("make"),
						Args: []ast.Expr{
							&ast.ArrayType{
								Elt: ast.NewIdent("error"),
							},
							ast.NewIdent("0"),
						},
					},
				},
			},
			&ast.RangeStmt{
				Key:   ast.NewIdent("_"),
				Tok:   token.DEFINE,
				Value: ast.NewIdent("node"),
				X:     ast.NewIdent("nodes"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{
								ast.NewIdent("ret"),
								ast.NewIdent("err"),
							},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("r"),
										Sel: ast.NewIdent("mutationExecutor"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("req"),
												Sel: ast.NewIdent("Context"),
											},
										},
										ast.NewIdent("node"),
										ast.NewIdent("variables"),
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
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											ast.NewIdent("graphqlErrors"),
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: ast.NewIdent("append"),
												Args: []ast.Expr{
													ast.NewIdent("graphqlErrors"),
													ast.NewIdent("err"),
												},
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Tok: token.ASSIGN,
							Lhs: []ast.Expr{
								&ast.IndexExpr{
									X: ast.NewIdent("data"),
									Index: &ast.CallExpr{
										Fun: ast.NewIdent("string"),
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("node"),
												Sel: ast.NewIdent("Name"),
											},
										},
									},
								},
							},
							Rhs: []ast.Expr{
								ast.NewIdent("ret"),
							},
						},
					},
				},
			}, generateResponseWrite())
	}

	subscriptionSwitchCases := []ast.Stmt{}
	if subscription != nil {
		subscriptionSwitchCases = append(subscriptionSwitchCases, &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{
				ast.NewIdent("data"),
			},
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
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					ast.NewIdent("graphqlErrors"),
				},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("make"),
						Args: []ast.Expr{
							&ast.ArrayType{
								Elt: ast.NewIdent("error"),
							},
							ast.NewIdent("0"),
						},
					},
				},
			},
			&ast.RangeStmt{
				Key:   ast.NewIdent("_"),
				Tok:   token.DEFINE,
				Value: ast.NewIdent("node"),
				X:     ast.NewIdent("nodes"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{
								ast.NewIdent("ret"),
								ast.NewIdent("err"),
							},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("r"),
										Sel: ast.NewIdent("subscriptionExecutor"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("req"),
												Sel: ast.NewIdent("Context"),
											},
										},
										ast.NewIdent("node"),
										ast.NewIdent("variables"),
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
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											ast.NewIdent("graphqlErrors"),
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: ast.NewIdent("append"),
												Args: []ast.Expr{
													ast.NewIdent("graphqlErrors"),
													ast.NewIdent("err"),
												},
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Tok: token.ASSIGN,
							Lhs: []ast.Expr{
								&ast.IndexExpr{
									X: ast.NewIdent("data"),
									Index: &ast.CallExpr{
										Fun: ast.NewIdent("string"),
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("node"),
												Sel: ast.NewIdent("Name"),
											},
										},
									},
								},
							},
							Rhs: []ast.Expr{
								ast.NewIdent("ret"),
							},
						},
					},
				},
			}, generateResponseWrite())
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
									{Names: []*ast.Ident{ast.NewIdent("Variables")}, Type: &ast.MapType{
										Key: ast.NewIdent("string"),
										Value: &ast.SelectorExpr{
											X:   ast.NewIdent("json"),
											Sel: ast.NewIdent("RawMessage"),
										},
									}},
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
									Lhs: []ast.Expr{
										ast.NewIdent("cacheMap"),
									},
									Tok: token.DEFINE,
									Rhs: []ast.Expr{
										&ast.TypeAssertExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("r"),
														Sel: ast.NewIdent("pool"),
													},
													Sel: ast.NewIdent("Get"),
												},
											},
											Type: &ast.SelectorExpr{
												X:   ast.NewIdent("executor"),
												Sel: ast.NewIdent("CacheMap"),
											},
										},
									},
								},

								&ast.DeferStmt{
									Call: &ast.CallExpr{
										Fun: &ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{
													List: []*ast.Field{},
												},
											},
											Body: &ast.BlockStmt{
												List: []ast.Stmt{
													&ast.ExprStmt{
														X: &ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.SelectorExpr{
																	X:   ast.NewIdent("r"),
																	Sel: ast.NewIdent("pool"),
																},
																Sel: ast.NewIdent("Put"),
															},
															Args: []ast.Expr{
																ast.NewIdent("cacheMap"),
															},
														},
													},
												},
											},
										},
									},
								},

								&ast.AssignStmt{
									Tok: token.DEFINE,
									Lhs: []ast.Expr{
										ast.NewIdent("nodes"),
									},
									Rhs: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("cacheMap"),
												Sel: ast.NewIdent("Get"),
											},
											Args: []ast.Expr{
												&ast.SelectorExpr{
													X:   ast.NewIdent("request"),
													Sel: ast.NewIdent("Query"),
												},
											},
										},
									},
								},

								&ast.IfStmt{
									Cond: &ast.BinaryExpr{
										X:  ast.NewIdent("nodes"),
										Op: token.EQL,
										Y:  ast.NewIdent("nil"),
									},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("nodes"),
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("PlanExecution"),
														},
														Args: []ast.Expr{
															ast.NewIdent("rootSelectionSet"),
														},
													},
												},
											},

											&ast.ExprStmt{
												X: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("cacheMap"),
														Sel: ast.NewIdent("Set"),
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X:   ast.NewIdent("request"),
															Sel: ast.NewIdent("Query"),
														},
														ast.NewIdent("nodes"),
														&ast.BinaryExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent("time"),
																Sel: ast.NewIdent("Minute"),
															},
															Op: token.MUL,
															Y: &ast.BasicLit{
																Kind:  token.INT,
																Value: "1",
															},
														},
													},
												},
											},
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
									Lhs: []ast.Expr{
										ast.NewIdent("cacheMap"),
									},
									Tok: token.DEFINE,
									Rhs: []ast.Expr{
										&ast.TypeAssertExpr{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X:   ast.NewIdent("r"),
														Sel: ast.NewIdent("pool"),
													},
													Sel: ast.NewIdent("Get"),
												},
											},
											Type: &ast.SelectorExpr{
												X:   ast.NewIdent("executor"),
												Sel: ast.NewIdent("CacheMap"),
											},
										},
									},
								},

								&ast.DeferStmt{
									Call: &ast.CallExpr{
										Fun: &ast.FuncLit{
											Type: &ast.FuncType{
												Params: &ast.FieldList{
													List: []*ast.Field{},
												},
											},
											Body: &ast.BlockStmt{
												List: []ast.Stmt{
													&ast.ExprStmt{
														X: &ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.SelectorExpr{
																	X:   ast.NewIdent("r"),
																	Sel: ast.NewIdent("pool"),
																},
																Sel: ast.NewIdent("Put"),
															},
															Args: []ast.Expr{
																ast.NewIdent("cacheMap"),
															},
														},
													},
												},
											},
										},
									},
								},

								&ast.AssignStmt{
									Tok: token.DEFINE,
									Lhs: []ast.Expr{
										ast.NewIdent("nodes"),
									},
									Rhs: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("cacheMap"),
												Sel: ast.NewIdent("Get"),
											},
											Args: []ast.Expr{
												&ast.SelectorExpr{
													X:   ast.NewIdent("request"),
													Sel: ast.NewIdent("Query"),
												},
											},
										},
									},
								},

								&ast.IfStmt{
									Cond: &ast.BinaryExpr{
										X:  ast.NewIdent("nodes"),
										Op: token.EQL,
										Y:  ast.NewIdent("nil"),
									},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("nodes"),
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("PlanExecution"),
														},
														Args: []ast.Expr{
															ast.NewIdent("rootSelectionSet"),
														},
													},
												},
											},

											&ast.ExprStmt{
												X: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("cacheMap"),
														Sel: ast.NewIdent("Set"),
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X:   ast.NewIdent("request"),
															Sel: ast.NewIdent("Query"),
														},
														ast.NewIdent("nodes"),
														&ast.BinaryExpr{
															X: &ast.SelectorExpr{
																X:   ast.NewIdent("time"),
																Sel: ast.NewIdent("Minute"),
															},
															Op: token.MUL,
															Y: &ast.BasicLit{
																Kind:  token.INT,
																Value: "1",
															},
														},
													},
												},
											},
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

func generateResponseWrite() ast.Stmt {
	return &ast.IfStmt{
		Init: &ast.AssignStmt{
			Lhs: []ast.Expr{ast.NewIdent("err")},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("json"),
								Sel: ast.NewIdent("NewEncoder"),
							},
							Args: []ast.Expr{
								ast.NewIdent("w"),
							},
						},
						Sel: ast.NewIdent("Encode"),
					},
					Args: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("executor"),
									Sel: ast.NewIdent("GraphQLResponse"),
								},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("Data"),
										Value: ast.NewIdent("data"),
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("Errors"),
										Value: ast.NewIdent("graphqlErrors"),
									},
								},
							},
						},
					},
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
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("w"),
							Sel: ast.NewIdent("WriteHeader"),
						},
						Args: []ast.Expr{
							ast.NewIdent("http.StatusInternalServerError"),
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

func generateTypeExprFromExpandedType(expandedType *introspection.FieldType) ast.Expr {
	if expandedType.Child != nil {
		if expandedType.NonNull {
			return generateTypeExprFromExpandedType(expandedType.Child)
		}

		if expandedType.IsList {
			return &ast.ArrayType{
				Elt: generateTypeExprFromExpandedType(expandedType.Child),
			}
		}

		return &ast.StarExpr{
			X: generateTypeExprFromExpandedType(expandedType.Child),
		}
	}

	if expandedType.IsPrimitive() {
		graphqlType := GraphQLType(expandedType.Name)
		return ast.NewIdent(graphqlType.golangType())
	}

	return &ast.SelectorExpr{
		X:   ast.NewIdent("model"),
		Sel: ast.NewIdent(string(expandedType.Name)),
	}
}

func generateResolverArgs(typePrefix string, field *schema.FieldDefinition) *ast.FieldList {
	ret := make([]*ast.Field, 0, len(field.Arguments)+1)
	ret = append(ret, &ast.Field{
		Names: []*ast.Ident{
			{
				Name: "ctx",
			},
		},
		Type: &ast.SelectorExpr{
			X:   ast.NewIdent("context"),
			Sel: ast.NewIdent("Context"),
		},
	})

	for _, arg := range field.Arguments {
		ret = append(ret, &ast.Field{
			Names: []*ast.Ident{
				{
					Name: string(arg.Name),
				},
			},
			Type: generateTypeExprFromFieldType(typePrefix, arg.Type),
		})
	}

	return &ast.FieldList{
		List: ret,
	}
}

func generateResolverReturns(typePrefix string, field *schema.FieldDefinition) *ast.FieldList {
	ret := make([]*ast.Field, 0, len(field.Arguments)+1)

	ret = append(ret, &ast.Field{
		Names: []*ast.Ident{},
		Type:  generateTypeExprFromFieldType(typePrefix, field.Type),
	})

	ret = append(ret, &ast.Field{
		Names: []*ast.Ident{},
		Type:  ast.NewIdent("error"),
	})

	return &ast.FieldList{
		List: ret,
	}
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

func generateInterfaceField(typePrefix string, operation *schema.OperationDefinition) *ast.GenDecl {
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
					Params:  generateResolverArgs(typePrefix, f),
					Results: generateResolverReturns(typePrefix, f),
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
								{
									Names: []*ast.Ident{
										ast.NewIdent("pool"),
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("sync"),
											Sel: ast.NewIdent("Pool"),
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
									&ast.KeyValueExpr{
										Key: ast.NewIdent("pool"),
										Value: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("executor"),
												Sel: ast.NewIdent("NewPool"),
											},
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

func generateResolverImplementation(typePrefix string, fields schema.FieldDefinitions) []ast.Decl {
	decls := make([]ast.Decl, 0, len(fields))

	for _, f := range fields {
		decls = append(decls, &ast.FuncDecl{
			Doc:  &ast.CommentGroup{},
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
				Params:  generateResolverArgs(typePrefix, f),
				Results: generateResolverReturns(typePrefix, f),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: ast.NewIdent("panic"),
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: `"` + string(f.Name) + " resolver is not implemented" + `"`,
								},
							},
						},
					},
				},
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

func generateOperationArgumentDecls(typePrefix string, operation *schema.OperationDefinition, indexes *schema.Indexes) []ast.Decl {
	decls := make([]ast.Decl, 0)

	if operation == nil {
		return decls
	}

	for _, field := range operation.Fields {
		if len(field.Arguments) == 0 {
			continue
		}

		decls = append(decls, generateExtractOperationArgumentsDecl(typePrefix, field, indexes))
	}

	return decls
}

func generateExtractOperationArgumentsDecl(typePrefix string, field *schema.FieldDefinition, indexes *schema.Indexes) ast.Decl {
	bodyStmts := make([]ast.Stmt, 0)
	bodyStmts = append(bodyStmts, generateDeclareStmts(typePrefix, field.Arguments)...)
	bodyStmts = append(bodyStmts, generateDefaultValueAssignmentStmts(field.Arguments, indexes)...)

	results := generateArgumentsParams(typePrefix, field.Arguments)
	results.List = append(results.List, &ast.Field{
		Type: ast.NewIdent("error"),
	})

	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("extract%sArgs", string(field.Name))),
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
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("node"),
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
							ast.NewIdent("variables"),
						},
						Type: &ast.MapType{
							Key: ast.NewIdent("string"),
							Value: &ast.SelectorExpr{
								X:   ast.NewIdent("json"),
								Sel: ast.NewIdent("RawMessage"),
							},
						},
					},
				},
			},
			Results: results,
		},
		Body: &ast.BlockStmt{
			List: bodyStmts,
		},
	}
}

func generateArgumentFieldExtractStmts(arg *schema.ArgumentDefinition) ast.Stmt {
	caseStmts := make([]ast.Stmt, 0)
	if arg.Type.IsBoolean() {
		caseStmts = append(caseStmts, generateBoolCaseStmt(arg))
	}

	if arg.Type.IsString() {
		caseStmts = append(caseStmts, generateStringCaseStmt(arg))
	}

	if arg.Type.IsID() {
		caseStmts = append(caseStmts, generateIDCaseStmt(arg))
	}

	// TODO handle other types like Int, Float, etc.

	switchStmt := &ast.TypeSwitchStmt{
		Assign: &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{
				ast.NewIdent("val"),
			},
			Rhs: []ast.Expr{
				&ast.TypeAssertExpr{
					X:    ast.NewIdent("ast"),
					Type: ast.NewIdent("type"),
				},
			},
		},
		Body: &ast.BlockStmt{
			List: caseStmts,
		},
	}

	return switchStmt
}

func generateBoolCaseStmt(arg *schema.ArgumentDefinition) ast.Stmt {
	var assignStmt ast.Stmt
	assignStmt = &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			ast.NewIdent(string(arg.Name)),
		},
		Rhs: []ast.Expr{
			ast.NewIdent("boolValue"),
		},
	}
	if arg.Type.Nullable {
		assignStmt = &ast.AssignStmt{
			Tok: token.ASSIGN,
			Lhs: []ast.Expr{ast.NewIdent(string(arg.Name))},
			Rhs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X:  ast.NewIdent("boolValue"),
				},
			},
		}
	}

	return &ast.CaseClause{
		List: []ast.Expr{
			&ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("goliteql"),
					Sel: ast.NewIdent("ValueParserLiteral"),
				},
			},
		},
		Body: []ast.Stmt{
			&ast.IfStmt{
				Cond: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("val"),
						Sel: ast.NewIdent("IsBool"),
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{ast.NewIdent("boolValue")},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("val"),
										Sel: ast.NewIdent("BoolValue"),
									},
									Args: []ast.Expr{},
								},
							},
						},
						assignStmt,
					},
				},
			},
		},
	}
}

func generateStringCaseStmt(arg *schema.ArgumentDefinition) ast.Stmt {
	var assignStmt ast.Stmt
	assignStmt = &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			ast.NewIdent(string(arg.Name)),
		},
		Rhs: []ast.Expr{
			ast.NewIdent("stringValue"),
		},
	}
	if arg.Type.Nullable {
		assignStmt = &ast.AssignStmt{
			Tok: token.ASSIGN,
			Lhs: []ast.Expr{ast.NewIdent(string(arg.Name))},
			Rhs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X:  ast.NewIdent("stringValue"),
				},
			},
		}
	}

	return &ast.CaseClause{
		List: []ast.Expr{
			&ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("goliteql"),
					Sel: ast.NewIdent("ValueParserLiteral"),
				},
			},
		},
		Body: []ast.Stmt{
			&ast.IfStmt{
				Cond: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("val"),
						Sel: ast.NewIdent("IsString"),
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{ast.NewIdent("stringValue")},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("val"),
										Sel: ast.NewIdent("StringValue"),
									},
									Args: []ast.Expr{},
								},
							},
						},
						assignStmt,
					},
				},
			},
		},
	}
}

func generateIDCaseStmt(arg *schema.ArgumentDefinition) ast.Stmt {
	var assignStmt ast.Stmt
	assignStmt = &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			ast.NewIdent(string(arg.Name)),
		},
		Rhs: []ast.Expr{
			ast.NewIdent("idValue"),
		},
	}
	if arg.Type.Nullable {
		assignStmt = &ast.AssignStmt{
			Tok: token.ASSIGN,
			Lhs: []ast.Expr{ast.NewIdent(string(arg.Name))},
			Rhs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X:  ast.NewIdent("idValue"),
				},
			},
		}
	}

	return &ast.CaseClause{
		List: []ast.Expr{
			&ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("goliteql"),
					Sel: ast.NewIdent("ValueParserLiteral"),
				},
			},
		},
		Body: []ast.Stmt{
			&ast.IfStmt{
				Cond: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("val"),
						Sel: ast.NewIdent("IsString"),
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Tok: token.DEFINE,
							Lhs: []ast.Expr{ast.NewIdent("idValue")},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("val"),
										Sel: ast.NewIdent("StringValue"),
									},
									Args: []ast.Expr{},
								},
							},
						},
						assignStmt,
					},
				},
			},
		},
	}
}

func generateArgumentsParams(typePrefix string, args schema.ArgumentDefinitions) *ast.FieldList {
	fields := make([]*ast.Field, 0, len(args))

	for _, arg := range args {
		fields = append(fields, &ast.Field{
			Type: generateTypeExprFromFieldType(typePrefix, arg.Type),
		})
	}

	return &ast.FieldList{
		List: fields,
	}
}

func generateDeclareStmts(typePrefix string, args schema.ArgumentDefinitions) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(args))

	stmts = append(stmts, &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok:   token.VAR,
			Specs: generateVarSpecs(typePrefix, args),
		},
	})

	return stmts
}

func generateValueCaseAssignStmt(arg *schema.ArgumentDefinition, indexes *schema.Indexes) ast.Stmt {
	caseSelector := "ValueParserLiteral"
	var rhs ast.Expr = &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent("val"),
			Sel: ast.NewIdent("StringValue"),
		},
	}

	if arg.Type.Nullable {
		rhs = generateStringPointerExpr(rhs)
	}

	if arg.Type.IsBoolean() {
		rhs = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("val"),
				Sel: ast.NewIdent("BoolValue"),
			},
		}

		if arg.Type.Nullable {
			rhs = generateBoolPointerExpr(rhs)
		}
	}

	if arg.Type.IsInt() {
		rhs = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("val"),
				Sel: ast.NewIdent("IntValue"),
			},
		}
	}

	if arg.Type.IsFloat() {
		rhs = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("val"),
				Sel: ast.NewIdent("FloatValue"),
			},
		}
	}

	if !arg.Type.IsID() && !arg.Type.IsString() && !arg.Type.IsBoolean() && !arg.Type.IsInt() && !arg.Type.IsFloat() {
		caseSelector = "ValueParserObject"
		rhs = generateObjectArgumentRhsType(arg, indexes)
	}

	return &ast.CaseClause{
		List: []ast.Expr{
			&ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("goliteql"),
					Sel: ast.NewIdent(caseSelector),
				},
			},
		},
		Body: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.ASSIGN,
				Lhs: []ast.Expr{
					ast.NewIdent(string(arg.Name)),
				},
				Rhs: []ast.Expr{
					rhs,
				},
			},
		},
	}
}

func generateObjectArgumentRhsType(arg *schema.ArgumentDefinition, indexes *schema.Indexes) ast.Expr {
	expandedType := introspection.ExpandType(arg.Type)
	return &ast.CompositeLit{
		Type: generateTypeExprFromExpandedType(expandedType),
		Elts: generateObjectArgumentElements(arg, indexes),
	}
}

func generateObjectArgumentElements(arg *schema.ArgumentDefinition, indexes *schema.Indexes) []ast.Expr {
	ret := make([]ast.Expr, 0)

	typeDefinition := indexes.TypeIndex[string(arg.Type.Name)]
	if typeDefinition != nil {
		for _, field := range typeDefinition.Fields {
			ret = append(ret, generateArgumentValue(field))
		}
	}

	inputDefinition := indexes.InputIndex[string(arg.Type.Name)]
	if inputDefinition != nil {
		for _, field := range inputDefinition.Fields {
			ret = append(ret, generateArgumentValue(field))
		}
	}

	return ret
}

func generateArgumentValue(field *schema.FieldDefinition) ast.Expr {
	// TODO: support object, list
	if !field.IsPrimitive() {
		return &ast.BasicLit{}
	}

	var v ast.Expr = &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.TypeAssertExpr{
				X: &ast.IndexExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("val"),
						Sel: ast.NewIdent("Fields"),
					},
					Index: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("\"%s\"", string(field.Name)),
					},
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("goliteql"),
						Sel: ast.NewIdent("ValueParserLiteral"),
					},
				},
			},
			Sel: ast.NewIdent(fmt.Sprintf("%sValue", string(field.Type.Name))),
		},
	}

	if field.Type.Nullable {
		v = generateStringPointerExpr(v)
	}

	return &ast.KeyValueExpr{
		Key:   ast.NewIdent(toUpperCase(string(field.Name))),
		Value: v,
	}
}

func generateDefaultValueAssignmentStmts(args schema.ArgumentDefinitions, indexes *schema.Indexes) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(args))

	argsCaseStmts := make([]ast.Stmt, 0, len(args))

	returns := make([]ast.Expr, 0, len(args))
	for i, arg := range args {
		returns = append(returns, ast.NewIdent(string(arg.Name)))
		if arg.Default != nil {
			stmts = append(stmts, generateAssignDefaultValueStmt(arg, indexes))
		}

		// default is object
		var bindStmt ast.Stmt = &ast.IfStmt{
			Init: &ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					ast.NewIdent("err"),
				},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("json"),
							Sel: ast.NewIdent("Unmarshal"),
						},
						Args: []ast.Expr{
							ast.NewIdent("rawJSONValue"),
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent(string(arg.Name)),
							},
						},
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
							ast.NewIdent(string(arg.Name)),
							ast.NewIdent("err"),
						},
					},
				},
			},
		}

		if arg.Type.IsBoolean() {
			var rh ast.Expr = &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: ast.NewIdent("string"),
					Args: []ast.Expr{
						ast.NewIdent("rawJSONValue"),
					},
				},
				Op: token.EQL,
				Y: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"true"`,
				},
			}

			rh = generateBoolPointerExpr(rh)

			bindStmt = &ast.AssignStmt{
				Tok: token.ASSIGN,
				Lhs: []ast.Expr{ast.NewIdent(string(arg.Name))},
				Rhs: []ast.Expr{
					rh,
				},
			}
		}

		if arg.Type.IsString() || arg.Type.IsID() {
			bindStmt = &ast.AssignStmt{
				Tok: token.ASSIGN,
				Lhs: []ast.Expr{ast.NewIdent(string(arg.Name))},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("string"),
						Args: []ast.Expr{
							ast.NewIdent("rawJSONValue"),
						},
					},
				},
			}
		}

		if arg.Type.IsInt() {
			bindStmt = &ast.AssignStmt{
				Tok: token.ASSIGN,
				Lhs: []ast.Expr{
					ast.NewIdent(string(arg.Name)),
					ast.NewIdent("_"),
				},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("strconv"),
							Sel: ast.NewIdent("Atoi"),
						},
						Args: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent("string"),
								Args: []ast.Expr{
									ast.NewIdent("rawJSONValue"),
								},
							},
						},
					},
				},
			}
		}

		argsCaseStmts = append(argsCaseStmts, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s\"", string(arg.Name)),
				},
			},
			Body: []ast.Stmt{
				&ast.IfStmt{
					Cond: &ast.UnaryExpr{
						Op: token.NOT,
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("arg"),
								Sel: ast.NewIdent("IsVariable"),
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("ast"),
									ast.NewIdent("err"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X: &ast.SelectorExpr{
													X:   ast.NewIdent("r"),
													Sel: ast.NewIdent("parser"),
												},
												Sel: ast.NewIdent("ValueParser"),
											},
											Sel: ast.NewIdent("Parse"),
										},
										Args: []ast.Expr{
											ast.NewIdent("arg.Value"),
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
										&ast.ReturnStmt{
											Results: []ast.Expr{
												ast.NewIdent(string(arg.Name)),
												ast.NewIdent("err"),
											},
										},
									},
								},
							},
							&ast.TypeSwitchStmt{
								Assign: &ast.AssignStmt{
									Tok: token.DEFINE,
									Lhs: []ast.Expr{
										ast.NewIdent("val"),
									},
									Rhs: []ast.Expr{
										&ast.TypeAssertExpr{
											X:    ast.NewIdent("ast"),
											Type: ast.NewIdent("type"),
										},
									},
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										generateValueCaseAssignStmt(arg, indexes),
									},
								},
							},
						},
					},
					Else: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Tok: token.DEFINE,
								Lhs: []ast.Expr{
									ast.NewIdent("rawJSONValue"),
									ast.NewIdent("ok"),
								},
								Rhs: []ast.Expr{
									&ast.IndexExpr{
										X: ast.NewIdent("variables"),
										Index: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("arg"),
												Sel: ast.NewIdent("VariableAnnotation"),
											},
										},
									},
								},
							},
							&ast.IfStmt{
								Cond: &ast.UnaryExpr{
									Op: token.NOT,
									X:  ast.NewIdent("ok"),
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.ReturnStmt{
											Results: []ast.Expr{
												ast.NewIdent(string(arg.Name)),
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("fmt"),
														Sel: ast.NewIdent("Errorf"),
													},
													Args: []ast.Expr{
														&ast.BasicLit{
															Kind:  token.STRING,
															Value: fmt.Sprintf(`"argument %s is not provided"`, string(arg.Name)),
														},
													},
												},
											},
										},
									},
								},
							},
							bindStmt,
						},
					},
				},
			},
		})

		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: ast.NewIdent("len"),
					Args: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("node"),
							Sel: ast.NewIdent("Arguments"),
						},
					},
				},
				Op: token.EQL,
				Y:  ast.NewIdent(fmt.Sprintf("%d", i)),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{},
			},
		})

		stmts = append(stmts, &ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Value: ast.NewIdent("arg"),
			Tok:   token.DEFINE,
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("node"),
				Sel: ast.NewIdent("Arguments"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.SwitchStmt{
						Tag: &ast.CallExpr{
							Fun: ast.NewIdent("string"),
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("arg"),
									Sel: ast.NewIdent("Name"),
								},
							},
						},
						Body: &ast.BlockStmt{
							List: argsCaseStmts,
						},
					},
				},
			},
		})
	}

	returns = append(returns, ast.NewIdent("nil"))

	stmts = append(stmts, &ast.ReturnStmt{
		Results: returns,
	})

	return stmts
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

	if arg.Type.IsID() {
		if arg.Type.Nullable {
			return generateStringPointerAST(string(arg.Default))
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
						Value: fmt.Sprintf("\\%q\\", value),
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

func generateIntPointerExpr(value ast.Expr) ast.Expr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.IndexExpr{
			X: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: ast.NewIdent("int"),
				},
				Elts: []ast.Expr{
					value,
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

func generateFloatPointerExpr(value ast.Expr) ast.Expr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.IndexExpr{
			X: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: ast.NewIdent("float64"),
				},
				Elts: []ast.Expr{
					value,
				},
			},
			Index: &ast.BasicLit{
				Kind:  token.INT,
				Value: "0",
			},
		},
	}
}

func generateObjectPointerExpr(objectTypeName string, expr ast.Expr) ast.Expr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.IndexExpr{
			X: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: ast.NewIdent(objectTypeName),
				},
				Elts: []ast.Expr{
					expr,
				},
			},
			Index: &ast.BasicLit{
				Kind:  token.INT,
				Value: "0",
			},
		},
	}
}

func generateVarSpecs(typePrefix string, args schema.ArgumentDefinitions) []ast.Spec {
	specs := make([]ast.Spec, 0)
	for _, arg := range args {
		specs = append(specs, &ast.ValueSpec{
			Names: []*ast.Ident{
				ast.NewIdent(string(arg.Name)),
			},
			Type: generateTypeExprFromFieldType(typePrefix, arg.Type),
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

func generateOperationResponseStructDecls(schema *schema.Schema) []ast.Decl {
	var decls []ast.Decl

	for _, typeDefinition := range schema.Types {
		if typeDefinition.IsIntrospection() {
			continue
		}

		if typeDefinition.IsPrimitive() {
			continue
		}

		decls = append(decls, generateResponseStructFromField(typeDefinition))
	}

	return decls
}

func generateResponseStructFromField(typeDefinition *schema.TypeDefinition) ast.Decl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(fmt.Sprintf("%sResponse", typeDefinition.Name)),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: generateAllPointerFieldStructFromField(typeDefinition),
					},
				},
			},
		},
	}
}

func generateAllPointerFieldStructFromField(typeDefinition *schema.TypeDefinition) []*ast.Field {
	fields := make([]*ast.Field, 0, len(typeDefinition.Fields))

	for _, field := range typeDefinition.Fields {
		var typeExpr ast.Expr = &ast.StarExpr{
			X: generateTypeExprFromFieldType("", field.Type),
		}

		if field.Type.Nullable {
			typeExpr = generateTypeExprFromFieldType("", field.Type)
		}

		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(toUpperCase(string(field.Name))),
			},
			Type: typeExpr,
		})
	}

	return fields
}

func generateFieldTypeForResponse(fieldType *schema.FieldType) *schema.FieldType {
	if fieldType == nil {
		return nil
	}

	if fieldType.IsList {
		return &schema.FieldType{
			Name:     fieldType.Name,
			IsList:   true,
			Nullable: true,
			ListType: generateFieldTypeForResponse(fieldType.ListType),
		}
	}

	return &schema.FieldType{
		Name:     []byte(string(fieldType.Name) + "Response"),
		IsList:   fieldType.IsList,
		Nullable: true,
		ListType: fieldType.ListType,
	}
}
