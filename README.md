## goliteql

goliteql is a lightweight schema first graphql code generator for golang.
It generates a graphql server code based on **http.Handler** interface and a schema files.

### Features

- Lightweight
- Schema first
- Dependency is only stdlib or goliteql

goliteql support for graphql specification on October 2021 Edition.

| Feature        | Status | Note |
|----------------|--------|------|
| Query          | ✅     | - |
| Mutation       | ✅     | - |
| Subscription   | ❌     | Parser supported, execution not implemented |
| Interface      | ❌     | Parser supported, execution not implemented |
| Union          | ❌     | Parser supported, execution not implemented |
| Enum           | ❌     | Parser supported |
| Input          | ✅     | - |
| Scalar         | ❌     | Parser supported (custom scalars unsupported) |
| Directive      | ❌     | Parser supported, directive execution not implemented |
| Fragment       | ❌     | Parser supported (inline fragments, named fragments planned) |
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

Write the following code in `resolver/mutation.resolver.go`:

```golang
package resolver

import (
	"encoding/json"
	"net/http"

	"<your-module-name>/graphql/model"
	"github.com/n9te9/goliteql/executor"
)

type MutationResolver interface {
	CreatePost(w http.ResponseWriter, req *http.Request)
}
// Read request body for CreatePostArgs
// Write response body for Post
func (r *resolver) CreatePost(w http.ResponseWriter, req *http.Request) {
	newPost := new(model.CreatePostArgs)

	if err := json.NewDecoder(req.Body).Decode(newPost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := model.Post{
		Id:   "1",
		Title: newPost.Data.Title,
		Content: newPost.Data.Content,
	}

	resp := executor.GraphQLResponse{
		Data: res,
		Errors: nil,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
```

#### Run

```bash
$ go run main.go
```

#### Sample Request
You can use `curl` to send a request to the server.

```bash
$ curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation ($data: NewPost!) { createPost(data: $data) { id title content } }",
    "variables": {
      "data": {
        "title": "title hoge",
        "content": "content fuga"
      }
    }
  }'
{
  "id": "1",
  "title": "title hoge",
  "content": "content fuga"
}
```

### Benchmark

I compared goliteql with other graphql code generator(gqlgen).
I took benchmark goliteql init schema.

```bash
$ go test -benchmem -v -bench .
go test -benchmem -v -bench .
goos: darwin
goarch: amd64
pkg: github.com/n9te9/goliteql/internal/generated/release_test/resolver
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkResolver
BenchmarkResolver-16    	   61759	     16392 ns/op	   11095 B/op	     173 allocs/op
BenchmarkGqlGen
BenchmarkGqlGen-16      	   32877	     35504 ns/op	   20661 B/op	     370 allocs/op
PASS
ok  	github.com/n9te9/goliteql/internal/generated/release_test/resolver	2.892s
```

goliteql is faster than gqlgen in this case.

### Contribution

If you want to contribute to goliteql, please fork the repository and create a pull request.
I will review your pull request and merge it if it is good.
If you want to add a new feature, please create an issue first and discuss it with me.
If you find a bug, please create an issue and I will fix it as soon as possible.

I welcome any contribution, but please make sure to follow the code style and conventions used in the project!!
