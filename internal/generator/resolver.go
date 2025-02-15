package generator

import (
	"go/ast"
	"go/token"

	"github.com/lkeix/gg-parser/schema"
)

func generateResolverInterface(query, mutation, subscription *schema.OperationDefinition) *ast.GenDecl {
	queryName := "Query"
	if query != nil {
		if len(query.Name) != 0 {
			queryName = string(query.Name)
		}
	}

	mutationName := "Mutation"
	if mutation != nil {
		if len(mutation.Name) != 0 {
			mutationName = string(mutation.Name)
		}
	}

	subscriptionName := "Subscription"
	if subscription != nil {
		if len(subscription.Name) != 0 {
			subscriptionName = string(subscription.Name)
		}
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
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: queryName,
									},
								},
								Type: &ast.Ident{
									Name: queryName,
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: mutationName,
									},
								},
								Type: &ast.Ident{
									Name: mutationName,
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: subscriptionName,
									},
								},
								Type: &ast.Ident{
									Name: subscriptionName,
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateInterfaceField(operation string, field schema.FieldDefinitions) *ast.FieldList {
	fields := make([]*ast.Field, 0, len(field))

	for _, f := range field {
		fieldType := GraphQLType(f.Type.Name)
		var fieldTypeIdent *ast.Ident
		if fieldType.IsPrimitive() {
			fieldTypeIdent = golangType(f.Type, fieldType)
		} else {
			fieldTypeIdent = golangType(f.Type, fieldType)
		}

		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{
				{
					Name: toUpperCase(string(f.Name)),
				},
			},
			Type: fieldTypeIdent,
		})
	}

	return nil
}

func generateResolver(operationDefinition *schema.OperationDefinition) *ast.File {
	return nil
}