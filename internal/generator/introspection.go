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
													generateStringPointerAST(string(i.Name)),
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
													ast.NewIdent("fields"),
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
													&ast.UnaryExpr{
														Op: token.AND,
														X: &ast.Ident{
															Name: "possibleTypes",
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

	stmts := make([]ast.Stmt, 0)
	stmts = append(stmts, generateIntrospectionTypesFieldSwitchStmts(typeDefinitions.WithoutMetaDefinition(), interfaceDefinitions, inputDefinitions, scalarDefinitions, enumDefinitions)...)
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

func generateIntrospectionTypesFieldSwitchStmts(typeDefinitions []*schema.TypeDefinition, interfaceDefinitions []*schema.InterfaceDefinition, inputDefinitions []*schema.InputDefinition, scalarDefinitions []*schema.ScalarDefinition, enumDefinitions schema.EnumDefinitions) []ast.Stmt {
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
					ast.NewIdent(fmt.Sprintf("%d", len(typeDefinitions)+len(interfaceDefinitions)+len(scalarDefinitions))),
				},
			},
		},
	})

	args := make([]ast.Expr, 0)
	args = append(args, ast.NewIdent("ret"))
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
		if t.IsIntrospection() {
			continue
		}
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
					&ast.UnaryExpr{
						Op: token.AND,
						X:  ast.NewIdent(fmt.Sprintf("interfaces%d", i)),
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
					// TODO
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
					// 					Elt: ast.NewIdent("__InputValue"),
					// 				},
					// 				ast.NewIdent("0"),
					// 				ast.NewIdent(fmt.Sprintf("%d", len(q.Fields))),
					// 			},
					// 		},
					// 	},
					// },
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
		if t.IsIntrospection() {
			continue
		}

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
				&ast.UnaryExpr{
					Op: token.AND,
					X:  ast.NewIdent("ret"),
				},
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
							Type: &ast.StarExpr{
								X: &ast.ArrayType{
									Elt: ast.NewIdent("__Field"),
								},
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
	generateIntrospectionFieldTypeAssignStmtFunc := func(fields schema.FieldDefinitions) []ast.Stmt {
		ret := make([]ast.Stmt, 0, len(fields))

		for _, field := range fields {
			assignStmt := &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("field" + string(field.Name)),
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
			}
			if field.Directives.Get([]byte("deprecated")) != nil {
				ret = append(ret, &ast.IfStmt{
					Cond: ast.NewIdent("includeDeprecated"),
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							assignStmt,
						},
					},
				})
			} else {
				ret = append(ret, assignStmt)
			}
		}

		return ret
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
							Value: `"name"`,
						},
					},
					Body: generateIntrospectionFieldTypeAssignStmtFunc(fieldDefinitions),
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
							Value: `"isDeprecated"`,
						},
					},
					Body: generateIntrospectionIsDeprecatedFieldStmts(fieldDefinitions),
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"deprecationReason"`,
						},
					},
					Body: generateIntrospectionDeprecationReasonFieldStmts(fieldDefinitions),
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"type"`,
						},
					},
					Body: generateIntrospectionFieldTypeBodyStmt(typeName, fieldDefinitions),
				},
			},
		},
	}
}

func generateIntrospectionIsDeprecatedFieldStmts(fieldDefinitions schema.FieldDefinitions) []ast.Stmt {
	ret := make([]ast.Stmt, 0, len(fieldDefinitions))
	for _, fieldDefinition := range fieldDefinitions {
		if fieldDefinition.Directives.Get([]byte("deprecated")) != nil {
			ret = append(ret, &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("field" + string(fieldDefinition.Name)),
						Sel: ast.NewIdent("IsDeprecated"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					ast.NewIdent("true"),
				},
			})
		} else {
			ret = append(ret, &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("field" + string(fieldDefinition.Name)),
						Sel: ast.NewIdent("IsDeprecated"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					ast.NewIdent("false"),
				},
			})
		}
	}

	return ret
}

func generateIntrospectionDeprecationReasonFieldStmts(fieldDefinitions schema.FieldDefinitions) []ast.Stmt {
	ret := make([]ast.Stmt, 0, len(fieldDefinitions))
	for _, fieldDefinition := range fieldDefinitions {
		if fieldDefinition.Directives.Get([]byte("deprecated")) != nil {
			directive := fieldDefinition.Directives.Get([]byte("deprecated"))
			ret = append(ret, &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("field" + string(fieldDefinition.Name)),
						Sel: ast.NewIdent("DeprecationReason"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					generateStringPointerExpr(ast.NewIdent(string(directive.Arguments[0].Value))),
				},
			})
		} else {
			ret = append(ret, &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("field" + string(fieldDefinition.Name)),
						Sel: ast.NewIdent("DeprecationReason"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					ast.NewIdent("nil"),
				},
			})
		}
	}

	return ret
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
						generateIntrospectionFieldSwitchStmt(string(t.Name), t.Fields),
					},
				},
			},
		}

		body = append(body, appendStmts...)
		body = append(body, &ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X:  ast.NewIdent("ret"),
				},
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
							Type: &ast.StarExpr{
								X: &ast.ArrayType{
									Elt: ast.NewIdent("__InputValue"),
								},
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
		if t.IsIntrospection() {
			continue
		}

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
				&ast.UnaryExpr{
					Op: token.AND,
					X:  ast.NewIdent("ret"),
				},
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
							Type: &ast.StarExpr{
								X: &ast.ArrayType{
									Elt: ast.NewIdent("__Field"),
								},
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

func generateIntrospectionTypeOfSwitchStmt(f *introspection.FieldType, callTypeOfFuncName string, indexes *schema.Indexes) ast.Stmt {
	var nameExpr, kindExpr ast.Expr
	var fieldAssignStmt, fieldAssignForPropertyStmt, errHandlingStmt ast.Stmt
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
	fieldAssignForPropertyStmt = &ast.EmptyStmt{}
	errHandlingStmt = &ast.EmptyStmt{}
	var extractArgsField ast.Stmt = &ast.EmptyStmt{}
	if f.IsObjectType() {
		extractArgsField = &ast.AssignStmt{
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
		}

		kindExpr = ast.NewIdent("__TypeKind_OBJECT")
		nameExpr = generateStringPointerAST(string(f.Name))
		fieldAssignStmt = &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("field"),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("__schema__" + string(f.Name) + "__fields"),
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
		}
		fieldAssignForPropertyStmt = &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret"),
					Sel: ast.NewIdent("Fields"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent("field"),
			},
		}
		errHandlingStmt = generateReturnErrorHandlingStmt([]ast.Expr{
			ast.NewIdent("nil"),
		})
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

	var ofTypeCaseStmt []ast.Stmt = []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret"),
					Sel: ast.NewIdent("OfType"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent("nil"),
			},
		},
	}
	if f.Child != nil {
		ofTypeCaseStmt = []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("t"),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent(callTypeOfFuncName),
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
				ast.NewIdent("t"),
			}),
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("ret"),
						Sel: ast.NewIdent("OfType"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					ast.NewIdent("t"),
				},
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
						extractArgsField,
						errHandlingStmt,
						fieldAssignStmt,
						errHandlingStmt,
						fieldAssignForPropertyStmt,
					},
				},
				&ast.CaseClause{
					List: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"interfaces"`,
						},
					},
					Body: generateIntrospectionInterfacesStmts(f, indexes),
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
					Body: ofTypeCaseStmt,
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

func generateIntrospectionInterfacesStmts(fieldType *introspection.FieldType, indexes *schema.Indexes) []ast.Stmt {
	if fieldType == nil || !fieldType.IsObjectType() {
		return []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("ret"),
						Sel: ast.NewIdent("Interfaces"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "nil",
					},
				},
			},
		}
	}

	ret := make([]ast.Stmt, 0)
	typeDefinition := indexes.GetTypeDefinition(string(fieldType.Name))
	if typeDefinition == nil {
		return ret
	}

	for _, interfaceType := range typeDefinition.Interfaces {
		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("interface" + string(interfaceType.Name)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("__schema__" + string(interfaceType.Name) + "__type"),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		})

		var prefixExpr ast.Expr = ast.NewIdent("ret")

		ret = append(ret, generateReturnErrorHandlingStmt([]ast.Expr{
			prefixExpr,
		}))

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("ret"),
						Sel: ast.NewIdent("Interfaces"),
					},
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("append"),
					Args: []ast.Expr{
						&ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("ret"),
								Sel: ast.NewIdent("Interfaces"),
							},
						},
						ast.NewIdent("interface" + string(interfaceType.Name)),
					},
				},
			},
		})
	}

	return ret
}

func generateIntrospectionTypeFieldSwitchStmt(typeName string, f *schema.FieldDefinition, indexes *schema.Indexes) ast.Stmt {
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

	var ofTypeCaseStmt []ast.Stmt = []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret"),
					Sel: ast.NewIdent("OfType"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				ast.NewIdent("nil"),
			},
		},
	}
	if f.Type.IsObject() || f.Type.IsList {
		if typeName == string(f.Type.Name) {
			typeName = ""
		} else {
			typeName = fmt.Sprintf("__%s", typeName)
		}

		ofTypeCaseStmt = []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("t"),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent(fmt.Sprintf("__schema%s__%s__typeof", typeName, string(f.Name))),
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
					&ast.SelectorExpr{
						X:   ast.NewIdent("ret"),
						Sel: ast.NewIdent("OfType"),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					ast.NewIdent("t"),
				},
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
									X:  ast.NewIdent("fields"),
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
					Body: generateIntrospectionInterfacesStmts(introspection.ExpandType(f.Type), indexes),
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
					Body: ofTypeCaseStmt,
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
			Body: generateIntrospectionOperationFieldsAST(s.GetQuery(), string(s.Definition.Query)),
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
			ast.NewIdent("fields"),
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
						Type: &ast.StarExpr{
							X: &ast.ArrayType{
								Elt: ast.NewIdent("__Field"),
							},
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
				&ast.UnaryExpr{
					Op: token.AND,
					X:  ast.NewIdent("ret"),
				},
				ast.NewIdent("nil"),
			},
		},
	}
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

func generateSchemaResponseWrite() ast.Stmt {
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
				&ast.ReturnStmt{},
			},
		},
	}
}

func generateTypeResponseWrite() ast.Stmt {
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
								Type: ast.NewIdent("__TypeResponse"),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: ast.NewIdent("Data"),
										Value: &ast.UnaryExpr{
											Op: token.AND,
											X: &ast.CompositeLit{
												Type: ast.NewIdent("__TypeResponseData"),
												Elts: []ast.Expr{
													&ast.KeyValueExpr{
														Key:   ast.NewIdent("Type"),
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

	ret = append(ret, generateIntrospectionTypeFuncDeclBodySwitchStmt(s.Types, s.Interfaces, s.Inputs, s.Scalars, s.Enums))
	ret = append(ret, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
			ast.NewIdent("nil"),
		},
	})

	return ret
}

func generateIntrospectionTypeFuncDeclBodySwitchStmt(typeDefinitions schema.TypeDefinitions, interfaceDefinitions []*schema.InterfaceDefinition, inputDefinitions []*schema.InputDefinition, scalarDefinitions []*schema.ScalarDefinition, enumDefinitions []*schema.EnumDefinition) ast.Stmt {
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

		for _, t := range typeDefinitions {
			if t.IsIntrospection() {
				continue
			}

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
		if t.IsIntrospection() {
			continue
		}

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
													generateStringPointerAST(string(t.Name)),
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
													ast.NewIdent("fields"),
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
													generateStringPointerAST(string(t.Name)),
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
															generateStringPointerAST(string(t.Name)),
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
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("fields"),
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
						&ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent("fields"),
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
					generateStringPointerAST(operationName),
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
	}

	return ret
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
													generateStringPointerAST(string(t.Name)),
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
				&ast.UnaryExpr{
					Op: token.AND,
					X:  ast.NewIdent("ret"),
				},
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
							Type: &ast.StarExpr{
								X: &ast.ArrayType{
									Elt: ast.NewIdent("__EnumValue"),
								},
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
					ast.NewIdent("nil"),
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
						generateStringPointerExpr(ast.NewIdent(reason)),
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
				ast.NewIdent(fmt.Sprintf(`"%s"`, string(elm.Name))),
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

		if elm.Directives.Get([]byte("deprecated")) != nil {
			nameBodyStmt = &ast.IfStmt{
				Cond: ast.NewIdent("includeDeprecated"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						nameBodyStmt,
					},
				},
			}

			descriptionStmt = &ast.IfStmt{
				Cond: ast.NewIdent("includeDeprecated"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						descriptionStmt,
					},
				},
			}
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
