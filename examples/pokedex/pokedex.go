package main

import (
	"encoding/json"
	"github.com/donseba/go-htmx"
	"github.com/donseba/go-htmx/middleware"
	"io"
	"log"
	"net/http"
	"strings"
)

type (
	App struct {
		htmx *htmx.HTMX
	}

	route struct {
		path    string
		handler http.Handler
	}

	PokemonResponse struct {
		Name    string `json:"name"`
		Sprites struct {
			FrontDefault string `json:"front_default"`
		} `json:"sprites"`
	}
)

func main() {
	// new app with htmx instance
	app := &App{
		htmx: htmx.New(),
	}

	mux := http.NewServeMux()

	// wrap routes with the middleware
	wrapRoutes(mux, middleware.MiddleWare, []route{
		{path: "/", handler: http.HandlerFunc(app.Home)},
		{path: "/search", handler: http.HandlerFunc(app.Search)},
	})

	err := http.ListenAndServe(":3210", mux)
	log.Fatal(err)
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func (a *App) Search(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.PostFormValue("pokemon"))
	if query == "" {
		_, _ = w.Write([]byte("Please enter a Pokemon name."))
		return
	}

	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + query)
	if err != nil || resp.StatusCode != 200 {
		_, _ = w.Write([]byte("Pokemon not found."))
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var pokemon PokemonResponse
	_ = json.Unmarshal(body, &pokemon)

	result := `<img src="` + pokemon.Sprites.FrontDefault + `" alt="` + pokemon.Name + `">`
	_, _ = w.Write([]byte(result))
}

// wrapRoutes takes a mux, middleware, and a list of routes to apply the middleware to.
func wrapRoutes(mux *http.ServeMux, mw func(http.Handler) http.Handler, routes []route) {
	for _, r := range routes {
		mux.Handle(r.path, mw(r.handler))
	}
}
