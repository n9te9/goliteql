## goliteql

goliteql is a lightweight schema first graphql code generator for golang.
It generates a graphql server code based on **http.Handler** interface and a schema files.

### Features

- Lightweight
- Schema first
- Generated code Dependency is only stdlib or goliteql

goliteql will support for graphql specification on October 2021 Edition.
But, it does not support all features of graphql specification.
The following table shows the current status of the features.

| Feature        | Status | Note |
|----------------|--------|------|
| Query          | ✅     | - |
| Mutation       | ✅     | - |
| Subscription   | ❌     | Parser supported, execution not implemented |
| Interface      | ⚙️     | Parser supported, execution is beta |
| Union          | ⚙️     | Parser supported, execution is beta |
| Enum           | ⚙️     | Parser supported, execution is beta |
| Input          | ✅     | - |
| Scalar         | ❌     | Parser supported (custom scalars unsupported) |
| Directive      | ❌     | Parser supported, directive execution not implemented |
| Fragment       | ⚙️     | Parser supported, execution is beta |
| Type           | ✅     | Object type definitions supported |
| extend         | ❌     | Parser supported, merging not yet implemented |
| Federation     | ❌     | Not supported |
| Introspection  | ❌     | Not supported |
| Validation     | ❌     | Parser structures exist, runtime validation WIP |

goliteql is not a full-featured graphql server.
If you want to full-featured graphql server, please use gqlgen.
goliteql will support for graphql specification on October 2021 Edition.

### Getting Started

#### Install

```bash
$ go install github.com/n9te9/goliteql/cmd/goliteql@latest
```

#### Generate

```bash
$ goliteql init
$ go mod init <your-module-name>
$ goliteql generate
$ go mod tidy
```

#### Example

```sh
$ touch main.go
```

Write the following code in `main.go`:

```go
package main

import (
	"net/http"

	"<your-module-name>/graphql/resolver"
)

func main() {
	r := resolver.NewResolver()

	http.ListenAndServe(":8080", r)
}
```

Write the following code in `resolver/query.resolver.go`:

```golang
package resolver

import (
	"context"

	"github.com/n9te9/goliteql/internal/generated/graphql/model"
)

type QueryResolver interface {
	Posts(ctx context.Context) ([]model.Node, error)
	Post(ctx context.Context, id string) (model.Post, error)
	Users(ctx context.Context) ([]model.User, error)
}

func (r *resolver) Posts(ctx context.Context) ([]model.Node, error) {
  return []model.Node{
		&model.Post{
			Id:          "1",
			Title:       "First Post",
			Content:     "This is the content of the first post.",
			Alt:         "First Post Alt",
			Description: nil,
		},
		&model.Post{
			Id:          "2",
			Title:       "Second Post",
			Content:     "This is the content of the second post.",
			Alt:         "Second Post Alt",
			Description: nil,
		},
	}, nil
}

func (r *resolver) Post(ctx context.Context, id string) (model.Post, error) {
	panic("post resolver is not implemented")
}

func (r *resolver) Users(ctx context.Context) ([]model.User, error) {
	panic("users resolver is not implemented")
}
```

#### Run

```bash
$ go run main.go
```

#### Sample Request
You can use `curl` to send a request to the server.

```bash
curl -s localhost:8080 -H "Content-Type: application/json" -d '{"query":"query { posts { ... on Post { id content } } }"}' | jq
```

You should see the following response:

```json
{
	"data": {
		"posts": [
			{
				"id": "1",
				"content": "This is the content of the first post."
			},
			{
				"id": "2",
				"content": "This is the content of the second post."
			}
		]
	}
}

### Benchmark

I compared goliteql with other graphql code generator(gqlgen).
I took benchmark goliteql init schema.
Benchmark code returns ten posts with id and content fields below implementation.

```go
package resolver

import (
	"context"

	"github.com/n9te9/goliteql/internal/generated/graphql/model"
)

type QueryResolver interface {
	Posts(ctx context.Context) ([]model.Node, error)
	Post(ctx context.Context, id string) (model.Post, error)
	Users(ctx context.Context) ([]model.User, error)
}

func (r *resolver) Posts(ctx context.Context) ([]model.Node, error) {
  return []model.Node{
		&model.Post{
			Id:          "1",
			Title:       "First Post",
			Content:     "This is the content of the first post.",
			Alt:         "First Post Alt",
			Description: nil,
		},
		&model.Post{
			Id:          "2",
			Title:       "Second Post",
			Content:     "This is the content of the second post.",
			Alt:         "Second Post Alt",
			Description: nil,
		},
		&model.Post{
			Id:          "3",
			Title:       "Third Post",
			Content:     "This is the content of the third post.",
			Alt:         "Third Post Alt",
			Description: nil,
		},
		&model.Post{
			Id:          "4",
			Title:       "Fourth Post",
			Content:     "This is the content of the fourth post.",
			Alt:         "Fourth Post Alt",
			Description: nil,
		},
		&model.Post{
			Id:          "5",
			Title:       "Fifth Post",
			Content:     "This is the content of the fifth post.",
			Alt:         "Fifth Post Alt",
			Description: nil,
		},
		&model.Post{
			Id:          "6",
			Title:       "Sixth Post",
			Content:     "This is the content of the sixth post.",
			Alt:         "Sixth Post Alt",
			Description: nil,
		},
		&model.Post{
			Id:          "7",
			Title:       "Seventh Post",
			Content:     "This is the content of the seventh post.",
			Alt:         "Seventh Post Alt",
			Description: nil,
		},
		&model.Post{
			Id:          "8",
			Title:       "Eighth Post",
			Content:     "This is the content of the eighth post.",
			Alt:         "Eighth Post Alt",
			Description: nil,
		},
		&model.Post{
			Id:          "9",
			Title:       "Ninth Post",
			Content:     "This is the content of the ninth post.",
			Alt:         "Ninth Post Alt",
			Description: nil,
		},
		&model.Post{
			Id:          "10",
			Title:       "Tenth Post",
			Content:     "This is the content of the tenth post.",
			Alt:         "Tenth Post Alt",
			Description: nil,
		},
	}, nil
}

func (r *resolver) Post(ctx context.Context, id string) (model.Post, error) {
	panic("post resolver is not implemented")
}

func (r *resolver) Users(ctx context.Context) ([]model.User, error) {
	panic("users resolver is not implemented")
}
```

```bash
$ go test -benchmem -v -bench .
goos: darwin
goarch: amd64
pkg: github.com/n9te9/goliteql/internal/benchmark_test
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkGqlgen
BenchmarkGqlgen-16      	   20308	     58628 ns/op	   33502 B/op	     491 allocs/op
BenchmarkGoliteql
BenchmarkGoliteql-16    	   63496	     19096 ns/op	   14299 B/op	     162 allocs/op
PASS
ok  	github.com/n9te9/goliteql/internal/benchmark_test	3.505s
```

### Getting started for parser

If you want to use goliteql as a parser, you can use the `goliteql` package.
Please run go get command to install goliteql package.

```bash
$ go get github.com/n9te9/goliteql
```

You can use the `schema.Parser.Parse` function to parse a schema file.

```golang
package main

import (
	"github.com/n9te9/goliteql/schema"
)

func main() {
	lexer := schema.NewLexer()
	parser := schema.NewParser(lexer)

	ast, err := parser.Parse([]byte(`${Your schema}`))
	if err != nil {
		panic(err)
	}

	// ast is the abstract syntax tree of the schema
}
```

### Contribution

If you want to contribute to goliteql, please fork the repository and create a pull request.
I will review your pull request and merge it if it is good.
If you want to add a new feature, please create an issue first and discuss it with me.
If you find a bug, please create an issue and I will fix it as soon as possible.

I welcome any contribution, but please make sure to follow the code style and conventions used in the project!!
