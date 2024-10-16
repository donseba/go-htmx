package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/donseba/go-htmx"
)

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
	htmx.UseTemplateCache = false

	// new app with htmx instance
	app := &App{
		htmx: htmx.New(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Index)

	err := http.ListenAndServe(":3210", mux)
	log.Fatal(err)
}

func (a *App) Index(w http.ResponseWriter, r *http.Request) {
	h := a.htmx.NewHandler(w, r)

	data := map[string]any{
		"Text": "Welcome to the home page",
	}

	page := htmx.NewComponent("index.html").
		SetData(data).
		With(htmx.NewComponent("block.html").SetData(map[string]any{"color": randomHexColor()}), "block1").
		With(htmx.NewComponent("block.html").SetData(map[string]any{"color": randomHexColor()}), "block2").
		With(htmx.NewComponent("block.html").SetData(map[string]any{"color": randomHexColor()}), "block3")

	_, err := h.Render(r.Context(), page)
	if err != nil {
		fmt.Printf("error rendering page: %v", err.Error())
	}
}

func randomHexColor() string {
	return fmt.Sprintf("#%06x", rand.Intn(0xffffff))
}
