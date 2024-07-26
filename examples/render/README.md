
## Getting Started

* If not already installed, please install the [gonew](https://pkg.go.dev/golang.org/x/tools/cmd/gonew) command.

```console
go install golang.org/x/tools/cmd/gonew@latest
```

* Create a new project using this template.
  - Second argument passed to `gonew` is a module path of your new app.

```console
gonew github.com/donseba/go-htmx/examples/render your.module/my-app # e.g. github.com/donseba/my-app
cd my-app
go mod tidy
go build

```

## Testing 

- Start your app

```console
./my-app
```

- Open your browser http://localhost:3210/
