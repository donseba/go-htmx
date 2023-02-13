# go-htmx
Fullstack example using [golang](https://go.dev), [htmx](https://htmx.org), [_hyperscript](https://hyperscript.org) & [tailwindcss](https://tailwindcss.com)

# Example
https://user-images.githubusercontent.com/2788923/218283740-b0c3e417-3629-41b3-86bd-e252a6b7f146.mp4

# Getting started
`go get github.com/donseba/go-htmx`

initialise the htmx service like so : 
```go
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"


	"github.com/donseba/go-htmx"
	"github.com/pkg/errors"
)

func main(){
    htmxService, err := htmx.NewService(&htmx.Config{
		ServerAddress:      "localhost:8888",
		TemplateDir:        "templates",
		TemplateFuncs:      nil,
		ErrorTemplate:      filepath.Join("error.gohtml"),
		DefaultTemplates:   []string{filepath.Join("index.gohtml")},
		DefaultTemplatesHx: []string{filepath.Join("hx", "index.gohtml")},
		Logger:             log.New(os.Stdout, "go-htmx | ", 0),
	})
	if err != nil {
		panic(errors.Wrap(err, "error loading .env file"))
	}
}
```
