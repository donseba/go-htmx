package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/donseba/go-htmx"
	"github.com/donseba/go-htmx/sse"
)

type (
	App struct {
		htmx *htmx.HTMX
	}
)

var (
	sseManager sse.Manager
)

func main() {
	app := &App{
		htmx: htmx.New(),
	}

	// the pages load htmx 4
	app.htmx.SetVersion(htmx.Version4)

	sseManager = sse.NewManager(5)

	go func() {
		for {
			time.Sleep(1 * time.Second) // Send a message every second
			// htmx 4 swaps the out-of-band element of a message without an event name into the element with its id
			sseManager.Send(sse.NewMessage(fmt.Sprintf(`<div id="time" hx-swap-oob="true">The current time is: %v</div>`, time.Now().Format(time.RFC850))))
		}
	}()

	go func() {
		for {
			clientsStr := ""
			clients := sseManager.Clients()
			for _, c := range clients {
				clientsStr += c + ", "
			}

			time.Sleep(1 * time.Second) // Send a message every seconds
			sseManager.Send(sse.NewMessage(fmt.Sprintf(`<div id="clients" hx-swap-oob="true">connected clients: %v</div>`, clientsStr)))
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
	cl := sse.NewClient(randStringRunes(10))

	sseManager.Handle(w, r, cl)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
