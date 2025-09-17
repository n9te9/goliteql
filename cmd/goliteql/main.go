package main

import (
	"log"
	"os"

	"github.com/n9te9/goliteql/internal/generator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "goliteql",
	Short: "A lightweight GraphQL codegen CLI",
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from GraphQL schema",
	Long:  `Generate code from GraphQL schema`,
	Run: func(cmd *cobra.Command, args []string) {
		yamlFile, err := os.ReadFile("goliteql.yaml")
		if err != nil {
			log.Fatalf("error reading config file: %v", err)
		}

		var config generator.Config
		if err := yaml.Unmarshal(yamlFile, &config); err != nil {
			log.Fatalf("error unmarshalling config file: %v", err)
		}

		g, err := generator.NewGenerator(&config)
		if err != nil {
			log.Fatalf("error creating generator: %v", err)
		}

		if err := g.Generate(); err != nil {
			log.Fatalf("error generating code: %v", err)
		}
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the project",
	Long:  `Initialize the project`,
	Run: func(cmd *cobra.Command, args []string) {
		initializeConfig()
	},
}

func main() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(generateCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error executing command: %v", err)
	}
}

var initConfig = generator.Config{
	SchemaDirectory:             "./graphql/schema",
	ModelOutputFile:             "./graphql/model/models.go",
	ScalarOutputFile:            "./graphql/model/scalar.go",
	QueryResolverOutputFile:     "./graphql/resolver/query.resolver.go",
	MutationResolverOutputFile:  "./graphql/resolver/mutate.resolver.go",
	RootResolverOutputFile:      "./graphql/resolver/resolver.go",
	ResolverGeneratedOutputFile: "./graphql/resolver/generated.go",
	EnumOutputFile:              "./graphql/model/enum.go",
	ModelPackageName:            "example.com/graphql/model",
	ResolverPackageName:         "example.com/graphql/resolver",
	Scalars: []generator.ScalarConfig{
		{
			Name:    "DateTime",
			Package: "time",
			Type:    "time.Time",
		},
	},
}

func initializeConfig() {
	if err := os.MkdirAll(initConfig.SchemaDirectory, 0755); err != nil && !os.IsExist(err) {
		log.Fatalf("error creating schema directory: %v", err)
	}

	schemaFilePath := initConfig.SchemaDirectory + "/schema.graphql"
	schemaFileIO, err := os.Create(schemaFilePath)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("error creating model output file: %v", err)
	}

	schemaFileIO.Write([]byte(schemaFile))
	schemaFileIO.Close()

	f, err := os.Create("./goliteql.yaml")
	if err != nil && os.IsExist(err) {
		log.Fatalf("error creating config file: %v", err)
	}

	if err := yaml.NewEncoder(f).Encode(initConfig); err != nil {
		log.Fatalf("error writing config file: %v", err)
	}
}

const schemaFile = `interface Node {
	id: ID!
}

type Post implements Node {
	id: ID!
	title: String!
	content: String!
	alt: String! @deprecated(reason: "use description")
	description: String
}

type User implements Node {
	id: ID!
	name: String!
	gender: Int
}

input NewPost {
	title: String!
	content: String!
	description: String
}

enum Role {
	ADMIN
	USER
	GUEST @deprecated(reason: "Use USER instead")
}

union SearchResult = Post | User

type Query {
	posts: [Node!]!
	post(id: ID!): Post!
	users: [User!]!
}

type Mutation {
	createPost(data: NewPost!): Post!
}
`
