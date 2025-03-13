package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/lkeix/gg-parser/schema"
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

func generateResolverStruct(query, mutation, subscription *schema.OperationDefinition) *ast.GenDecl {
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
		Tok:   token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "Resolver",
				},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
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
					Type:  &ast.StarExpr{X: ast.NewIdent("Resolver")},
				},
			},
		},
		Type: &ast.FuncType{
			Params: generateServeHTTPArgs(),
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
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				generateServeHTTPBody(query, mutation, subscription),
			},
		},
	}
}

func generateServeHTTPBody(query, mutation, subscription *schema.OperationDefinition) *ast.BlockStmt {
	querySwitchCases := []ast.Stmt{}

	if query != nil {
		for _, field := range query.Fields {
			querySwitchCases = append(querySwitchCases, &ast.CaseClause{
				List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", field.Name)}},
				Body: []ast.Stmt{
					&ast.ExprStmt{X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r." + newQueryName(query)),
							Sel: ast.NewIdent(toUpperCase(string(field.Name))),
						},
						Args: []ast.Expr{
							ast.NewIdent("req"),
							ast.NewIdent("w"),
						},
					}},
					&ast.ReturnStmt{},
				},
			})
		}
	}

	mutationSwitchCases := []ast.Stmt{}
	if mutation != nil {
		for _, field := range mutation.Fields {
			mutationSwitchCases = append(mutationSwitchCases, &ast.CaseClause{
				List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", field.Name)}},
				Body: []ast.Stmt{
					&ast.ExprStmt{X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r." + newMutationName(mutation)),
							Sel: ast.NewIdent(toUpperCase(string(field.Name))),
						},
						Args: []ast.Expr{
							ast.NewIdent("req"),
							ast.NewIdent("w"),
						},
					}},
					&ast.ReturnStmt{},
				},
			})
		}
	}

	subscriptionSwitchCases := []ast.Stmt{}
	if subscription != nil {
		for _, field := range subscription.Fields {
			subscriptionSwitchCases = append(subscriptionSwitchCases, &ast.CaseClause{
				List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", field.Name)}},
				Body: []ast.Stmt{
					&ast.ExprStmt{X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r." + newSubscriptionName(subscription)),
							Sel: ast.NewIdent(toUpperCase(string(field.Name))),
						},
						Args: []ast.Expr{
							ast.NewIdent("req"),
							ast.NewIdent("w"),
						},
					}},
					&ast.ReturnStmt{},
				},
			})
		}
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
								&ast.BasicLit{Kind: token.STRING, Value: "\"Invalid JSON\""},
								ast.NewIdent("http.StatusBadRequest"),
							},
						}},
						&ast.ReturnStmt{},
					},
				},
			},
			
			&ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("operationType")},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("detectOperationType"),
					},
				},
			},

			&ast.SwitchStmt{
				Tag: ast.NewIdent("operationType"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.CaseClause{
							List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "\"query\""}},
							Body: []ast.Stmt{
								&ast.SwitchStmt{
									Tag: ast.NewIdent("request.OperationName"),
									Body: &ast.BlockStmt{
										List: querySwitchCases,
									},
								},
							},
						},
						
						&ast.CaseClause{
							List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "\"mutation\""}},
							Body: []ast.Stmt{
								&ast.SwitchStmt{
									Tag: ast.NewIdent("request.OperationName"),
									Body: &ast.BlockStmt{
										List: mutationSwitchCases,
									},
								},
							},
						},
						
						&ast.CaseClause{
							List: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: "\"subscription\""}},
							Body: []ast.Stmt{
								&ast.SwitchStmt{
									Tag: ast.NewIdent("request.OperationName"),
									Body: &ast.BlockStmt{
										List: subscriptionSwitchCases,
									},
								},
							},
						},
					},
				},
			},
			
			&ast.ExprStmt{X: &ast.CallExpr{
				Fun: ast.NewIdent("http.Error"),
				Args: []ast.Expr{
					ast.NewIdent("w"),
					&ast.BasicLit{Kind: token.STRING, Value: "\"Unknown operation\""},
					ast.NewIdent("http.StatusBadRequest"),
				},
			}},
		},
	}
}

func generateServeHTTPArgs() *ast.FieldList {
	return &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{
					{
						Name: "req",
					},
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X: ast.NewIdent("http"),
						Sel: ast.NewIdent("Request"),
					},
				},
			},
			{
				Names: []*ast.Ident{
					{
						Name: "w",
					},
				},
				Type: &ast.ParenExpr{
					X: &ast.SelectorExpr{
						X: ast.NewIdent("http"),
						Sel: ast.NewIdent("ResponseWriter"),
					},
				},
			},
		},
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

func generateInterfaceField(operation *schema.OperationDefinition, modelPackagePath string) *ast.GenDecl {
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
					Params: generateServeHTTPArgs(),
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
