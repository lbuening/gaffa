# Gaffa

Gaffa (German for Duct Tape) is a dependency injection framework for Go.
So far it is just a stripped down version of [Service Weaver](https://serviceweaver.dev/)'s dependency injection.

## Installation

```bash
go install github.com/lbuening/gaffa/cmd/gaffa@latest
```

`go install` installs the gaffa command to `$GOBIN`, which defaults to `$HOME/go/bin`.
Make sure this directory is included in your PATH.
You can accomplish this, for example, by adding the following to your `.bashrc` and running `source ~/.bashrc`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

### Create a new project

```bash
mkdir myproject
cd myproject
go mod init myproject
go get github.com/lbuening/gaffa
```

### Create main.go

```go
package main

import (
	"context"
	"github.com/lbuening/gaffa"
)

type app struct {
	gaffa.Implements[gaffa.Main]
	exampleService gaffa.Ref[IExampleService]
}

func run(ctx context.Context, a *app) error {
	return a.exampleService.DoSomething(ctx)
}

func main() {
	err := gaffa.Run[app, *app](context.Background(), run)
	if err != nil {
		panic(err)
	}
}
```

### Create example_service.go

```go
package main

import (
    "context"
    "fmt"
    "github.com/lbuening/gaffa"
)

type IExampleService interface {
    DoSomething(ctx context.Context) error
}

type exampleService struct {
    gaffa.Implements[IExampleService]
}

func (s *exampleService) DoSomething(_ context.Context) error {
    fmt.Println("Hello World!")
    return nil
}
```

### Generate code

```bash
gaffa generate ./...
```

### Run

```bash
go run .
```

Output:

```bash
Hello World!
```

## Credits

This project is completely based on [Service Weaver](https://serviceweaver.dev/).
The only difference is that I have just stripped the part out for the dependency injection.
