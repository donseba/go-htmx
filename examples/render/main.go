package main

import (
	"fmt"
	"github.com/donseba/go-htmx"
	"github.com/donseba/go-htmx/middleware"
	"log"
	"net/http"
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

	// wrap routes with the middleware
	wrapRoutes(mux, middleware.MiddleWare, []route{
		{path: "/", handler: http.HandlerFunc(app.Home)},
		{path: "/child", handler: http.HandlerFunc(app.Child)},
	})

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

// wrapRoutes takes a mux, middleware, and a list of routes to apply the middleware to.
func wrapRoutes(mux *http.ServeMux, mw func(http.Handler) http.Handler, routes []route) {
	for _, r := range routes {
		mux.Handle(r.path, mw(r.handler))
	}
}
