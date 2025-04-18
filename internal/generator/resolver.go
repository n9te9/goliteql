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

func generateWrapResponseWriter(op *schema.OperationDefinition, index map[string]*schema.TypeDefinition) []ast.Decl {
	res := make([]ast.Decl, 0, len(op.Fields))

	for _, field := range op.Fields {
		res = append(res, generateWrapResponseWriterStruct(field))
		res = append(res, generateWrapResponseWriterFunc(field))
		res = append(res, generateWrapResponseWriterWrite(field, index))
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

func generateSliceMapType(typ *schema.FieldType) ast.Expr {
	if !typ.IsList {
		return &ast.MapType{
			Key:   &ast.Ident{Name: "string"},
			Value: &ast.Ident{Name: "json.RawMessage"},
		}
	}

	return &ast.ArrayType{
		Elt: generateSliceMapType(typ.ListType),
	}
}

func generateWrapResponseWriterWriteRangeStmt(field *schema.FieldDefinition, index map[string]*schema.TypeDefinition, typeExpr ast.Expr, k int) []ast.Stmt {
	ft := field.Type
	for ft.IsList {
		ft = ft.ListType
	}

	excludeStmt := &ast.AssignStmt{
		Tok: token.ASSIGN,
		Lhs: []ast.Expr{
			ast.NewIdent("resp"),
		},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("executor"),
					Sel: ast.NewIdent("ExcludeSelectFields"),
				},
				Args: []ast.Expr{
					ast.NewIdent("resp"),
					&ast.SelectorExpr{
						X:   ast.NewIdent("w"),
						Sel: ast.NewIdent("selections"),
					},
				},
			},
		},
	}

	xName := "resp"
	itr := ""
	if k > 0 {
		xName = fmt.Sprintf("v%d", k-1)
		for i := 0; i < k; i++ {
			itr += fmt.Sprintf("[k%d]", i)
		}

		for i := 0; i < k; i++ {
			excludeStmt = &ast.AssignStmt{
				Tok: token.ASSIGN,
				Lhs: []ast.Expr{
					ast.NewIdent(fmt.Sprintf("resp%s", itr)),
				},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("ExcludeSelectFields"),
						},
						Args: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("resp%s", itr)),
							&ast.SelectorExpr{
								X:   ast.NewIdent("w"),
								Sel: ast.NewIdent("selections"),
							},
						},
					},
				},
			}
		}
	}

	body := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.IfStmt{
				Init: &ast.AssignStmt{
					Tok: token.DEFINE,
					Lhs: []ast.Expr{
						ast.NewIdent("_"),
						ast.NewIdent("ok"),
					},
					Rhs: []ast.Expr{
						&ast.IndexExpr{
							X: ast.NewIdent(fmt.Sprintf("resp%s", itr)),
							Index: ast.NewIdent(fmt.Sprintf("k%d", k)),
						},
					},
				},
				Cond: ast.NewIdent("!ok"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.SelectorExpr{
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
										X:   &ast.SelectorExpr{
											X:   ast.NewIdent("w"),
											Sel: ast.NewIdent("ResponseWriter"),
										},
										Sel: ast.NewIdent("Write"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: ast.NewIdent("[]byte"),
											Args: []ast.Expr{
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("fmt"),
														Sel: ast.NewIdent("Sprintf"),
													},
													Args: []ast.Expr{
														&ast.BasicLit{
															Kind:  token.STRING,
															Value: "\"unknown field: %s\"",
														},
														&ast.BasicLit{
															Kind:  token.STRING,
															Value: fmt.Sprintf("k%d", k),
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
				},
			},
		},
	}

	if _, ok := typeExpr.(*ast.MapType); ok {
		return []ast.Stmt{
			&ast.RangeStmt{
				Key:  ast.NewIdent(fmt.Sprintf("k%d", k)),
				Tok:  token.DEFINE,
				X:    ast.NewIdent("selected"),
				Body: body,
			},
			excludeStmt,
		}
	}

	if t, ok := typeExpr.(*ast.ArrayType); ok {
		return []ast.Stmt{
			&ast.RangeStmt{
				Key:   ast.NewIdent(fmt.Sprintf("k%d", k)),
				Tok:   token.DEFINE,
				X:     ast.NewIdent(xName),
				Body: &ast.BlockStmt{
					List: generateWrapResponseWriterWriteRangeStmt(field, index, t.Elt, k+1),
				},
			},
		}
	}

	panic("unknown type")
}

func generateWrapResponseWriterWrite(field *schema.FieldDefinition, index map[string]*schema.TypeDefinition) *ast.FuncDecl {
	typeExpr := generateSliceMapType(field.Type)

	
	rangeStmts := generateWrapResponseWriterWriteRangeStmt(field, index, typeExpr, 0)
	var validationStmt, excludeStmt ast.Stmt
	if len(rangeStmts) == 1 {
		validationStmt = rangeStmts[0]
		excludeStmt = &ast.ExprStmt{
			X: &ast.BasicLit{},
		}
	}

	if len(rangeStmts) == 2 {
		validationStmt = rangeStmts[0]
		excludeStmt = rangeStmts[1]
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
						&ast.CompositeLit{
							Type: typeExpr,
							Elts: nil,
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
									ast.NewIdent("0"),
									ast.NewIdent("err"),
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
						&ast.CallExpr{
							Fun: &ast.BasicLit{
								Kind:  token.IDENT,
								Value: "make",
							},
							Args: []ast.Expr{
								&ast.MapType{
									Key:   &ast.Ident{Name: "string"},
									Value: &ast.StructType{
										Fields: &ast.FieldList{},
									},
								},
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.BasicLit{},
				},
				&ast.RangeStmt{
					Key: ast.NewIdent("_"),
					Value: ast.NewIdent("selection"),
					Tok: token.DEFINE,
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("w"),
						Sel: ast.NewIdent("selections"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
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
							},
							&ast.AssignStmt{
								Tok: token.ASSIGN,
								Lhs: []ast.Expr{
									&ast.IndexExpr{
										X: ast.NewIdent("selected"),
										Index: &ast.CallExpr{
											Fun: &ast.BasicLit{
												Kind:  token.IDENT,
												Value: "string",
											},
											Args: []ast.Expr{
												&ast.SelectorExpr{
													X:   ast.NewIdent("sel"),
													Sel: ast.NewIdent("Name"),
												},
											},
										},
									},
								},
								Rhs: []ast.Expr{
									&ast.CompositeLit{
										Type: &ast.StructType{
											Fields: &ast.FieldList{},
										},
										Elts: []ast.Expr{},
									},
								},
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.BasicLit{},
				},

				validationStmt,
				excludeStmt,

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
								ast.NewIdent("resp"),
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
		argsStr := make([]string, 0, len(f.Arguments))
		for _, arg := range f.Arguments {
			s := recv(arg.Type)
			argsStr = append(argsStr, s)
		}

		returnsStr := recv(f.Type)

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
