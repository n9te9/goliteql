package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/internal/generator/introspection"
	"github.com/n9te9/goliteql/schema"
)

func generateIntrospectionModelAST(types []*schema.TypeDefinition) []ast.Decl {
	decls := make([]ast.Decl, 0, len(types))

	for _, t := range types {
		if t.IsIntrospection() {
			if t.PrimitiveTypeName != nil {
				decls = append(decls, &ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: ast.NewIdent(string(t.Name)),
							Type: ast.NewIdent(string(t.PrimitiveTypeName)),
						},
					},
				})

				continue
			}

			decls = append(decls, &ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent(string(t.Name)),
						Type: &ast.StructType{
							Fields: generateModelFieldWithOmitempty(t.Fields),
						},
					},
				},
			})
		}
	}

	return decls
}

func generateIntrospectionSchemaResponseDataModelAST() ast.Decl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("__SchemaResponseData"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("Schema"),
								},
								Type: &ast.StarExpr{
									X: ast.NewIdent("__Schema"),
								},
								Tag: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "`json:\"__schema\"`",
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateIntrospectionSchemaResponseModelAST() ast.Decl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("__SchemaResponse"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("Data"),
								},
								Type: &ast.StarExpr{
									X: ast.NewIdent("__SchemaResponseData"),
								},
								Tag: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "`json:\"data\"`",
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("Errors"),
								},
								Type: &ast.ArrayType{
									Elt: ast.NewIdent("error"),
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateIntrospectionModelFieldCaseAST(s *schema.Schema, field *schema.FieldDefinition) ast.Stmt {
	var stmts []ast.Stmt
	switch string(field.Name) {
	case "description":
		// TODO
	case "queryType":
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret"),
					Sel: ast.NewIdent("QueryType"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("__schema_queryType"),
					},
					Args: []ast.Expr{
						ast.NewIdent("child"),
					},
				},
			},
		})
	}

	return &ast.CaseClause{
		List: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s"`, string(field.Name)),
			},
		},
		Body: stmts,
	}
}

func generateNodeWalkerArgs() *ast.FieldList {
	return &ast.FieldList{
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
		},
	}
}

func generateIntrospectionQueryTypeMethodAST(s *schema.Schema) ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("resolver"),
					},
				},
			},
		},
		Name: ast.NewIdent("__schema_queryType"),
		Type: &ast.FuncType{
			Params: generateNodeWalkerArgs(),
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("__Type"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: generateIntrospectionQueryTypeMethodBodyAST(s),
		},
	}
}

func generateIntrospectionTypeMethodDecls(s *schema.Schema) []ast.Decl {
	ret := make([]ast.Decl, 0)
	q := s.GetQuery()
	if q == nil {
		return ret
	}

	for _, field := range q.Fields {
		ret = append(ret, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("r"),
						},
						Type: &ast.StarExpr{
							X: ast.NewIdent("resolver"),
						},
					},
				},
			},
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(field.Name))),
			Type: &ast.FuncType{
				Params: generateNodeWalkerArgs(),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("__Type"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("ret"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CompositeLit{
								Type: ast.NewIdent("__Type"),
								Elts: []ast.Expr{},
							},
						},
					},
					&ast.RangeStmt{
						Key:   ast.NewIdent("_"),
						Tok:   token.DEFINE,
						Value: ast.NewIdent("child"),
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("node"),
							Sel: ast.NewIdent("Children"),
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								generateIntrospectionTypeFieldSwitchStmt(string(field.Type.Name), field),
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("ret"),
						},
					},
				},
			},
		})
	}

	return ret
}

func generateIntrospectionFieldTypeTypeOfDecls(s *schema.Schema) []ast.Decl {
	ret := make([]ast.Decl, 0)

	q := s.GetQuery()
	if q == nil {
		return ret
	}

	for _, field := range q.Fields {
		if field.Type.IsList && field.Type.Nullable {
			ret = append(ret, generateIntrospectionRecursiveFieldTypeOfDecls(string(field.Name), introspection.ExpandType(field.Type).Unwrap(), 0)...)
		} else {
			ret = append(ret, generateIntrospectionRecursiveFieldTypeOfDecls(string(field.Name), introspection.ExpandType(field.Type).Unwrap(), 0)...)
		}
	}

	for _, t := range s.Types {
		if t.IsIntrospection() {
			continue
		}

		for _, field := range t.Fields {
			if field.Type.IsList && field.Type.Nullable {
				ret = append(ret, generateIntrospectionRecursiveFieldTypeOfDecls(fmt.Sprintf("%s__%s", t.Name, field.Name), introspection.ExpandType(field.Type).Unwrap(), 0)...)
			} else {
				ret = append(ret, generateIntrospectionRecursiveFieldTypeOfDecls(fmt.Sprintf("%s__%s", t.Name, field.Name), introspection.ExpandType(field.Type).Unwrap(), 0)...)
			}
		}
	}

	return ret
}

func generateIntrospectionTypeResolverDeclsFromTypeDefinitions(typeDefinitions []*schema.TypeDefinition) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, t := range typeDefinitions {
		if t.IsIntrospection() {
			continue
		}

		for _, field := range t.Fields {
			ret = append(ret, &ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("r"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("resolver"),
							},
						},
					},
				},
				Name: ast.NewIdent(fmt.Sprintf("__schema__%s__%s__type", string(t.Name), string(field.Name))),
				Type: &ast.FuncType{
					Params: generateNodeWalkerArgs(),
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("__Type"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("ret"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CompositeLit{
									Type: ast.NewIdent("__Type"),
									Elts: []ast.Expr{},
								},
							},
						},
						&ast.RangeStmt{
							Key:   ast.NewIdent("_"),
							Tok:   token.DEFINE,
							Value: ast.NewIdent("child"),
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("node"),
								Sel: ast.NewIdent("Children"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									generateIntrospectionTypeFieldSwitchStmt(string(t.Name), field),
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("ret"),
							},
						},
					},
				},
			})
		}
	}

	return ret
}

func generateIntrospectionTypeFieldsDecls(typeDefinitions []*schema.TypeDefinition) []ast.Decl {
	ret := make([]ast.Decl, 0)
	generateIntrospectionFieldTypeAssignStmtFunc := func(t *schema.TypeDefinition) []ast.Stmt {
		ret := make([]ast.Stmt, 0, len(t.Fields))

		for i, field := range t.Fields {
			ret = append(ret, &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.IndexExpr{
							X: ast.NewIdent("ret"),
							Index: &ast.BasicLit{
								Kind:  token.INT,
								Value: fmt.Sprintf("%d", i),
							},
						},
						Sel: ast.NewIdent("Name"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s"`, string(field.Name)),
					},
				},
			})
		}

		return ret
	}

	for _, t := range typeDefinitions {
		if t.IsIntrospection() {
			continue
		}

		ret = append(ret, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("r"),
						},
						Type: &ast.StarExpr{
							X: ast.NewIdent("resolver"),
						},
					},
				},
			},
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", string(t.Name))),
			Type: &ast.FuncType{
				Params: generateNodeWalkerArgs(),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.StarExpr{
								X: &ast.ArrayType{
									Elt: ast.NewIdent("__Field"),
								},
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("ret"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent("make"),
								Args: []ast.Expr{
									&ast.ArrayType{
										Elt: ast.NewIdent("__Field"),
									},
									ast.NewIdent(fmt.Sprintf("%d", len(t.Fields))),
								},
							},
						},
					},
					&ast.RangeStmt{
						Key:   ast.NewIdent("_"),
						Tok:   token.DEFINE,
						Value: ast.NewIdent("child"),
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("node"),
							Sel: ast.NewIdent("Children"),
						},
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
										List: []ast.Stmt{
											&ast.CaseClause{
												List: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"name"`,
													},
												},
												Body: generateIntrospectionFieldTypeAssignStmtFunc(t),
											},
											&ast.CaseClause{
												List: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"description"`,
													},
												},
												Body: []ast.Stmt{
													&ast.ExprStmt{
														X: &ast.BasicLit{
															Kind:  token.STRING,
															Value: "// TODO",
														},
													},
												},
											},
											&ast.CaseClause{
												List: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"type"`,
													},
												},
												Body: generateIntrospectionFieldTypeBodyStmt(string(t.Name), t.Fields),
											},
										},
									},
								},
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("ret"),
							},
						},
					},
				},
			},
		})
	}

	return ret
}

func generateIntrospectionRecursiveFieldTypeOfDecls(fieldDefinitionName string, field *introspection.FieldType, nestCount int) []ast.Decl {
	typeOfSuffix := "__typeof"
	for range nestCount {
		typeOfSuffix += "__typeof"
	}

	decls := make([]ast.Decl, 0)

	var bodyStmts []ast.Stmt
	if field != nil {
		bodyStmts = append(bodyStmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("ret"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("new"),
					Args: []ast.Expr{
						ast.NewIdent("__Type"),
					},
				},
			},
		})

		bodyStmts = append(bodyStmts, &ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Tok:   token.DEFINE,
			Value: ast.NewIdent("child"),
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("node"),
				Sel: ast.NewIdent("Children"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					generateIntrospectionTypeOfSwitchStmt(field, fmt.Sprintf("__schema__%s%s__typeof", fieldDefinitionName, typeOfSuffix)),
				},
			},
		})

		bodyStmts = append(bodyStmts, &ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
			},
		})
	} else {
		bodyStmts = append(bodyStmts, &ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: "nil",
				},
			},
		})
	}

	decls = append(decls, &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("resolver"),
					},
				},
			},
		},
		Name: ast.NewIdent(fmt.Sprintf("__schema__%s%s", fieldDefinitionName, typeOfSuffix)),
		Type: &ast.FuncType{
			Params: generateNodeWalkerArgs(),
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: ast.NewIdent("__Type"),
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: bodyStmts,
		},
	})

	if field != nil && field.Child != nil {
		decls = append(decls, generateIntrospectionRecursiveFieldTypeOfDecls(fieldDefinitionName, field.Child, nestCount+1)...)
	}

	return decls
}

func generateIntrospectionTypeOfSwitchStmt(f *introspection.FieldType, callTypeOfFuncName string) ast.Stmt {
	var nameExpr, kindExpr ast.Expr
	var fieldAssignStmt ast.Stmt
	if f.IsPrimitive() {
		kindExpr = ast.NewIdent("__TypeKind_SCALAR")
		nameExpr = generateStringPointerAST(string(f.Name))
	}

	fieldAssignStmt = &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent("ret"),
				Sel: ast.NewIdent("Fields"),
			},
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: "nil",
			},
		},
	}

	if f.IsObject() && !f.IsPrimitive() {
		kindExpr = ast.NewIdent("__TypeKind_OBJECT")
		nameExpr = generateStringPointerAST(string(f.Name))
		fieldAssignStmt = &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret"),
					Sel: ast.NewIdent("Fields"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("__schema__" + string(f.Name) + "__fields"),
					},
					Args: []ast.Expr{
						ast.NewIdent("child"),
					},
				},
			},
		}
	}

	if f.IsList {
		kindExpr = ast.NewIdent("__TypeKind_LIST")
		nameExpr = &ast.BasicLit{
			Kind:  token.STRING,
			Value: "nil",
		}
	}

	if f.NonNull {
		kindExpr = ast.NewIdent("__TypeKind_NON_NULL")
		nameExpr = &ast.BasicLit{
			Kind:  token.STRING,
			Value: "nil",
		}
	}

	var ofTypeAssignRhsExpr ast.Expr = ast.NewIdent("nil")
	if f.Child != nil{
		ofTypeAssignRhsExpr = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("r"),
				Sel: ast.NewIdent(callTypeOfFuncName),
			},
			Args: []ast.Expr{
				ast.NewIdent("child"),
			},
		}
	}

	return &ast.SwitchStmt{
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
			List: []ast.Stmt{
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"kind"`,
						},
					},
					Body: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("Kind"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								kindExpr,
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"name"`,
						},
					},
					Body: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("Name"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								nameExpr,
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"description"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"fields"`,
						},
					},
					Body: []ast.Stmt{
						fieldAssignStmt,
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"interfaces"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"enumValues"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"possibleTypes"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"inputFields"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"ofType"`,
						},
					},
					Body: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("OfType"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								ofTypeAssignRhsExpr,
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"specifiedByURL"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
			},
		},
	}
}

func generateIntrospectionTypeFieldSwitchStmt(typeName string, f *schema.FieldDefinition) ast.Stmt {
	var nameExpr ast.Expr
	kindValue := ""
	if f.IsPrimitive() {
		kindValue = "__TypeKind_SCALAR"
		nameExpr = generateStringPointerAST(string(f.Type.Name))
	} else {
		kindValue = "__TypeKind_OBJECT"
		nameExpr = generateStringPointerAST(string(f.Type.Name))
	}

	if f.Type.IsList {
		kindValue = "__TypeKind_LIST"
		nameExpr = &ast.BasicLit{
			Kind:  token.STRING,
			Value: "nil",
		}
	}

	if !f.Type.Nullable {
		kindValue = "__TypeKind_NON_NULL"
		nameExpr = &ast.BasicLit{
			Kind:  token.STRING,
			Value: "nil",
		}
	}

	var ofTypeAssignRhsExpr ast.Expr = ast.NewIdent("nil")
	if f.Type.IsObject() || f.Type.IsList {
		if typeName == string(f.Type.Name) {
			typeName = ""
		} else {
			typeName = fmt.Sprintf("__%s", typeName)
		}

		ofTypeAssignRhsExpr = &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("r"),
				Sel: ast.NewIdent(fmt.Sprintf("__schema%s__%s__typeof", typeName, string(f.Name))),
			},
			Args: []ast.Expr{
				ast.NewIdent("child"),
			},
		}
	}

	return &ast.SwitchStmt{
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
			List: []ast.Stmt{
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"kind"`,
						},
					},
					Body: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("Kind"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								ast.NewIdent(kindValue),
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"name"`,
						},
					},
					Body: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("Name"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								nameExpr,
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"description"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"fields"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"interfaces"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"enumValues"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"possibleTypes"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"inputFields"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"ofType"`,
						},
					},
					Body: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("OfType"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								ofTypeAssignRhsExpr,
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"specifiedByURL"`,
						},
					},
					Body: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.BasicLit{
								Kind:  token.STRING,
								Value: "// TODO",
							},
						},
					},
				},
			},
		},
	}
}

func generateIntrospectionQueryTypeMethodBodyAST(s *schema.Schema) []ast.Stmt {
	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("ret"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: ast.NewIdent("__Type"),
					Elts: []ast.Expr{},
				},
			},
		},
		&ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Tok:   token.DEFINE,
			Value: ast.NewIdent("child"),
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("node"),
				Sel: ast.NewIdent("Children"),
			},
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
							List: generateQueryTypeSwitchBodyAST(s),
						},
					},
				},
			},
		},
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
			},
		},
	}
}

func generateQueryTypeSwitchBodyAST(s *schema.Schema) []ast.Stmt {
	return []ast.Stmt{
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"name"`,
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("Name"),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						generateStringPointerAST(string(s.Definition.Query)),
					},
				},
			},
		},
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"kind"`,
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("Kind"),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent("__TypeKind_OBJECT"),
					},
				},
			},
		},
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"fields"`,
				},
			},
			Body: generateIntrospectionOperationFieldsAST(s.GetQuery()),
		},
	}
}

func generateIntrospectionOperationFieldsAST(fieldDefinitions *schema.OperationDefinition) []ast.Stmt {
	if fieldDefinitions == nil {
		return []ast.Stmt{}
	}

	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret"),
					Sel: ast.NewIdent("Fields"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("__schema_fields"),
					},
					Args: []ast.Expr{
						ast.NewIdent("child"),
					},
				},
			},
		},
	}
}

func generateIntrospectionFieldsFuncsAST(attributeName string, fieldDefinitions schema.FieldDefinitions) []ast.Decl {
	decls := make([]ast.Decl, 0)

	decls = append(decls, generateIntrospectionFieldsFuncAST(attributeName, fieldDefinitions))

	return decls
}

func generateIntrospectionFieldsFuncAST(attributeName string, fields schema.FieldDefinitions) ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("resolver"),
					},
				},
			},
		},
		Name: ast.NewIdent("__schema_fields"),
		Type: &ast.FuncType{
			Params: generateNodeWalkerArgs(),
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.ArrayType{
								Elt: ast.NewIdent("__Field"),
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: generateIntrospectionFieldFuncBodyStmts(attributeName, fields),
		},
	}
}

func generateIntrospectionFieldFuncBodyStmts(attributeName string, fields schema.FieldDefinitions) []ast.Stmt {
	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("ret"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("make"),
					Args: []ast.Expr{
						&ast.ArrayType{
							Elt: ast.NewIdent("__Field"),
						},
						ast.NewIdent(fmt.Sprintf("%d", len(fields))),
					},
				},
			},
		},
		&ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Tok:   token.DEFINE,
			Value: ast.NewIdent("child"),
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("node"),
				Sel: ast.NewIdent("Children"),
			},
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
							List: []ast.Stmt{
								&ast.CaseClause{
									List: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"name"`,
										},
									},
									Body: generateIntrospectionFieldNameBodyStmt(fields),
								},
								&ast.CaseClause{
									List: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"description"`,
										},
									},
								},
								&ast.CaseClause{
									List: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"type"`,
										},
									},
									Body: generateIntrospectionFieldTypeBodyStmt(attributeName, fields),
								},
								&ast.CaseClause{
									List: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"args"`,
										},
									},
								},
								&ast.CaseClause{
									List: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"isDeprecated"`,
										},
									},
									Body: []ast.Stmt{},
								},
								&ast.CaseClause{
									List: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"deprecationReason"`,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		&ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X:  ast.NewIdent("ret"),
				},
			},
		},
	}
}

func generateIntrospectionFieldNameBodyStmt(fields schema.FieldDefinitions) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(fields))
	for i, field := range fields {
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X: &ast.IndexExpr{
						X: ast.NewIdent("ret"),
						Index: &ast.BasicLit{
							Kind:  token.INT,
							Value: fmt.Sprintf("%d", i),
						},
					},
					Sel: ast.NewIdent("Name"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf(`"%s"`, string(field.Name))),
			},
		})
	}

	return stmts
}

func generateIntrospectionFieldTypeBodyStmt(attributeName string, fields schema.FieldDefinitions) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(fields))
	for i, field := range fields {
		prefix := fmt.Sprintf("__schema__%s__%s", attributeName, string(field.Name))
		if attributeName == "" {
			prefix = fmt.Sprintf("__schema__%s", string(field.Name))
		}

		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X: &ast.IndexExpr{
						X: ast.NewIdent("ret"),
						Index: &ast.BasicLit{
							Kind:  token.INT,
							Value: fmt.Sprintf("%d", i),
						},
					},
					Sel: ast.NewIdent("Type"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("%s__type", prefix)),
					},
					Args: []ast.Expr{
						ast.NewIdent("child"),
					},
				},
			},
		})
	}

	return stmts
}

func generateStringPointerAST(value string) ast.Expr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.IndexExpr{
			X: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: ast.NewIdent("string"),
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

func generateModelFieldCaseASTs(s *schema.Schema, fields []*schema.FieldDefinition) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(fields))
	for _, f := range fields {
		stmts = append(stmts, generateIntrospectionModelFieldCaseAST(s, f))
	}

	return stmts
}

func generateIntrospectionSchemaQueryAST(s *schema.Schema) ast.Decl {
	stmts := make([]ast.Stmt, 0)
	for _, t := range s.Types {
		if string(t.Name) == "__Schema" {
			stmts = append(stmts, generateModelFieldCaseASTs(s, t.Fields)...)
		}
	}

	body := make([]ast.Stmt, 0)

	body = append(body, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("ret"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: ast.NewIdent("new"),
				Args: []ast.Expr{
					ast.NewIdent("__Schema"),
				},
			},
		},
	})
	body = append(body, &ast.RangeStmt{
		Key:   ast.NewIdent("_"),
		Tok:   token.DEFINE,
		Value: ast.NewIdent("child"),
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("node"),
			Sel: ast.NewIdent("Children"),
		},
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
						List: stmts,
					},
				},
			},
		},
	})

	body = append(body, generateResponseWrite())

	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("resolver"),
					},
				},
			},
		},
		Name: ast.NewIdent("__schema"),
		Type: &ast.FuncType{
			Params: generateOperationExecutorArgs(),
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

func generateResponseWrite() ast.Stmt {
	return &ast.IfStmt{
		Init: &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("err"),
			},
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
						Sel: &ast.Ident{
							Name: "Encode",
						},
					},
					Args: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent("__SchemaResponse"),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: ast.NewIdent("Data"),
										Value: &ast.UnaryExpr{
											Op: token.AND,
											X: &ast.CompositeLit{
												Type: ast.NewIdent("__SchemaResponseData"),
												Elts: []ast.Expr{
													&ast.KeyValueExpr{
														Key:   ast.NewIdent("Schema"),
														Value: ast.NewIdent("ret"),
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
							&ast.SelectorExpr{
								X:   ast.NewIdent("http"),
								Sel: ast.NewIdent("StatusInternalServerError"),
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("w"),
							Sel: ast.NewIdent("Write"),
						},
						Args: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent("[]byte"),
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("err"),
											Sel: ast.NewIdent("Error"),
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
