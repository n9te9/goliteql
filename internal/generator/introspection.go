package generator

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

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
							Fields: generateModelFieldForResponse(t.Fields),
						},
					},
				},
			})
		}
	}

	return decls
}

func generateIntrospectionTypeResponseDataModelAST() ast.Decl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("__TypeResponseData"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("Type"),
								},
								Type: &ast.StarExpr{
									X: ast.NewIdent("__Type"),
								},
								Tag: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "`json:\"__type\"`",
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateIntrospectionTypeResponseModelAST() ast.Decl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("__TypeResponse"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("Data"),
								},
								Type: &ast.StarExpr{
									X: ast.NewIdent("__TypeResponseData"),
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
								Tag: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "`json:\"errors,omitempty\"`",
								},
							},
						},
					},
				},
			},
		},
	}
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

func generateIntrospectionInterfaceTypeFuncDecls(interfaces []*schema.InterfaceDefinition, indexes *schema.Indexes) []ast.Decl {
	ret := make([]ast.Decl, 0, len(interfaces))

	for _, i := range interfaces {
		body := []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{
								ast.NewIdent("ret"),
							},
							Type: ast.NewIdent("__Type"),
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
													ast.NewIdent("__TypeKind_INTERFACE"),
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"name"`),
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
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent(fmt.Sprintf("%q", string(i.Name))),
														},
													},
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
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("includeDeprecated"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("r"),
															Sel: ast.NewIdent("extract__fieldsArgs"),
														},
														Args: []ast.Expr{
															ast.NewIdent("child"),
															ast.NewIdent("variables"),
														},
													},
												},
											},
											generateReturnErrorHandlingStmt([]ast.Expr{
												&ast.CompositeLit{
													Type: ast.NewIdent("__Type"),
												},
											}),
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("fields"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("r"),
															Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", string(i.Name))),
														},
														Args: []ast.Expr{
															ast.NewIdent("ctx"),
															ast.NewIdent("child"),
															ast.NewIdent("variables"),
															&ast.StarExpr{
																X: ast.NewIdent("includeDeprecated"),
															},
														},
													},
												},
											},
											generateReturnErrorHandlingStmt([]ast.Expr{
												&ast.CompositeLit{
													Type: ast.NewIdent("__Type"),
												},
											}),
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "ret",
														},
														Sel: ast.NewIdent("Fields"),
													},
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent("fields"),
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"enumValues"`),
										},
										Body: []ast.Stmt{
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("ret"),
														Sel: ast.NewIdent("EnumValues"),
													},
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent("nil"),
														},
													},
												},
											},
										},
									}, &ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"inputFields"`),
										},
										Body: []ast.Stmt{
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("ret"),
														Sel: ast.NewIdent("InputFields"),
													},
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent("nil"),
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"possibleTypes"`),
										},
										Body: func() []ast.Stmt {
											ret := make([]ast.Stmt, 0, len(indexes.TypeIndex))
											ret = append(ret, &ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("possibleTypes"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: ast.NewIdent("make"),
														Args: []ast.Expr{
															&ast.ArrayType{
																Elt: ast.NewIdent("__Type"),
															},
															ast.NewIdent("0"),
														},
													},
												},
											})
											for k, typeDefinition := range indexes.TypeIndex {
												for _, interfaceDefinition := range typeDefinition.Interfaces {
													if string(interfaceDefinition.Name) == string(i.Name) {
														ret = append(ret, &ast.AssignStmt{
															Lhs: []ast.Expr{
																ast.NewIdent(fmt.Sprintf("interfaceType%s", k)),
																ast.NewIdent("err"),
															},
															Tok: token.DEFINE,
															Rhs: []ast.Expr{
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X:   ast.NewIdent("r"),
																		Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(typeDefinition.Name))),
																	},
																	Args: []ast.Expr{
																		ast.NewIdent("ctx"),
																		ast.NewIdent("child"),
																		ast.NewIdent("variables"),
																	},
																},
															},
														},
															generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("ret")}),
															&ast.AssignStmt{
																Lhs: []ast.Expr{
																	ast.NewIdent("possibleTypes"),
																},
																Tok: token.ASSIGN,
																Rhs: []ast.Expr{
																	&ast.CallExpr{
																		Fun: ast.NewIdent("append"),
																		Args: []ast.Expr{
																			ast.NewIdent("possibleTypes"),
																			ast.NewIdent(fmt.Sprintf("interfaceType%s", k)),
																		},
																	},
																},
															})
													}
												}
											}

											ret = append(ret, &ast.AssignStmt{
												Lhs: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("ret"),
														Sel: ast.NewIdent("PossibleTypes"),
													},
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent("possibleTypes"),
														},
													},
												},
											})

											return ret
										}(),
									},
								},
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
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(i.Name))),
			Type: &ast.FuncType{
				Params: generateNodeWalkerArgs(),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("__Type"),
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: body,
			},
		})
	}

	return ret
}

func generateIntrospectionTypesFuncDecl(schema *schema.Schema) ast.Decl {
	typeDefinitions := schema.Types
	interfaceDefinitions := schema.Interfaces
	inputDefinitions := schema.Inputs
	scalarDefinitions := schema.Scalars
	enumDefinitions := schema.Enums
	unionDefinitions := schema.Unions

	stmts := make([]ast.Stmt, 0)
	stmts = append(stmts, generateIntrospectionTypesFieldSwitchStmts(schema, typeDefinitions, interfaceDefinitions, inputDefinitions, scalarDefinitions, enumDefinitions, unionDefinitions)...)
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("ret"),
			ast.NewIdent("nil"),
		},
	})

	return &ast.FuncDecl{
		Name: ast.NewIdent("__schema_types"),
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
		Type: &ast.FuncType{
			Params: generateNodeWalkerArgs(),
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: ast.NewIdent("__Type"),
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}
}

func generateIntrospectionTypesFieldSwitchStmts(schema *schema.Schema, typeDefinitions []*schema.TypeDefinition, interfaceDefinitions []*schema.InterfaceDefinition, inputDefinitions []*schema.InputDefinition, scalarDefinitions []*schema.ScalarDefinition, enumDefinitions schema.EnumDefinitions, unionDefinitions schema.UnionDefinitions) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(typeDefinitions))

	stmts = append(stmts, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("ret"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: ast.NewIdent("make"),
				Args: []ast.Expr{
					&ast.ArrayType{
						Elt: ast.NewIdent("__Type"),
					},
					ast.NewIdent("0"),
					ast.NewIdent(fmt.Sprintf("%d", len(typeDefinitions)+len(interfaceDefinitions)+len(scalarDefinitions)+len(inputDefinitions)+len(enumDefinitions)+len(unionDefinitions)+3)),
				},
			},
		},
	})

	args := make([]ast.Expr, 0)
	args = append(args, ast.NewIdent("ret"))
	if q := schema.GetQuery(); q != nil {
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("queryType"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("__schema_queryType"),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("node"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		args = append(args, ast.NewIdent("queryType"))
	}

	if m := schema.GetMutation(); m != nil {
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("mutationType"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("__schema_mutationType"),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("node"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		args = append(args, ast.NewIdent("mutationType"))
	}

	for _, scalarDefinition := range scalarDefinitions {
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("scalar%s", scalarDefinition.Name)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(scalarDefinition.Name))),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("node"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}))
		args = append(args, ast.NewIdent(fmt.Sprintf("scalar%s", scalarDefinition.Name)))
	}

	for _, t := range typeDefinitions {
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("type%s", t.Name)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(t.Name))),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("node"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}))
		args = append(args, ast.NewIdent(fmt.Sprintf("type%s", t.Name)))
	}

	for _, interfaceDefinition := range interfaceDefinitions {
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("interface%s", interfaceDefinition.Name)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(interfaceDefinition.Name))),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("node"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}))
		args = append(args, ast.NewIdent(fmt.Sprintf("interface%s", interfaceDefinition.Name)))
	}

	for _, inputDefinition := range inputDefinitions {
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("input%s", inputDefinition.Name)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(inputDefinition.Name))),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("node"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}))
		args = append(args, ast.NewIdent(fmt.Sprintf("input%s", inputDefinition.Name)))
	}

	for _, enumDefinition := range enumDefinitions {
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("enum%s", enumDefinition.Name)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(enumDefinition.Name))),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("node"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}))
		args = append(args, ast.NewIdent(fmt.Sprintf("enum%s", enumDefinition.Name)))
	}

	for _, unionDefinition := range unionDefinitions {
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("union%s", unionDefinition.Name)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(unionDefinition.Name))),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("node"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}))
		args = append(args, ast.NewIdent(fmt.Sprintf("union%s", unionDefinition.Name)))
	}

	stmts = append(stmts, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("ret"),
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun:  ast.NewIdent("append"),
				Args: args,
			},
		},
	})

	return stmts
}

func generateIntrospectionTypeFieldInterfacesCallStmts(interfaces []*schema.InterfaceDefinition, i int) []ast.Stmt {
	var interfaceStmts []ast.Stmt
	for j, interfaceDefinition := range interfaces {
		interfaceStmts = append(interfaceStmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("interfaces%d", i)),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("make"),
					Args: []ast.Expr{
						&ast.ArrayType{
							Elt: ast.NewIdent("__Type"),
						},
						ast.NewIdent("0"),
					},
				},
			},
		})
		interfaceStmts = append(interfaceStmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("interface%d%d", i, j)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(interfaceDefinition.Name))),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		},
			generateReturnErrorHandlingStmt([]ast.Expr{
				&ast.CompositeLit{
					Type: ast.NewIdent("__Type"),
					Elts: []ast.Expr{},
				},
			}),
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(fmt.Sprintf("interfaces%d", i)),
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("append"),
						Args: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("interfaces%d", i)),
							ast.NewIdent(fmt.Sprintf("interface%d%d", i, j)),
						},
					},
				},
			}, &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("ret"),
						Sel: ast.NewIdent("Interfaces"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("NewNullable"),
						},
						Args: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("interfaces%d", i)),
						},
					},
				},
			})
	}

	return interfaceStmts
}

func generateIntrospectionModelFieldCaseAST(s *schema.Schema, field *schema.FieldDefinition, index int) ast.Stmt {
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
						ast.NewIdent("ctx"),
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
	case "mutationType":
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("mutationType"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("__schema_mutationType"),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret"),
					Sel: ast.NewIdent("MutationType"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("executor"),
						Sel: ast.NewIdent("NewNullable"),
					},
					Args: []ast.Expr{
						ast.NewIdent("mutationType"),
					},
				},
			},
		})
	case "types":
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("type%d", index)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("__schema_types"),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.Ident{
					Name: "err",
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
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
		})
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret"),
					Sel: ast.NewIdent("Types"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("type%d", index)),
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
					ast.NewIdent("ctx"),
				},
				Type: &ast.SelectorExpr{
					X:   ast.NewIdent("context"),
					Sel: ast.NewIdent("Context"),
				},
			}, {
				Names: []*ast.Ident{
					ast.NewIdent("node"),
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("executor"),
						Sel: ast.NewIdent("Node"),
					},
				},
			}, {
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

func generateIntrospectionMutationTypeMethodAST(s *schema.Schema) ast.Decl {
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
		Name: ast.NewIdent("__schema_mutationType"),
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
			List: generateIntrospectionMutationTypeMethodBodyAST(s),
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
						{
							Type: ast.NewIdent("error"),
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
								generateIntrospectionTypeFieldSwitchStmt(string(field.Type.Name), field, s.Indexes),
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
			},
		})
	}

	m := s.GetMutation()
	if m == nil {
		return ret
	}

	for _, field := range m.Fields {
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
						{
							Type: ast.NewIdent("error"),
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
								generateIntrospectionTypeFieldSwitchStmt(string(field.Type.Name), field, s.Indexes),
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
			ret = append(ret, generateIntrospectionRecursiveFieldTypeOfDecls(string(field.Name), introspection.ExpandType(field.Type).Unwrap(), 0, s.Indexes)...)
		} else {
			ret = append(ret, generateIntrospectionRecursiveFieldTypeOfDecls(string(field.Name), introspection.ExpandType(field.Type).Unwrap(), 0, s.Indexes)...)
		}
	}

	for _, t := range s.Types {
		for _, field := range t.Fields {
			ret = append(ret, generateIntrospectionRecursiveFieldTypeOfDecls(fmt.Sprintf("%s__%s", t.Name, field.Name), introspection.ExpandType(field.Type).Unwrap(), 0, s.Indexes)...)

		}
	}

	for _, i := range s.Interfaces {
		for _, field := range i.Fields {
			ret = append(ret, generateIntrospectionRecursiveFieldTypeOfDecls(fmt.Sprintf("%s__%s", i.Name, field.Name), introspection.ExpandType(field.Type).Unwrap(), 0, s.Indexes)...)
		}
	}

	for _, i := range s.Inputs {
		for _, field := range i.Fields {
			ret = append(ret, generateIntrospectionRecursiveFieldTypeOfDecls(fmt.Sprintf("%s__%s", i.Name, field.Name), introspection.ExpandType(field.Type).Unwrap(), 0, s.Indexes)...)
		}
	}

	m := s.GetMutation()
	if m == nil {
		return ret
	}

	for _, field := range m.Fields {
		if field.Type.IsList && field.Type.Nullable {
			ret = append(ret, generateIntrospectionRecursiveFieldTypeOfDecls(string(field.Name), introspection.ExpandType(field.Type).Unwrap(), 0, s.Indexes)...)
		} else {
			ret = append(ret, generateIntrospectionRecursiveFieldTypeOfDecls(string(field.Name), introspection.ExpandType(field.Type).Unwrap(), 0, s.Indexes)...)
		}
	}

	return ret
}

func generateIntrospectionTypeFieldsFuncDecls(typeDefinitions []*schema.TypeDefinition, indexes *schema.Indexes) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, t := range typeDefinitions {
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
							{
								Type: ast.NewIdent("error"),
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
									generateIntrospectionTypeFieldSwitchStmt(string(t.Name), field, indexes),
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
				},
			})
		}
	}

	return ret
}

func generateIntrospectionInputFieldsFuncDecls(typeDefinitions []*schema.InputDefinition, indexes *schema.Indexes) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, t := range typeDefinitions {
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
							{
								Type: ast.NewIdent("error"),
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
						// &ast.AssignStmt{
						// 	Lhs: []ast.Expr{
						// 		ast.NewIdent("fields"),
						// 	},
						// 	Tok: token.DEFINE,
						// 	Rhs: []ast.Expr{
						// 		&ast.CallExpr{
						// 			Fun: ast.NewIdent("make"),
						// 			Args: []ast.Expr{
						// 				&ast.ArrayType{
						// 					Elt: ast.NewIdent("__Field"),
						// 				},
						// 				ast.NewIdent("0"),
						// 				ast.NewIdent(fmt.Sprintf("%d", len(t.Fields))),
						// 			},
						// 		},
						// 	},
						// },
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
									generateIntrospectionTypeFieldSwitchStmt(string(t.Name), field, indexes),
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
				},
			})
		}
	}

	return ret
}

func generateIntrospectionTypeResolverDeclsFromInterfaces(interfaces []*schema.InterfaceDefinition, indexes *schema.Indexes) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, t := range interfaces {
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
							{
								Type: ast.NewIdent("error"),
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
									generateIntrospectionTypeFieldSwitchStmt(string(t.Name), field, indexes),
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
				},
			})
		}
	}

	return ret
}

func generateIntrospectionInterfaceFieldsDecls(interfaces []*schema.InterfaceDefinition) []ast.Decl {
	ret := make([]ast.Decl, 0)
	for _, t := range interfaces {
		args := generateNodeWalkerArgs()

		args.List = append(args.List, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent("includeDeprecated"),
			},
			Type: &ast.Ident{
				Name: "bool",
			},
		})

		names := make([]*ast.Ident, 0, len(t.Fields))
		appendStmts := make([]ast.Stmt, 0, len(t.Fields))
		for _, field := range t.Fields {
			names = append(names, ast.NewIdent("field"+string(field.Name)))

			if field.Directives.Get([]byte("deprecated")) != nil {
				appendStmts = append(appendStmts, &ast.IfStmt{
					Cond: ast.NewIdent("includeDeprecated"),
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
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
											ast.NewIdent("field" + string(field.Name)),
										},
									},
								},
							},
						},
					},
				})
			} else {
				appendStmts = append(appendStmts, &ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("ret"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("append"),
							Args: []ast.Expr{
								ast.NewIdent("ret"),
								ast.NewIdent("field" + string(field.Name)),
							},
						},
					},
				})
			}
		}

		var varFieldDeclStmt ast.Stmt = &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: names,
						Type:  ast.NewIdent("__Field"),
					},
				},
			},
		}

		var body []ast.Stmt = []ast.Stmt{
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
							ast.NewIdent("0"),
							ast.NewIdent(fmt.Sprintf("%d", len(t.Fields))),
						},
					},
				},
			},
			varFieldDeclStmt,
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
						generateIntrospectionFieldSwitchStmt(string(t.Name), t.Fields),
					},
				},
			},
		}

		body = append(body, appendStmts...)
		body = append(body, &ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
				ast.NewIdent("nil"),
			},
		})

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
				Params: args,
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.ArrayType{
								Elt: ast.NewIdent("__Field"),
							},
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: body,
			},
		})
	}

	return ret
}

func generateIntrospectionFieldSwitchStmt(typeName string, fieldDefinitions schema.FieldDefinitions) ast.Stmt {
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
			List: introspection.GenerateFieldsCaseStmts(fieldDefinitions),
		},
	}
}

func generateIntrospectionInputValueSwitchStmt(typeName string, fieldDefinitions schema.FieldDefinitions) ast.Stmt {
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
			List: introspection.GenerateInputValuesCaseStmts(fieldDefinitions),
		},
	}
}

func generateIntrospectionInputFieldsDecls(inputDefinitions []*schema.InputDefinition) []ast.Decl {
	ret := make([]ast.Decl, 0)
	for _, t := range inputDefinitions {
		args := generateNodeWalkerArgs()

		args.List = append(args.List, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent("includeDeprecated"),
			},
			Type: &ast.Ident{
				Name: "bool",
			},
		})

		names := make([]*ast.Ident, 0, len(t.Fields))
		appendStmts := make([]ast.Stmt, 0, len(t.Fields))
		for _, field := range t.Fields {
			names = append(names, ast.NewIdent("field"+string(field.Name)))

			if field.Directives.Get([]byte("deprecated")) != nil {
				appendStmts = append(appendStmts, &ast.IfStmt{
					Cond: ast.NewIdent("includeDeprecated"),
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
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
											ast.NewIdent("field" + string(field.Name)),
										},
									},
								},
							},
						},
					},
				})
			} else {
				appendStmts = append(appendStmts, &ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("ret"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("append"),
							Args: []ast.Expr{
								ast.NewIdent("ret"),
								ast.NewIdent("field" + string(field.Name)),
							},
						},
					},
				})
			}
		}

		var varFieldDeclStmt ast.Stmt = &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: names,
						Type:  ast.NewIdent("__InputValue"),
					},
				},
			},
		}

		var body []ast.Stmt = []ast.Stmt{
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
								Elt: ast.NewIdent("__InputValue"),
							},
							ast.NewIdent("0"),
							ast.NewIdent(fmt.Sprintf("%d", len(t.Fields))),
						},
					},
				},
			},
			varFieldDeclStmt,
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
						generateIntrospectionInputValueSwitchStmt(string(t.Name), t.Fields),
					},
				},
			},
		}

		body = append(body, appendStmts...)
		body = append(body, &ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
				ast.NewIdent("nil"),
			},
		})

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
				Params: args,
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.ArrayType{
								Elt: ast.NewIdent("__InputValue"),
							},
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: body,
			},
		})
	}

	return ret
}

func generateIntrospectionTypeFieldsDecls(typeDefinitions []*schema.TypeDefinition) []ast.Decl {
	ret := make([]ast.Decl, 0)
	for _, t := range typeDefinitions {
		// if t.IsIntrospection() {
		// 	continue
		// }

		args := generateNodeWalkerArgs()

		args.List = append(args.List, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent("includeDeprecated"),
			},
			Type: &ast.Ident{
				Name: "bool",
			},
		})

		names := make([]*ast.Ident, 0, len(t.Fields))
		appendStmts := make([]ast.Stmt, 0, len(t.Fields))
		for _, field := range t.Fields {
			names = append(names, ast.NewIdent("field"+string(field.Name)))

			if field.Directives.Get([]byte("deprecated")) != nil {
				appendStmts = append(appendStmts, &ast.IfStmt{
					Cond: ast.NewIdent("includeDeprecated"),
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
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
											ast.NewIdent("field" + string(field.Name)),
										},
									},
								},
							},
						},
					},
				})
			} else {
				appendStmts = append(appendStmts, &ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("ret"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("append"),
							Args: []ast.Expr{
								ast.NewIdent("ret"),
								ast.NewIdent("field" + string(field.Name)),
							},
						},
					},
				})
			}
		}

		var varFieldDeclStmt ast.Stmt = &ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: names,
						Type:  ast.NewIdent("__Field"),
					},
				},
			},
		}

		var body []ast.Stmt = []ast.Stmt{
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
							ast.NewIdent("0"),
							ast.NewIdent(fmt.Sprintf("%d", len(t.Fields))),
						},
					},
				},
			},
			varFieldDeclStmt,
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
						generateIntrospectionFieldSwitchStmt(string(t.Name), t.Fields),
					},
				},
			},
		}

		body = append(body, appendStmts...)
		body = append(body, &ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
				ast.NewIdent("nil"),
			},
		})

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
				Params: args,
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.ArrayType{
								Elt: ast.NewIdent("__Field"),
							},
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: body,
			},
		})
	}

	return ret
}

func generateIntrospectionRecursiveFieldTypeOfDecls(fieldDefinitionName string, field *introspection.FieldType, nestCount int, indexes *schema.Indexes) []ast.Decl {
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
					generateIntrospectionTypeOfSwitchStmt(field, fmt.Sprintf("__schema__%s%s__typeof", fieldDefinitionName, typeOfSuffix), indexes),
				},
			},
		})

		bodyStmts = append(bodyStmts, &ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
				ast.NewIdent("nil"),
			},
		})
	} else {
		bodyStmts = append(bodyStmts, &ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: "nil",
				},
				ast.NewIdent("nil"),
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
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: bodyStmts,
		},
	})

	if field != nil && field.Child != nil {
		decls = append(decls, generateIntrospectionRecursiveFieldTypeOfDecls(fieldDefinitionName, field.Child, nestCount+1, indexes)...)
	}

	return decls
}

func generateIntrospectionTypeCaseStmts(f *introspection.FieldType, callTypeOfFuncName string, indexes *schema.Indexes) []ast.Stmt {
	if f.NonNull {
		return introspection.GenerateNonNullTypeCaseStmts(callTypeOfFuncName)
	}

	if f.IsList {
		return introspection.GenerateListTypeCaseStmts(callTypeOfFuncName)
	}

	if enumDefinition, ok := indexes.EnumIndex[string(f.Name)]; ok {
		return introspection.GenerateEnumTypeCaseStmts(enumDefinition)
	}

	if unionDefinition, ok := indexes.UnionIndex[string(f.Name)]; ok {
		return introspection.GenerateUnionTypeCaseStmts(unionDefinition)
	}

	if inputDefinition, ok := indexes.InputIndex[string(f.Name)]; ok {
		return introspection.GenerateInputTypeCaseStmts(inputDefinition)
	}

	if interfaceDefinition, ok := indexes.InterfaceIndex[string(f.Name)]; ok {
		return introspection.GenerateInterfaceTypeCaseStmts(interfaceDefinition, indexes)
	}

	if typeDefinition, ok := indexes.TypeIndex[string(f.Name)]; ok {
		return introspection.GenerateTypeObjectCaseStmts(typeDefinition)
	}

	return introspection.GenerateScalarCaseStmts(f)
}

func generateIntrospectionTypeOfSwitchStmt(f *introspection.FieldType, callTypeOfFuncName string, indexes *schema.Indexes) ast.Stmt {
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
			List: generateIntrospectionTypeCaseStmts(f, callTypeOfFuncName, indexes),
		},
	}
}

func generateIntrospectionTypeFieldSwitchStmt(typeName string, f *schema.FieldDefinition, indexes *schema.Indexes) ast.Stmt {
	if typeName == string(f.Type.Name) {
		typeName = ""
	} else {
		typeName = fmt.Sprintf("__%s", typeName)
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
			List: generateIntrospectionTypeCaseStmts(introspection.ExpandType(f.Type), fmt.Sprintf("__schema%s__%s__typeof", typeName, string(f.Name)), indexes),
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

func generateIntrospectionMutationTypeMethodBodyAST(s *schema.Schema) []ast.Stmt {
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
							List: generateMutationTypeSwitchBodyAST(s),
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
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("NewNullable"),
							},
							Args: []ast.Expr{
								ast.NewIdent(fmt.Sprintf("%q", string(s.Definition.Query))),
							},
						},
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
			Body: generateIntrospectionOperationFieldsAST(s.GetQuery(), string(s.Definition.Query)),
		},
	}
}

func generateMutationTypeSwitchBodyAST(s *schema.Schema) []ast.Stmt {
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
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("NewNullable"),
							},
							Args: []ast.Expr{
								ast.NewIdent(fmt.Sprintf("%q", string(s.Definition.Mutation))),
							},
						},
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
			Body: generateIntrospectionOperationFieldsAST(s.GetMutation(), string(s.Definition.Mutation)),
		},
	}
}

func generateIntrospectionOperationFieldsAST(fieldDefinitions *schema.OperationDefinition, operationName string) []ast.Stmt {
	if fieldDefinitions == nil {
		return []ast.Stmt{}
	}

	ret := make([]ast.Stmt, 0)
	ret = append(ret, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("fields"),
			ast.NewIdent("err"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("r"),
					Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", string(operationName))),
				},
				Args: []ast.Expr{
					ast.NewIdent("ctx"),
					ast.NewIdent("child"),
					ast.NewIdent("variables"),
				},
			},
		},
	})
	ret = append(ret, &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  ast.NewIdent("err"),
			Op: token.NEQ,
			Y:  &ast.BasicLit{Kind: token.STRING, Value: "nil"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{},
		},
	})
	ret = append(ret, &ast.AssignStmt{
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
					X:   ast.NewIdent("executor"),
					Sel: ast.NewIdent("NewNullable"),
				},
				Args: []ast.Expr{
					ast.NewIdent("fields"),
				},
			},
		},
	})

	return ret
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
		Name: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", attributeName)),
		Type: &ast.FuncType{
			Params: generateNodeWalkerArgs(),
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: ast.NewIdent("__Field"),
						},
					},
					{
						Type: ast.NewIdent("error"),
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
	if attributeName == "Query" || attributeName == "Mutation" || attributeName == "Subscription" {
		attributeName = ""
	}

	args := make([]ast.Expr, 0)
	args = append(args, ast.NewIdent("ret"))

	names := make([]*ast.Ident, 0)
	for _, field := range fields {
		names = append(names, ast.NewIdent("ret"+string(field.Name)))
		args = append(args, ast.NewIdent("ret"+string(field.Name)))
	}
	varSpec := &ast.ValueSpec{
		Names: names,
		Type:  ast.NewIdent("__Field"),
	}

	appendStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("ret"),
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun:  ast.NewIdent("append"),
				Args: args,
			},
		},
	}

	rangeBodyStmts := []ast.Stmt{}

	if fields.HasDeprecatedDirective() {
		rangeBodyStmts = append(rangeBodyStmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("includeDeprecated"),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("extract__fieldsArgs"),
					},
					Args: []ast.Expr{
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		rangeBodyStmts = append(rangeBodyStmts, generateReturnErrorHandlingStmt([]ast.Expr{
			ast.NewIdent("nil"),
		}))
	}

	rangeBodyStmts = append(rangeBodyStmts, &ast.SwitchStmt{
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
					Body: generateIntrospectionSchemaFieldTypeBodyStmt(attributeName, fields),
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"args"`,
						},
					},
					Body: generateArgsAssignStmt(fields),
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
	})

	return []ast.Stmt{
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					varSpec,
				},
			},
		},
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
						ast.NewIdent("0"),
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
				List: rangeBodyStmts,
			},
		},
		appendStmt,
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
				ast.NewIdent("nil"),
			},
		},
	}
}

func generateArgsAssignStmt(fields schema.FieldDefinitions) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(fields))
	for _, field := range fields {
		if len(field.Arguments) == 0 {
			continue
		}

		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("ret" + string(field.Name) + "Args"),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__args", string(field.Name))),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		})

		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{
			ast.NewIdent("nil"),
		}))

		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret" + string(field.Name)),
					Sel: ast.NewIdent("Args"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent("ret" + string(field.Name) + "Args"),
			},
		})
	}

	return stmts
}

func generateIntrospectionFieldNameBodyStmt(fields schema.FieldDefinitions) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(fields))
	for _, field := range fields {
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent(fmt.Sprintf("ret%s", field.Name)),
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

func generateReturnErrorHandlingStmt(prefixReturnExpr []ast.Expr) ast.Stmt {
	prefixReturnExpr = append(prefixReturnExpr, ast.NewIdent("err"))
	return &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  ast.NewIdent("err"),
			Op: token.NEQ,
			Y:  &ast.BasicLit{Kind: token.STRING, Value: "nil"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: prefixReturnExpr,
				},
			},
		},
	}
}

func generateIntrospectionSchemaFieldTypeBodyStmt(attributeName string, fields schema.FieldDefinitions) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(fields))
	for i, field := range fields {
		prefix := fmt.Sprintf("__schema__%s__%s", attributeName, string(field.Name))
		if attributeName == "" {
			prefix = fmt.Sprintf("__schema__%s", string(field.Name))
		}

		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(string(fmt.Sprintf("t%d", i))),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("%s__type", prefix)),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		})

		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{
			ast.NewIdent("nil"),
		}))

		assignStmt := &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret" + string(field.Name)),
					Sel: ast.NewIdent("Type"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent(string(fmt.Sprintf("t%d", i))),
			},
		}

		if field.Directives.Get([]byte("deprecated")) != nil {
			stmts = append(stmts, &ast.IfStmt{
				Cond: &ast.StarExpr{
					X: ast.NewIdent("includeDeprecated"),
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						assignStmt,
					},
				},
			})
		} else {
			stmts = append(stmts, assignStmt)
		}
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
				ast.NewIdent(string(fmt.Sprintf("t%d", i))),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("%s__type", prefix)),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		})

		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{
			ast.NewIdent("nil"),
		}))

		assignStmt := &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("field" + string(field.Name)),
					Sel: ast.NewIdent("Type"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent(string(fmt.Sprintf("t%d", i))),
			},
		}

		if field.Directives.Get([]byte("deprecated")) != nil {
			stmts = append(stmts, &ast.IfStmt{
				Cond: ast.NewIdent("includeDeprecated"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						assignStmt,
					},
				},
			})
		} else {
			stmts = append(stmts, assignStmt)
		}
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

func generateStringPointerExpr(value ast.Expr) ast.Expr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.IndexExpr{
			X: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: ast.NewIdent("string"),
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

func generateBoolPointerExpr(value ast.Expr) ast.Expr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.IndexExpr{
			X: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: ast.NewIdent("bool"),
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

func generateBoolPointerAST(value string) ast.Expr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.IndexExpr{
			X: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: ast.NewIdent("bool"),
				},
				Elts: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
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

func generateModelFieldCaseASTs(s *schema.Schema, fields []*schema.FieldDefinition) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(fields))
	for i, f := range fields {
		stmts = append(stmts, generateIntrospectionModelFieldCaseAST(s, f, i))
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

	body = append(body, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("ret"),
			ast.NewIdent("nil"),
		},
	})
	params := generateNodeWalkerArgs()

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
			Params: params,
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: ast.NewIdent("__Schema"),
						},
					}, {
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

func generateSchemaErrorResponseWrite() ast.Stmt {
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
									&ast.KeyValueExpr{
										Key: ast.NewIdent("Errors"),
										Value: &ast.CompositeLit{
											Type: &ast.ArrayType{
												Elt: ast.NewIdent("error"),
											},
											Elts: []ast.Expr{
												ast.NewIdent("err"),
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
				&ast.ReturnStmt{},
			},
		},
	}
}

func generateIntrospectionTypeFuncDecl(s *schema.Schema) ast.Decl {
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
		Name: ast.NewIdent("__type"),
		Type: &ast.FuncType{
			Params: generateNodeWalkerArgs(),
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: ast.NewIdent("__Type"),
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: generateIntrospectionTypeFuncDeclBody(s),
		},
	}
}

func generateIntrospectionTypeFuncDeclBody(s *schema.Schema) []ast.Stmt {
	ret := make([]ast.Stmt, 0)
	ret = append(ret, &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X: &ast.CallExpr{
				Fun: &ast.Ident{Name: "len"},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   &ast.Ident{Name: "node"},
						Sel: &ast.Ident{Name: "Arguments"},
					},
				},
			},
			Op: token.EQL,
			Y:  &ast.BasicLit{Kind: token.INT, Value: "0"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{Name: "nil"},
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "fmt"},
								Sel: &ast.Ident{Name: "Errorf"},
							},
							Args: []ast.Expr{
								&ast.BasicLit{Kind: token.STRING, Value: `"type must have a name argument"`},
							},
						},
					},
				},
			},
		},
	},
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{Name: "arg"},
				&ast.Ident{Name: "err"},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "r"},
								Sel: &ast.Ident{Name: "parser"},
							},
							Sel: &ast.Ident{Name: "ValueParser"},
						},
						Sel: &ast.Ident{Name: "Parse"},
					},
					Args: []ast.Expr{
						&ast.SelectorExpr{
							X: &ast.IndexExpr{
								X: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "node"},
									Sel: &ast.Ident{Name: "Arguments"},
								},
								Index: &ast.BasicLit{Kind: token.INT, Value: "0"},
							},
							Sel: &ast.Ident{Name: "Value"},
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  &ast.Ident{Name: "err"},
				Op: token.NEQ,
				Y:  &ast.Ident{Name: "nil"},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{Name: "nil"},
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "fmt"},
									Sel: &ast.Ident{Name: "Errorf"},
								},
								Args: []ast.Expr{
									&ast.BasicLit{Kind: token.STRING, Value: `"failed to parse type name argument: %w"`},
									&ast.Ident{Name: "err"},
								},
							},
						},
					},
				},
			},
		})

	ret = append(ret, generateIntrospectionTypeFuncDeclBodySwitchStmt(s.Types, s.Interfaces, s.Inputs, s.Scalars, s.Enums, s.Unions))
	ret = append(ret, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
			ast.NewIdent("nil"),
		},
	})

	return ret
}

func generateIntrospectionTypeFuncDeclBodySwitchStmt(typeDefinitions schema.TypeDefinitions, interfaceDefinitions []*schema.InterfaceDefinition, inputDefinitions []*schema.InputDefinition, scalarDefinitions []*schema.ScalarDefinition, enumDefinitions []*schema.EnumDefinition, unionDefinitions schema.UnionDefinitions) ast.Stmt {
	caseStmts := func(typeDefinitions schema.TypeDefinitions) []ast.Stmt {
		stmts := make([]ast.Stmt, 0)
		stmts = append(stmts, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"\"Query\""`,
				},
			},
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
								Sel: ast.NewIdent("__schema__Query__type"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("node"),
								ast.NewIdent("variables"),
							},
						},
					},
				},
				generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent("ret"),
						},
						ast.NewIdent("nil"),
					},
				},
			},
		})

		stmts = append(stmts, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"\"Mutation\""`,
				},
			},
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
								Sel: ast.NewIdent("__schema__Mutation__type"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("node"),
								ast.NewIdent("variables"),
							},
						},
					},
				},
				generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent("ret"),
						},
						ast.NewIdent("nil"),
					},
				},
			},
		})

		for _, t := range typeDefinitions {
			stmts = append(stmts, &ast.CaseClause{
				List: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"\"%s\""`, string(t.Name)),
					},
				},
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
									Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(t.Name))),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("node"),
									ast.NewIdent("variables"),
								},
							},
						},
					},
					generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("ret"),
							},
							ast.NewIdent("nil"),
						},
					},
				},
			})
		}

		for _, i := range interfaceDefinitions {
			stmts = append(stmts, &ast.CaseClause{
				List: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"\"%s\""`, string(i.Name)),
					},
				},
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
									Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(i.Name))),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("node"),
									ast.NewIdent("variables"),
								},
							},
						},
					},
					generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("ret"),
							},
							ast.NewIdent("nil"),
						},
					},
				},
			})
		}

		for _, i := range inputDefinitions {
			stmts = append(stmts, &ast.CaseClause{
				List: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"\"%s\""`, string(i.Name)),
					},
				},
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
									Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(i.Name))),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("node"),
									ast.NewIdent("variables"),
								},
							},
						},
					},
					generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("ret"),
							},
							ast.NewIdent("nil"),
						},
					},
				},
			})
		}

		for _, s := range scalarDefinitions {
			stmts = append(stmts, &ast.CaseClause{
				List: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"\"%s\""`, string(s.Name)),
					},
				},
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
									Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(s.Name))),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("node"),
									ast.NewIdent("variables"),
								},
							},
						},
					},
					generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("ret"),
							},
							ast.NewIdent("nil"),
						},
					},
				},
			})
		}

		for _, s := range enumDefinitions {
			stmts = append(stmts, &ast.CaseClause{
				List: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"\"%s\""`, string(s.Name)),
					},
				},
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
									Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(s.Name))),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("node"),
									ast.NewIdent("variables"),
								},
							},
						},
					},
					generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("ret"),
							},
							ast.NewIdent("nil"),
						},
					},
				},
			})
		}

		for _, u := range unionDefinitions {
			stmts = append(stmts, &ast.CaseClause{
				List: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"\"%s\""`, string(u.Name)),
					},
				},
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
									Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(u.Name))),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("node"),
									ast.NewIdent("variables"),
								},
							},
						},
					},
					generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("ret"),
							},
							ast.NewIdent("nil"),
						},
					},
				},
			})
		}

		return stmts
	}

	return &ast.TypeSwitchStmt{
		Assign: &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("val"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.TypeAssertExpr{
					X:    ast.NewIdent("arg"),
					Type: ast.NewIdent("type"),
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("goliteql"),
								Sel: ast.NewIdent("ValueParserLiteral"),
							},
						},
					},
					Body: []ast.Stmt{
						&ast.SwitchStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("name"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("val"),
											Sel: ast.NewIdent("StringValue"),
										},
									},
								},
							},
							Tag: ast.NewIdent("name"),
							Body: &ast.BlockStmt{
								List: caseStmts(typeDefinitions),
							},
						},
					},
				},
				&ast.CaseClause{
					List: nil, // default case
					Body: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("nil"),
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("fmt"),
										Sel: ast.NewIdent("Errorf"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"type must be a string, got %T"`,
										},
										ast.NewIdent("val"),
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

func generateIntrospectionTypeFuncDecls(typeDefinitions schema.TypeDefinitions) []ast.Decl {
	ret := make([]ast.Decl, 0, len(typeDefinitions))
	for i, t := range typeDefinitions {
		interfacesStmts := generateIntrospectionTypeFieldInterfacesCallStmts(t.Interfaces, i)
		body := []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{
								ast.NewIdent("ret"),
							},
							Type: ast.NewIdent("__Type"),
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
											ast.NewIdent(`"name"`),
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
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent(fmt.Sprintf("%q", string(t.Name))),
														},
													},
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
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("includeDeprecated"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("r"),
															Sel: ast.NewIdent("extract__fieldsArgs"),
														},
														Args: []ast.Expr{
															ast.NewIdent("child"),
															ast.NewIdent("variables"),
														},
													},
												},
											},
											generateReturnErrorHandlingStmt([]ast.Expr{
												&ast.CompositeLit{
													Type: ast.NewIdent("__Type"),
												},
											}),
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("fields"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("r"),
															Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", string(t.Name))),
														},
														Args: []ast.Expr{
															ast.NewIdent("ctx"),
															ast.NewIdent("child"),
															ast.NewIdent("variables"),
															&ast.StarExpr{
																X: ast.NewIdent("includeDeprecated"),
															},
														},
													},
												},
											},
											generateReturnErrorHandlingStmt([]ast.Expr{
												&ast.CompositeLit{
													Type: ast.NewIdent("__Type"),
												},
											}),
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "ret",
														},
														Sel: ast.NewIdent("Fields"),
													},
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent("fields"),
														},
													},
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
										Body: interfacesStmts,
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"enumValues"`),
										},
										Body: []ast.Stmt{
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("ret"),
														Sel: ast.NewIdent("EnumValues"),
													},
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															&ast.CompositeLit{
																Type: &ast.ArrayType{
																	Elt: ast.NewIdent("__EnumValue"),
																},
															},
														},
													},
												},
											},
										},
									}, &ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"inputFields"`),
										},
										Body: []ast.Stmt{
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("ret"),
														Sel: ast.NewIdent("InputFields"),
													},
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															&ast.CompositeLit{
																Type: &ast.ArrayType{
																	Elt: ast.NewIdent("__InputValue"),
																},
															},
														},
													},
												},
											},
										},
									}, &ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"possibleTypes"`),
										},
										Body: []ast.Stmt{
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("ret"),
														Sel: ast.NewIdent("PossibleTypes"),
													},
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															&ast.CompositeLit{
																Type: &ast.ArrayType{
																	Elt: ast.NewIdent("__Type"),
																},
															},
														},
													},
												},
											},
										},
									}, &ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"ogType"`),
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
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent("nil"),
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
			&ast.ReturnStmt{
				Results: []ast.Expr{
					ast.NewIdent("ret"),
					ast.NewIdent("nil"),
				},
			},
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
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(t.Name))),
			Type: &ast.FuncType{
				Params: generateNodeWalkerArgs(),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("__Type"),
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: body,
			},
		})
	}

	return ret
}

func generateIntrospectionInputFuncDecls(inputDefinitions []*schema.InputDefinition) []ast.Decl {
	ret := make([]ast.Decl, 0, len(inputDefinitions))
	for _, t := range inputDefinitions {
		body := []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{
								ast.NewIdent("ret"),
							},
							Type: ast.NewIdent("__Type"),
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
													ast.NewIdent("__TypeKind_INPUT_OBJECT"),
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"name"`),
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
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent(fmt.Sprintf("%q", string(t.Name))),
														},
													},
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
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("includeDeprecated"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("r"),
															Sel: ast.NewIdent("extract__fieldsArgs"),
														},
														Args: []ast.Expr{
															ast.NewIdent("child"),
															ast.NewIdent("variables"),
														},
													},
												},
											},
											generateReturnErrorHandlingStmt([]ast.Expr{
												&ast.CompositeLit{
													Type: ast.NewIdent("__Type"),
												},
											}),
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("fields"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("r"),
															Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", string(t.Name))),
														},
														Args: []ast.Expr{
															ast.NewIdent("ctx"),
															ast.NewIdent("child"),
															ast.NewIdent("variables"),
															&ast.StarExpr{
																X: ast.NewIdent("includeDeprecated"),
															},
														},
													},
												},
											},
											generateReturnErrorHandlingStmt([]ast.Expr{
												&ast.CompositeLit{
													Type: ast.NewIdent("__Type"),
												},
											}),
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "ret",
														},
														Sel: ast.NewIdent("InputFields"),
													},
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent("fields"),
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
			&ast.ReturnStmt{
				Results: []ast.Expr{
					ast.NewIdent("ret"),
					ast.NewIdent("nil"),
				},
			},
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
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(t.Name))),
			Type: &ast.FuncType{
				Params: generateNodeWalkerArgs(),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("__Type"),
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: body,
			},
		})
	}

	return ret
}

func generateIntrospectionScalarFuncBodySwitchStmt(t *schema.ScalarDefinition) ast.Stmt {
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
								ast.NewIdent("__TypeKind_SCALAR"),
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						ast.NewIdent(`"name"`),
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
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("executor"),
										Sel: ast.NewIdent("NewNullable"),
									},
									Args: []ast.Expr{
										ast.NewIdent(fmt.Sprintf("%q", string(t.Name))),
									},
								},
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						ast.NewIdent(`"description"`),
					},
					Body: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("Description"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("executor"),
										Sel: ast.NewIdent("NewNullable"),
									},
									Args: []ast.Expr{
										ast.NewIdent(fmt.Sprintf("%q", string(t.Name))),
									},
								},
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						ast.NewIdent(`"fields"`),
					},
					Body: []ast.Stmt{
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
										X:   ast.NewIdent("executor"),
										Sel: ast.NewIdent("NewNullable"),
									},
									Args: []ast.Expr{
										&ast.CompositeLit{
											Type: &ast.ArrayType{
												Elt: ast.NewIdent("__Field"),
											},
											Elts: []ast.Expr{},
										},
									},
								},
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						ast.NewIdent(`"interfaces"`),
					},
					Body: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("Interfaces"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("executor"),
										Sel: ast.NewIdent("NewNullable"),
									},
									Args: []ast.Expr{
										&ast.CompositeLit{
											Type: &ast.ArrayType{
												Elt: ast.NewIdent("__Type"),
											},
											Elts: []ast.Expr{},
										},
									},
								},
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						ast.NewIdent(`"possibleTypes"`),
					},
					Body: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("PossibleTypes"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("executor"),
										Sel: ast.NewIdent("NewNullable"),
									},
									Args: []ast.Expr{
										&ast.CompositeLit{
											Type: &ast.ArrayType{
												Elt: ast.NewIdent("__Type"),
											},
											Elts: []ast.Expr{},
										},
									},
								},
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						ast.NewIdent(`"enumValues"`),
					},
					Body: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("EnumValues"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("executor"),
										Sel: ast.NewIdent("NewNullable"),
									},
									Args: []ast.Expr{
										&ast.CompositeLit{
											Type: &ast.ArrayType{
												Elt: ast.NewIdent("__EnumValue"),
											},
											Elts: []ast.Expr{},
										},
									},
								},
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						ast.NewIdent(`"inputFields"`),
					},
					Body: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("InputFields"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("executor"),
										Sel: ast.NewIdent("NewNullable"),
									},
									Args: []ast.Expr{
										&ast.CompositeLit{
											Type: &ast.ArrayType{
												Elt: ast.NewIdent("__InputValue"),
											},
											Elts: []ast.Expr{},
										},
									},
								},
							},
						},
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						ast.NewIdent(`"ofType"`),
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
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("executor"),
										Sel: ast.NewIdent("NewNullable"),
									},
									Args: []ast.Expr{
										ast.NewIdent("nil"),
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

func generateIntrospectionScalarFuncDecls(scalarDefinitions []*schema.ScalarDefinition) []ast.Decl {
	ret := make([]ast.Decl, 0, len(scalarDefinitions))

	for _, t := range scalarDefinitions {
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
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(t.Name))),
			Type: &ast.FuncType{
				Params: generateNodeWalkerArgs(),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("__Type"),
						},
						{
							Type: ast.NewIdent("error"),
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
								generateIntrospectionScalarFuncBodySwitchStmt(t),
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
			},
		})
	}

	return ret
}

func generateIntrospectionOperationFuncDecls(s *schema.Schema) []ast.Decl {
	ret := make([]ast.Decl, 0, len(s.Operations))

	q := s.GetQuery()
	if q != nil {
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
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(s.Definition.Query))),
			Type: &ast.FuncType{
				Params: generateNodeWalkerArgs(),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("__Type"),
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("fields"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent("make"),
								Args: []ast.Expr{
									&ast.ArrayType{
										Elt: ast.NewIdent("__Field"),
									},
									ast.NewIdent("0"),
									ast.NewIdent(fmt.Sprintf("%d", len(q.Fields))),
								},
							},
						},
					},
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
								generateIntrospectionOperationSwitchStmt(string(s.Definition.Query), q),
							},
						},
					},
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
									X:   ast.NewIdent("executor"),
									Sel: ast.NewIdent("NewNullable"),
								},
								Args: []ast.Expr{
									ast.NewIdent("fields"),
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
			},
		})
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
		Name: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(s.Definition.Mutation))),
		Type: &ast.FuncType{
			Params: generateNodeWalkerArgs(),
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("__Type"),
					},
					{
						Type: ast.NewIdent("error"),
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
				}, &ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("fields"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("make"),
							Args: []ast.Expr{
								&ast.ArrayType{
									Elt: ast.NewIdent("__Field"),
								},
								ast.NewIdent("0"),
								ast.NewIdent(fmt.Sprintf("%d", len(q.Fields))),
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
							generateIntrospectionOperationSwitchStmt(string(s.Definition.Mutation), s.GetMutation()),
						},
					},
				},
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
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("NewNullable"),
							},
							Args: []ast.Expr{
								ast.NewIdent("fields"),
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
		},
	})

	return ret
}

func generateIntrospectionOperationSwitchStmt(operationName string, op *schema.OperationDefinition) ast.Stmt {
	caseStmts := make([]ast.Stmt, 0, len(op.Fields)+1)
	caseStmts = append(caseStmts, &ast.CaseClause{
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
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("NewNullable"),
						},
						Args: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("%q", string(operationName))),
						},
					},
				},
			},
		},
	})

	caseStmts = append(caseStmts, &ast.CaseClause{
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
	})

	fieldBodyStmts := make([]ast.Stmt, 0)

	if op.Fields.HasDeprecatedDirective() {
		fieldBodyStmts = append(fieldBodyStmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("includeDeprecated"),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("extract__fieldsArgs"),
					},
					Args: []ast.Expr{
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
	}

	for _, f := range op.Fields {
		if f.Directives.Get([]byte("deprecated")) != nil {
			ifStmt := &ast.IfStmt{
				Cond: &ast.StarExpr{
					X: ast.NewIdent("includeDeprecated"),
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent(fmt.Sprintf("%sfield", f.Name)),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("r"),
										Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", string(f.Name))),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										ast.NewIdent("child"),
										ast.NewIdent("variables"),
									},
								},
							},
						},
						generateReturnErrorHandlingStmt([]ast.Expr{
							&ast.CompositeLit{
								Type: ast.NewIdent("__Type"),
								Elts: []ast.Expr{},
							},
						}), &ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("fields"),
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("append"),
									Args: []ast.Expr{
										ast.NewIdent("fields"),
										ast.NewIdent(fmt.Sprintf("%sfield", f.Name)),
									},
								},
							},
						},
					},
				},
			}
			fieldBodyStmts = append(fieldBodyStmts, ifStmt)
		} else {
			fieldBodyStmts = append(fieldBodyStmts, &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(fmt.Sprintf("%sfield", f.Name)),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", string(f.Name))),
						},
						Args: []ast.Expr{
							ast.NewIdent("ctx"),
							ast.NewIdent("child"),
							ast.NewIdent("variables"),
						},
					},
				},
			})
			fieldBodyStmts = append(fieldBodyStmts, generateReturnErrorHandlingStmt([]ast.Expr{
				&ast.CompositeLit{
					Type: ast.NewIdent("__Type"),
					Elts: []ast.Expr{},
				},
			}))
			fieldBodyStmts = append(fieldBodyStmts, &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("fields"),
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("append"),
						Args: []ast.Expr{
							ast.NewIdent("fields"),
							ast.NewIdent(fmt.Sprintf("%sfield", f.Name)),
						},
					},
				},
			})
		}
	}

	caseStmts = append(caseStmts, &ast.CaseClause{
		List: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: `"fields"`,
			},
		},
		Body: fieldBodyStmts,
	})

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
			List: caseStmts,
		},
	}
}

func generateIntrospectionOperationFieldFuncDecls(operationDefinition *schema.OperationDefinition) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, f := range operationDefinition.Fields {
		caseStmts := make([]ast.Stmt, 0)
		caseStmts = append(caseStmts, &ast.CaseClause{
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
						ast.NewIdent(fmt.Sprintf(`"%s"`, string(f.Name))),
					},
				},
			},
		})

		caseStmts = append(caseStmts, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"type"`,
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(fmt.Sprintf("%sType", f.Name)),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(f.Name))),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("child"),
								ast.NewIdent("variables"),
							},
						},
					},
				},
				generateReturnErrorHandlingStmt([]ast.Expr{
					&ast.CompositeLit{
						Type: ast.NewIdent("__Field"),
					},
				}),
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("Type"),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent(fmt.Sprintf("%sType", f.Name)),
					},
				},
			},
		})

		if len(f.Arguments) > 0 {
			caseStmts = append(caseStmts, &ast.CaseClause{
				List: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: `"args"`,
					},
				},
				Body: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("%sArgs", f.Name)),
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("r"),
									Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__args", string(f.Name))),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("child"),
									ast.NewIdent("variables"),
								},
							},
						},
					},
					generateReturnErrorHandlingStmt([]ast.Expr{
						&ast.CompositeLit{
							Type: ast.NewIdent("__Field"),
						},
					}),
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("ret"),
								Sel: ast.NewIdent("Args"),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("%sArgs", f.Name)),
						},
					},
				},
			})
		}

		stmts := make([]ast.Stmt, 0)
		stmts = append(stmts, &ast.SwitchStmt{
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
				List: caseStmts,
			},
		})

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
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", string(f.Name))),
			Type: &ast.FuncType{
				Params: generateNodeWalkerArgs(),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("__Field"),
						},
						{
							Type: ast.NewIdent("error"),
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
								Type: ast.NewIdent("__Field"),
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
							List: stmts,
						},
					},

					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("ret"),
							ast.NewIdent("nil"),
						},
					},
				},
			},
		})

		if len(f.Arguments) > 0 {
			ret = append(ret, generateIntrospectionOperationArgsFuncDecl(f))
		}
	}

	return ret
}

func generateIntrospectionOperationArgsFuncDecl(fieldDefinition *schema.FieldDefinition) ast.Decl {
	stmts := make([]ast.Stmt, 0)
	stmts = append(stmts, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("ret"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: ast.NewIdent("make"),
				Args: []ast.Expr{
					&ast.ArrayType{
						Elt: ast.NewIdent("__InputValue"),
					},
					ast.NewIdent("0"),
					ast.NewIdent(fmt.Sprintf("%d", len(fieldDefinition.Arguments))),
				},
			},
		},
	})

	appendExprs := make([]ast.Expr, 0, len(fieldDefinition.Arguments)+1)
	appendExprs = append(appendExprs, ast.NewIdent("ret"))
	valueSpecs := make([]ast.Spec, 0, len(fieldDefinition.Arguments))
	nameAssignStmts := make([]ast.Stmt, 0, len(fieldDefinition.Arguments))
	descriptionStmts := make([]ast.Stmt, 0, len(fieldDefinition.Arguments))
	typeStmts := make([]ast.Stmt, 0, len(fieldDefinition.Arguments))
	defaultValueStmts := make([]ast.Stmt, 0, len(fieldDefinition.Arguments))
	for _, arg := range fieldDefinition.Arguments {
		valueSpecs = append(valueSpecs, &ast.ValueSpec{
			Names: []*ast.Ident{
				ast.NewIdent(fmt.Sprintf("%sRet", string(arg.Name))),
			},
			Type: ast.NewIdent("__InputValue"),
		})
		appendExprs = append(appendExprs, ast.NewIdent(fmt.Sprintf("%sRet", string(arg.Name))))
		nameAssignStmts = append(nameAssignStmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent(fmt.Sprintf("%sRet", string(arg.Name))),
					Sel: ast.NewIdent("Name"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf(`"%s"`, string(arg.Name))),
			},
		})
		descriptionStmts = append(descriptionStmts, &ast.ExprStmt{
			X: ast.NewIdent("// TODO"),
		})
		typeStmts = append(typeStmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sType", string(arg.Name))),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(arg.Type.Name))),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		}, generateReturnErrorHandlingStmt([]ast.Expr{
			ast.NewIdent("nil"),
		}), &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent(fmt.Sprintf("%sRet", string(arg.Name))),
					Sel: ast.NewIdent("Type"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sType", string(arg.Name))),
			},
		})

		if arg.Default != nil {
			defaultValueStmts = append(defaultValueStmts, &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent(fmt.Sprintf("%sRet", string(arg.Name))),
						Sel: ast.NewIdent("DefaultValue"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("NewNullable"),
						},
						Args: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("%q", string(arg.Default))),
						},
					},
				},
			})
		} else {
			defaultValueStmts = append(defaultValueStmts, &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent(fmt.Sprintf("%sRet", string(arg.Name))),
						Sel: ast.NewIdent("DefaultValue"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					ast.NewIdent("nil"),
				},
			})
		}
	}

	stmts = append(stmts, &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok:   token.VAR,
			Specs: valueSpecs,
		},
	})

	stmts = append(stmts, &ast.RangeStmt{
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
								Body: nameAssignStmts,
							},
							&ast.CaseClause{
								List: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: `"description"`,
									},
								},
								Body: descriptionStmts,
							},
							&ast.CaseClause{
								List: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: `"type"`,
									},
								},
								Body: typeStmts,
							},
							&ast.CaseClause{
								List: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: `"defaultValue"`,
									},
								},
								Body: defaultValueStmts,
							},
						},
					},
				},
			},
		},
	})

	stmts = append(stmts, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("ret"),
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun:  ast.NewIdent("append"),
				Args: appendExprs,
			},
		},
	})

	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("ret"),
			ast.NewIdent("nil"),
		},
	})

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
		Name: ast.NewIdent(fmt.Sprintf("__schema__%s__args", string(fieldDefinition.Name))),
		Type: &ast.FuncType{
			Params: generateNodeWalkerArgs(),
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: ast.NewIdent("__InputValue"),
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}
}

func generateIntrospectionEnumFuncDecls(enumDefinitions []*schema.EnumDefinition) []ast.Decl {
	ret := make([]ast.Decl, 0, len(enumDefinitions))

	for _, t := range enumDefinitions {
		body := []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("ret"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CompositeLit{
						Type: ast.NewIdent("__Type"),
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
													ast.NewIdent("__TypeKind_ENUM"),
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"name"`),
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
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent(fmt.Sprintf(`"%s"`, string(t.Name))),
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"description"`),
										},
										Body: []ast.Stmt{
											&ast.ExprStmt{
												X: &ast.BasicLit{
													Kind:  token.STRING,
													Value: `// TODO`,
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent(`"enumValues"`),
										},
										Body: []ast.Stmt{
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("includeDeprecated"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("r"),
															Sel: ast.NewIdent("extract__fieldsArgs"),
														},
														Args: []ast.Expr{
															ast.NewIdent("child"),
															ast.NewIdent("variables"),
														},
													},
												},
											},
											generateReturnErrorHandlingStmt([]ast.Expr{
												&ast.CompositeLit{
													Type: ast.NewIdent("__Type"),
													Elts: []ast.Expr{},
												},
											}),
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													ast.NewIdent("enumValues"),
													ast.NewIdent("err"),
												},
												Tok: token.DEFINE,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("r"),
															Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__enumValues", string(t.Name))),
														},
														Args: []ast.Expr{
															ast.NewIdent("ctx"),
															ast.NewIdent("child"),
															&ast.StarExpr{
																X: ast.NewIdent("includeDeprecated"),
															},
														},
													},
												},
											},
											generateReturnErrorHandlingStmt([]ast.Expr{
												&ast.CompositeLit{
													Type: ast.NewIdent("__Type"),
													Elts: []ast.Expr{},
												},
											}),
											&ast.AssignStmt{
												Lhs: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("ret"),
														Sel: ast.NewIdent("EnumValues"),
													},
												},
												Tok: token.ASSIGN,
												Rhs: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("executor"),
															Sel: ast.NewIdent("NewNullable"),
														},
														Args: []ast.Expr{
															ast.NewIdent("enumValues"),
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

		body = append(body, &ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
				ast.NewIdent("nil"),
			},
		})

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
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(t.Name))),
			Type: &ast.FuncType{
				Params: generateNodeWalkerArgs(),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("__Type"),
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: body,
			},
		})
	}

	return ret
}

func generateIntrospectionEnumValuesFuncDecl(enums schema.EnumDefinitions) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, enum := range enums {
		ret = append(ret, generateIntrospectionEnumValuesFieldFuncDecls(enum)...)

		bodyStmt := make([]ast.Stmt, 0)
		bodyStmt = append(bodyStmt, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("ret"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("make"),
					Args: []ast.Expr{
						&ast.ArrayType{
							Elt: ast.NewIdent("__EnumValue"),
						},
						ast.NewIdent("0"),
						ast.NewIdent(fmt.Sprintf("%d", len(enum.Values))),
					},
				},
			},
		})

		for _, value := range enum.Values {
			bodyStmt = append(bodyStmt, &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(fmt.Sprintf("%sRet", value.Name)),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__%s__enumValue", enum.Name, value.Name)),
						},
						Args: []ast.Expr{
							ast.NewIdent("ctx"),
							ast.NewIdent("node"),
							ast.NewIdent("includeDeprecated"),
						},
					},
				},
			}, generateReturnErrorHandlingStmt([]ast.Expr{
				ast.NewIdent("nil"),
			}))

			if value.Directives.Get([]byte("deprecated")) != nil {
				bodyStmt = append(bodyStmt, &ast.IfStmt{
					Cond: ast.NewIdent("includeDeprecated"),
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
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
											ast.NewIdent(fmt.Sprintf("%sRet", value.Name)),
										},
									},
								},
							},
						},
					},
				})
			} else {
				bodyStmt = append(bodyStmt, &ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("ret"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("append"),
							Args: []ast.Expr{
								ast.NewIdent("ret"),
								ast.NewIdent(fmt.Sprintf("%sRet", value.Name)),
							},
						},
					},
				})
			}
		}

		bodyStmt = append(bodyStmt, &ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
				ast.NewIdent("nil"),
			},
		})

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
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__enumValues", enum.Name)),
			Body: &ast.BlockStmt{
				List: bodyStmt,
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("ctx"),
							},
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("Context"),
							},
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
						{
							Names: []*ast.Ident{
								ast.NewIdent("includeDeprecated"),
							},
							Type: ast.NewIdent("bool"),
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.ArrayType{
								Elt: ast.NewIdent("__EnumValue"),
							},
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

func generateIntrospectionEnumValuesFieldFuncDecls(enum *schema.EnumDefinition) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, elm := range enum.Values {
		var isDeprecatedStmts, deprecatedReasonStmts []ast.Stmt
		isDeprecatedStmts = []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("ret"),
						Sel: ast.NewIdent("IsDeprecated"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					ast.NewIdent("false"),
				},
			},
		}

		deprecatedReasonStmts = []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("ret"),
						Sel: ast.NewIdent("DeprecationReason"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("NewNullable"),
						},
						Args: []ast.Expr{
							ast.NewIdent("nil"),
						},
					},
				},
			},
		}

		deprecated := elm.Directives.Get([]byte("deprecated"))
		if deprecated != nil {
			isDeprecatedStmts = []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("IsDeprecated"),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent("true"),
					},
				},
			}

			reason := string(deprecated.Arguments[0].Value)
			deprecatedReasonStmts = []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("DeprecationReason"),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("NewNullable"),
							},
							Args: []ast.Expr{
								ast.NewIdent(fmt.Sprintf("%q", string(reason))),
							},
						},
					},
				},
			}
		}

		var nameBodyStmt ast.Stmt = &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret"),
					Sel: ast.NewIdent("Name"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf(`"%s"`, strings.ReplaceAll(strings.ReplaceAll(string(elm.Name), "__TypeKind_", ""), "__DirectiveLocation_", ""))),
			},
		}

		var descriptionStmt ast.Stmt = &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret"),
					Sel: ast.NewIdent("Description"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent("nil"),
			},
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
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__%s__enumValue", enum.Name, elm.Name)),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("ctx"),
							},
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("Context"),
							},
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
						{
							Names: []*ast.Ident{
								ast.NewIdent("includeDeprecated"),
							},
							Type: ast.NewIdent("bool"),
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("__EnumValue"),
						},
						{
							Type: ast.NewIdent("error"),
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
								Type: ast.NewIdent("__EnumValue"),
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
										List: []ast.Stmt{
											&ast.CaseClause{
												List: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"name"`,
													},
												},
												Body: []ast.Stmt{
													nameBodyStmt,
												},
											}, &ast.CaseClause{
												List: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"description"`,
													},
												},
												Body: []ast.Stmt{
													descriptionStmt,
												},
											}, &ast.CaseClause{
												List: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"isDeprecated"`,
													},
												},
												Body: isDeprecatedStmts,
											}, &ast.CaseClause{
												List: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"deprecationReason"`,
													},
												},
												Body: deprecatedReasonStmts,
											},
										},
									},
								},
							},
						},
					}, &ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("ret"),
							ast.NewIdent("nil"),
						},
					},
				},
			},
		})
	}

	return ret
}

func generateIntrospectionUnionTypeFuncDecls(unionDefinitions schema.UnionDefinitions, indexes *schema.Indexes) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, u := range unionDefinitions {
		stmts := make([]ast.Stmt, 0)
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("ret"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: ast.NewIdent("__Type"),
				},
			},
		})
		stmts = append(stmts, &ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Tok:   token.DEFINE,
			Value: ast.NewIdent("child"),
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("node"),
				Sel: ast.NewIdent("Children"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					generateIntrospectionUnionSwitchStmt(u, indexes),
				},
			},
		})

		stmts = append(stmts, &ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
				ast.NewIdent("nil"),
			},
		})

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
			Name: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(u.Name))),
			Type: &ast.FuncType{
				Params: generateNodeWalkerArgs(),
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: ast.NewIdent("__Type"),
						},
						{
							Type: ast.NewIdent("error"),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: stmts,
			},
		})
	}

	return ret
}

func generateIntrospectionUnionSwitchStmt(unionDefinition *schema.UnionDefinition, indexes *schema.Indexes) *ast.SwitchStmt {
	caseStmts := make([]ast.Stmt, 0)

	// kind
	caseStmts = append(caseStmts, &ast.CaseClause{
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
					ast.NewIdent("__TypeKind_UNION"),
				},
			},
		},
	})

	// name
	caseStmts = append(caseStmts, &ast.CaseClause{
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
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("NewNullable"),
						},
						Args: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("%q", string(unionDefinition.Name))),
						},
					},
				},
			},
		},
	})

	// description
	caseStmts = append(caseStmts, &ast.CaseClause{
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
					Value: `// TODO`,
				},
			},
		},
	})

	// possibleTypes
	possibleTypesStmts := make([]ast.Stmt, 0)
	possibleTypesStmts = append(possibleTypesStmts, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("possibleTypes"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: ast.NewIdent("make"),
				Args: []ast.Expr{
					&ast.ArrayType{
						Elt: ast.NewIdent("__Type"),
					},
					ast.NewIdent("0"),
					ast.NewIdent(fmt.Sprintf("%d", len(unionDefinition.Types))),
				},
			},
		},
	})
	for _, possibleType := range unionDefinition.Types {
		possibleTypesStmts = append(possibleTypesStmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sRet", string(possibleType))),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", string(possibleType))),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
		possibleTypesStmts = append(possibleTypesStmts, generateReturnErrorHandlingStmt([]ast.Expr{
			&ast.CompositeLit{
				Type: ast.NewIdent("__Type"),
			},
		}))
		possibleTypesStmts = append(possibleTypesStmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("possibleTypes"),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("append"),
					Args: []ast.Expr{
						ast.NewIdent("possibleTypes"),
						ast.NewIdent(fmt.Sprintf("%sRet", string(possibleType))),
					},
				},
			},
		})
	}
	possibleTypesStmts = append(possibleTypesStmts, &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent("ret"),
				Sel: ast.NewIdent("PossibleTypes"),
			},
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("executor"),
					Sel: ast.NewIdent("NewNullable"),
				},
				Args: []ast.Expr{
					ast.NewIdent("possibleTypes"),
				},
			},
		},
	})

	caseStmts = append(caseStmts, &ast.CaseClause{
		List: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: `"possibleTypes"`,
			},
		},
		Body: possibleTypesStmts,
	})

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
			List: caseStmts,
		},
	}
}
