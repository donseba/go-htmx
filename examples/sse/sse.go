package main

import (
	"fmt"
	"github.com/donseba/go-htmx"
	"github.com/donseba/go-htmx/sse"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type (
	App struct {
		htmx *htmx.HTMX
	}

	client struct {
		id string
		ch chan sse.Message
	}
)

func (c *client) ID() string {
	return c.id
}
func (c *client) Chan() chan sse.Message {
	return c.ch
}

var (
	sseManager *sse.Manager
)

func main() {
	app := &App{
		htmx: htmx.New(),
	}

	sseManager = sse.NewManager(5)

	go func() {
		for {
			time.Sleep(1 * time.Second) // Send a message every seconds
			sseManager.Send(sse.NewMessage(fmt.Sprintf("The current time is: %v", time.Now().Format(time.RFC850))).WithEvent("Time"))
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("GET /", http.HandlerFunc(app.Home))
	mux.Handle("GET /sse", http.HandlerFunc(app.SSE))

	err := http.ListenAndServe(":3210", mux)
	log.Fatal(err)
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func (a *App) SSE(w http.ResponseWriter, r *http.Request) {
	cl := &client{
		id: randStringRunes(10),
		ch: make(chan sse.Message),
	}

	sseManager.Handle(w, cl)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
