package example

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/donseba/go-htmx/pkg/htmx"
	"github.com/go-chi/chi/v5"
)

type (
	Controller struct {
		app *htmx.Service
	}
)

func NewController(app *htmx.Service) *Controller {
	return &Controller{
		app: app,
	}
}

func (c *Controller) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", c.getHome)
	r.Get("/example", c.getExample)
	r.Get("/who-are-you", c.getWhoAreYou)
	r.Post("/pokemon", c.postPokemon)

	return r
}

func (c *Controller) getHome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	td, err := c.app.TemplateData(ctx)
	if err != nil {
		c.app.Error(w, r, err, http.StatusInternalServerError)

	}

	c.app.Serve(td, []string{"home.gohtml"}).ServeHTTP(w, r)
}

func (c *Controller) getExample(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	td, err := c.app.TemplateData(ctx)
	if err != nil {
		c.app.Error(w, r, err, http.StatusInternalServerError)

	}

	c.app.Serve(td, []string{"example.gohtml"}).ServeHTTP(w, r)
}

func (c *Controller) getWhoAreYou(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hxh := c.app.HxHeader(ctx)

	w.WriteHeader(http.StatusOK)
	if hxh.HxPrompt == "" {
		_, _ = w.Write([]byte("hello, dont be shy!"))
		return
	}

	_, _ = w.Write([]byte(fmt.Sprintf("hello, %s!", hxh.HxPrompt)))
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	Height int `json:"height"`
	ID     int `json:"id"`
	Moves  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name    string `json:"name"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
	} `json:"sprites"`

	Weight int `json:"weight"`
}

func (c *Controller) postPokemon(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	search := r.FormValue("search")

	resp, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", search))
	if err != nil {
		c.app.Error(w, r, err, http.StatusFailedDependency)
		return
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		c.app.Error(w, r, errors.New(fmt.Sprintf("unable to find pokemon : %s", search)), http.StatusNotFound)
		return
	}

	var pokemon Pokemon
	err = json.NewDecoder(resp.Body).Decode(&pokemon)
	if err != nil {
		c.app.Error(w, r, err, http.StatusInternalServerError)
		return
	}

	td, err := c.app.TemplateData(ctx)
	if err != nil {
		c.app.Error(w, r, err, http.StatusInternalServerError)
		return
	}

	td.Data["Pokemon"] = pokemon

	c.app.Serve(td, []string{"pokemon.gohtml"}).ServeHTTP(w, r)
}
