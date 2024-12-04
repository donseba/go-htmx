package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/donseba/go-htmx"
)

//go:embed home.html
var templates embed.FS

type (
	App struct {
		htmx *htmx.HTMX
	}

	route struct {
		path    string
		handler http.Handler
	}
)

func main() {
	// new app with htmx instance
	app := &App{
		htmx: htmx.New(),
	}

	mux := http.NewServeMux()

	htmx.UseTemplateCache = false

	mux.HandleFunc("/", app.Home)

	err := http.ListenAndServe(":3210", mux)
	log.Fatal(err)
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	h := a.htmx.NewHandler(w, r)

	data := map[string]any{
		"Text": "Welcome to the home page",
	}

	page := htmx.NewComponent("home.html").FS(templates).SetData(data)

	_, err := h.Render(r.Context(), page)
	if err != nil {
		fmt.Printf("error rendering page: %v", err.Error())
	}
}
