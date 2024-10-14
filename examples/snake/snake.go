package main

import (
	"context"
	"github.com/donseba/go-htmx"
	"github.com/donseba/go-htmx/sse"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type (
	App struct {
		htmx *htmx.HTMX
		game SnakeGame
	}

	Position struct {
		X, Y int
	}

	SnakeGame struct {
		Board  [20][20]string
		Snake  []Position
		Food   Position
		Dir    Position
		Mu     sync.Mutex
		Active bool
	}
)

var (
	sseManager sse.Manager
)

func main() {

	app := &App{
		htmx: htmx.New(),
		game: SnakeGame{
			Dir:    Position{X: 1, Y: 0},
			Snake:  []Position{{X: 10, Y: 10}, {X: 9, Y: 10}, {X: 8, Y: 10}},
			Active: true,
		},
	}

	placeFood(&app.game)

	sseManager = sse.NewManager(5)

	go func() {
		for {
			app.game.Mu.Lock()
			if app.game.Active {
				moveSnake(&app.game)
			}
			app.game.Mu.Unlock()

			page := htmx.NewComponent("board.gohtml").SetData(map[string]any{
				"game": &app.game,
			})

			out, err := page.Render(context.Background())
			if err != nil {
				log.Printf("error rendering page: %v", err.Error())
			}

			re := regexp.MustCompile(`\s+`)
			stringOut := re.ReplaceAllString(string(out), "")

			sseManager.Send(sse.NewMessage(stringOut).WithEvent("board"))
			time.Sleep(150 * time.Millisecond)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("GET /", http.HandlerFunc(app.Home))
	mux.Handle("GET /sse", http.HandlerFunc(app.SSE))
	mux.Handle("PUT /move/{dir}", http.HandlerFunc(app.Move))
	mux.Handle("PUT /pause", http.HandlerFunc(app.Pause))

	err := http.ListenAndServe(":3210", mux)
	log.Fatal(err)
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	a.game.Mu.Lock()
	defer a.game.Mu.Unlock()

	h := a.htmx.NewHandler(w, r)

	page := htmx.NewComponent("board.gohtml").SetData(map[string]any{
		"game": &a.game,
	}).Wrap(mainContent(), "board")

	_, err := h.Render(r.Context(), page)
	if err != nil {
		log.Printf("error rendering page: %v", err.Error())
	}
}

func (a *App) Move(w http.ResponseWriter, r *http.Request) {
	a.game.Mu.Lock()
	defer a.game.Mu.Unlock()

	dir := r.PathValue("dir")
	switch dir {
	case "up":
		a.game.Dir = Position{X: -1, Y: 0}
	case "down":
		a.game.Dir = Position{X: 1, Y: 0}
	case "left":
		a.game.Dir = Position{X: 0, Y: -1}
	case "right":
		a.game.Dir = Position{X: 0, Y: 1}
	}

	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write(nil)
}

func (a *App) Pause(w http.ResponseWriter, r *http.Request) {
	a.game.Mu.Lock()
	defer a.game.Mu.Unlock()

	a.game.Active = !a.game.Active

	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write(nil)
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

func mainContent() htmx.RenderableComponent {
	return htmx.NewComponent("index.gohtml")
}
