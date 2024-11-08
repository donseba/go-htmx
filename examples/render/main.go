package main

import (
	"fmt"
	partial "github.com/donseba/go-htmx/partials"
	"log"
	"net/http"

	"github.com/donseba/go-htmx"
)

type (
	App struct {
		htmx *htmx.HTMX
	}
)

func main() {
	// new app with htmx instance
	app := &App{
		htmx: htmx.New(),
	}

	mux := http.NewServeMux()

	partial.UseTemplateCache = false
	partial.DefaultPartialHeader = "Hx-Target"

	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/child", app.Child)

	err := http.ListenAndServe(":3210", mux)
	log.Fatal(err)
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	h := a.htmx.NewHandler(w, r)

	data := map[string]any{
		"Text": "Welcome to the home page",
	}

	page := partial.NewID("content").Templates("content.html").SetData(data).Wrap(mainContent())

	out, err := page.RenderWithRequest(r.Context(), r)
	if err != nil {
		fmt.Printf("error rendering page: %v\n", err.Error())
	}

	_, _ = h.WriteHTML(out)
}

func (a *App) Child(w http.ResponseWriter, r *http.Request) {
	h := a.htmx.NewHandler(w, r)

	data := map[string]any{
		"Text": "Welcome to the child page",
	}

	page := partial.NewID("content", "content.html").SetData(data).Wrap(mainContent())

	out, err := page.RenderWithRequest(r.Context(), r)
	if err != nil {
		fmt.Printf("error rendering child page: %v\n", err.Error())
	}

	_, _ = h.WriteHTML(out)
}

func mainContent() *partial.Partial {
	menuItems := []struct {
		Name string
		Link string
	}{
		{"Home", "/"},
		{"Child", "/child"},
	}

	data := map[string]any{
		"Title":     "Home",
		"MenuItems": menuItems,
	}

	sidebar := partial.NewID("sidebar", "sidebar.html")
	return partial.New("index.html").SetGlobalData(data).With(sidebar)
}
