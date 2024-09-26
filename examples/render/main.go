package main

import (
	"fmt"
	"log"
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
	// new app with htmx instance
	app := &App{
		htmx: htmx.New(),
	}

	mux := http.NewServeMux()

	htmx.UseTemplateCache = false

	http.HandleFunc("/", app.Home)
	http.HandleFunc("/child", app.Child)

	err := http.ListenAndServe(":3210", mux)
	log.Fatal(err)
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	h := a.htmx.NewHandler(w, r)

	data := map[string]any{
		"Text": "Welcome to the home page",
	}

	page := htmx.NewComponent("home.html").SetData(data).Wrap(mainContent(), "Content")

	_, err := h.Render(r.Context(), page)
	if err != nil {
		fmt.Printf("error rendering page: %v", err.Error())
	}
}

func (a *App) Child(w http.ResponseWriter, r *http.Request) {
	h := a.htmx.NewHandler(w, r)

	data := map[string]any{
		"Text": "Welcome to the child page",
	}

	page := htmx.NewComponent("child.html").SetData(data).Wrap(mainContent(), "Content")

	_, err := h.Render(r.Context(), page)
	if err != nil {
		fmt.Printf("error rendering page: %v", err.Error())
	}
}

func mainContent() htmx.RenderableComponent {
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

	sidebar := htmx.NewComponent("sidebar.html")
	return htmx.NewComponent("index.html").SetData(data).With(sidebar, "Sidebar")
}
