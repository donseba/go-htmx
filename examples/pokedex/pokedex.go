package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

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

	mux.HandleFunc("GET /", app.Home)
	mux.HandleFunc("POST /search", app.Search)

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
