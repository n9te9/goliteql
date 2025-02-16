package generator

import (
	"go/ast"
	"go/token"

	"github.com/lkeix/gg-parser/schema"
)

func newQueryName(query *schema.OperationDefinition) *ast.Ident {
	queryName := "Query"
	if query != nil {
		if len(query.Name) != 0 {
			queryName = string(query.Name)
		}
	}

	return ast.NewIdent(queryName)
}

func newMutationName(mutation *schema.OperationDefinition) *ast.Ident {
	mutationName := "Mutation"
	if mutation != nil {
		if len(mutation.Name) != 0 {
			mutationName = string(mutation.Name)
		}
	}

	return ast.NewIdent(mutationName)
}

func newSubscriptionName(subscription *schema.OperationDefinition) *ast.Ident {
	subscriptionName := "Subscription"
	if subscription != nil {
		if len(subscription.Name) != 0 {
			subscriptionName = string(subscription.Name)
		}
	}

	return ast.NewIdent(subscriptionName)
}

func generateResolverStruct(query, mutation, subscription *schema.OperationDefinition) *ast.GenDecl {
	generateField := func(query, mutation, subscription *schema.OperationDefinition) []*ast.Field {
		fields := make([]*ast.Field, 0, 3)
		if query != nil {
			fields = append(fields, &ast.Field{
				Type: newQueryName(query),
			})
		}

		if mutation != nil {
			fields = append(fields, &ast.Field{
				Type: newMutationName(mutation),
			})
		}

		if subscription != nil {
			fields = append(fields, &ast.Field{
				Type: newSubscriptionName(subscription),
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
					Params: &ast.FieldList{
						List: generateInterfaceMethodArgs(f.Arguments, modelPackagePath),
					},
					Results: &ast.FieldList{
						List: generateInterfaceMethodResults(f.Type, modelPackagePath),
					},
				},
			})
		}

		return &ast.FieldList{
			List: fields,
		}
	}

	var ident *ast.Ident
	if operation.OperationType.IsQuery() {
		ident = newQueryName(operation)
	}

	if operation.OperationType.IsMutation() {
		ident = newMutationName(operation)
	}

	if operation.OperationType.IsSubscription() {
		ident = newSubscriptionName(operation)
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

func generateInterfaceMethodArgs(args schema.ArgumentDefinitions, modelPackagePath string) []*ast.Field {
	fields := make([]*ast.Field, 0, len(args))

	for _, arg := range args {
		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{
				{
					Name: string(arg.Name),
				},
			},
			Type: golangType(arg.Type, GraphQLType(arg.Type.Name), modelPackagePath),
		})
	}

	return fields
}

func generateInterfaceMethodResults(fieldType *schema.FieldType, modelPackagePath string) []*ast.Field {
	return []*ast.Field{
		{
			Type: golangType(fieldType, GraphQLType(fieldType.Name), modelPackagePath),
		},
		{
			Type: ast.NewIdent("error"),
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
