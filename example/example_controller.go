package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	go_htmx "github.com/donseba/go-htmx"
	"github.com/go-chi/chi/v5"
	"golang.org/x/net/websocket"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type (
	Controller struct {
		app *go_htmx.Service
	}
)

func NewController(app *go_htmx.Service) *Controller {
	return &Controller{
		app: app,
	}
}

func (c *Controller) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", c.getHome)
	r.Get("/example", c.getExample)
	r.Get("/greeting", c.getGreeting)
	r.Get("/who-are-you", c.getWhoAreYou)
	r.Post("/pokemon", c.postPokemon)

	r.Get("/echo", c.wsEcho)
	r.Get("/heartbeat", c.wsHeartbeat)

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
		} `json:"ability"`
	} `json:"abilities"`
	Height  int    `json:"height"`
	ID      int    `json:"id"`
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

	if resp.StatusCode != http.StatusOK {
		c.app.Error(w, r, errors.New(fmt.Sprintf("unable to find pokemon : %s", search)), resp.StatusCode)
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

func (c *Controller) wsHeartbeat(w http.ResponseWriter, r *http.Request) {
	handler := websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		for i := 0; ; i = i + 1 {
			time.Sleep(1 * time.Second)

			random := rand.Int()
			message := `<div id="idMessageHeartbeat" hx-swap-oob="true">Message ` + strconv.Itoa(i) + `: ` + strconv.Itoa(random) + `</div>`

			if err := websocket.Message.Send(ws, message); err != nil {
				c.app.Logger().Print("send", err)
				return
			}
		}
	})

	handler.ServeHTTP(w, r)
}

func (c *Controller) wsEcho(w http.ResponseWriter, r *http.Request) {
	handler := websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		for {
			msg := ""

			if err := websocket.Message.Receive(ws, &msg); err != nil {
				c.app.Logger().Print("receive", err)
				return
			}

			response := `<div id="idMessageEcho" hx-swap-oob="true">` + msg + `</div>`

			if err := websocket.Message.Send(ws, response); err != nil {
				c.app.Logger().Print("send", err)
				return
			}
		}
	})

	handler.ServeHTTP(w, r)
}

func (c *Controller) getGreeting(w http.ResponseWriter, r *http.Request) {
	formats := []string{
		"Hi, %s. Welcome!",
		"Great to see you, %s!",
		"Hail, %s! Well met!",
		"Hello, full stack %s!",
	}

	_, _ = w.Write([]byte(fmt.Sprintf(formats[rand.Intn(len(formats))], "gopher")))
}
