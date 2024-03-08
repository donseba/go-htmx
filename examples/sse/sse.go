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
)

var (
	sseManager sse.Manager
)

func main() {
	app := &App{
		htmx: htmx.New(),
	}

	sseManager = sse.NewManager(5)

	go func() {
		for {
			time.Sleep(1 * time.Second) // Send a message every second
			sseManager.Send(sse.NewMessage(fmt.Sprintf("<div>The current time is: %v</div>", time.Now().Format(time.RFC850))).WithEvent("time"))
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
			sseManager.Send(sse.NewMessage(fmt.Sprintf("connected clients: %v", clientsStr)).WithEvent("clients"))
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
