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

		var config Config
		if err := yaml.Unmarshal(yamlFile, &config); err != nil {
			log.Fatalf("error unmarshalling config file: %v", err)
		}

		generateorConfig := &generator.Config{
			SchemaDirectory:            config.SchemaDirectory,
			ModelOutputFile:            config.ModelOutputFile,
			QueryResolverOutputFile:    config.QueryResolverOutputFile,
			MutationResolverOutputFile: config.MutationResolverOutputFile,
			RootResolverOutputFile:     config.RootResolverOutputFile,
			EnumOutputFile:             config.EnumOutputFile,
			ModelPackageName:           config.ModelPackageName,
			ResolverPackageName:        config.ResolverPackageName,
		}
		g, err := generator.NewGenerator(generateorConfig)
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

type Config struct {
	SchemaDirectory            string `yaml:"schema_directory"`
	ModelOutputFile            string `yaml:"model_output_file"`
	QueryResolverOutputFile    string `yaml:"query_resolver_output_file"`
	MutationResolverOutputFile string `yaml:"mutation_resolver_output_file"`
	RootResolverOutputFile     string `yaml:"root_resolver_output_file"`
	EnumOutputFile             string `yaml:"enum_output_file"`
	ModelPackageName           string `yaml:"model_package_name"`
	ResolverPackageName        string `yaml:"resolver_package_name"`
}

var initConfig = Config{
	SchemaDirectory:            "./graphql/schema",
	ModelOutputFile:            "./graphql/model/models.go",
	QueryResolverOutputFile:    "./graphql/resolver/query.resolver.go",
	MutationResolverOutputFile: "./graphql/resolver/mutate.resolver.go",
	RootResolverOutputFile:     "./graphql/resolver/resolver.go",
	EnumOutputFile:             "./graphql/model/enum.go",
	ModelPackageName:           "example/graphql/model",
	ResolverPackageName:        "example/graphql/resolver",
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

const schemaFile = `type Post {
	id: ID!
	title: String!
	content: String!
	description: String
}

input NewPost {
	title: String!
	content: String!
	description: String
}

type Query {
	posts: [Post!]!
	post(id: ID!): Post!
}

type Mutation {
	createPost(data: NewPost!): Post!
}
`
