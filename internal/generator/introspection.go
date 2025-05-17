package generator

import (
	"fmt"
	"go/ast"
	"go/token"

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

func generateSchemaResponseDataModelAST() ast.Decl {
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

func generateSchemaResponseModelAST() ast.Decl {
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

func generateModelFieldCaseAST(s *schema.Schema, field *schema.FieldDefinition) ast.Stmt {
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

func generateQueryTypeMethodAST(s *schema.Schema) ast.Decl {
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
			List: generateQueryTypeMethodBodyAST(s),
		},
	}
}

func generateQueryTypeMethodBodyAST(s *schema.Schema) []ast.Stmt {
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
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: &ast.ArrayType{
							Elt: ast.NewIdent("__Field"),
						},
						Elts: generateIntrospectionFieldsAST(fieldDefinitions.Fields),
					},
				},
			},
		},
	}
}

func generateIntrospectionFieldsAST(fieldDefinitions []*schema.FieldDefinition) []ast.Expr {
	ret := make([]ast.Expr, 0, len(fieldDefinitions))
	for _, f := range fieldDefinitions {
		elm := make([]ast.Expr, 0, 5)

		elm = append(elm, &ast.KeyValueExpr{
			Key:   ast.NewIdent("Name"),
			Value: ast.NewIdent(fmt.Sprintf(`"%s"`, string(f.Name))),
		})
		// TODO' implement args, description
		elm = append(elm, &ast.KeyValueExpr{
			Key:   ast.NewIdent("Type"),
			Value: generateIntrospectionFieldTypeAST(f.Type),
		})
		elm = append(elm, &ast.KeyValueExpr{
			Key: ast.NewIdent("IsDeprecated"),
			Value: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("%t", f.IsDeprecated()),
			},
		})
		elm = append(elm, &ast.KeyValueExpr{
			Key: ast.NewIdent("DeprecationReason"),
			Value: generateStringPointerAST(f.DeprecatedReason()),
		})

		ret = append(ret, &ast.CompositeLit{
			Elts: elm,
		})
	}

	return ret
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

func generateIntrospectionFieldTypeAST(fieldType *schema.FieldType) ast.Expr {
	nameKeyValue := &ast.KeyValueExpr{
		Key: ast.NewIdent("Name"),
		Value: ast.NewIdent("nil"),
	}
	if fieldType.Name != nil {
		nameKeyValue.Value = generateStringPointerAST(string(fieldType.Name))
	}
	kindKeyValue := &ast.KeyValueExpr{
		Key: ast.NewIdent("Kind"),
		Value: ast.NewIdent("__TypeKind_NON_NULL"),
	}

	return &ast.CompositeLit{
		Type: ast.NewIdent("__Type"),
		Elts: []ast.Expr{
			nameKeyValue,
			kindKeyValue,
			generateIntrospectionFieldOfTypeAST(fieldType),
		},
	}
}

func generateIntrospectionFieldOfTypeAST(fieldType *schema.FieldType) ast.Expr {
	nameKeyValue := &ast.KeyValueExpr{
		Key: ast.NewIdent("Name"),
		Value: ast.NewIdent("nil"),
	}
	if fieldType.Name != nil {
		nameKeyValue.Value = generateStringPointerAST(string(fieldType.Name))
	}

	typeKindValue := &ast.KeyValueExpr{
		Key: ast.NewIdent("Kind"),
		Value: ast.NewIdent("__TypeKind_SCALAR"),
	}

	if !fieldType.IsPrimitive() {
		typeKindValue.Value = ast.NewIdent("__TypeKind_OBJECT")
	}

	if !fieldType.Nullable {
		if fieldType.IsList {
			return &ast.KeyValueExpr{
				Key: ast.NewIdent("OfType"),
				Value: &ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: ast.NewIdent("__Type"),
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: ast.NewIdent("Kind"),
								Value: ast.NewIdent("__TypeKind_NON_NULL"),
							},
							&ast.KeyValueExpr{
								Key: ast.NewIdent("OfType"),
								Value: &ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: ast.NewIdent("__Type"),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: ast.NewIdent("Kind"),
												Value: ast.NewIdent("__TypeKind_LIST"),
											},
											generateIntrospectionFieldOfTypeAST(fieldType.ListType),
										},
									},
								},
							},
						},
					},
				},
			}
		} else {
			return &ast.KeyValueExpr{
				Key: ast.NewIdent("OfType"),
				Value: &ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: ast.NewIdent("__Type"),
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key: ast.NewIdent("Kind"),
								Value: ast.NewIdent("__TypeKind_NON_NULL"),
							},
							&ast.KeyValueExpr{
								Key: ast.NewIdent("OfType"),
								Value: &ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: ast.NewIdent("__Type"),
										Elts: []ast.Expr{
											nameKeyValue,
											typeKindValue,
										},
									},
								},
							},
						},
					},
				},
			}
		}
	}

	return &ast.KeyValueExpr{
		Key: ast.NewIdent("OfType"),
		Value: &ast.UnaryExpr{
			Op: token.AND,
			X: &ast.CompositeLit{
				Type: ast.NewIdent("__Type"),
				Elts: []ast.Expr{
					nameKeyValue,
					typeKindValue,
				},
			},
		},
	}
}

func generateModelFieldCaseASTs(s *schema.Schema, fields []*schema.FieldDefinition) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(fields))
	for _, f := range fields {
		stmts = append(stmts, generateModelFieldCaseAST(s, f))
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
